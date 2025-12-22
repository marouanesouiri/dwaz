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

// GuildScheduledEventPrivacyLevel represents Guild Scheduled Event Privacy Levels types.
//
// Reference: https://discord.com/developers/docs/resources/invite#invite-object-invite-types
type GuildScheduledEventPrivacyLevel int

const (
	// The scheduled event is only accessible to guild members.
	GuildScheduledEventPrivacyLevelGuildOnly GuildScheduledEventPrivacyLevel = 2
)

// Is returns true if the guild scheduled event's level matches the provided one.
func (l GuildScheduledEventPrivacyLevel) Is(privacyLevel GuildScheduledEventPrivacyLevel) bool {
	return l == privacyLevel
}

// GuildScheduledEventEntityType represents invite target types.
//
// Reference: https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
type GuildScheduledEventEntityType int

const (
	GuildScheduledEventEntityTypeStageInstance GuildScheduledEventEntityType = 1
	GuildScheduledEventEntityTypeVoice         GuildScheduledEventEntityType = 2
	GuildScheduledEventEntityTypeExternal      GuildScheduledEventEntityType = 3
)

// Is returns true if the fuild scheduled event entity's Type matches the provided one.
func (t GuildScheduledEventEntityType) Is(typ GuildScheduledEventEntityType) bool {
	return t == typ
}

// GuildScheduledEvent is a representation of a scheduled event in a guild.
type GuildScheduledEvent struct {
	// ID is the id of the scheduled event.
	ID Snowflake `json:"id"`

	// GuildID is the guild id which the scheduled event belongs to.
	GuildID Snowflake `json:"guild_id"`

	// ChannelID is the channel id in which the scheduled event will be hosted, or null if scheduled entity type is EXTERNAL
	ChannelID Snowflake `json:"channel_id"`
}

// TODO: continue guild_sheduled_event.go
