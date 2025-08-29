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
	// BotID is the ID of the bot this role belongs to.
	//
	// Always present, 0 if not set.
	BotID Snowflake `json:"bot_id"`

	// IntegrationID is the ID of the integration this role belongs to.
	//
	// Always present, 0 if not set.
	IntegrationID Snowflake `json:"integration_id"`

	// PremiumSubscriber indicates whether this is the guild's Booster role.
	//
	// True if present, false otherwise.
	PremiumSubscriber *struct{} `json:"premium_subscriber,omitempty"`

	// SubscriptionListingID is the ID of this role's subscription SKU and listing.
	//
	// Always present, 0 if not set.
	SubscriptionListingID Snowflake `json:"subscription_listing_id"`

	// AvailableForPurchase indicates whether this role is available for purchase.
	//
	// True if present, false otherwise.
	AvailableForPurchase *struct{} `json:"available_for_purchase,omitempty"`

	// GuildConnections indicates whether this role is a guild's linked role.
	//
	// True if present, false otherwise.
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
	// Always present, 0 if not set.
	SecondaryColor Color `json:"secondary_color"`

	// TertiaryColor is the tertiary color for the role.
	//
	// Always present, 0 if not set.
	TertiaryColor Color `json:"tertiary_color"`
}

// Role represents a Discord role.
//
// Reference: https://discord.com/developers/docs/resources/guild#role-object-role-structure
type Role struct {
	// ID is the role ID.
	//
	// Always present.
	ID Snowflake `json:"id"`

	// Name is the role name.
	//
	// Always present.
	Name string `json:"name"`

	// Colors contains the role's color definitions.
	//
	// Always present.
	Colors RoleColors `json:"colors"`

	// Hoist indicates if this role is pinned in the user listing.
	//
	// Always present.
	Hoist bool `json:"hoist"`

	// Icon is the role's icon hash.
	//
	// Always present, may be empty string if no icon.
	Icon string `json:"icon"`

	// UnicodeEmoji is the role's unicode emoji.
	//
	// Always present, may be empty string if not set.
	UnicodeEmoji string `json:"unicode_emoji"`

	// Position is the position of this role (roles with same position are sorted by ID).
	Position int `json:"position"`

	// Permissions is the permission bit set for this role.
	//
	// Always present.
	Permissions string `json:"permissions"`

	// Managed indicates whether this role is managed by an integration.
	//
	// Always present.
	Managed bool `json:"managed"`

	// Mentionable indicates whether this role is mentionable.
	//
	// Always present.
	Mentionable bool `json:"mentionable"`

	// Tags contains the tags this role has.
	//
	// Optional; may be nil if no tags.
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
