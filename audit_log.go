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

// AuditLogEvent represents the type of action that triggered an audit log entry.
//
// Reference: https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-audit-log-events
type AuditLogEvent int

const (
	AuditLogEventGuildUpdate                        AuditLogEvent = 1
	AuditLogEventChannelCreate                      AuditLogEvent = 10
	AuditLogEventChannelUpdate                      AuditLogEvent = 11
	AuditLogEventChannelDelete                      AuditLogEvent = 12
	AuditLogEventChannelOverwriteCreate             AuditLogEvent = 13
	AuditLogEventChannelOverwriteUpdate             AuditLogEvent = 14
	AuditLogEventChannelOverwriteDelete             AuditLogEvent = 15
	AuditLogEventMemberKick                         AuditLogEvent = 20
	AuditLogEventMemberPrune                        AuditLogEvent = 21
	AuditLogEventMemberBanAdd                       AuditLogEvent = 22
	AuditLogEventMemberBanRemove                    AuditLogEvent = 23
	AuditLogEventMemberUpdate                       AuditLogEvent = 24
	AuditLogEventMemberRoleUpdate                   AuditLogEvent = 25
	AuditLogEventMemberMove                         AuditLogEvent = 26
	AuditLogEventMemberDisconnect                   AuditLogEvent = 27
	AuditLogEventBotAdd                             AuditLogEvent = 28
	AuditLogEventRoleCreate                         AuditLogEvent = 30
	AuditLogEventRoleUpdate                         AuditLogEvent = 31
	AuditLogEventRoleDelete                         AuditLogEvent = 32
	AuditLogEventInviteCreate                       AuditLogEvent = 40
	AuditLogEventInviteUpdate                       AuditLogEvent = 41
	AuditLogEventInviteDelete                       AuditLogEvent = 42
	AuditLogEventWebhookCreate                      AuditLogEvent = 50
	AuditLogEventWebhookUpdate                      AuditLogEvent = 51
	AuditLogEventWebhookDelete                      AuditLogEvent = 52
	AuditLogEventEmojiCreate                        AuditLogEvent = 60
	AuditLogEventEmojiUpdate                        AuditLogEvent = 61
	AuditLogEventEmojiDelete                        AuditLogEvent = 62
	AuditLogEventMessageDelete                      AuditLogEvent = 72
	AuditLogEventMessageBulkDelete                  AuditLogEvent = 73
	AuditLogEventMessagePin                         AuditLogEvent = 74
	AuditLogEventMessageUnpin                       AuditLogEvent = 75
	AuditLogEventIntegrationCreate                  AuditLogEvent = 80
	AuditLogEventIntegrationUpdate                  AuditLogEvent = 81
	AuditLogEventIntegrationDelete                  AuditLogEvent = 82
	AuditLogEventStageInstanceCreate                AuditLogEvent = 83
	AuditLogEventStageInstanceUpdate                AuditLogEvent = 84
	AuditLogEventStageInstanceDelete                AuditLogEvent = 85
	AuditLogEventStickerCreate                      AuditLogEvent = 90
	AuditLogEventStickerUpdate                      AuditLogEvent = 91
	AuditLogEventStickerDelete                      AuditLogEvent = 92
	AuditLogEventGuildScheduledEventCreate          AuditLogEvent = 100
	AuditLogEventGuildScheduledEventUpdate          AuditLogEvent = 101
	AuditLogEventGuildScheduledEventDelete          AuditLogEvent = 102
	AuditLogEventThreadCreate                       AuditLogEvent = 110
	AuditLogEventThreadUpdate                       AuditLogEvent = 111
	AuditLogEventThreadDelete                       AuditLogEvent = 112
	AuditLogEventApplicationCommandPermissionUpdate AuditLogEvent = 121
	AuditLogEventSoundboardSoundCreate              AuditLogEvent = 130
	AuditLogEventSoundboardSoundUpdate              AuditLogEvent = 131
	AuditLogEventSoundboardSoundDelete              AuditLogEvent = 132
	AuditLogEventAutoModerationRuleCreate           AuditLogEvent = 140
	AuditLogEventAutoModerationRuleUpdate           AuditLogEvent = 141
	AuditLogEventAutoModerationRuleDelete           AuditLogEvent = 142
	AuditLogEventAutoModerationBlockMessage         AuditLogEvent = 143
	AuditLogEventAutoModerationFlagToChannel        AuditLogEvent = 144
	AuditLogEventAutoModerationUserCommunicationDis AuditLogEvent = 145
	AuditLogEventCreatorMonetizationRequestCreated  AuditLogEvent = 150
	AuditLogEventCreatorMonetizationTermsAccepted   AuditLogEvent = 151
	AuditLogEventOnboardingPromptCreate             AuditLogEvent = 163
	AuditLogEventOnboardingPromptUpdate             AuditLogEvent = 164
	AuditLogEventOnboardingPromptDelete             AuditLogEvent = 165
	AuditLogEventOnboardingCreate                   AuditLogEvent = 166
	AuditLogEventOnboardingUpdate                   AuditLogEvent = 167
	AuditLogEventHomeSettingsCreate                 AuditLogEvent = 190
	AuditLogEventHomeSettingsUpdate                 AuditLogEvent = 191
	AuditLogEventVoiceChannelStatusCreate           AuditLogEvent = 192
	AuditLogEventVoiceChannelStatusDelete           AuditLogEvent = 193
)

