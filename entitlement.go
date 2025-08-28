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

// EntitlementType represents the type of an entitlement in Discord.
//
// Reference: https://discord.com/developers/docs/resources/entitlement#entitlement-object-entitlement-types
type EntitlementType int

const (
	// EntitlementTypePurchase indicates the entitlement was purchased by the user.
	EntitlementTypePurchase EntitlementType = 1 + iota

	// EntitlementTypePremiumSubscription indicates the entitlement is for a Discord Nitro subscription.
	EntitlementTypePremiumSubscription

	// EntitlementTypeDeveloperGift indicates the entitlement was gifted by a developer.
	EntitlementTypeDeveloperGift

	// EntitlementTypeTestModePurchase indicates the entitlement was purchased by a developer in application test mode.
	EntitlementTypeTestModePurchase

	// EntitlementTypeFreePurchase indicates the entitlement was granted when the SKU was free.
	EntitlementTypeFreePurchase

	// EntitlementTypeUserGift indicates the entitlement was gifted by another user.
	EntitlementTypeUserGift

	// EntitlementTypePremiumPurchase indicates the entitlement was claimed for free by a Nitro subscriber.
	EntitlementTypePremiumPurchase

	// EntitlementTypeApplicationSubscription indicates the entitlement was purchased as an app subscription.
	EntitlementTypeApplicationSubscription
)

// Is returns true if the entitlement's Type matches the provided one.
func (t EntitlementType) Is(entitlementType EntitlementType) bool {
	return t == entitlementType
}

// Entitlement represents a Discord Entitlement.
//
// Reference: https://discord.com/developers/docs/resources/entitlement#entitlement-object
type Entitlement struct {
	// ID is the unique identifier of the entitlement.
	ID Snowflake `json:"id"`

	// SkuID is the ID of the SKU associated with this entitlement.
	SkuID Snowflake `json:"sku_id"`

	// ApplicationID is the ID of the application this entitlement belongs to.
	ApplicationID Snowflake `json:"application_id"`

	// UserID is the id of the user that is granted access to the entitlement's SKU.
	//
	// Optional:
	//   - Will be 0 if the entitlement is not associated with a specific user.
	UserID Snowflake `json:"user_id"`

	// Type is the type of entitlement.
	Type EntitlementType `json:"type"`

	// Deleted indicates whether the entitlement has been deleted.
	Deleted bool `json:"deleted"`

	// StartsAt is the start date at which the entitlement is valid.
	//
	// Optional.
	StartsAt *time.Time `json:"starts_at"`

	// EndsAt is the optional date at which the entitlement is no longer valid.
	//
	// Optional.
	EndsAt *time.Time `json:"ends_at"`

	// GuildID is the id of the guild that is granted access to the entitlement's SKU.
	//
	// Optional:
	//   - Will be 0 if the entitlement is not associated with a guild.
	GuildID Snowflake `json:"guild_id"`

	// Consumed indicates whether the entitlement for a consumable item has been consumed.
	//
	// Optional:
	//   - Will be null for non-consumable entitlements.
	Consumed *bool `json:"consumed"`
}
