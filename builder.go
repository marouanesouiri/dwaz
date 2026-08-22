/************************************************************************************
 *
 * wumpkit, A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD-3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package wumpkit

// ReusableBuilder represents a builder whose internal state can be reset
// and reused for subsequent object constructions to avoid allocations.
type ReusableBuilder interface {
	Reset()
}

// Compile-time assertions verifying that all builder types implement ReusableBuilder.
var (
	_ ReusableBuilder = (*EmbedBuilder)(nil)
	_ ReusableBuilder = (*ActionRowBuilder)(nil)
	_ ReusableBuilder = (*ButtonBuilder)(nil)
	_ ReusableBuilder = (*StringSelectMenuBuilder)(nil)
	_ ReusableBuilder = (*TextInputBuilder)(nil)
	_ ReusableBuilder = (*UserSelectMenuBuilder)(nil)
	_ ReusableBuilder = (*RoleSelectMenuBuilder)(nil)
	_ ReusableBuilder = (*MentionableSelectMenuBuilder)(nil)
	_ ReusableBuilder = (*ChannelSelectMenuBuilder)(nil)
	_ ReusableBuilder = (*SectionBuilder)(nil)
	_ ReusableBuilder = (*TextDisplayBuilder)(nil)
	_ ReusableBuilder = (*ThumbnailBuilder)(nil)
	_ ReusableBuilder = (*MediaGalleryBuilder)(nil)
	_ ReusableBuilder = (*FileBuilder)(nil)
	_ ReusableBuilder = (*SeparatorBuilder)(nil)
	_ ReusableBuilder = (*ContainerBuilder)(nil)
	_ ReusableBuilder = (*LabelBuilder)(nil)
	_ ReusableBuilder = (*SlashCommandBuilder)(nil)
	_ ReusableBuilder = (*UserCommandBuilder)(nil)
	_ ReusableBuilder = (*MessageCommandBuilder)(nil)
)
