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

// RoleFlags represents flags on a Discord guild role.
//
// Reference: https://discord.com/developers/docs/topics/permissions#role-object-role-flags
type RoleFlags int

const (
	// Role can be selected by members in an onboarding prompt.
	RoleFlagInPrompt RoleFlags = 1 << 0
)

// Has returns true if all provided flags are set.
func (f RoleFlags) Has(flag RoleFlags) bool {
	return BitFieldHas(f, flag)
}

// RoleTags represents the tags object attached to a role.
//
// Reference: https://discord.com/developers/docs/topics/permissions#role-object-role-tags-structure
type RoleTags struct {
	// BotID is the ID of the bot that this role belongs to.
	// It is set for roles automatically created when adding a bot
	// to a guild with specific permissions.
	//
	// Optional:
	//   - Will be 0 if the role is not associated with a bot.
	BotID Snowflake `json:"bot_id"`

	// IntegrationID is the ID of the integration that this role belongs to.
	//
	// Optional:
	//   - Will be 0 if the role is not associated with an integration.
	IntegrationID Snowflake `json:"integration_id"`

	// PremiumSubscriber indicates whether this is the guild's Booster role.
	//
	// True if present (not nil), false otherwise (nil).
	PremiumSubscriber *struct{} `json:"premium_subscriber,omitempty"`

	// SubscriptionListingID is the ID of this role's subscription SKU and listing.
	//
	// Optional:
	//   - Will be 0 if the role is not linked to a subscription.
	SubscriptionListingID Snowflake `json:"subscription_listing_id"`

	// AvailableForPurchase indicates whether this role is available for purchase.
	//
	// True if present (not nil), false otherwise (nil).
	AvailableForPurchase *struct{} `json:"available_for_purchase,omitempty"`

	// GuildConnections indicates whether this role is a guild's linked role.
	//
	// True if present (not nil), false otherwise (nil).
	GuildConnections *struct{} `json:"guild_connections,omitempty"`
}

// RoleColors represents a role's color definitions.
//
// Reference: https://discord.com/developers/docs/resources/guild#role-object-role-colors-object
type RoleColors struct {
	// PrimaryColor is the primary color for the role.
	PrimaryColor Color `json:"primary_color"`

	// SecondaryColor is the secondary color for the role.
	//
	// Optional:
	//   - Will be nil if not set.
	SecondaryColor *Color `json:"secondary_color"`

	// TertiaryColor is the tertiary color for the role.
	//
	// Optional:
	//   - Will be nil if not set.
	TertiaryColor *Color `json:"tertiary_color"`
}

// Role represents a Discord role.
//
// Reference: https://discord.com/developers/docs/resources/guild#role-object-role-structure
type Role struct {
	// ID is the role ID.
	ID Snowflake `json:"id"`

	// Name is the role name.
	Name string `json:"name"`

	// Colors contains the role's color definitions.
	Colors RoleColors `json:"colors"`

	// Hoist indicates if this role is pinned in the user listing.
	Hoist bool `json:"hoist"`

	// Icon is the role's icon hash.
	//
	// Optional:
	//   - Will be empty string if no icon.
	Icon string `json:"icon"`

	// UnicodeEmoji is the role's unicode emoji.
	//
	// Optional:
	//   - Will be empty string if not set.
	UnicodeEmoji string `json:"unicode_emoji"`

	// Position is the position of this role (roles with same position are sorted by ID).
	//
	// Note:
	//   - Roles with same position are sorted by ID.
	Position int `json:"position"`

	// Permissions is the permission bit set for this role.
	Permissions Permissions `json:"permissions"`

	// Managed indicates whether this role is managed by an integration.
	Managed bool `json:"managed"`

	// Mentionable indicates whether this role is mentionable.
	Mentionable bool `json:"mentionable"`

	// Tags contains the tags this role has.
	//
	// Optional:
	//   - Will be nil if no tags.
	Tags *RoleTags `json:"tags,omitempty"`

	// Flags are role flags combined as a bitfield.
	Flags RoleFlags `json:"flags"`
}

// Mention returns a Discord mention string for the role.
//
// Example output: "<@&123456789012345678>"
func (r *Role) Mention() string {
	return "<@&" + r.ID.String() + ">"
}

// IconURL returns the URL to the role's icon image in PNG format.
//
// If the role has a custom icon set, it returns the URL to that icon,
// Otherwise it returns an empty string.
//
// Example usage:
//
//	url := role.IconURL()
func (u *Role) IconURL() string {
	if u.Icon != "" {
		return RoleIconURL(u.ID, u.Icon, ImageFormatDefault, ImageSizeDefault)
	}
	return ""
}

// IconURLWith returns the URL to the role's icon image,
// allowing explicit specification of image format and size.
//
// If the role has a custom icon set, it returns the URL to that icon
// using the provided format and size, Otherwise it returns an empty string.
//
// Example usage:
//
//	url := role.IconURLWith(ImageFormatWebP, ImageSize512)
func (u *Role) IconURLWith(format ImageFormat, size ImageSize) string {
	if u.Icon != "" {
		return RoleIconURL(u.ID, u.Icon, format, size)
	}
	return ""
}
