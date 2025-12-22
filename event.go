/************************************************************************************
 *
 * dwaz (Discord Wrapper API for Zwafriya), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package dwaz

import "encoding/json"

// ReadyCreateEvent Shard is ready
type ReadyEvent struct {
	ShardsID int // shard that dispatched this event
	Guilds   []Guild
}

// GuildCreateEvent Guild was created
type GuildCreateEvent struct {
	ShardsID int // shard that dispatched this event
	Guild    GatewayGuild
}

// MessageCreateEvent Message was created
type MessageCreateEvent struct {
	ShardsID int // shard that dispatched this event
	Message Message
}

// MessageCreateEvent Message was created
type MessageUpdateEvent struct {
	ShardsID   int // shard that dispatched this event
	OldMessage Message
	NewMessage Message
}

// MessageDeleteEvent Message was deleted
type MessageDeleteEvent struct {
	ShardsID int // shard that dispatched this event
	Message  Message
}

// InteractionCreateEvent Interaction created
type InteractionCreateEvent struct {
	ShardsID    int // shard that dispatched this event
	Interaction Interaction
}

var _ json.Unmarshaler = (*InteractionCreateEvent)(nil)

// UnmarshalJSON implements json.Unmarshaler for InteractionCreateEvent.
func (c *InteractionCreateEvent) UnmarshalJSON(buf []byte) error {
	interaction, err := UnmarshalInteraction(buf)
	if err == nil {
		c.Interaction = interaction
	}
	return err
}

// VoiceStateUpdateEvent VoiceState was updated
type VoiceStateUpdateEvent struct {
	ShardsID int // shard that dispatched this event
	OldState GatewayVoiceState
	NewState GatewayVoiceState
}

// ApplicationCommandPermissionsUpdateEvent Application command permission was updated
type ApplicationCommandPermissionsUpdateEvent struct {
	// TODO: complete this struct
}

// AutoModerationRuleCreateEvent Auto Moderation rule was created
type AutoModerationRuleCreateEvent struct {
	// TODO: complete this struct
}

// AutoModerationRuleUpdateEvent Auto Moderation rule was updated
type AutoModerationRuleUpdateEvent struct {
	// TODO: complete this struct
}

// AutoModerationRuleDeleteEvent Auto Moderation rule was deleted
type AutoModerationRuleDeleteEvent struct {
	// TODO: complete this struct
}

// AutoModerationActionExecutionEvent Auto Moderation rule was triggered and an action was executed
type AutoModerationActionExecutionEvent struct {
	// TODO: complete this struct
}

// ChannelCreateEvent New guild channel created
type ChannelCreateEvent struct {
	// TODO: complete this struct
}

// ChannelUpdateEvent Channel was updated
type ChannelUpdateEvent struct {
	// TODO: complete this struct
}

// ChannelDeleteEvent Channel was deleted
type ChannelDeleteEvent struct {
	// TODO: complete this struct
}

// ChannelPinsUpdateEvent Message was pinned or unpinned
type ChannelPinsUpdateEvent struct {
	// TODO: complete this struct
}

// ThreadCreateEvent Thread created
type ThreadCreateEvent struct {
	// TODO: complete this struct
}

// ThreadUpdateEvent Thread was updated
type ThreadUpdateEvent struct {
	// TODO: complete this struct
}

// ThreadDeleteEvent Thread was deleted
type ThreadDeleteEvent struct {
	// TODO: complete this struct
}

// ThreadListSyncEvent Sent when gaining access to a channel, contains all active threads
type ThreadListSyncEvent struct {
	// TODO: complete this struct
}

// ThreadMemberUpdateEvent Thread member for the current user was updated
type ThreadMemberUpdateEvent struct {
	// TODO: complete this struct
}

// ThreadMembersUpdateEvent Some user(s) were added to or removed from a thread
type ThreadMembersUpdateEvent struct {
	// TODO: complete this struct
}

// EntitlementCreateEvent Entitlement was created
type EntitlementCreateEvent struct {
	// TODO: complete this struct
}

// EntitlementUpdateEvent Entitlement was updated or renewed
type EntitlementUpdateEvent struct {
	// TODO: complete this struct
}

// EntitlementDeleteEvent Entitlement was deleted
type EntitlementDeleteEvent struct {
	// TODO: complete this struct
}

// GuildUpdateEvent Guild was updated
type GuildUpdateEvent struct {
	// TODO: complete this struct
}

// GuildDeleteEvent Guild became unavailable, or user left/was removed from a guild
type GuildDeleteEvent struct {
	// TODO: complete this struct
}

// GuildAuditLogEntryCreateEvent A guild audit log entry was created
type GuildAuditLogEntryCreateEvent struct {
	// TODO: complete this struct
}

// GuildBanAddEvent User was banned from a guild
type GuildBanAddEvent struct {
	// TODO: complete this struct
}

// GuildBanRemoveEvent User was unbanned from a guild
type GuildBanRemoveEvent struct {
	// TODO: complete this struct
}

// GuildEmojisUpdateEvent Guild emojis were updated
type GuildEmojisUpdateEvent struct {
	// TODO: complete this struct
}

// GuildStickersUpdateEvent Guild stickers were updated
type GuildStickersUpdateEvent struct {
	// TODO: complete this struct
}

// GuildIntegrationsUpdateEvent Guild integration was updated
type GuildIntegrationsUpdateEvent struct {
	// TODO: complete this struct
}

// GuildMemberAddEvent New user joined a guild
type GuildMemberAddEvent struct {
	// TODO: complete this struct
}

// GuildMemberRemoveEvent User was removed from a guild
type GuildMemberRemoveEvent struct {
	// TODO: complete this struct
}

// GuildMemberUpdateEvent Guild member was updated
type GuildMemberUpdateEvent struct {
	// TODO: complete this struct
}

// GuildMembersChunkEvent Response to Request Guild Members
type GuildMembersChunkEvent struct {
	// TODO: complete this struct
}

// GuildRoleCreateEvent Guild role was created
type GuildRoleCreateEvent struct {
	// TODO: complete this struct
}

// GuildRoleUpdateEvent Guild role was updated
type GuildRoleUpdateEvent struct {
	// TODO: complete this struct
}

// GuildRoleDeleteEvent Guild role was deleted
type GuildRoleDeleteEvent struct {
	// TODO: complete this struct
}

// GuildScheduledEventCreateEvent Guild scheduled event was created
type GuildScheduledEventCreateEvent struct {
	// TODO: complete this struct
}

// GuildScheduledEventUpdateEvent Guild scheduled event was updated
type GuildScheduledEventUpdateEvent struct {
	// TODO: complete this struct
}

// GuildScheduledEventDeleteEvent Guild scheduled event was deleted
type GuildScheduledEventDeleteEvent struct {
	// TODO: complete this struct
}

// GuildScheduledEventUserAddEvent User subscribed to a guild scheduled event
type GuildScheduledEventUserAddEvent struct {
	// TODO: complete this struct
}

// GuildScheduledEventUserRemoveEvent User unsubscribed from a guild scheduled event
type GuildScheduledEventUserRemoveEvent struct {
	// TODO: complete this struct
}

// GuildSoundboardSoundCreateEvent Guild soundboard sound was created
type GuildSoundboardSoundCreateEvent struct {
	// TODO: complete this struct
}

// GuildSoundboardSoundUpdateEvent Guild soundboard sound was updated
type GuildSoundboardSoundUpdateEvent struct {
	// TODO: complete this struct
}

// GuildSoundboardSoundDeleteEvent Guild soundboard sound was deleted
type GuildSoundboardSoundDeleteEvent struct {
	// TODO: complete this struct
}

// GuildSoundboardSoundsUpdateEvent Guild soundboard sounds were updated
type GuildSoundboardSoundsUpdateEvent struct {
	// TODO: complete this struct
}

// SoundboardSoundsEvent Response to Request Soundboard Sounds
type SoundboardSoundsEvent struct {
	// TODO: complete this struct
}

// IntegrationCreateEvent Guild integration was created
type IntegrationCreateEvent struct {
	// TODO: complete this struct
}

// IntegrationUpdateEvent Guild integration was updated
type IntegrationUpdateEvent struct {
	// TODO: complete this struct
}

// IntegrationDeleteEvent Guild integration was deleted
type IntegrationDeleteEvent struct {
	// TODO: complete this struct
}

// InviteCreateEvent Invite to a channel was created
type InviteCreateEvent struct {
	// TODO: complete this struct
}

// InviteDeleteEvent Invite to a channel was deleted
type InviteDeleteEvent struct {
	// TODO: complete this struct
}

// MessageDeleteBulkEvent Multiple messages were deleted at once
type MessageDeleteBulkEvent struct {
	// TODO: complete this struct
}

// MessageReactionAddEvent User reacted to a message
type MessageReactionAddEvent struct {
	// TODO: complete this struct
}

// MessageReactionRemoveEvent User removed a reaction from a message
type MessageReactionRemoveEvent struct {
	// TODO: complete this struct
}

// MessageReactionRemoveAllEvent All reactions were explicitly removed from a message
type MessageReactionRemoveAllEvent struct {
	// TODO: complete this struct
}

// MessageReactionRemoveEmojiEvent All reactions for a given emoji were explicitly removed from a message
type MessageReactionRemoveEmojiEvent struct {
	// TODO: complete this struct
}

// PresenceUpdateEvent User was updated
type PresenceUpdateEvent struct {
	// TODO: complete this struct
}

// StageInstanceCreateEvent Stage instance was created
type StageInstanceCreateEvent struct {
	// TODO: complete this struct
}

// StageInstanceUpdateEvent Stage instance was updated
type StageInstanceUpdateEvent struct {
	// TODO: complete this struct
}

// StageInstanceDeleteEvent Stage instance was deleted or closed
type StageInstanceDeleteEvent struct {
	// TODO: complete this struct
}

// SubscriptionCreateEvent Premium App Subscription was created
type SubscriptionCreateEvent struct {
	// TODO: complete this struct
}

// SubscriptionUpdateEvent Premium App Subscription was updated
type SubscriptionUpdateEvent struct {
	// TODO: complete this struct
}

// SubscriptionDeleteEvent Premium App Subscription was deleted
type SubscriptionDeleteEvent struct {
	// TODO: complete this struct
}

// TypingStartEvent User started typing in a channel
type TypingStartEvent struct {
	// TODO: complete this struct
}

// UserUpdateEvent Properties about the user changed
type UserUpdateEvent struct {
	// TODO: complete this struct
}

// VoiceChannelEffectSendEvent Someone sent an effect in a voice channel
type VoiceChannelEffectSendEvent struct {
	// TODO: complete this struct
}

// VoiceServerUpdateEvent Guild's voice server was updated
type VoiceServerUpdateEvent struct {
	// TODO: complete this struct
}

// WebhooksUpdateEvent Guild channel webhook was created, update, or deleted
type WebhooksUpdateEvent struct {
	// TODO: complete this struct
}

// MessagePollVoteAddEvent User voted on a poll
type MessagePollVoteAddEvent struct {
	// TODO: complete this struct
}

// MessagePollVoteRemoveEvent User removed a vote on a poll
type MessagePollVoteRemoveEvent struct {
	// TODO: complete this struct
}
