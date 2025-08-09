/************************************************************************************
 *
 * yada (yet another discord api), A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package yada

import "time"

// MemberFlags represents flags of a guild member.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-member-object-guild-member-flags
type MemberFlags int

const (
	// Member has left and rejoined the guild
	//  - Editable: false
	MemberFlagDidRejoin MemberFlags = 1 << iota

	// Member has completed onboarding
	//  - Editable: false
	MemberFlagCompletedOnboarding

	// Member is exempt from guild verification requirements
	//  - Editable: true
	MemberFlagBypassesVerification

	// Member has started onboarding
	//  - Editable: false
	MemberFlagStartedOnboarding

	// Member is a guest and can only access the voice channel they were invited to
	//  - Editable: false
	MemberFlagIsGuest

	// Member has started Server Guide new member actions
	//  - Editable: false
	MemberFlagStartedHomeActions

	// Member has completed Server Guide new member actions
	//  - Editable: false
	MemberFlagCompletedHomeActions

	// Member's username, display name, or nickname is blocked by AutoMod
	//  - Editable: false
	MemberFlagQuarantinedUsername
	_

	// Member has dismissed the DM settings upsell
	//  - Editable: false
	MemberFlagDMSettingsUpsellAcknowledged

	// Member's guild tag is blocked by AutoMod
	//  - Editable: false
	MemberFlagQuarantinedGuildTag
)

// Has returns true if all provided flags are set.
func (f MemberFlags) Has(flags ...MemberFlags) bool {
	for _, flag := range flags {
		if f&flag != flag {
			return false
		}
	}
	return true
}

// Member is a discord GuildMember
type Member struct {
	// GuildID is the member's guild id.
	GuildID Snowflake `json:"guild_id"`

	// User is the member's user object.
	User User `json:"user"`

	// Nickname is the user's nickname.
	Nickname string `json:"nick"`

	// Avatar is the member's avatar hash.
	// Note:
	//  - this is difrent from the user avatar, this one is spesific to this guild
	//
	// Optional:
	//  - May be empty string if no avatar.
	Avatar string `json:"avatar"`

	// Banner is the member's banner hash.
	// Note:
	//  - this is difrent from the user banner, this one is spesific to this guild
	//
	// Optional:
	//  - May be empty string if no banner.
	Banner string `json:"banner"`

	// RoleIDs is the ids of roles this member have
	RoleIDs []Snowflake `json:"roles,omitempty"`

	// JoinedAt when the user joined the guild
	//
	// Optional:
	//  - Nil in VoiceStateUpdate event if the member was invited as a guest.
	JoinedAt *time.Time `json:"joined_at"`

	// PremiumSince when the user started boosting the guild
	//
	// Optional:
	//  - Nil if member is not a server booster
	PremiumSince *time.Time `json:"premium_since,omitempty"`

	// Deaf is whether the user is deafened in voice channels
	Deaf bool `json:"deaf,omitempty"`

	// Mute is whether the user is muted in voice channels
	Mute bool `json:"mute,omitempty"`

	// Flags guild member flags represented as a bit set, defaults to 0
	Flags MemberFlags `json:"flags"`

	// Pending is whether the user has not yet passed the guild's Membership Screening requirements
	Pending bool `json:"pending"`

	// CommunicationDisabledUntil is when the user's timeout will expire and the user will be able to communicate in the guild again, null or a time in the past if the user is not timed out
	CommunicationDisabledUntil *time.Time `json:"communication_disabled_until"`

	// AvatarDecorationData is the data for the member's guild avatar decoration
	AvatarDecorationData *AvatarDecorationData `json:"avatar_decoration_data"`
}

// Mention returns a Discord mention string for the user.
//
// Example output: "<@123456789012345678>"
func (m *Member) Mention() string {
	return m.User.Mention()
}

// CreatedAt returns the time when this member account is created.
func (m *Member) CreatedAt() time.Time {
	return m.User.CreatedAt()
}

// DisplayName returns the member's nickname if set,
// otherwise it returns their global display name if set,
// otherwise it falls back to their username.
//
// - Nickname: a guild-specific name set by the user or server mods.
// - Globalname: the name shown across Discord (can differ from username).
// - Username: the original account username.
//
// Example usage:
//
//	name := member.DisplayName()
func (m *Member) DisplayName() string {
	if m.Nickname != "" {
		return m.Nickname
	}
	return m.User.DisplayName()
}

// AvatarURL returns the URL to the members's avatar image.
//
// If the member has a custom avatar set, it returns the URL to that avatar.
// Otherwise it returns their global user avatar URL,
// By default, it uses GIF format if the avatar is animated, otherwise PNG,
// with a default size of 1024.
//
// Example usage:
//
//	url := member.AvatarURL()
func (m *Member) AvatarURL() string {
	if m.Avatar != "" {
		return GuildMemberAvatarURL(m.GuildID, m.User.ID, m.Avatar, GuildMemberAvatarFormatGIF, ImageSize1024)
	}
	return m.User.AvatarURL()
}

// AvatarURLWith returns the URL to the member's avatar image,
// allowing explicit specification of image format and size.
//
// If the user has a custom avatar set, it returns the URL to that avatar.
// Otherwise it returns their global user avatar URL using the provided format and size.
//
// Example usage:
//
//	url := member.AvatarURLWith(GuildMemberAvatarFormatWebP, ImageSize512)
func (m *Member) AvatarURLWith(format GuildMemberAvatarFormat, size ImageSize) string {
	if m.Avatar != "" {
		return GuildMemberAvatarURL(m.GuildID, m.User.ID, m.Avatar, format, size)
	}
	return m.User.AvatarURLWith(UserAvatarFormat(format), size)
}

// BannerURL returns the URL to the members's banner image.
//
// If the member has a custom banner set, it returns the URL to that banner.
// Otherwise it returns their global user banner URL,
// By default, it uses GIF format if the banner is animated, otherwise PNG,
// with a default size of 1024.
//
// Example usage:
//
//	url := member.BannerURL()
func (m *Member) BannerURL() string {
	if m.Avatar != "" {
		return GuildMemberBannerURL(m.GuildID, m.User.ID, m.Avatar, GuildMemberBannerFormatGIF, ImageSize1024)
	}
	return m.User.BannerURL()
}

// BannerURLWith returns the URL to the member's banner image,
// allowing explicit specification of image format and size.
//
// If the user has a custom banner set, it returns the URL to that avatar.
// Otherwise it returns their global user banner URL using the provided format and size.
//
// Example usage:
//
//	url := member.BannerURLWith()
func (m *Member) BannerURLWith(format GuildMemberBannerFormat, size ImageSize) string {
	if m.Avatar != "" {
		return GuildMemberBannerURL(m.GuildID, m.User.ID, m.Avatar, format, size)
	}
	return m.User.BannerURLWith(UserBannerFormat(format), size)
}
