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

import (
	"errors"
	"time"

	"github.com/bytedance/sonic"
)

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
	return BitFieldHas(f, flags...)
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

// Is returns true if the overWrite's Type matches the provided one.
func (t PermissionOverwriteType) Is(overWriteType PermissionOverwriteType) bool {
	return t == overWriteType
}

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

// ForumTag represents a tag that can be applied to a thread
// in a GuildForum or GuildMedia channel.
//
// Reference: https://discord.com/developers/docs/resources/channel#forum-tag-object
type ForumTag struct {
	// ID is the id of the tag.
	ID Snowflake `json:"id"`

	// Name is the name of the tag (0-20 characters).
	Name string `json:"name"`

	// Moderated indicates whether this tag can only be added to or removed from
	// threads by a member with the ManageThreads permission.
	Moderated bool `json:"moderated"`

	// EmojiID is the ID of a guild's custom emoji.
	//
	// Optional:
	//  - May be equal 0.
	//
	// Note:
	//  - If EmojiName is empty (not set), then EmojiID must be set (non-zero).
	EmojiID Snowflake `json:"emoji_id"`

	// EmojiName is the Unicode character of the emoji.
	//
	// Optional:
	//  - May be empty string.
	//
	// Note:
	//  - If EmojiName is empty (not set), then EmojiID must be set (non-zero).
	EmojiName string `json:"emoji_name"`
}

// DefaultReactionEmoji represents a default reaction emoji for forum channels.
type DefaultReactionEmoji struct {
	// EmojiID is the ID of a guild's custom emoji.
	//
	// Optional:
	//  - May be equal to 0.
	//
	// Info:
	//  - If 0, EmojiName will be set instead.
	EmojiID Snowflake `json:"emoji_id"`

	// EmojiName is the Unicode character of the emoji.
	//
	// Optional:
	//  - May be empty string.
	//
	// Info:
	//  - If empty, EmojiID will be set instead.
	EmojiName string `json:"emoji_name"`
}

// AutoArchiveDuration represents the auto archive duration of a thread channel
//
// Reference: https://discord.com/developers/docs/resources/channel#thread-metadata-object
type AutoArchiveDuration int

const (
	AutoArchiveDuration1h  AutoArchiveDuration = 60
	AutoArchiveDuration24h AutoArchiveDuration = 1440
	AutoArchiveDuration3d  AutoArchiveDuration = 4320
	AutoArchiveDuration1w  AutoArchiveDuration = 10080
)

// Is returns true if the thread's auto archive duration matches the provided auto archive duration.
func (d AutoArchiveDuration) Is(duration AutoArchiveDuration) bool {
	return d == duration
}

