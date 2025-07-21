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

// UserFlags represents flags on a Discord user account.
//
// Reference: https://discord.com/developers/docs/resources/user#user-object-user-flags
type UserFlags int

const (
	// Discord Employee
	UserFlagStaff UserFlags = 1 << 0

	// Partnered Server Owner
	UserFlagPartner UserFlags = 1 << 1

	// HypeSquad Events Member
	UserFlagHypeSquad UserFlags = 1 << 2

	// Bug Hunter Level 1
	UserFlagBugHunterLevel1 UserFlags = 1 << 3

	// House Bravery Member
	UserFlagHypeSquadOnlineHouse1 UserFlags = 1 << 6

	// House Brilliance Member
	UserFlagHypeSquadOnlineHouse2 UserFlags = 1 << 7

	// House Balance Member
	UserFlagHypeSquadOnlineHouse3 UserFlags = 1 << 8

	// Early Nitro Supporter
	UserFlagPremiumEarlySupporter UserFlags = 1 << 9

	// User is a team
	UserFlagTeamPseudoUser UserFlags = 1 << 10

	// Bug Hunter Level 2
	UserFlagBugHunterLevel2 UserFlags = 1 << 14

	// Verified Bot
	UserFlagVerifiedBot UserFlags = 1 << 16

	// Early Verified Bot Developer
	UserFlagVerifiedDeveloper UserFlags = 1 << 17

	// Moderator Programs Alumni
	UserFlagCertifiedModerator UserFlags = 1 << 18

	// Bot uses only HTTP interactions and is shown in the online member list
	UserFlagBotHTTPInteractions UserFlags = 1 << 19

	// User is an Active Developer
	UserFlagActiveDeveloper UserFlags = 1 << 22
)

// Has returns true if all provided flags are set.
func (f UserFlags) Has(flags ...UserFlags) bool {
	for _, flag := range flags {
		if f&flag != flag {
			return false
		}
	}
	return true
}

// Nameplate represents the nameplate the user has.
//
// Reference: https://discord.com/developers/docs/resources/user#nameplate
type Nameplate struct {
	// SkuID is the Discord snowflake ID of the nameplate SKU.
	//
	// Always present.
	SkuID Snowflake `json:"sku_id"`

	// Asset is the path to the nameplate asset.
	//
	// Always present.
	Asset string `json:"asset"`

	// Label is the label of this nameplate.
	//
	// Optional and currently unused by Discord, may be empty string.
	Label string `json:"label"`

	// Palette is the background color of the nameplate.
	//
	// Always present.
	// Allowed values:
	// "crimson", "berry", "sky", "teal", "forest",
	// "bubble_gum", "violet", "cobalt", "clover", "lemon", "white"
	Palette string `json:"palette"`
}

// Collectibles represents collectibles the user owns,
// excluding avatar decorations and profile effects.
//
// Reference: https://discord.com/developers/docs/resources/user#collectibles
type Collectibles struct {
	// Nameplate is the user's nameplate collectible data.
	//
	// Optional, may be nil if the user has no nameplate collectible.
	Nameplate *Nameplate `json:"nameplate,omitempty"`
}

// UserPrimaryGuild represents the user's primary guild info.
//
// Optionally included by Discord API.
//
// Reference: https://discord.com/developers/docs/resources/user#user-primary-guild-object
type UserPrimaryGuild struct {
	// IdentityGuildID is the Discord snowflake ID of the user's primary guild.
	//
	// Optional:
	// - May be nil if the user has no primary guild set.
	// - May be nil if the system cleared the identity due to guild tag support removal or privacy.
	IdentityGuildID *Snowflake `json:"identity_guild_id,omitempty"`

	// IdentityEnabled indicates if the user currently displays the primary guild's server tag.
	//
	// Optional:
	// - May be nil if the identity was cleared by the system (e.g., guild tag disabled).
	// - May be false if the user explicitly disabled showing the tag.
	IdentityEnabled *bool `json:"identity_enabled,omitempty"`

	// Tag is the text of the user's server tag.
	//
	// Optional:
	// - May be nil or empty string if no tag is set.
	// - Limited to 4 characters.
	// - May be cleared if tag data is invalid or unavailable.
	Tag *string `json:"tag,omitempty"`

	// Badge is the hash string of the user's server tag badge.
	//
	// Optional:
	// - May be nil if user has no badge or badge info unavailable.
	// - Appearance depends on guild config or Discord rollout.
	Badge *string `json:"badge,omitempty"`
}

