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
	"fmt"
	"time"
)

// WebhookType represents the type of a Discord webhook.
//
// Reference: https://discord.com/developers/docs/resources/webhook#webhook-object-webhook-types
type WebhookType int

const (
	// WebhookTypeIncoming is an Incoming Webhook that can post messages to channels with a token.
	WebhookTypeIncoming WebhookType = 1

	// WebhookTypeChannelFollower is a Channel Follower Webhook used for Announcement Channel following.
	WebhookTypeChannelFollower WebhookType = 2

	// WebhookTypeApplication is an Application Webhook used with Interactions.
	WebhookTypeApplication WebhookType = 3
)

// Is returns true if this WebhookType matches the provided one.
func (v WebhookType) Is(other WebhookType) bool {
	return v == other
}

// Webhook represents a Discord webhook object.
//
// Reference: https://discord.com/developers/docs/resources/webhook#webhook-object
type Webhook struct {
	// ID is the unique ID of the webhook.
	ID Snowflake `json:"id"`

	// Type is the type of the webhook.
	Type WebhookType `json:"type"`

	// GuildID is the guild ID this webhook is for, if any.
	GuildID Snowflake `json:"guild_id,omitempty"`

	// ChannelID is the channel ID this webhook is for, if any.
	ChannelID Snowflake `json:"channel_id,omitempty"`

	// User is the user that created this webhook (not included when getting webhook with token).
	User *User `json:"user,omitempty"`

	// Name is the default name of the webhook.
	Name string `json:"name,omitempty"`

	// Avatar is the default user avatar hash of the webhook.
	Avatar string `json:"avatar,omitempty"`

	// Token is the secure token of the webhook (returned for Incoming Webhooks).
	Token string `json:"token,omitempty"`

	// ApplicationID is the bot/OAuth2 application that created this webhook.
	ApplicationID Snowflake `json:"application_id,omitempty"`

	// SourceGuild is the guild of the channel that this webhook is following (Channel Follower Webhooks).
	SourceGuild *Guild `json:"source_guild,omitempty"`

	// SourceChannel is the channel that this webhook is following (Channel Follower Webhooks).
	SourceChannel Channel `json:"source_channel,omitempty"`

	// URL is the url used for executing the webhook (returned by the webhooks OAuth2 flow).
	URL string `json:"url,omitempty"`
}

// CreatedAt returns the creation timestamp of the webhook.
func (w *Webhook) CreatedAt() time.Time {
	return w.ID.Timestamp()
}

// String implements the fmt.Stringer interface.
func (w *Webhook) String() string {
	return w.Name
}

var _ fmt.Stringer = (*Webhook)(nil)
