# v2.0.2 (2026-03-11)

## ✨ Features

- feat: Added internationalization (i18n) support. (e7f4c50)

## 🐛 Bug Fixes

- fix: 选择单播优先时部分频道地址没有正确显示 (refs #8) (1e21cc8)
- fix: i18n display issues. (8f6909c)
- fix: 关于系统中补充本项目的仓库链接。 (382f0a6)
- fix: 优化EPG源的节目单展示。 (9baec15)
- fix: 优化直播源频道列表展示和检测逻辑 #7 (d93490b)
- fix: 修改IPTV直播源参数配置时，未同步更新关联的EPG源配置。 (a05719d)
- fix: 优化查看直播源和EPG源的频道列表的排序顺序。 (3054718)
- fix: update web favicon icon (c1cb622)

# v2.0.1 (2026-03-08)

## 🐛 Bug Fixes

- fix: M3U格式直播接口的tvg-name属性优先展示频道别名。 (05b0a9c)
- fix: 优化修改定时检测逻辑。 (c95b9ce)
- fix: 修改优化定时任务执行逻辑。 (f26493f)
- fix: 修改初始化、登录接口问题 (326e413)

# v2.0.0 (2026-03-08)

## ✨ Features

- feat: 补充关于信息 (e1c3505)
- feat: 增加M3U直播源时可自动创建对应的XMLTV格式EPG源 (74cc240)
- feat: 增加直播源有效性探测和过滤功能。 (50e316e)
- feat: 新增/修改网络订阅URL直播源，支持自定义HTTP请求头。 (6e52066)
- feat: 优化直播源、EPG源列表展示，支持数据同步时自动刷新 (94d4dad)
- feat: 增强登录接口的安全性。 (8789863)

## 🐛 Bug Fixes

- fix: 增加版本信息展示。 (871b2f7)
- fix: 优化和修改直播频道检测。 (98bc01e)
- fix: 修改直播源参数展示问题 (9f697d0)
- fix: 增强初始化、登录和修改密码接口的安全性 (0690cd0)
- fix: 修改EPG源管理查看节目单的时间展示问题。 (0f07234)
- fix: 优化日志打印 (b26390e)
- fix: 修改直播接口M3U格式的catchup参数问题。 (8c5d094)
- fix: 将接口返回的错误提示信息统一修改为中文。 (f6dc686)
