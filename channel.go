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

// ChannelType represents Discord channel types.
//
// Reference: https://discord.com/developers/docs/resources/channel#channel-object-channel-types
type ChannelType int

const (
	// GuildText is a text channel within a server.
	ChannelTypeGuildText ChannelType = 0

	// DM is a direct message between users.
	ChannelTypeDM ChannelType = 1

	// GuildVoice is a voice channel within a server.
	ChannelTypeGuildVoice ChannelType = 2

	// GroupDM is a direct message between multiple users.
	ChannelTypeGroupDM ChannelType = 3

	// GuildCategory is an organizational category that contains up to 50 channels.
	ChannelTypeGuildCategory ChannelType = 4

	// GuildAnnouncement is a channel that users can follow and crosspost into their own server (formerly news channels).
	ChannelTypeGuildAnnouncement ChannelType = 5

	// AnnouncementThread is a temporary sub-channel within a GuildAnnouncement channel.
	ChannelTypeAnnouncementThread ChannelType = 10

	// PublicThread is a temporary sub-channel within a GuildText or GuildForum channel.
	ChannelTypePublicThread ChannelType = 11

	// PrivateThread is a temporary sub-channel within a GuildText channel that is only viewable by those invited and those with the MANAGE_THREADS permission.
	ChannelTypePrivateThread ChannelType = 12

	// GuildStageVoice is a voice channel for hosting events with an audience.
	ChannelTypeGuildStageVoice ChannelType = 13

	// GuildDirectory is the channel in a hub containing the listed servers.
	ChannelTypeGuildDirectory ChannelType = 14

	// GuildForum is a channel that can only contain threads.
	ChannelTypeGuildForum ChannelType = 15

	// GuildMedia is a channel that can only contain threads, similar to GuildForum channels.
	ChannelTypeGuildMedia ChannelType = 16
)

// Is returns true if the channel's Type matches the provided one.
func (t ChannelType) Is(channelType ChannelType) bool {
	return t == channelType
}

// ChannelFlags represents Discord channel flags combined as a bitfield.
//
// Reference: https://discord.com/developers/docs/resources/channel#channel-object-channel-flags
type ChannelFlags int

const (
	// ChannelFlagPinned indicates that this thread is pinned to the top of its parent
	// GUILD_FORUM or GUILD_MEDIA channel.
	//
	// Applicable only to threads within forum or media channels.
	ChannelFlagPinned ChannelFlags = 1 << 1

	// ChannelFlagRequireTag indicates whether a tag is required to be specified when
	// creating a thread in a GUILD_FORUM or GUILD_MEDIA channel.
	//
	// Tags are specified in the AppliedTags field.
	ChannelFlagRequireTag ChannelFlags = 1 << 4

	// ChannelFlagHideMediaDownloadOptions, when set, hides the embedded media download options
	// for media channel attachments.
	//
	// Available only for media channels.
	ChannelFlagHideMediaDownloadOptions ChannelFlags = 1 << 15
)

// Has returns true if all provided flags are set.
func (f ChannelFlags) Has(flags ...ChannelFlags) bool {
	for _, flag := range flags {
		if f&flag != flag {
			return false
		}
	}
	return true
}

// PermissionOverwriteType defines the type of permission overwrite target.
//
// Reference: https://discord.com/developers/docs/resources/channel#overwrite-object-overwrite-structure
type PermissionOverwriteType int

const (
	// PermissionOverwriteTypeRole indicates the overwrite applies to a role.
	PermissionOverwriteTypeRole PermissionOverwriteType = 0

	// PermissionOverwriteTypeMember indicates the overwrite applies to a member.
	PermissionOverwriteTypeMember PermissionOverwriteType = 1
)

// ForumPostsSortOrder defines the sort order type used to order posts in forum/media channels.
//
// Reference: https://discord.com/developers/docs/resources/channel#channel-object-sort-order-types
type ForumPostsSortOrder int

