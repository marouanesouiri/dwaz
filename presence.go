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

// UserStatus represents a user's status on Discord.
//
// Reference: https://discord.com/developers/docs/events/gateway-events#update-presence-status-types
type UserStatus string

const (
	UserStatusOnline       UserStatus = "online"
	UserStatusDoNotDisturb UserStatus = "dnd"
	UserStatusIdle         UserStatus = "idle"
	UserStatusInvisible    UserStatus = "invisible"
	UserStatusOffline      UserStatus = "offline"
)

// Is returns true if this UserStatus matches the provided one.
func (v UserStatus) Is(other UserStatus) bool {
	return v == other
}

// ClientStatus represents active sessions on a per-platform basis.
//
// Reference: https://discord.com/developers/docs/events/gateway-events#client-status-object
type ClientStatus struct {
	Desktop UserStatus `json:"desktop,omitempty"`
	Mobile  UserStatus `json:"mobile,omitempty"`
	Web     UserStatus `json:"web,omitempty"`
}

// ActivityType represents the type of an activity.
//
// Reference: https://discord.com/developers/docs/events/gateway-events#activity-object-activity-types
type ActivityType int

const (
	ActivityTypeGame      ActivityType = 0
	ActivityTypeStreaming ActivityType = 1
	ActivityTypeListening ActivityType = 2
	ActivityTypeWatching  ActivityType = 3
	ActivityTypeCustom    ActivityType = 4
	ActivityTypeCompeting ActivityType = 5
)

// Is returns true if this ActivityType matches the provided one.
func (v ActivityType) Is(other ActivityType) bool {
	return v == other
}

// ActivityFlags represents flags for an activity.
//
// Reference: https://discord.com/developers/docs/events/gateway-events#activity-object-activity-flags
type ActivityFlags int

const (
	ActivityFlagInstance                 ActivityFlags = 1 << 0
	ActivityFlagJoin                     ActivityFlags = 1 << 1
	ActivityFlagSpectate                 ActivityFlags = 1 << 2
	ActivityFlagJoinRequest              ActivityFlags = 1 << 3
	ActivityFlagSync                     ActivityFlags = 1 << 4
	ActivityFlagPlay                     ActivityFlags = 1 << 5
	ActivityFlagPartyPrivacyFriends      ActivityFlags = 1 << 6
	ActivityFlagPartyPrivacyVoiceChannel ActivityFlags = 1 << 7
	ActivityFlagEmbedded                 ActivityFlags = 1 << 8
)

// Has returns true if the flags contain all of the provided flags.
func (f ActivityFlags) Has(flags ...ActivityFlags) bool {
	return BitFieldHas(f, flags...)
}

// ActivityTimestamps represents unix timestamps for start and/or end of the activity.
type ActivityTimestamps struct {
	Start int64 `json:"start,omitempty"`
	End   int64 `json:"end,omitempty"`
}

// ActivityParty represents party information for the activity.
type ActivityParty struct {
	ID   string `json:"id,omitempty"`
	Size [2]int `json:"size,omitempty"`
}

// ActivityAssets represents images and text for the presence.
type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

// ActivitySecrets represents secrets for Rich Presence joining and spectating.
type ActivitySecrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

// ActivityButton represents a custom button in a Rich Presence.
type ActivityButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// ActivityEmoji represents an emoji used in a custom status activity.
type ActivityEmoji struct {
	Name     string    `json:"name"`
	ID       Snowflake `json:"id,omitempty"`
	Animated bool      `json:"animated,omitempty"`
}

// Activity represents a Discord user activity/presence object.
//
// Reference: https://discord.com/developers/docs/events/gateway-events#activity-object
type Activity struct {
	Name          string              `json:"name"`
	Type          ActivityType        `json:"type"`
	URL           string              `json:"url,omitempty"`
	CreatedAt     int64               `json:"created_at,omitempty"`
	Timestamps    *ActivityTimestamps `json:"timestamps,omitempty"`
	ApplicationID Snowflake           `json:"application_id,omitempty"`
	Details       string              `json:"details,omitempty"`
	State         string              `json:"state,omitempty"`
	Emoji         *ActivityEmoji      `json:"emoji,omitempty"`
	Party         *ActivityParty      `json:"party,omitempty"`
	Assets        *ActivityAssets     `json:"assets,omitempty"`
	Secrets       *ActivitySecrets    `json:"secrets,omitempty"`
	Instance      bool                `json:"instance,omitempty"`
	Flags         ActivityFlags       `json:"flags,omitempty"`
	Buttons       []ActivityButton    `json:"buttons,omitempty"`
}

// GatewayPresence represents a presence update sent over the Gateway.
//
// Reference: https://discord.com/developers/docs/events/gateway-events#presence-update-presence-update-event-fields
type GatewayPresence struct {
	User         User         `json:"user"`
	GuildID      Snowflake    `json:"guild_id"`
	Status       UserStatus   `json:"status"`
	Activities   []Activity   `json:"activities"`
	ClientStatus ClientStatus `json:"client_status"`
}