// Is returns true if this AuditLogEvent matches the provided one.
func (v AuditLogEvent) Is(other AuditLogEvent) bool {
	return v == other
}

// AuditLogChange represents a single change in an audit log entry.
//
// Reference: https://discord.com/developers/docs/resources/audit-log#audit-log-change-object
type AuditLogChange struct {
	// NewValue is the new value of the key.
	NewValue any `json:"new_value,omitempty"`

	// OldValue is the old value of the key.
	OldValue any `json:"old_value,omitempty"`

	// Key is the name of audit log change key.
	Key string `json:"key"`
}

// OptionalAuditEntryInfo contains additional audit log options depending on the action type.
//
// Reference: https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object-optional-audit-entry-info
type OptionalAuditEntryInfo struct {
	// ApplicationID is the ID of the app whose permissions were targeted.
	ApplicationID Snowflake `json:"application_id,omitempty"`

	// AutoModerationRuleName is name of the AutoMod rule that was triggered.
	AutoModerationRuleName string `json:"auto_moderation_rule_name,omitempty"`

	// AutoModerationRuleTriggerType is trigger type of the AutoMod rule.
	AutoModerationRuleTriggerType string `json:"auto_moderation_rule_trigger_type,omitempty"`

	// ChannelID is the channel in which the entities were targeted.
	ChannelID Snowflake `json:"channel_id,omitempty"`

	// Count is the number of entities targeted.
	Count string `json:"count,omitempty"`

	// DeleteMemberDays is number of days of inactivity for prune.
	DeleteMemberDays string `json:"delete_member_days,omitempty"`

	// ID is ID of the overwritten entity.
	ID Snowflake `json:"id,omitempty"`

	// MembersRemoved is number of members removed by the prune.
	MembersRemoved string `json:"members_removed,omitempty"`

	// MessageID is ID of the message that was targeted.
	MessageID Snowflake `json:"message_id,omitempty"`

	// RoleName is name of the role if type was "0" (role).
	RoleName string `json:"role_name,omitempty"`

	// Type is type of overwritten entity - role ("0") or member ("1").
	Type string `json:"type,omitempty"`

	// Status is the status of the stage instance.
	Status string `json:"status,omitempty"`
}

// AuditLogEntry represents a single audit log entry.
//
// Reference: https://discord.com/developers/docs/resources/audit-log#audit-log-entry-object
type AuditLogEntry struct {
	// TargetID is the ID of the affected entity (webhook, user, role, etc.).
	TargetID Snowflake `json:"target_id,omitempty"`

	// Changes is changes made to the target_id.
	Changes []AuditLogChange `json:"changes,omitempty"`

	// UserID is user or app that made the changes.
	UserID Snowflake `json:"user_id,omitempty"`

	// ID is the unique ID of the entry.
	ID Snowflake `json:"id"`

	// ActionType is the type of action that occurred.
	ActionType AuditLogEvent `json:"action_type"`

	// Options contains additional information for certain action types.
	Options *OptionalAuditEntryInfo `json:"options,omitempty"`

	// Reason is the reason for the change (1-512 characters).
	Reason string `json:"reason,omitempty"`
}

// AuditLog represents a Discord Guild Audit Log.
//
// Reference: https://discord.com/developers/docs/resources/audit-log#audit-log-object
type AuditLog struct {
	// ApplicationCommands is list of application commands referenced in the audit log.
	ApplicationCommands []ApplicationCommand `json:"application_commands"`

	// AuditLogEntries is list of audit log entries, sorted from most to least recent.
	AuditLogEntries []AuditLogEntry `json:"audit_log_entries"`

	// AutoModerationRules is list of auto moderation rules referenced in the audit log.
	AutoModerationRules []AutoModerationRule `json:"auto_moderation_rules"`

	// GuildScheduledEvents is list of guild scheduled events referenced in the audit log.
	GuildScheduledEvents []GuildScheduledEvent `json:"guild_scheduled_events"`

	// Integrations is list of partial integration objects.
	Integrations []Integration `json:"integrations"`

	// Threads is list of threads referenced in the audit log.
	Threads []Channel `json:"threads"`

	// Users is list of users referenced in the audit log.
	Users []User `json:"users"`

	// Webhooks is list of webhooks referenced in the audit log.
	Webhooks []Webhook `json:"webhooks"`
}
