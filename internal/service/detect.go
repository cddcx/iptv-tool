package service

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"iptv-tool-v2/internal/model"
)

const (
	DefaultDetectConcurrency = 3
	DefaultDetectTimeout     = 5
)

// DetectService handles channel latency detection using ffmpeg
type DetectService struct {
	dataDir string
}

func NewDetectService(dataDir string) *DetectService {
	return &DetectService{dataDir: dataDir}
}

// GetFFmpegPath returns the path to the ffmpeg executable, or error if not found
func (s *DetectService) GetFFmpegPath() (string, error) {
	name := "ffmpeg"
	if runtime.GOOS == "windows" {
		name = "ffmpeg.exe"
	}
	ffmpegPath := filepath.Join(s.dataDir, "detect", name)

	// Check existence by trying to stat
	cmd := exec.Command(ffmpegPath, "-version")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg 可执行文件未找到或无法运行: %w", err)
	}
	return ffmpegPath, nil
}

// GetFFmpegVersion returns the version string of the installed ffmpeg
func (s *DetectService) GetFFmpegVersion() (string, error) {
	ffmpegPath, err := s.GetFFmpegPath()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, ffmpegPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取 ffmpeg 版本失败: %w", err)
	}

	// Parse first line: "ffmpeg version N-xxxxx-gxxxxxxx Copyright ..."
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "unknown", nil
}

// getDetectConfig reads concurrency and timeout settings from the database
func (s *DetectService) getDetectConfig() (concurrency int, timeout int) {
	concurrency = DefaultDetectConcurrency
	timeout = DefaultDetectTimeout

	var settings []model.SystemSetting
	model.DB.Where("key IN ?", []string{"detect_concurrency", "detect_timeout"}).Find(&settings)

	for _, setting := range settings {
		switch setting.Key {
		case "detect_concurrency":
			if v, err := strconv.Atoi(setting.Value); err == nil && v >= 1 && v <= 10 {
				concurrency = v
			}
		case "detect_timeout":
			if v, err := strconv.Atoi(setting.Value); err == nil && v >= 1 && v <= 30 {
				timeout = v
			}
		}
	}
	return
}

// DetectChannels performs latency detection on all parsed channels for a given source.
// manual=true: fails immediately if the source is syncing.
// manual=false: waits for syncing to finish (up to 10 minutes) before starting detection.
func (s *DetectService) DetectChannels(sourceID uint, manual bool) error {
	var source model.LiveSource
	if err := model.DB.First(&source, sourceID).Error; err != nil {
		return fmt.Errorf("直播源 %d 未找到: %w", sourceID, err)
	}

	if !source.Status {
		return nil // Source is disabled, skip
	}

	// Check if already detecting
	if source.IsDetecting {
		return fmt.Errorf("该直播源正在检测中，请勿重复触发")
	}

	// Check syncing status
	if manual {
		if source.IsSyncing {
			return fmt.Errorf("该直播源正在刷新中，请等待刷新完成后再执行检测")
		}
	} else {
		// Wait for syncing to finish (auto/scheduled mode)
		if err := s.waitForSyncComplete(sourceID, 10*time.Minute); err != nil {
			return err
		}
	}

	// Get ffmpeg path
	ffmpegPath, err := s.GetFFmpegPath()
	if err != nil {
		return err
	}

	// Mark as detecting
	model.DB.Model(&model.LiveSource{}).Where("id = ?", sourceID).Update("is_detecting", true)
	defer func() {
		model.DB.Model(&model.LiveSource{}).Where("id = ?", sourceID).Update("is_detecting", false)
	}()

	// Reset latency and detected_at for all channels in this source before detecting
	model.DB.Model(&model.ParsedChannel{}).Where("source_id = ?", sourceID).Updates(map[string]interface{}{
		"latency":     nil,
		"detected_at": nil,
	})

	// Load channels
	var channels []model.ParsedChannel
	if err := model.DB.Where("source_id = ?", sourceID).Find(&channels).Error; err != nil {
		return fmt.Errorf("加载频道列表失败: %w", err)
	}

	if len(channels) == 0 {
		return nil
	}

	concurrency, timeout := s.getDetectConfig()

	slog.Info("Starting channel detection",
		"source_id", sourceID,
		"channels", len(channels),
		"concurrency", concurrency,
		"timeout_seconds", timeout,
	)

	// Concurrent detection with semaphore
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i := range channels {
		wg.Add(1)
		go func(ch *model.ParsedChannel) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire
			defer func() { <-sem }() // Release

			// For channels with multiple URLs (pipe-separated), test the first one
			testURL := ch.URL
			if idx := strings.Index(testURL, "|"); idx > 0 {
				testURL = strings.TrimSpace(testURL[:idx])
			}

			latency, detectErr := s.detectSingleChannel(ffmpegPath, testURL, timeout)
			now := time.Now()

			if detectErr != nil {
				// Timeout or error
				timeoutVal := -1
				ch.Latency = &timeoutVal
			} else {
				ch.Latency = &latency
			}
			ch.DetectedAt = &now

			// Update single channel result in DB
			model.DB.Model(&model.ParsedChannel{}).Where("id = ?", ch.ID).Updates(map[string]interface{}{
				"latency":     ch.Latency,
				"detected_at": ch.DetectedAt,
			})
		}(&channels[i])
	}

	wg.Wait()

	slog.Info("Channel detection completed", "source_id", sourceID, "channels", len(channels))
	return nil
}

// detectSingleChannel runs ffmpeg to probe a single URL and returns the latency in milliseconds
func (s *DetectService) detectSingleChannel(ffmpegPath string, url string, timeoutSec int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	start := time.Now()

	// ffmpeg -i <url> -t 1 -f null - : try to read 1 second of data and discard
	cmd := exec.CommandContext(ctx, ffmpegPath, "-i", url, "-t", "1", "-f", "null", "-")

	// Combine stdout and stderr (ffmpeg writes progress to stderr)
	err := cmd.Run()

	elapsed := time.Since(start)
	latencyMs := int(elapsed.Milliseconds())

	if ctx.Err() == context.DeadlineExceeded {
		return 0, fmt.Errorf("timeout after %ds", timeoutSec)
	}

	if err != nil {
		// ffmpeg returned non-zero exit code — stream is unreachable or invalid
		return 0, fmt.Errorf("ffmpeg error: %w", err)
	}

	return latencyMs, nil
}

// waitForSyncComplete polls the source's is_syncing status until it becomes false
func (s *DetectService) waitForSyncComplete(sourceID uint, maxWait time.Duration) error {
	deadline := time.Now().Add(maxWait)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		var source model.LiveSource
		if err := model.DB.First(&source, sourceID).Error; err != nil {
			return fmt.Errorf("直播源 %d 未找到: %w", sourceID, err)
		}
		if !source.IsSyncing {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("等待直播源刷新完成超时（超过 %v）", maxWait)
		}
		<-ticker.C
	}
}