// ThreadMetaData represents the metadata object that contains a number of thread-specific channel fields.
//
// Reference: https://discord.com/developers/docs/resources/channel#thread-metadata-object
type ThreadMetaData struct {
	// Archived is whether the thread is archived
	Archived bool `json:"archived"`

	// AutoArchiveDuration is the duration will thread need to stop showing in the channel list.
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`

	// ArchiveTimestamp is the timestamp when the thread's archive status was last changed,
	// used for calculating recent activity
	ArchiveTimestamp time.Time `json:"archive_timestamp"`

	// Locked is whether the thread is locked; when a thread is locked,
	// only users with MANAGE_THREADS can unarchive it
	Locked bool `json:"locked"`

	// Invitable is whether non-moderators can add other non-moderators to a thread.
	Invitable bool `json:"invitable"`
}

// Channel is the interface representing a Discord channel.
//
// This interface can represent any type of channel returned by Discord,
// including text channels, voice channels, thread channels, forum channels, etc.
//
// Use this interface when you want to handle channels generically without knowing
// the specific concrete type in advance.
//
// You can convert (assert) it to a specific channel type using a type assertion or
// a type switch, as described in the official Go documentation:
//   - https://go.dev/ref/spec#Type_assertions
//   - https://go.dev/doc/effective_go#type_switch
//
// Example usage:
//
//	var myChannel Channel
//
//	switch c := ch.(type) {
//	case *TextChannel:
//	    fmt.Println("Text channel name:", c.Name)
//	case *VoiceChannel:
//	    fmt.Println("Voice channel bitrate:", c.Bitrate)
//	case *ForumChannel:
//	    fmt.Println("Forum channel tags:", c.AvailableTags)
//	default:
//	    fmt.Println("Other channel type:", c.GetType())
//	}
//
// You can also use an if-condition to check a specific type:
//
//	if textCh, ok := ch.(*TextChannel); ok {
//	    fmt.Println("Text channel:", textCh.Name)
//	}
type Channel interface {
	GetID() Snowflake
	GetType() ChannelType
	CreatedAt() time.Time
	Mention() string
}

// GuildChannel represents a guild-specific Discord channel.
//
// This interface extends the Channel interface and adds guild-specific fields,
// such as the guild ID, channel name, permission overwrites, flags, and jump URL.
//
// Use this interface when you want to handle guild channels generically without
// knowing the specific concrete type (TextChannel, VoiceChannel, ForumChannel, etc.).
//
// You can convert (assert) it to a specific guild channel type using a type assertion
// or a type switch, as described in the official Go documentation:
//   - https://go.dev/ref/spec#Type_assertions
//   - https://go.dev/doc/effective_go#type_switch
//
// Example usage:
//
//	var myGuildChannel GuildChannel
//
//	switch c := ch.(type) {
//	case *TextChannel:
//	    fmt.Println("Text channel name:", c.Name)
//	case *VoiceChannel:
//	    fmt.Println("Voice channel bitrate:", c.Bitrate)
//	case *ForumChannel:
//	    fmt.Println("Forum channel tags:", c.AvailableTags)
//	default:
//	    fmt.Println("Other guild channel type:", c.GetType())
//	}
//
// You can also use an if-condition to check a specific type:
//
//	if textCh, ok := ch.(*TextChannel); ok {
//	    fmt.Println("Text channel:", textCh.Name)
//	}
type GuildChannel interface {
	Channel
	GetGuildID() Snowflake
	GetName() string
	GetPermissionOverwrites() []PermissionOverwrite
	GetFlags() ChannelFlags
	JumpURL() string
}

// BaseChannel contains only fields present in all channel types.
//
// Reference: https://discord.com/developers/docs/resources/channel#channel-object-channel-structure
type BaseChannel struct {
	// ID is the unique Discord snowflake ID of the channel.
	ID Snowflake `json:"id"`

	// Type is the type of the channel.
	Type ChannelType `json:"type"`
}

func (c *BaseChannel) GetID() Snowflake     { return c.ID }
func (c *BaseChannel) GetType() ChannelType { return c.Type }
func (c *BaseChannel) CreatedAt() time.Time { return c.ID.Timestamp() }

// Mention returns a Discord mention string for the channel.
//
// Example output: "<#123456789012345678>"
func (c *BaseChannel) Mention() string {
	return "<#" + c.ID.String() + ">"
}

// BaseGuildChannel embeds BaseChannel and adds fields common to guild channels except threads.
//
// Used by guild-specific channel types like TextChannel, VoiceChannel, ForumChannel, etc.
type BaseGuildChannel struct {
	BaseChannel

	// GuildID is the id of the guild.
	GuildID Snowflake `json:"guild_id"`

	// Name is the name of the channel.
	//
	// Info:
	//  - can be 1 to 100 characters.
	Name string `json:"name,omitempty"`

	// Position is the sorting position of the channel.
	Position int `json:"position,omitempty"`

	// PermissionOverwrites are explicit permission overwrites for members and roles.
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`

	// Flags are combined channel flags.
	Flags ChannelFlags `json:"flags,omitempty"`
}

