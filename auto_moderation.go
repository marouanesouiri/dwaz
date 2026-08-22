/************************************************************************************
 *
 * wumpkit, A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package wumpkit

import (
	"time"
)

// AutoModerationEventType indicates in what event context a rule should be checked.
//
// Reference: https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-event-types
type AutoModerationEventType int

const (
	// when a member sends or edits a message in the guild
	AutoModerationEventTypeMessageSend AutoModerationEventType = 1
	// when a member edits their profile
	AutoModerationEventTypeMemberUpdate AutoModerationEventType = 2
)

// Is returns true if the event type matches the provided one.
func (t AutoModerationEventType) Is(eventType AutoModerationEventType) bool {
	return t == eventType
}

// AutoModerationTriggerType characterizes the type of content which can trigger the rule.
//
// Reference: https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-trigger-types
type AutoModerationTriggerType int

const (
	// check if content contains words from a user defined list of keywords
	AutoModerationTriggerTypeKeyword AutoModerationTriggerType = 1
	// check if content represents generic spam
	AutoModerationTriggerTypeSpam AutoModerationTriggerType = 3
	// check if content contains words from internal pre-defined wordsets
	AutoModerationTriggerTypeKeywordPreset AutoModerationTriggerType = 4
	// check if content contains more unique mentions than allowed
	AutoModerationTriggerTypeMentionSpam AutoModerationTriggerType = 5
	// check if member profile contains words from a user defined list of keywords
	AutoModerationTriggerTypeMemberProfile AutoModerationTriggerType = 6
)

// Is returns true if the trigger type matches the provided one.
func (t AutoModerationTriggerType) Is(triggerType AutoModerationTriggerType) bool {
	return t == triggerType
}

// AutoModerationKeywordPresetType represents the internally pre-defined wordsets.
//
// Reference: https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-keyword-preset-types
type AutoModerationKeywordPresetType int

const (
	// words that may be considered forms of swearing or cursing
	AutoModerationKeywordPresetTypeProfanity AutoModerationKeywordPresetType = 1
	// words that refer to sexually explicit behavior or activity
	AutoModerationKeywordPresetTypeSexualContent AutoModerationKeywordPresetType = 2
	// personal insults or words that may be considered hate speech
	AutoModerationKeywordPresetTypeSlurs AutoModerationKeywordPresetType = 3
)

// Is returns true if the preset type matches the provided one.
func (t AutoModerationKeywordPresetType) Is(presetType AutoModerationKeywordPresetType) bool {
	return t == presetType
}

// AutoModerationTriggerMetadata represents additional data used to determine whether a rule should be triggered.
//
// Reference: https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object-trigger-metadata
type AutoModerationTriggerMetadata struct {
	// keyword_filter: substrings which will be searched for in content (Maximum of 1000)
	// Associated Trigger Types: KEYWORD, MEMBER_PROFILE
	KeywordFilter []string `json:"keyword_filter,omitempty"`

	// regex_patterns: regular expression patterns which will be matched against content (Maximum of 10)
	// Associated Trigger Types: KEYWORD, MEMBER_PROFILE
	RegexPatterns []string `json:"regex_patterns,omitempty"`

	// presets: the internally pre-defined wordsets which will be searched for in content
	// Associated Trigger Types: KEYWORD_PRESET
	Presets []AutoModerationKeywordPresetType `json:"presets,omitempty"`

	// allow_list: substrings which should not trigger the rule (Maximum of 100 or 1000)
	// Associated Trigger Types: KEYWORD, KEYWORD_PRESET, MEMBER_PROFILE
	AllowList []string `json:"allow_list,omitempty"`

	// mention_total_limit: total number of unique role and user mentions allowed per message (Maximum of 50)
	// Associated Trigger Types: MENTION_SPAM
	MentionTotalLimit int `json:"mention_total_limit,omitempty"`

	// mention_raid_protection_enabled: whether to automatically detect mention raids
	// Associated Trigger Types: MENTION_SPAM
	MentionRaidProtectionEnabled bool `json:"mention_raid_protection_enabled,omitempty"`
}

// AutoModerationActionType represents the type of action that will execute whenever a rule is triggered.
//
// Reference: https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-action-object-action-types
type AutoModerationActionType int

const (
	// blocks a member's message and prevents it from being posted
	AutoModerationActionTypeBlockMessage AutoModerationActionType = 1
	// logs user content to a specified channel
	AutoModerationActionTypeSendAlertMessage AutoModerationActionType = 2
	// timeout user for a specified duration
	AutoModerationActionTypeTimeout AutoModerationActionType = 3
	// prevents a member from using text, voice, or other interactions
	AutoModerationActionTypeBlockMemberInteraction AutoModerationActionType = 4
)

// Is returns true if the action type matches the provided one.
func (t AutoModerationActionType) Is(actionType AutoModerationActionType) bool {
	return t == actionType
}

// AutoModerationActionMetadata represents additional data used when an action is executed.
//
// Reference: https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-action-object-action-metadata
type AutoModerationActionMetadata struct {
	// channel to which user content should be logged
	// Associated Action Types: SEND_ALERT_MESSAGE
	ChannelID Snowflake `json:"channel_id,omitempty"`

	// timeout duration in seconds
	// Associated Action Types: TIMEOUT
	DurationSeconds int `json:"duration_seconds,omitempty"`

	// additional explanation that will be shown to members whenever their message is blocked
	// Associated Action Types: BLOCK_MESSAGE
	// Max 150 characters
	CustomMessage string `json:"custom_message,omitempty"`
}

// AutoModerationAction represents an action which will execute whenever a rule is triggered.
//
// Reference: https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-action-object
type AutoModerationAction struct {
	// the type of action
	Type AutoModerationActionType `json:"type"`

	// additional metadata needed during execution for this specific action type
	Metadata *AutoModerationActionMetadata `json:"metadata,omitempty"`
}

// AutoModerationRule represents a rule which allows each guild to set up rules that trigger based on some criteria.
//
// Reference: https://discord.com/developers/docs/resources/auto-moderation#auto-moderation-rule-object
type AutoModerationRule struct {
	// the id of this rule
	ID Snowflake `json:"id"`

	// the id of the guild which this rule belongs to
	GuildID Snowflake `json:"guild_id"`

	// the rule name
	Name string `json:"name"`

	// the user which first created this rule
	CreatorID Snowflake `json:"creator_id"`

	// the rule event type
	EventType AutoModerationEventType `json:"event_type"`

	// the rule trigger type
	TriggerType AutoModerationTriggerType `json:"trigger_type"`

	// the rule trigger metadata
	TriggerMetadata AutoModerationTriggerMetadata `json:"trigger_metadata"`

	// the actions which will execute when the rule is triggered
	Actions []AutoModerationAction `json:"actions"`

	// whether the rule is enabled
	Enabled bool `json:"enabled"`

	// the role ids that should not be affected by the rule (Maximum of 20)
	ExemptRoles []Snowflake `json:"exempt_roles"`

	// the channel ids that should not be affected by the rule (Maximum of 50)
	ExemptChannels []Snowflake `json:"exempt_channels"`
}

// CreatedAt returns the creation timestamp of the auto moderation rule.
func (r *AutoModerationRule) CreatedAt() time.Time {
	return r.ID.Timestamp()
}

// String implements the fmt.Stringer interface.
func (r *AutoModerationRule) String() string {
	return r.Name
}
