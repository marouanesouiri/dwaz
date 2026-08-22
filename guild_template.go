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

// GuildTemplate represents a code that when used, creates a guild based on a snapshot of an existing guild.
//
// Reference: https://discord.com/developers/docs/resources/guild-template#guild-template-object
type GuildTemplate struct {
	// Code is the template code (unique ID).
	Code string `json:"code"`

	// Name is the template name.
	Name string `json:"name"`

	// Description is the description for the template.
	Description *string `json:"description,omitempty"`

	// UsageCount is the number of times this template has been used.
	UsageCount int `json:"usage_count"`

	// CreatorID is the ID of the user who created the template.
	CreatorID Snowflake `json:"creator_id"`

	// Creator is the user who created the template.
	Creator User `json:"creator"`

	// CreatedAt is when this template was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when this template was last synced to the source guild.
	UpdatedAt time.Time `json:"updated_at"`

	// SourceGuildID is the ID of the guild this template is based on.
	SourceGuildID Snowflake `json:"source_guild_id"`

	// SerializedSourceGuild is the guild snapshot this template contains.
	SerializedSourceGuild RestGuild `json:"serialized_source_guild"`

	// IsDirty is whether the template has unsynced changes.
	IsDirty *bool `json:"is_dirty,omitempty"`
}