const (
	// ForumPostsSortOrderLatestActivity sorts posts by latest activity (default).
	ForumPostsSortOrderLatestActivity ForumPostsSortOrder = 0

	// ForumPostsSortOrderCreationDate sorts posts by creation time (most recent to oldest).
	ForumPostsSortOrderCreationDate ForumPostsSortOrder = 1
)

// Is returns true if the channel's SortOrder type matches the provided one.
func (t ForumPostsSortOrder) Is(sortOrderType ForumPostsSortOrder) bool {
	return t == sortOrderType
}

// ForumLayout defines the layout type used to place posts in forum/media channels.
//
// Reference: https://discord.com/developers/docs/resources/channel#channel-object-forum-layout-types
type ForumLayout int

const (
	// ForumLayoutNotSet indicates no default has been set for forum channel.
	ForumLayoutNotSet ForumLayout = 0

	// ForumLayoutListView displays posts as a list.
	ForumLayoutListView ForumLayout = 1

	// ForumLayoutGalleryView displays posts as a collection of tiles.
	ForumLayoutGalleryView ForumLayout = 2
)

// Is returns true if the channel's PostsLayout type matches the provided one.
func (t ForumLayout) Is(layoutType ForumLayout) bool {
	return t == layoutType
}

// PermissionOverwrite represents a permission overwrite for a role or member.
//
// Used to grant or deny specific permissions in a channel.
//
// Reference: https://discord.com/developers/docs/resources/channel#overwrite-object-overwrite-structure
type PermissionOverwrite struct {
	// ID is the role or user ID the overwrite applies to.
	ID Snowflake `json:"id"`

	// Type specifies whether this overwrite is for a role or a member.
	Type PermissionOverwriteType `json:"type"`

	// Allow is the permission bit set explicitly allowed.
	Allow Permissions `json:"allow"`

	// Deny is the permission bit set explicitly denied.
	Deny Permissions `json:"deny"`
}

// Channel is the interface representing a Discord channel.
type Channel interface {
	GetID() Snowflake
	GetType() ChannelType
	Mention() string
}

// GuildChannel represents a guild-specific channel.
type GuildChannel interface {
	Channel
	GetGuildID() Snowflake
	GetName() string
	GetPosition() int
	GetPermissionOverwrites() []PermissionOverwrite
	GetFlags() ChannelFlags
	JumpURL() string
}

// baseChannel contains only fields present in all channel types.
//
// Reference: https://discord.com/developers/docs/resources/channel#channel-object-channel-structure
type baseChannel struct {
	// ID is the unique Discord snowflake ID of the channel.
	ID Snowflake `json:"id"`

	// Type is the type of the channel.
	Type ChannelType `json:"type"`
}

func (c *baseChannel) GetID() Snowflake     { return c.ID }
func (c *baseChannel) GetType() ChannelType { return c.Type }