func (c *BaseGuildChannel) GetGuildID() Snowflake { return c.GuildID }
func (c *BaseGuildChannel) GetName() string       { return c.Name }
func (c *BaseGuildChannel) GetPosition() int      { return c.Position }
func (c *BaseGuildChannel) GetPermissionOverwrites() []PermissionOverwrite {
	return c.PermissionOverwrites
}
func (c *BaseGuildChannel) GetFlags() ChannelFlags { return c.Flags }
func (c *BaseGuildChannel) JumpURL() string {
	return "https://discord.com/channels/" + c.GuildID.String() + "/" + c.ID.String()
}

// BaseThreadChannel embeds BaseChannel and adds fields common to thread channels.
type BaseThreadChannel struct {
	BaseChannel

	// GuildID is the id of the guild.
	GuildID Snowflake `json:"guild_id"`

	// Name is the name of the channel.
	//
	// Info:
	//  - can be 1 to 100 characters.
	Name string `json:"name,omitempty"`

	// PermissionOverwrites are explicit permission overwrites for members and roles.
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`

	// Flags are combined channel flags.
	Flags ChannelFlags `json:"flags,omitempty"`
}

func (c *BaseThreadChannel) GetGuildID() Snowflake { return c.GuildID }
func (c *BaseThreadChannel) GetName() string       { return c.Name }
func (c *BaseThreadChannel) GetPermissionOverwrites() []PermissionOverwrite {
	return c.PermissionOverwrites
}
func (c *BaseThreadChannel) GetFlags() ChannelFlags { return c.Flags }
func (c *BaseThreadChannel) JumpURL() string {
	return "https://discord.com/channels/" + c.GuildID.String() + "/" + c.ID.String()
}

// CategorizedFields holds the parent category field for categorized guild channels.
type CategorizedFields struct {
	// ParentID is the id of the parent category for this channel.
	//
	// Info:
	//  - Each parent category can contain up to 50 channels.
	//
	// Optional:
	//  - May be equal 0 if the channel is not in a category.
	ParentID Snowflake `json:"parent_id"`
}

// TextBasedFields holds fields related to text-based features like messaging.
type TextBasedFields struct {
	// LastMessageID is the id of the last message sent in this channel.
	LastMessageID Snowflake `json:"last_message_id"`

	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message.
	// Bots, as well as users with the permission manageMessages or manageChannel, are unaffected.
	RateLimitPerUser time.Duration `json:"rate_limit_per_user"`
}

// NsfwFields holds the NSFW indicator field.
type NsfwFields struct {
	// Nsfw indicates whether the channel is NSFW.
	Nsfw bool `json:"nsfw"`
}

// TopicFields holds the topic field.
type TopicFields struct {
	// Topic is the channel topic.
	//
	// Length:
	//  - 0-1024 characters for text, announcement, and stage voice channels.
	//  - 0-4096 characters for forum and media channels.
	//
	// Optional:
	//  - May be empty string if the channel has no topic.
	Topic string `json:"topic"`
}

// VoiceFields holds voice-related configuration fields.
type VoiceFields struct {
	// Bitrate is the bitrate (in bits) of the voice channel.
	Bitrate int `json:"bitrate"`

	// UserLimit is the user limit of the voice channel.
	UserLimit int `json:"user_limit"`

	// RtcRegion is the voice region id for the voice channel. Automatic when set to empty string.
	RtcRegion string `json:"rtc_region"`
}

// ForumFields holds forum and media channel specific fields.
type ForumFields struct {
	// AvailableTags is the set of tags that can be used in this channel.
	AvailableTags []ForumTag `json:"available_tags"`

	// DefaultReactionEmoji specifies the emoji used as the default way to react to a forum post.
	DefaultReactionEmoji DefaultReactionEmoji `json:"default_reaction_emoji"`

	// DefaultSortOrder is the default sort order type used to order posts
	// in GuildForum and GuildMedia channels. Defaults to PostsSortOrderLatestActivity.
	DefaultSortOrder ForumPostsSortOrder `json:"default_sort_order"`

	// DefaultForumLayout is the default forum layout view used to display posts
	// in GuildForum channels. Defaults to ForumLayoutNotSet.
	DefaultForumLayout ForumLayout `json:"default_forum_layout"`
}

// CategoryChannel represents a guild category channel.
type CategoryChannel struct {
	BaseGuildChannel
}

// TextChannel represents a guild text channel.
type TextChannel struct {
	BaseGuildChannel
	CategorizedFields
	TextBasedFields
	NsfwFields
	TopicFields
}

// VoiceChannel represents a guild voice channel.
type VoiceChannel struct {
	BaseGuildChannel
	CategorizedFields
	TextBasedFields
	NsfwFields
	VoiceFields
}

// AnnouncementChannel represents an announcement channel.
type AnnouncementChannel struct {
	BaseGuildChannel
	CategorizedFields
	TextBasedFields
	NsfwFields
	TopicFields
}

// StageVoiceChannel represents a stage voice channel.
type StageVoiceChannel struct {
	BaseGuildChannel
	CategorizedFields
	TextBasedFields
	NsfwFields
	VoiceFields
	TopicFields
}

// ForumChannel represents a guild forum channel.
type ForumChannel struct {
	BaseGuildChannel
	CategorizedFields
	TextBasedFields
	NsfwFields
	TopicFields
	ForumFields
}

// MediaChannel represents a media channel.
type MediaChannel struct {
	ForumChannel
}

// ThreadChannel represents the base for thread channels.
type ThreadChannel struct {
	BaseThreadChannel
	CategorizedFields
	TextBasedFields
	// OwnerID is the id of this thread owner
	OwnerID Snowflake `json:"owner_id"`
	// ThreadMetadata is the metadata that contains a number of thread-specific channel fields.
	ThreadMetadata ThreadMetaData `json:"thread_metadata"`
}

var (
	// All channels must implement Channel.
	_ Channel = (*CategoryChannel)(nil)
	_ Channel = (*TextChannel)(nil)
	_ Channel = (*VoiceChannel)(nil)
	_ Channel = (*AnnouncementChannel)(nil)
	_ Channel = (*StageVoiceChannel)(nil)
	_ Channel = (*ForumChannel)(nil)
	_ Channel = (*MediaChannel)(nil)
	_ Channel = (*ThreadChannel)(nil)
	// All guild channels must implement GuildChannel.
	_ GuildChannel = (*CategoryChannel)(nil)
	_ GuildChannel = (*TextChannel)(nil)
	_ GuildChannel = (*VoiceChannel)(nil)
	_ GuildChannel = (*AnnouncementChannel)(nil)
	_ GuildChannel = (*StageVoiceChannel)(nil)
	_ GuildChannel = (*ForumChannel)(nil)
	_ GuildChannel = (*MediaChannel)(nil)
)

func channelFromJson(buf []byte) (Channel, error) {
	var meta struct {
		Type ChannelType `json:"type"`
	}
	if err := sonic.Unmarshal(buf, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case ChannelTypeGuildCategory:
		var c CategoryChannel
		return &c, sonic.Unmarshal(buf, &c)
	case ChannelTypeGuildText:
		var c TextChannel
		return &c, sonic.Unmarshal(buf, &c)
	case ChannelTypeGuildVoice:
		var c VoiceChannel
		return &c, sonic.Unmarshal(buf, &c)
	case ChannelTypeGuildAnnouncement:
		var c AnnouncementChannel
		return &c, sonic.Unmarshal(buf, &c)
	case ChannelTypeGuildStageVoice:
		var c StageVoiceChannel
		return &c, sonic.Unmarshal(buf, &c)
	case ChannelTypeGuildForum:
		var c ForumChannel
		return &c, sonic.Unmarshal(buf, &c)
	case ChannelTypeGuildMedia:
		var c MediaChannel
		return &c, sonic.Unmarshal(buf, &c)
	case ChannelTypeAnnouncementThread, ChannelTypePrivateThread, ChannelTypePublicThread:
		var c ThreadChannel
		return &c, sonic.Unmarshal(buf, &c)
	default:
		return nil, errors.New("unknown channel type")
	}
}
