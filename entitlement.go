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

import (
	"time"
)

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
	UserID Snowflake `json:"user_id,omitempty"`

	// Type is the type of entitlement.
	Type EntitlementType `json:"type"`

	// Deleted indicates whether the entitlement has been deleted.
	Deleted bool `json:"deleted"`

	// StartsAt is the start date at which the entitlement is valid.
	StartsAt *time.Time `json:"starts_at,omitempty"`

	// EndsAt is the optional date at which the entitlement is no longer valid.
	EndsAt *time.Time `json:"ends_at,omitempty"`

	// GuildID is the id of the guild that is granted access to the entitlement's SKU.
	GuildID Snowflake `json:"guild_id,omitempty"`

	// Consumed indicates whether the entitlement for a consumable item has been consumed.
	Consumed *bool `json:"consumed,omitempty"`
}

// SKUType represents the type of a SKU.
//
// Reference: https://discord.com/developers/docs/resources/sku#sku-object-sku-types
type SKUType int

const (
	// SKUTypeDurable represents a durable one-time purchase.
	SKUTypeDurable SKUType = 2

	// SKUTypeConsumable represents a consumable one-time purchase.
	SKUTypeConsumable SKUType = 3

	// SKUTypeSubscription represents a recurring subscription.
	SKUTypeSubscription SKUType = 5

	// SKUTypeSubscriptionGroup represents a system-generated group for each subscription SKU.
	SKUTypeSubscriptionGroup SKUType = 6
)

// Is returns true if this SKUType matches the provided one.
func (v SKUType) Is(other SKUType) bool {
	return v == other
}

// SKUFlags represents flags for a SKU combined as a bitfield.
//
// Reference: https://discord.com/developers/docs/resources/sku#sku-object-sku-flags
type SKUFlags int

const (
	SKUFlagAvailable         SKUFlags = 1 << 2
	SKUFlagGuildSubscription SKUFlags = 1 << 7
	SKUFlagUserSubscription  SKUFlags = 1 << 8
)

// Has returns true if the flags contain all of the provided flags.
func (f SKUFlags) Has(flags ...SKUFlags) bool {
	return BitFieldHas(f, flags...)
}

// SKU represents a premium offering that can be made available to an application's users or guilds.
//
// Reference: https://discord.com/developers/docs/resources/sku#sku-object
type SKU struct {
	// ID is the ID of the SKU.
	ID Snowflake `json:"id"`

	// Type is the type of SKU.
	Type SKUType `json:"type"`

	// ApplicationID is the ID of the parent application.
	ApplicationID Snowflake `json:"application_id"`

	// Name is the customer-facing name of the premium offering.
	Name string `json:"name"`

	// Slug is the system-generated URL slug based on the SKU's name.
	Slug string `json:"slug"`

	// Flags is the SKU flags combined as a bitfield.
	Flags SKUFlags `json:"flags"`
}

// SubscriptionStatus represents the current status of a subscription.
//
// Reference: https://discord.com/developers/docs/resources/subscription#subscription-statuses
type SubscriptionStatus int

const (
	// SubscriptionStatusActive indicates the subscription is active and scheduled to renew.
	SubscriptionStatusActive SubscriptionStatus = 0

	// SubscriptionStatusEnding indicates the subscription is active but will not renew.
	SubscriptionStatusEnding SubscriptionStatus = 1

	// SubscriptionStatusInactive indicates the subscription is inactive and not being charged.
	SubscriptionStatusInactive SubscriptionStatus = 2
)

// Is returns true if this SubscriptionStatus matches the provided one.
func (v SubscriptionStatus) Is(other SubscriptionStatus) bool {
	return v == other
}

// Subscription represents a user making recurring payments for a SKU.
//
// Reference: https://discord.com/developers/docs/resources/subscription#subscription-object
type Subscription struct {
	// ID is the ID of the subscription.
	ID Snowflake `json:"id"`

	// UserID is the ID of the user who is subscribed.
	UserID Snowflake `json:"user_id"`

	// SKUIDs is the list of SKUs subscribed to.
	SKUIDs []Snowflake `json:"sku_ids"`

	// EntitlementIDs is the list of entitlements granted for this subscription.
	EntitlementIDs []Snowflake `json:"entitlement_ids"`

	// RenewalSKUIDs is the list of SKUs that this user will be subscribed to at renewal.
	RenewalSKUIDs []Snowflake `json:"renewal_sku_ids,omitempty"`

	// CurrentPeriodStart is the start of the current subscription period.
	CurrentPeriodStart time.Time `json:"current_period_start"`

	// CurrentPeriodEnd is the end of the current subscription period.
	CurrentPeriodEnd time.Time `json:"current_period_end"`

	// Status is the current status of the subscription.
	Status SubscriptionStatus `json:"status"`

	// CanceledAt is when the subscription was canceled.
	CanceledAt *time.Time `json:"canceled_at,omitempty"`

	// Country is the ISO3166-1 alpha-2 country code of the payment source.
	Country *string `json:"country,omitempty"`
}
