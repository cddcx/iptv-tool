/**
 * Prompt Template System
 *
 * Extensible module for building AI prompts. Each template exposes:
 *   id          – unique identifier
 *   name        – human-readable name (i18n key)
 *   description – short description (i18n key)
 *   build(channelNames: string[]) => string
 */

// ---------------------------------------------------------------------------
// Template: Channel Group Rules
// ---------------------------------------------------------------------------

const GROUP_RULES_TEMPLATE = {
  id: 'group_rules',
  name: 'rules.ai_template_group_rules',
  description: 'rules.ai_template_group_rules_desc',

  /**
   * Build a prompt that asks the AI to generate channel grouping rules.
   * Uses Few-shot prompting with concrete input/output examples.
   *
   * @param {string[]} channelNames – deduplicated channel names
   * @returns {string} the full prompt text
   */
  build(channelNames) {
    const channelList = channelNames.join('\n')

    return `You are an IPTV channel management expert. I need you to generate a set of channel grouping rules based on the provided channel name list.

## Task

Analyze all the channel names below, categorize them into appropriate groups, and write regex matching rules for each group.

## Requirements

1. No more than 7 groups in total.
2. Group names must be concise and descriptive.
3. **IMPORTANT: Group names must be in the same language as the channel names.** For example, if the channel names are in Chinese, the group names must also be in Chinese (e.g. 央视, 卫视, 地方, 国际). If the channel names are in English, use English group names.
4. Each group should contain one or more regex rules to match channel names.
5. Try to ensure every channel can be matched by at least one group's rules.
6. Regex patterns should be as concise and accurate as possible.

## Output Format

Return the result strictly in the following JSON format. Do not include any extra text or explanation — only output the raw JSON array:

\`\`\`json
[
  {
    "group_name": "Group Name",
    "rules": [
      { "target": "name", "match_mode": "regex", "pattern": "regex_pattern" }
    ]
  }
]
\`\`\`

## Example

### Input channel list:
CCTV1
CCTV2
CCTV5
湖南卫视
浙江卫视
北京卫视
凤凰中文
凤凰资讯
CGTN
东方卫视

### Output:
\`\`\`json
[
  {
    "group_name": "央视",
    "rules": [
      { "target": "name", "match_mode": "regex", "pattern": "^CCTV" }
    ]
  },
  {
    "group_name": "卫视",
    "rules": [
      { "target": "name", "match_mode": "regex", "pattern": "卫视" }
    ]
  },
  {
    "group_name": "国际",
    "rules": [
      { "target": "name", "match_mode": "regex", "pattern": "凤凰|CGTN" }
    ]
  }
]
\`\`\`

Note: In the example above, since the channel names are in Chinese, the group names are also in Chinese. Always follow this principle.

## Now generate grouping rules for the following channel list:

${channelList}
`
  },
}

// ---------------------------------------------------------------------------
// Registry – add new templates here
// ---------------------------------------------------------------------------

const templates = [GROUP_RULES_TEMPLATE]

/**
 * Get all registered prompt templates.
 * @returns {Array} template objects
 */
export function getPromptTemplates() {
  return templates
}

/**
 * Get a template by id.
 * @param {string} id
 * @returns {object|undefined}
 */
export function getPromptTemplate(id) {
  return templates.find((t) => t.id === id)
}

/**
 * Validate AI-returned JSON against the group rules schema.
 * Returns { valid: boolean, data?: Array, error?: string }.
 *
 * Expected shape:
 *   [ { group_name: string, rules: [ { target, match_mode, pattern } ] } ]
 */
export function validateGroupRulesJSON(jsonString) {
  // Try to extract JSON from markdown code block if present
  let cleanedJson = jsonString.trim()
  const codeBlockMatch = cleanedJson.match(/```(?:json)?\s*\n?([\s\S]*?)\n?\s*```/)
  if (codeBlockMatch) {
    cleanedJson = codeBlockMatch[1].trim()
  }

  let parsed
  try {
    parsed = JSON.parse(cleanedJson)
  } catch {
    return { valid: false, error: 'invalid_json' }
  }

  if (!Array.isArray(parsed) || parsed.length === 0) {
    return { valid: false, error: 'not_array' }
  }

  for (let i = 0; i < parsed.length; i++) {
    const group = parsed[i]

    if (!group.group_name || typeof group.group_name !== 'string') {
      return { valid: false, error: 'missing_group_name', index: i }
    }

    if (!Array.isArray(group.rules) || group.rules.length === 0) {
      return { valid: false, error: 'missing_rules', index: i }
    }

    for (let j = 0; j < group.rules.length; j++) {
      const rule = group.rules[j]
      if (!rule.pattern || typeof rule.pattern !== 'string') {
        return { valid: false, error: 'missing_pattern', groupIndex: i, ruleIndex: j }
      }
      // Normalise: ensure target and match_mode have defaults
      if (!rule.target) rule.target = 'name'
      if (!rule.match_mode) rule.match_mode = 'regex'
    }
  }

  return { valid: true, data: parsed }
}