// Mention returns a Discord mention string for the channel.
//
// Example output: "<#123456789012345678>"
func (c *baseChannel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

// baseGuildChannel embeds baseChannel and adds fields present in guild channels only.
//
// Used by all guild-specific channel types like TextChannel, VoiceChannel, ForumChannel, etc.
type baseGuildChannel struct {
	baseChannel

	// GuildID is the id of the guild.
	//
	// Always present.
	GuildID Snowflake `json:"guild_id"`

	// Name is the name of the channel.
	//
	// Always present. 1-100 characters.
	Name string `json:"name,omitempty"`

	// Position is the sorting position of the channel.
	//
	// Always present.
	Position int `json:"position,omitempty"`

	// PermissionOverwrites are explicit permission overwrites for members and roles.
	//
	// Always present.
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`

	// Flags are combined channel flags.
	//
	// Always present.
	Flags ChannelFlags `json:"flags,omitempty"`
}

func (c *baseGuildChannel) GetGuildID() Snowflake { return c.GuildID }
func (c *baseGuildChannel) GetName() string       { return c.Name }
func (c *baseGuildChannel) GetPosition() int      { return c.Position }
func (c *baseGuildChannel) GetPermissionOverwrites() []PermissionOverwrite {
	return c.PermissionOverwrites
}
func (c *baseGuildChannel) GetFlags() ChannelFlags { return c.Flags }
func (c *baseGuildChannel) JumpURL() string {
	return "https://discord.com/channels/" + c.GuildID.String() + "/" + c.ID.String()
}

// CategoryChannel represents a guild category channel.
type CategoryChannel struct {
	baseGuildChannel
}

// TextChannel represents a guild text channel.
type TextChannel struct {
	baseGuildChannel

	// ParentID is the id of the parent category for this channel.
	//
	// Note:
	//  Each parent category can contain up to 50 channels.
	//
	// Always present. If ParentID == 0, the channel is not in a category.
	ParentID Snowflake `json:"parent_id"`

	// LastMessageID is the id of the last message sent in this channel.
	//
	// Always present.
	LastMessageID Snowflake `json:"last_message_id"`

	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message.
	// Bots, as well as users with the permission manageMessages or manageChannel, are unaffected.
	//
	// Always present.
	RateLimitPerUser time.Duration `json:"rate_limit_per_user"`

	// Topic is the channel topic (can be 0-1024 characters).
	//
	// Always present. Can be empty string if the channel has no topic.
	Topic string `json:"topic"`

	// Nsfw indicates whether the channel is NSFW.
	//
	// Always present.
	Nsfw bool `json:"nsfw"`
}

// VoiceChannel represents a guild voice channel.
type VoiceChannel struct {
	baseGuildChannel

	// ParentID is the id of the parent category for this channel.
	//
	// Note:
	//  Each parent category can contain up to 50 channels.
	//
	// Always present. If ParentID == 0, the channel is not in a category.
	ParentID Snowflake `json:"parent_id"`

	// LastMessageID is the id of the last message sent in this channel.
	//
	// Always present.
	LastMessageID Snowflake `json:"last_message_id"`

	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message.
	// Bots, as well as users with the permission manageMessages or manageChannel, are unaffected.
	//
	// Always present.
	RateLimitPerUser time.Duration `json:"rate_limit_per_user"`

	// Nsfw indicates whether the channel is NSFW.
	//
	// Always present.
	Nsfw bool `json:"nsfw"`

	// Bitrate is the bitrate (in bits) of the voice channel.
	//
	// Always present.
	Bitrate int `json:"bitrate"`

	// UserLimit is the user limit of the voice channel.
	//
	// Always present.
	UserLimit int `json:"user_limit"`

	// RtcRegion is the voice region id for the voice channel. Automatic when set to empty string.
	//
	// Always present.
	RtcRegion string `json:"rtc_region"`
}

// ForumTag represents a tag that can be applied to a thread
// in a GuildForum or GuildMedia channel.
//
// Reference: https://discord.com/developers/docs/resources/channel#forum-tag-object
type ForumTag struct {
	// ID is the id of the tag.
	//
	// Always present.
	ID Snowflake `json:"id"`

	// Name is the name of the tag (0-20 characters).
	//
	// Always present. Can be empty string.
	Name string `json:"name"`

	// Moderated indicates whether this tag can only be added to or removed from
	// threads by a member with the ManageThreads permission.
	//
	// Always present.
	Moderated bool `json:"moderated"`

	// EmojiID is the ID of a guild's custom emoji.
	//
	// Optional:
	// - Zero value (0) means it is not set.
	// - If EmojiName is empty (not set), then EmojiID must be set (non-zero).
	EmojiID Snowflake `json:"emoji_id"`

	// EmojiName is the Unicode character of the emoji.
	//
	// Optional:
	// - Empty string means it is not set.
	// - If EmojiID is zero (not set), then EmojiName must be set (non-empty).
	EmojiName string `json:"emoji_name"`
}

// ForumChannel represents a guild forum channel.
type ForumChannel struct {
	baseGuildChannel

	// ParentID is the id of the parent category for this channel.
	//
	// Note:
	//  Each parent category can contain up to 50 channels.
	//
	// Always present. If ParentID == 0, the channel is not in a category.
	ParentID Snowflake `json:"parent_id"`

	// LastMessageID is the id of the last message sent in this channel.
	//
	// Always present.
	LastMessageID Snowflake `json:"last_message_id"`

	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message.
	// Bots, as well as users with the permission manageMessages or manageChannel, are unaffected.
	//
	// Always present.
	RateLimitPerUser time.Duration `json:"rate_limit_per_user"`

	// Nsfw indicates whether the channel is NSFW.
	//
	// Always present.
	Nsfw bool `json:"nsfw"`

	// Topic is the channel topic (can be 0-4096 characters).
	//
	// Always present. Can be empty string if the channel has no topic.
	Topic string `json:"topic"`

	// AvailableTags is the set of tags that can be used in this channel.
	//
	// Always present. Can be empty if this channel has no tags.
	AvailableTags []ForumTag `json:"available_tags"`

	// DefaultReactionEmoji specifies the emoji used as the default way to react to a forum post.
	//
	// Exactly one of EmojiID and EmojiName must be set.
	//
	// Always present. If EmojiID != 0, it refers to a guild custom emoji.
	// If EmojiID == 0, EmojiName will contain a Unicode emoji character.
	DefaultReactionEmoji struct {
		// EmojiID is the ID of a guild's custom emoji.
		//
		// Optional. If 0, EmojiName will be set instead.
		EmojiID Snowflake `json:"emoji_id"`

		// EmojiName is the Unicode character of the emoji.
		//
		// Optional. If empty, EmojiID will be set instead.
		EmojiName string `json:"emoji_name"`
	} `json:"default_reaction_emoji"`

	// DefaultSortOrder is the default sort order type used to order posts
	// in GuildForum and GuildMedia channels. Defaults to PostsSortOrderLatestActivity.
	//
	// Always present.
	DefaultSortOrder ForumPostsSortOrder `json:"default_sort_order"`

	// DefaultForumLayout is the default forum layout view used to display posts
	// in GuildForum channels. Defaults to ForumLayoutNotSet.
	//
	// Always present.
	DefaultForumLayout ForumLayout `json:"default_forum_layout"`
}

// MediaChannel represents a media channel.
type MediaChannel struct {
	baseGuildChannel

	// ParentID is the id of the parent category for this channel.
	//
	// Note:
	//  Each parent category can contain up to 50 channels.
	//
	// Always present. If ParentID == 0, the channel is not in a category.
	ParentID Snowflake `json:"parent_id"`

	// LastMessageID is the id of the last message sent in this channel.
	//
	// Always present.
	LastMessageID Snowflake `json:"last_message_id"`

	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message.
	// Bots, as well as users with the permission manageMessages or manageChannel, are unaffected.
	//
	// Always present.
	RateLimitPerUser time.Duration `json:"rate_limit_per_user"`

	// Nsfw indicates whether the channel is NSFW.
	//
	// Always present.
	Nsfw bool `json:"nsfw"`

	// Topic is the channel topic (can be 0-4096 characters).
	//
	// Always present. Can be empty string if the channel has no topic.
	Topic string `json:"topic"`

	// AvailableTags is the set of tags that can be used in this channel.
	//
	// Always present. Can be empty if this channel has no tags.
	AvailableTags []ForumTag `json:"available_tags"`

	// DefaultReactionEmoji specifies the emoji used as the default way to react to a forum post.
	//
	// Exactly one of EmojiID and EmojiName must be set.
	//
	// Always present. If EmojiID != 0, it refers to a guild custom emoji.
	// If EmojiID == 0, EmojiName will contain a Unicode emoji character.
	DefaultReactionEmoji struct {
		// EmojiID is the ID of a guild's custom emoji.
		//
		// Optional. If 0, EmojiName will be set instead.
		EmojiID Snowflake `json:"emoji_id"`

		// EmojiName is the Unicode character of the emoji.
		//
		// Optional. If empty, EmojiID will be set instead.
		EmojiName string `json:"emoji_name"`
	} `json:"default_reaction_emoji"`

	// DefaultSortOrder is the default sort order type used to order posts
	// in GuildForum and GuildMedia channels. Defaults to PostsSortOrderLatestActivity.
	//
	// Always present.
	DefaultSortOrder ForumPostsSortOrder `json:"default_sort_order"`

	// DefaultForumLayout is the default forum layout view used to display posts
	// in GuildForum channels. Defaults to ForumLayoutNotSet.
	//
	// Always present.
	DefaultForumLayout ForumLayout `json:"default_forum_layout"`
}

// AnnouncementChannel represents an announcement channel.
type AnnouncementChannel struct {
	baseGuildChannel

	// ParentID is the id of the parent category for this channel.
	//
	// Note:
	//  Each parent category can contain up to 50 channels.
	//
	// Always present. If ParentID == 0, the channel is not in a category.
	ParentID Snowflake `json:"parent_id"`

	// LastMessageID is the id of the last message sent in this channel.
	//
	// Always present.
	LastMessageID Snowflake `json:"last_message_id"`

	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message.
	// Bots, as well as users with the permission manageMessages or manageChannel, are unaffected.
	//
	// Always present.
	RateLimitPerUser time.Duration `json:"rate_limit_per_user"`

	// Topic is the channel topic (can be 0-1024 characters).
	//
	// Always present. Can be empty string if the channel has no topic.
	Topic string `json:"topic"`

	// Nsfw indicates whether the channel is NSFW.
	//
	// Always present.
	Nsfw bool `json:"nsfw"`
}

// StageVoiceChannel represents a stage voice channel.
type StageVoiceChannel struct {
	baseGuildChannel

	// ParentID is the id of the parent category for this channel.
	//
	// Note:
	//  Each parent category can contain up to 50 channels.
	//
	// Always present. If ParentID == 0, the channel is not in a category.
	ParentID Snowflake `json:"parent_id"`

	// LastMessageID is the id of the last message sent in this channel.
	//
	// Always present.
	LastMessageID Snowflake `json:"last_message_id"`

	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message.
	// Bots, as well as users with the permission manageMessages or manageChannel, are unaffected.
	//
	// Always present.
	RateLimitPerUser time.Duration `json:"rate_limit_per_user"`

	// Nsfw indicates whether the channel is NSFW.
	//
	// Always present.
	Nsfw bool `json:"nsfw"`

	// Bitrate is the bitrate (in bits) of the voice channel.
	//
	// Always present.
	Bitrate int `json:"bitrate"`

	// UserLimit is the user limit of the voice channel.
	//
	// Always present.
	UserLimit int `json:"user_limit"`

	// RtcRegion is the voice region id for the voice channel. Automatic when set to empty string.
	//
	// Always present.
	RtcRegion string `json:"rtc_region"`

	// Topic is the channel topic (can be 0-1024 characters).
	//
	// Always present. Can be empty string if the channel has no topic.
	Topic string `json:"topic"`
}

// TODO: continue the thread channels

// AnnouncementThreadChannel represents an announcement thread channel.
type AnnouncementThreadChannel struct {
	baseGuildChannel
}

// PublicThreadChannel represents a public thread channel.
type PublicThreadChannel struct {
	baseGuildChannel
}

// PrivateThreadChannel represents a private thread channel.
type PrivateThreadChannel struct {
	baseGuildChannel
}
