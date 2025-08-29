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

import "time"

// StagePrivacyLevel represents the privacy level of a Discord stage instance.
//
// Reference: https://discord.com/developers/docs/resources/stage-instance#stage-instance-object-privacy-level
type StagePrivacyLevel int

const (
	// The Stage instance is visible publicly. (deprecated)
	StagePrivacyLevelPublic StagePrivacyLevel = iota + 1

	// The Stage instance is visible to only guild members.
	StagePrivacyLevelGuildOnly
)

// StageInstance represent a Discord stage instance.
//
// Reference: https://discord.com/developers/docs/resources/stage-instance#stage-instance-object
type StageInstance struct {
	// ID is the stageInstance's unique Discord snowflake ID.
	ID Snowflake `json:"id"`

	// GuildID is the guild id of the associated Stage channel
	GuildID Snowflake `json:"guild_id"`

	// ChannelID is the id of the associated Stage channel
	ChannelID Snowflake `json:"channel_id"`

	// Topic is the topic of the Stage instance (1-120 characters)
	Topic string `json:"topic"`

	// PrivacyLevel is the privacy level of the Stage instance
	PrivacyLevel StagePrivacyLevel `json:"privacy_level"`

	// DiscoverableDisabled is whether or not Stage Discovery is disabled (deprecated)
	DiscoverableDisabled bool `json:"discoverable_disabled"`
}

// CreatedAt returns the time when this stage instance is created.
func (s *StageInstance) CreatedAt() time.Time {
	return s.ID.Timestamp()
}