// AvatarDecorationData represents avatar decoration info.
//
// Reference: https://discord.com/developers/docs/resources/user#avatar-decoration-object
type AvatarDecorationData struct {
	// Asset is the avatar decoration hash.
	//
	// Always present.
	Asset string `json:"asset"`

	// SkuID is the Discord snowflake ID of the avatar decoration SKU.
	//
	// Always present.
	SkuID Snowflake `json:"sku_id"`
}

// User represents a Discord user object.
//
// Reference: https://discord.com/developers/docs/resources/user#user-object-user-structure
type User struct {
	// ID is the user's unique Discord snowflake ID.
	//
	// Always present.
	ID Snowflake `json:"id"`

	// Username is the user's username (not unique).
	//
	// Always present.
	Username string `json:"username"`

	// Discriminator is the user's 4-digit Discord tag suffix.
	//
	// Always present.
	Discriminator string `json:"discriminator"`

	// GlobalName is the user's display name.
	//
	// Always present, may be empty string if unset.
	// For bots, this is the application name.
	GlobalName string `json:"global_name"`

	// Avatar is the user's avatar hash.
	//
	// Always present, may be empty string if no avatar.
	Avatar string `json:"avatar"`

	// Bot indicates if the user is a bot account.
	//
	// Optional:
	// - Omitted or false for normal users.
	// - Present and true for bot accounts.
	// When true, bot-only fields may be accessed safely without nil checks.
	Bot bool `json:"bot,omitempty"`

	// System indicates if the user is an official Discord system user.
	//
	// Optional:
	// - Omitted if false.
	// Only applicable for special system users.
	System bool `json:"system,omitempty"`

	// MFAEnabled indicates if the user has two-factor authentication enabled.
	//
	// Optional:
	// - Omitted if false or not applicable.
	MFAEnabled bool `json:"mfa_enabled,omitempty"`

	// Banner is the user's banner hash.
	//
	// Always present, may be empty string if no banner.
	Banner string `json:"banner"`

	// AccentColor is the user's banner color encoded as an integer.
	//
	// Optional:
	// - May be nil if no accent color is set.
	AccentColor *int `json:"accent_color"`

	// Locale is the user's chosen language/locale.
	//
	// Optional:
	// - May be omitted for bots or partial user objects.
	Locale *string `json:"locale,omitempty"`

	// Verified indicates if the user's email is verified.
	//
	// Optional:
	// - Present only in OAuth2 user objects with `email` scope.
	// - Nil if not present or scope not granted.
	Verified *bool `json:"verified,omitempty"`

	// Email is the user's email address.
	//
	// Optional:
	// - Present only in OAuth2 user objects with `email` scope.
	// - Nil if not present or scope not granted.
	Email *string `json:"email"`

	// Flags are internal user account flags.
	//
	// Optional:
	// - May be nil if flags are not present.
	Flags *UserFlags `json:"flags,omitempty"`

	// PremiumType is the Nitro subscription type.
	//
	// Optional:
	// - Nil if user has no Nitro.
	PremiumType *int `json:"premium_type,omitempty"`

	// PublicFlags are the public flags on the user account.
	//
	// Optional:
	// - May be nil if no public flags.
	PublicFlags *UserFlags `json:"public_flags,omitempty"`

	// AvatarDecorationData holds avatar decoration info.
	//
	// Optional:
	// - May be nil if user has no avatar decoration.
	AvatarDecorationData *AvatarDecorationData `json:"avatar_decoration_data,omitempty"`

	// Collectibles holds user's collectibles.
	//
	// Optional:
	// - May be nil if user has no collectibles.
	Collectibles *Collectibles `json:"collectibles,omitempty"`

	// PrimaryGuild holds the user's primary guild info.
	//
	// Optional:
	// - May be nil if no primary guild set.
	// - May be nil if identity cleared due to guild tag or privacy settings.
	PrimaryGuild *UserPrimaryGuild `json:"primary_guild,omitempty"`
}
