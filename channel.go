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

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"

	"github.com/marouanesouiri/stdx/optional"
	"github.com/marouanesouiri/stdx/result"
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
	Allow Permissions `json:"allow,omitempty"`

	// Deny is the permission bit set explicitly denied.
	Deny Permissions `json:"deny,omitempty"`
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
	EmojiID Snowflake `json:"emoji_id,omitempty"`

	// EmojiName is the Unicode character of the emoji.
	//
	// Optional:
	//  - May be empty string.
	//
	// Note:
	//  - If EmojiName is empty (not set), then EmojiID must be set (non-zero).
	EmojiName string `json:"emoji_name,omitempty"`
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
	ArchiveTimestamp time.Time `json:"archive_timestamp,omitzero"`

	// Locked is whether the thread is locked; when a thread is locked,
	// only users with MANAGE_THREADS can unarchive it
	Locked bool `json:"locked"`

	// Invitable is whether non-moderators can add other non-moderators to a thread.
	Invitable bool `json:"invitable"`
}

// ChannelFields contains only fields present in all channel types.
//
// Reference: https://discord.com/developers/docs/resources/channel#channel-object-channel-structure
type ChannelFields struct {
	// ID is the unique Discord snowflake ID of the channel.
	ID Snowflake `json:"id"`

	// Type is the type of the channel.
	Type ChannelType `json:"type"`
}

func (c *ChannelFields) GetID() Snowflake {
	return c.ID
}

func (c *ChannelFields) GetType() ChannelType {
	return c.Type
}

func (c *ChannelFields) CreatedAt() time.Time {
	return c.ID.Timestamp()
}

// Mention returns a Discord mention string for the channel.
//
// Example output: "<#123456789012345678>"
func (c *ChannelFields) Mention() string {
	return "<#" + c.ID.String() + ">"
}

// String implements the fmt.Stringer interface.
func (c *ChannelFields) String() string {
	return c.Mention()
}

// GuildChannelFields embeds BaseChannel and adds fields common to guild channels except threads.
//
// Used by guild-specific channel types like TextChannel, VoiceChannel, ForumChannel, etc.
type GuildChannelFields struct {
	ChannelFields

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

func (c *GuildChannelFields) GetGuildID() Snowflake {
	return c.GuildID
}

func (c *GuildChannelFields) GetName() string {
	return c.Name
}

func (c *GuildChannelFields) GetPosition() int {
	return c.Position
}

func (c *GuildChannelFields) GetPermissionOverwrites() []PermissionOverwrite {
	return c.PermissionOverwrites
}

func (c *GuildChannelFields) GetFlags() ChannelFlags {
	return c.Flags
}

func (c *GuildChannelFields) JumpURL() string {
	return "https://discord.com/channels/" + c.GuildID.String() + "/" + c.ID.String()
}

// ThreadChannelFields embeds BaseChannel and adds fields common to thread channels.
type ThreadChannelFields struct {
	ChannelFields

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

func (c *ThreadChannelFields) GetGuildID() Snowflake {
	return c.GuildID
}

func (c *ThreadChannelFields) GetName() string {
	return c.Name
}

func (c *ThreadChannelFields) GetPermissionOverwrites() []PermissionOverwrite {
	return c.PermissionOverwrites
}

func (c *ThreadChannelFields) GetFlags() ChannelFlags {
	return c.Flags
}

func (c *ThreadChannelFields) JumpURL() string {
	return "https://discord.com/channels/" + c.GuildID.String() + "/" + c.ID.String()
}

// CategorizedChannelFields holds the parent category field for categorized guild channels.
type CategorizedChannelFields struct {
	// ParentID is the id of the parent category for this channel.
	//
	// Info:
	//  - Each parent category can contain up to 50 channels.
	//
	// Optional:
	//  - May be equal 0 if the channel is not in a category.
	ParentID Snowflake `json:"parent_id"`
}

func (c *CategorizedChannelFields) GetParentID() Snowflake {
	return c.ParentID
}

// MessageChannelFields holds fields related to text-based features like messaging.
type MessageChannelFields struct {
	// LastMessageID is the id of the last message sent in this channel.
	LastMessageID Snowflake `json:"last_message_id"`
}

func (t *MessageChannelFields) GetLastMessageID() Snowflake {
	return t.LastMessageID
}

// GuildMessageChannelFields holds fields related to text-based features like messaging.
type GuildMessageChannelFields struct {
	MessageChannelFields
	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message.
	// Bots, as well as users with the permission manageMessages or manageChannel, are unaffected.
	RateLimitPerUser time.Duration `json:"rate_limit_per_user"`
}

func (t *GuildMessageChannelFields) GetRateLimitPerUser() time.Duration {
	return t.RateLimitPerUser
}

// NsfwChannelFields holds the NSFW indicator field.
type NsfwChannelFields struct {
	// Nsfw indicates whether the channel is NSFW.
	Nsfw bool `json:"nsfw"`
}

// TopicChannelFields holds the topic field.
type TopicChannelFields struct {
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

type Bitrate int

const BitrateMin Bitrate = 8000
const BitrateMaxForStageChannels Bitrate = 64000
const BitrateMaxForNormalGuilds Bitrate = 96000
const BitrateMaxForLevel1Guilds Bitrate = 128000
const BitrateMaxForLevel2Guilds Bitrate = 256000
const BitrateMaxForLevel3Guilds Bitrate = 384000

// AudioChannelFields holds voice-related configuration fields.
type AudioChannelFields struct {
	// Bitrate is the bitrate (in bits) of the voice channel.
	Bitrate int `json:"bitrate"`

	// UserLimit is the user limit of the voice channel.
	UserLimit int `json:"user_limit"`

	// RtcRegion is the voice region id for the voice channel. Automatic when set to empty string.
	RtcRegion string `json:"rtc_region"`
}

func (c *AudioChannelFields) GetBitrate() int {
	return c.Bitrate
}

func (c *AudioChannelFields) GetUserLimit() int {
	return c.UserLimit
}

func (c *AudioChannelFields) GetRtcRegion() string {
	return c.RtcRegion
}

// ForumChannelFields holds forum and media channel specific fields.
type ForumChannelFields struct {
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
	GuildChannelFields
}

func (c *CategoryChannel) MarshalJSON() ([]byte, error) {
	type NoMethod CategoryChannel
	return json.Marshal((*NoMethod)(c))
}

// TextChannel represents a guild text channel.
type TextChannel struct {
	GuildChannelFields
	CategorizedChannelFields
	GuildMessageChannelFields
	NsfwChannelFields
	TopicChannelFields
}

func (c *TextChannel) MarshalJSON() ([]byte, error) {
	type NoMethod TextChannel
	return json.Marshal((*NoMethod)(c))
}

// VoiceChannel represents a guild voice channel.
type VoiceChannel struct {
	GuildChannelFields
	CategorizedChannelFields
	GuildMessageChannelFields
	NsfwChannelFields
	AudioChannelFields
}

func (c *VoiceChannel) MarshalJSON() ([]byte, error) {
	type NoMethod VoiceChannel
	return json.Marshal((*NoMethod)(c))
}

// AnnouncementChannel represents an announcement channel.
type AnnouncementChannel struct {
	GuildChannelFields
	CategorizedChannelFields
	GuildMessageChannelFields
	NsfwChannelFields
	TopicChannelFields
}

func (c *AnnouncementChannel) MarshalJSON() ([]byte, error) {
	type NoMethod AnnouncementChannel
	return json.Marshal((*NoMethod)(c))
}

// StageVoiceChannel represents a stage voice channel.
type StageVoiceChannel struct {
	GuildChannelFields
	CategorizedChannelFields
	GuildMessageChannelFields
	NsfwChannelFields
	AudioChannelFields
	TopicChannelFields
}

func (c *StageVoiceChannel) MarshalJSON() ([]byte, error) {
	type NoMethod StageVoiceChannel
	return json.Marshal((*NoMethod)(c))
}

// ForumChannel represents a guild forum channel.
type ForumChannel struct {
	GuildChannelFields
	CategorizedChannelFields
	GuildMessageChannelFields
	NsfwChannelFields
	TopicChannelFields
	ForumChannelFields
}

func (c *ForumChannel) MarshalJSON() ([]byte, error) {
	type NoMethod ForumChannel
	return json.Marshal((*NoMethod)(c))
}

// MediaChannel represents a media channel.
type MediaChannel struct {
	ForumChannel
}

type ThreadMemberFlags int

// ThreadMember represents Discord thread channel member.
//
// Reference: https://discord.com/developers/docs/resources/channel#channel-object-channel-types
type ThreadMember struct {
	// ThreadID is the id of the thread.
	ThreadID Snowflake `json:"id"`

	// UserID is the id of the member.
	UserID Snowflake `json:"user_id"`

	// JoinTimestamp is the time the user last joined the thread.
	JoinTimestamp time.Time `json:"join_timestamp,omitzero"`

	// Flags are any user-thread settings, currently only used for notifications.
	Flags ThreadMemberFlags `json:"flags"`

	// Member is the guild member object of this thread member.
	//
	// Optional:
	//   - This field is only present when 'with_member' is set to true when calling [ListThreadMembers] or [GetThreadMember].
	//
	// [ListThreadMembers]: https://discord.com/developers/docs/resources/channel#list-thread-members
	// [GetThreadMember]: https://discord.com/developers/docs/resources/channel#get-thread-member
	Member *Member `json:"member"`
}

// ThreadChannel represents the base for thread channels.
type ThreadChannel struct {
	ThreadChannelFields
	CategorizedChannelFields
	GuildMessageChannelFields
	// OwnerID is the id of this thread owner
	OwnerID Snowflake `json:"owner_id"`
	// ThreadMetadata is the metadata that contains a number of thread-specific channel fields.
	ThreadMetadata ThreadMetaData `json:"thread_metadata"`
}

func (c *ThreadChannel) MarshalJSON() ([]byte, error) {
	type NoMethod ThreadChannel
	return json.Marshal((*NoMethod)(c))
}

// DMChannelFields contains fields common to DM and Group DM channels.
type DMChannelFields struct {
	ChannelFields
	MessageChannelFields
}

// ThreadChannel represents a DM channel between the currect user and other user.
type DMChannel struct {
	DMChannelFields
	// Recipients is the list of users participating in the group DM channel.
	//
	// Info:
	//   - Contains the users involved in the group DM, excluding the current user or bot.
	Recipients []User `json:"recipients"`
}

func (c *DMChannel) MarshalJSON() ([]byte, error) {
	type NoMethod DMChannel
	return json.Marshal((*NoMethod)(c))
}

// ThreadChannel represents a DM channel between the currect user and other user.
type GroupDMChannel struct {
	DMChannelFields
	// Icon is the custom icon for the group DM channel.
	//
	// Optional:
	//   - Will be empty string if no icon.
	Icon string `json:"icon"`
}

func (c *GroupDMChannel) MarshalJSON() ([]byte, error) {
	type NoMethod GroupDMChannel
	return json.Marshal((*NoMethod)(c))
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
	json.Marshaler
	GetID() Snowflake
	GetType() ChannelType
	CreatedAt() time.Time
	Mention() string
}

var (
	_ Channel = (*CategoryChannel)(nil)
	_ Channel = (*TextChannel)(nil)
	_ Channel = (*VoiceChannel)(nil)
	_ Channel = (*AnnouncementChannel)(nil)
	_ Channel = (*StageVoiceChannel)(nil)
	_ Channel = (*ForumChannel)(nil)
	_ Channel = (*MediaChannel)(nil)
	_ Channel = (*ThreadChannel)(nil)
	_ Channel = (*DMChannel)(nil)
	_ Channel = (*GroupDMChannel)(nil)
)

// MessageChannel represents a Discord text channel.
//
// This interface extends the Channel interface and adds text-channel-specific fields,
// such as the ID of the last message and the rate limit (slowmode) per user.
//
// Use this interface when you want to handle text channels specifically.
//
// You can convert (assert) it to a concrete type using a type assertion or type switch:
//
// Example usage:
//
//	var ch MessageChannel
//
//	switch c := ch.(type) {
//	case *TextChannel:
//	    fmt.Println("Text channel name:", c.GetName())
//	    fmt.Println("Last message ID:", c.GetLastMessageID())
//	case *VoiceChannel:
//	    fmt.Println("Voiec channel name:", c.GetName())
//	    fmt.Println("Last message ID:", c.GetLastMessageID())
//	case *DMChannel:
//	    fmt.Println("DM channel name:", c.GetName())
//	    fmt.Println("Last message ID:", c.GetLastMessageID())
//	default:
//	    fmt.Println("Other text channel type:", c.GetType())
//	}
//
// You can also use an if-condition to check a specific type:
//
//	if textCh, ok := ch.(*TextChannel); ok {
//	    fmt.Println("Text channel:", textCh.GetName())
//	}
type MessageChannel interface {
	Channel
	// GetLastMessageID returns the Snowflake ID to the last message sent in this channel.
	//
	// Note:
	//   - Will always return 0 if no Message has been sent yet.
	GetLastMessageID() Snowflake
}

var (
	_ MessageChannel = (*TextChannel)(nil)
	_ MessageChannel = (*VoiceChannel)(nil)
	_ MessageChannel = (*AnnouncementChannel)(nil)
	_ MessageChannel = (*StageVoiceChannel)(nil)
	_ MessageChannel = (*ForumChannel)(nil)
	_ MessageChannel = (*MediaChannel)(nil)
	_ MessageChannel = (*ThreadChannel)(nil)
	_ MessageChannel = (*DMChannel)(nil)
	_ MessageChannel = (*GroupDMChannel)(nil)
)

// NamedChannel represents a Discord channel that has a name.
//
// This interface is used for channel types that expose a name, such as text channels,
// voice channels, forum channels, thread channels, DM channels, and Group DM channels.
//
// Use this interface when you want to handle channels generically by their name without
// knowing the specific concrete type in advance.
//
// You can convert (assert) it to a specific channel type using a type assertion or a type
// switch, as described in the official Go documentation:
//   - https://go.dev/ref/spec#Type_assertions
//   - https://go.dev/doc/effective_go#type_switch
//
// Example usage:
//
//	var ch NamedChannel
//
//	// Using a type switch to handle specific channel types
//	switch c := ch.(type) {
//	case *TextChannel:
//	    fmt.Println("Text channel name:", c.GetName())
//	case *VoiceChannel:
//	    fmt.Println("Voice channel name:", c.GetName())
//	default:
//	    fmt.Println("Other named channel type:", c.GetType())
//	}
//
//	// Using a type assertion to check a specific type
//	if textCh, ok := ch.(*TextChannel); ok {
//	    fmt.Println("Text channel name:", textCh.GetName())
//	}
type NamedChannel interface {
	Channel
	GetName() string
}

var (
	_ NamedChannel = (*CategoryChannel)(nil)
	_ NamedChannel = (*TextChannel)(nil)
	_ NamedChannel = (*VoiceChannel)(nil)
	_ NamedChannel = (*AnnouncementChannel)(nil)
	_ NamedChannel = (*StageVoiceChannel)(nil)
	_ NamedChannel = (*ForumChannel)(nil)
	_ NamedChannel = (*MediaChannel)(nil)
	_ NamedChannel = (*ThreadChannel)(nil)
)

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
	NamedChannel
	GetGuildID() Snowflake
	GetPermissionOverwrites() []PermissionOverwrite
	GetFlags() ChannelFlags
	JumpURL() string
}

var (
	_ GuildChannel = (*CategoryChannel)(nil)
	_ GuildChannel = (*TextChannel)(nil)
	_ GuildChannel = (*VoiceChannel)(nil)
	_ GuildChannel = (*AnnouncementChannel)(nil)
	_ GuildChannel = (*StageVoiceChannel)(nil)
	_ GuildChannel = (*ForumChannel)(nil)
	_ GuildChannel = (*MediaChannel)(nil)
	_ GuildChannel = (*ThreadChannel)(nil)
)

// GuildMessageChannel represents a Discord text channel.
//
// This interface extends the Channel interface and adds text-channel-specific fields,
// such as the ID of the last message and the rate limit (slowmode) per user.
//
// Use this interface when you want to handle text channels specifically.
//
// You can convert (assert) it to a concrete type using a type assertion or type switch:
//
// Example usage:
//
//	var ch GuildMessageChannel
//
//	switch c := ch.(type) {
//	case *TextChannel:
//	    fmt.Println("Text channel name:", c.GetName())
//	    fmt.Println("Last message ID:", c.GetLastMessageID())
//	    fmt.Println("Rate limit per user:", c.GetRateLimitPerUser())
//	case *VoiceChannel:
//	    fmt.Println("Voiec channel name:", c.GetName())
//	    fmt.Println("Last message ID:", c.GetLastMessageID())
//	    fmt.Println("Rate limit per user:", c.GetRateLimitPerUser())
//	default:
//	    fmt.Println("Other text channel type:", c.GetType())
//	}
//
// You can also use an if-condition to check a specific type:
//
//	if textCh, ok := ch.(*TextChannel); ok {
//	    fmt.Println("Text channel:", textCh.GetName())
//	}
type GuildMessageChannel interface {
	GuildChannel
	MessageChannel
	GetRateLimitPerUser() time.Duration
}

var (
	_ GuildMessageChannel = (*TextChannel)(nil)
	_ GuildMessageChannel = (*VoiceChannel)(nil)
	_ GuildMessageChannel = (*AnnouncementChannel)(nil)
	_ GuildMessageChannel = (*StageVoiceChannel)(nil)
	_ GuildMessageChannel = (*ForumChannel)(nil)
	_ GuildMessageChannel = (*MediaChannel)(nil)
	_ GuildMessageChannel = (*ThreadChannel)(nil)
)

// PositionedChannel represents a Discord channel that has a sorting position within its parent category.
//
// This interface is used for guild channels that have a defined position, such as category channels, text channels,
// voice channels, announcement channels, stage voice channels, forum channels, and media channels.
// The position determines the order in which channels appear within their parent category in the
// Discord client. If the channel is not under a parent category, the position is relative to other
// top-level channels in the guild.
//
// Use this interface when you want to handle channels generically by their position without knowing
// the specific concrete type in advance.
//
// You can convert (assert) it to a specific channel type using a type assertion or a type switch,
// as described in the official Go documentation:
//   - https://go.dev/ref/spec#Type_assertions
//   - https://go.dev/doc/effective_go#type_switch
//
// Example usage:
//
//	var ch PositionedChannel
//
//	// Using a type switch to handle specific channel types
//	switch c := ch.(type) {
//	case *TextChannel:
//	    fmt.Println("Text channel position:", c.GetPosition())
//	case *VoiceChannel:
//	    fmt.Println("Voice channel position:", c.GetPosition())
//	case *ForumChannel:
//	    fmt.Println("Forum channel position:", c.GetPosition())
//	default:
//	    fmt.Println("Other positioned channel type:", c.GetType())
//	}
//
//	// Using a type assertion to check a specific type
//	if textCh, ok := ch.(*TextChannel); ok {
//	    fmt.Println("Text channel position:", textCh.GetPosition())
//	}
type PositionedChannel interface {
	NamedChannel
	GetPosition() int
}

var (
	_ PositionedChannel = (*CategoryChannel)(nil)
	_ PositionedChannel = (*TextChannel)(nil)
	_ PositionedChannel = (*VoiceChannel)(nil)
	_ PositionedChannel = (*AnnouncementChannel)(nil)
	_ PositionedChannel = (*StageVoiceChannel)(nil)
	_ PositionedChannel = (*ForumChannel)(nil)
	_ PositionedChannel = (*MediaChannel)(nil)
)

// CategorizedChannel represents a Discord channel that can be placed under a parent category channel within a guild.
//
// This interface is used for guild channels that can be organized under a category, such as text channels,
// voice channels, announcement channels, stage voice channels, forum channels, media channels, and thread channels.
//
// Use this interface when you want to handle channels generically by their parent category without knowing
// the specific concrete type in advance.
//
// You can convert (assert) it to a specific channel type using a type assertion or a type switch,
// as described in the official Go documentation:
//   - https://go.dev/ref/spec#Type_assertions
//   - https://go.dev/doc/effective_go#type_switch
//
// Example usage:
//
//	var ch CategorizedChannel
//
//	// Using a type switch to handle specific channel types
//	switch c := ch.(type) {
//	case *TextChannel:
//	    fmt.Println("Text channel parent ID:", c.GetParentID())
//	case *VoiceChannel:
//	    fmt.Println("Voice channel parent ID:", c.GetParentID())
//	case *ThreadChannel:
//	    fmt.Println("Thread channel parent ID:", c.GetParentID())
//	default:
//	    fmt.Println("Other categorized channel type:", c.GetType())
//	}
//
//	// Using a type assertion to check a specific type
//	if textCh, ok := ch.(*TextChannel); ok {
//	    fmt.Println("Text channel parent ID:", textCh.GetParentID())
//	}
type CategorizedChannel interface {
	NamedChannel
	GetParentID() Snowflake
}

var (
	_ CategorizedChannel = (*TextChannel)(nil)
	_ CategorizedChannel = (*VoiceChannel)(nil)
	_ CategorizedChannel = (*AnnouncementChannel)(nil)
	_ CategorizedChannel = (*StageVoiceChannel)(nil)
	_ CategorizedChannel = (*ForumChannel)(nil)
	_ CategorizedChannel = (*MediaChannel)(nil)
	_ CategorizedChannel = (*ThreadChannel)(nil)
)

// AudioChannel represents a Discord channel that supports voice or audio functionality.
//
// This interface is used for guild channels that have voice-related features, such as voice channels
// and stage voice channels. It provides access to audio-specific properties like bitrate, user limit,
// and RTC region.
//
// Note:
//   - DM channels (ChannelTypeDM) and Group DM channels (ChannelTypeGroupDM) support audio features
//     like calls, streams, and webcams for users. However, for bots, these channels are treated as
//     text channels, as bots cannot interact with their audio features (e.g., bots cannot initiate calls in them).
//
// Use this interface when you want to handle audio channels generically without knowing
// the specific concrete type in advance.
//
// You can convert (assert) it to a specific channel type using a type assertion or a type switch,
// as described in the official Go documentation:
//   - https://go.dev/ref/spec#Type_assertions
//   - https://go.dev/doc/effective_go#type_switch
//
// Example usage:
//
//	var ch AudioChannel
//
//	// Using a type switch to handle specific channel types
//	switch c := ch.(type) {
//	case *VoiceChannel:
//	    fmt.Println("Voice channel bitrate:", c.GetBitrate())
//	    fmt.Println("Voice channel user limit:", c.GetUserLimit())
//	    fmt.Println("Voice channel RTC region:", c.GetRtcRegion())
//	case *StageVoiceChannel:
//	    fmt.Println("Stage voice channel bitrate:", c.GetBitrate())
//	    fmt.Println("Stage voice channel user limit:", c.GetUserLimit())
//	    fmt.Println("Stage voice channel RTC region:", c.GetRtcRegion())
//	}
//
//	// Using a type assertion to check a specific type
//	if voiceCh, ok := ch.(*VoiceChannel); ok {
//	    fmt.Println("Voice channel bitrate:", voiceCh.GetBitrate())
//	}
type AudioChannel interface {
	GuildChannel
	GuildMessageChannel
	GetBitrate() int
	GetUserLimit() int
	GetRtcRegion() string
}

var (
	_ AudioChannel = (*VoiceChannel)(nil)
	_ AudioChannel = (*StageVoiceChannel)(nil)
)

// Helper func to Unmarshal any channel type to a Channel interface.
func UnmarshalChannel(buf []byte) (Channel, error) {
	var meta struct {
		Type ChannelType `json:"type"`
	}
	if err := json.Unmarshal(buf, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case ChannelTypeGuildCategory:
		var c CategoryChannel
		return &c, json.Unmarshal(buf, &c)
	case ChannelTypeGuildText:
		var c TextChannel
		return &c, json.Unmarshal(buf, &c)
	case ChannelTypeGuildVoice:
		var c VoiceChannel
		return &c, json.Unmarshal(buf, &c)
	case ChannelTypeGuildAnnouncement:
		var c AnnouncementChannel
		return &c, json.Unmarshal(buf, &c)
	case ChannelTypeGuildStageVoice:
		var c StageVoiceChannel
		return &c, json.Unmarshal(buf, &c)
	case ChannelTypeGuildForum:
		var c ForumChannel
		return &c, json.Unmarshal(buf, &c)
	case ChannelTypeGuildMedia:
		var c MediaChannel
		return &c, json.Unmarshal(buf, &c)
	case ChannelTypeAnnouncementThread,
		ChannelTypePrivateThread,
		ChannelTypePublicThread:
		var c ThreadChannel
		return &c, json.Unmarshal(buf, &c)
	case ChannelTypeDM:
		var c DMChannel
		return &c, json.Unmarshal(buf, &c)
	case ChannelTypeGroupDM:
		var c GroupDMChannel
		return &c, json.Unmarshal(buf, &c)
	default:
		return nil, errors.New("unknown channel type")
	}
}

type ResolvedChannel struct {
	Channel
	Permissions Permissions `json:"permissions"`
}

var _ json.Unmarshaler = (*ResolvedChannel)(nil)

// UnmarshalJSON implements json.Unmarshaler for ResolvedChannel.
func (c *ResolvedChannel) UnmarshalJSON(buf []byte) error {
	var t struct {
		Permissions Permissions `json:"permissions"`
	}
	if err := json.Unmarshal(buf, &t); err != nil {
		return err
	}
	c.Permissions = t.Permissions

	channel, err := UnmarshalChannel(buf)
	if err != nil {
		return err
	}
	c.Channel = channel

	return nil
}

type ResolvedMessageChannel struct {
	MessageChannel
	Permissions Permissions `json:"permissions"`
}

var _ json.Unmarshaler = (*ResolvedMessageChannel)(nil)

// UnmarshalJSON implements json.Unmarshaler for ResolvedMessageChannel.
func (c *ResolvedMessageChannel) UnmarshalJSON(buf []byte) error {
	var t struct {
		Permissions Permissions `json:"permissions"`
	}
	if err := json.Unmarshal(buf, &t); err != nil {
		return err
	}
	c.Permissions = t.Permissions

	channel, err := UnmarshalChannel(buf)
	if err != nil {
		return err
	}
	if messageCh, ok := channel.(MessageChannel); ok {
		c.MessageChannel = messageCh
	} else {
		return errors.New("cannot unmarshal non-MessageChannel into ResolvedMessageChannel")
	}

	return nil
}

type ResolvedThread struct {
	ThreadChannel
	Permissions Permissions `json:"permissions"`
}

type PartialChannel struct {
	ChannelFields

	// Name is the name of the channel.
	//
	// Info:
	//  - can be 1 to 100 characters.
	Name string `json:"name,omitempty"`
}

func (c *PartialChannel) GetName() string {
	return c.Name
}

func (c *PartialChannel) MarshalJSON() ([]byte, error) {
	return json.Marshal(c)
}

var _ NamedChannel = (*PartialChannel)(nil)

type VideoQualityModes int

const (
	VideoQualityModesAuto VideoQualityModes = iota + 1
	VideoQualityModesFull
)

// FetchChannel retrieves a channel by its ID.
//
// Returns a Channel interface which can be asserted to specific types.
func (r *requester) FetchChannel(channelID Snowflake) result.Result[Channel] {
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    "/channels/" + channelID.String(),
	})
	if res.IsErr() {
		return result.Err[Channel](res.Err())
	}
	body := res.Value()
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return result.Err[Channel](err)
	}

	channel, err := UnmarshalChannel(bodyBytes)
	if err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/channels/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Channel](err)
	}
	return result.Ok(channel)
}

func (r *requester) modifyChannel(channelID Snowflake, reqBody []byte, reason string) (Channel, error) {
	res := r.DoRequest(Request{
		Method: "PATCH",
		URL:    "/channels/" + channelID.String(),
		Body:   reqBody,
		Reason: reason,
	})
	if res.IsErr() {
		return nil, res.Err()
	}
	body := res.Value()
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, res.Err()
	}

	channel, err := UnmarshalChannel(bodyBytes)
	if err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PATCH",
			"url":    "/channels/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return nil, res.Err()
	}
	return channel, nil
}

// ModifyGroupDMOptions contains parameters for modifying a Group DM channel.
type ModifyGroupDMOptions struct {
	// Name is the channel's name (1-100 characters).
	Name string `json:"name,omitempty"`

	// Icon sets the icon for Group DM.
	Icon Base64Image `json:"icon,omitempty"`

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// ModifyGroupDMChannel updates a group DM channel's settings.
func (r *requester) ModifyGroupDMChannel(channelID Snowflake, opts ModifyGroupDMOptions) result.Result[*GroupDMChannel] {
	reqBody, _ := json.Marshal(opts)
	channel, err := r.modifyChannel(channelID, reqBody, opts.Reason)
	if err != nil {
		return result.Err[*GroupDMChannel](err)
	}
	if groupdmChannel, ok := channel.(*GroupDMChannel); ok {
		return result.From(groupdmChannel, err)
	}
	return result.Err[*GroupDMChannel](fmt.Errorf("ModifyGroupDMChannel: channel ID %v is not a Group DM", channelID))
}

// ModifyGuildChannelOptions contains parameters for modifying a Guild channel.
type ModifyGuildChannelOptions struct {
	// Name is the channel's name (1-100 characters).
	Name string `json:"name,omitempty"`

	// Type specifies the type of channel to update.
	//
	// Note:
	//  - Defaults to ChannelTypeGuildText if unset.
	//  - Valid values include ChannelTypeGuildText, ChannelTypeGuildVoice, ChannelTypeGuildForum, etc.
	//
	// Applies to All Channels.
	Type optional.Option[ChannelType] `json:"type,omitzero"`

	// Position determines the channel’s position in the server’s channel list (lower numbers appear higher).
	//
	// Note:
	//  - Channels with the same position are sorted by their internal ID.
	//
	// Applies to All Channels.
	Position optional.Option[int] `json:"position,omitzero"`

	// Topic is a description of the channel (0-1024 characters).
	//
	// Note:
	//  - This field is optional.
	//
	// Applies to Channels of Type: Text, Announcement, Forum, Media.
	Topic optional.Option[string] `json:"topic,omitzero"`

	// Nsfw marks the channel as Not Safe For Work, restricting it to 18+ users.
	//
	// Note:
	//  - Set to true to enable the age restriction.
	//
	// Applies to Channels of Type: Text, Voice, Announcement, Stage, Forum.
	Nsfw optional.Option[bool] `json:"nsfw,omitzero"`

	// RateLimitPerUser sets the seconds a user must wait before sending another message (0-21600).
	//
	// Note:
	//  - Bots and users with manage_messages or manage_channel permissions are unaffected.
	//
	// Applies to Channels of Type: Text, Voice, Stage, Forum, Media.
	RateLimitPerUser optional.Option[int] `json:"rate_limit_per_user,omitzero"`

	// Bitrate sets the audio quality for voice or stage channels (in bits, minimum 8000).
	//
	// Note:
	//  - This field is ignored for non-voice channels.
	//
	// Applies to Channels of Type: Voice, Stage.
	Bitrate optional.Option[int] `json:"bitrate,omitzero"`

	// UserLimit caps the number of users in a voice or stage channel (0 for unlimited, 1-99 for a limit).
	//
	// Note:
	//  - Set to 0 to allow unlimited users.
	//
	// Applies to Channels of Type: Voice, Stage.
	UserLimit optional.Option[int] `json:"user_limit,omitzero"`

	// PermissionOverwrites defines custom permissions for specific roles or users.
	//
	// Note:
	//  - This field requires valid overwrite objects.
	//
	// Applies to All Channels.
	PermissionOverwrites optional.Option[[]PermissionOverwrite] `json:"permission_overwrites,omitzero"`

	// ParentID is the ID of the category to nest the channel under.
	//
	// Note:
	//  - This field is ignored for category channels.
	//
	// Applies to Channels of Type: Text, Voice, Announcement, Stage, Forum, Media.
	ParentID optional.Option[Snowflake] `json:"parent_id,omitzero"`

	// RtcRegion sets the voice region id for the voice channel.
	//
	// Note:
	//  - Automatic when set to null (or empty string in Option).
	//
	// Applies to Channels of Type: Voice, Stage.
	RtcRegion optional.Option[string] `json:"rtc_region,omitzero"`

	// VideoQualityMode sets the camera video quality for voice or stage channels.
	//
	// Note:
	//  - Valid options are defined in VideoQualityModes.
	//
	// Applies to Channels of Type: Voice, Stage.
	VideoQualityMode optional.Option[VideoQualityModes] `json:"video_quality_mode,omitzero"`

	// DefaultAutoArchiveDuration sets the default auto archive duration for new threads in this channel.
	DefaultAutoArchiveDuration optional.Option[AutoArchiveDuration] `json:"default_auto_archive_duration,omitzero"`

	// Flags sets the channel flags.
	Flags optional.Option[ChannelFlags] `json:"flags,omitzero"`

	// AvailableTags sets the available tags for a forum/media channel.
	AvailableTags optional.Option[[]ForumTag] `json:"available_tags,omitzero"`

	// DefaultReactionEmoji sets the default reaction emoji for a forum/media channel.
	DefaultReactionEmoji optional.Option[DefaultReactionEmoji] `json:"default_reaction_emoji,omitzero"`

	// DefaultThreadRateLimitPerUser sets the default thread slowmode.
	DefaultThreadRateLimitPerUser optional.Option[int] `json:"default_thread_rate_limit_per_user,omitzero"`

	// DefaultSortOrder sets the default sort order for a forum/media channel.
	DefaultSortOrder optional.Option[ForumPostsSortOrder] `json:"default_sort_order,omitzero"`

	// DefaultForumLayout sets the default layout for a forum channel.
	DefaultForumLayout optional.Option[ForumLayout] `json:"default_forum_layout,omitzero"`

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// ModifyGuildChannel updates a guild channel's settings.
//
// Requires the PermissionManageChannels permission.
func (r *requester) ModifyGuildChannel(channelID Snowflake, opts ModifyGuildChannelOptions) result.Result[GuildChannel] {
	reqBody, _ := json.Marshal(opts)
	channel, err := r.modifyChannel(channelID, reqBody, opts.Reason)
	if err != nil {
		return result.Err[GuildChannel](err)
	}
	if gChannel, ok := channel.(GuildChannel); ok {
		return result.From(gChannel, err)
	}
	return result.Err[GuildChannel](fmt.Errorf("ModifyGuildChannel: channel ID %v is not a Guild channel", channelID))
}

// ModifyGuildThreadOptions contains parameters for modifying a Guild thread's settings.
type ModifyGuildThreadOptions struct {
	// Name is the channel's name (1-100 characters).
	Name string `json:"name,omitempty"`

	// Archived sets the archived state of the thread.
	Archived optional.Option[bool] `json:"archived,omitzero"`

	// AutoArchiveDuration sets the duration after which the thread will automatically archive.
	AutoArchiveDuration optional.Option[AutoArchiveDuration] `json:"auto_archive_duration,omitzero"`

	// Locked sets the locked state of the thread.
	Locked optional.Option[bool] `json:"locked,omitzero"`

	// Invitable sets whether non-moderators can invite others to the thread (private threads only).
	Invitable optional.Option[bool] `json:"invitable,omitzero"`

	// RateLimitPerUser sets the seconds a user must wait before sending another message (0-21600).
	//
	// Note:
	//  - Bots and users with manage_messages or manage_channel permissions are unaffected.
	RateLimitPerUser optional.Option[int] `json:"rate_limit_per_user,omitzero"`

	// Flags sets the channel flags.
	Flags optional.Option[ChannelFlags] `json:"flags,omitzero"`

	// AppliedTags sets the tags applied to a thread.
	AppliedTags optional.Option[[]Snowflake] `json:"applied_tags,omitzero"`

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// ModifyGuildThread updates a guild thread channel's settings.
//
// Requires the PermissionManageThreads permission.
func (r *requester) ModifyGuildThread(channelID Snowflake, opts ModifyGuildThreadOptions) result.Result[*ThreadChannel] {
	reqBody, _ := json.Marshal(opts)
	channel, err := r.modifyChannel(channelID, reqBody, opts.Reason)
	if err != nil {
		return result.Err[*ThreadChannel](err)
	}
	if threadChannel, ok := channel.(*ThreadChannel); ok {
		return result.From(threadChannel, err)
	}
	return result.Err[*ThreadChannel](fmt.Errorf("ModifyGuildThread: channel ID %v is not a Thread channel", channelID))
}

// DeleteChannelOptions contains parameters for deleting a channel, or closing a private message.
type DeleteChannelOptions struct {
	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// DeleteChannel deletes/Close a channel.
//
// Note:
//   - For Community guilds, the Rules or Guidelines channelID
//     and the Community Updates channel cannot be deleted.
//   - Deleting a category does not delete its child channels;
//
// Requires the PermissionManageChannels permission for the guild, or PermissionManageThreads if the channel is a thread.
func (r *requester) DeleteChannel(channelID Snowflake, opts DeleteChannelOptions) result.Result[Channel] {
	res := r.DoRequest(Request{
		Method: "DELETE",
		URL:    "/channels/" + channelID.String(),
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[Channel](res.Err())
	}
	body := res.Value()
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return result.Err[Channel](err)
	}

	channel, err := UnmarshalChannel(bodyBytes)
	if err != nil {
		r.logger.WithFields(map[string]any{
			"method": "DELETE",
			"url":    "/channels/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Channel](err)
	}
	return result.Ok(channel)
}

// EditChannelPermissionsOptions contains parameters for updating a channel overwrite permissions.
type EditChannelPermissionsOptions struct {
	// Allow is the permissions to allow for the overwite.
	Allow optional.Option[Permissions] `json:"allow,omitzero"`

	// Deny is the permissions to deny for the overwite.
	Deny optional.Option[Permissions] `json:"deny,omitzero"`

	// Type is the type of the overwite.
	Type PermissionOverwriteType `json:"type"`

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// EditChannelPermissions edits the channel permission overwrites for a user or role in a channel.
func (r *requester) EditChannelPermissions(channelID Snowflake, overwriteID Snowflake, opts EditChannelPermissionsOptions) result.Void {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "PUT",
		URL:    "/channels/" + channelID.String() + "/permissions/" + overwriteID.String(),
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// FetchChannelInvites Returns a list of invite objects for the channel
//
// Note:
//   - Only usable for guild channels.
//
// Requires the PermissionManageChannels permission.
func (r *requester) FetchChannelInvites(channelID Snowflake) result.Result[[]FullInvite] {
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    "/channels/" + channelID.String() + "/invites",
	})
	if res.IsErr() {
		return result.Err[[]FullInvite](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var invites []FullInvite
	if err := json.NewDecoder(body).Decode(&invites); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/channels/{id}/invites",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]FullInvite](err)
	}
	return result.Ok(invites)
}

// CreateChannelInviteOptions contains parameters for creating a new invite for a channel.
type CreateChannelInviteOptions struct {
	// MaxAge is the duration of invite in seconds before expiry, or 0 for never. between 0 and 604800 (7 days).
	MaxAge optional.Option[int] `json:"max_age,omitzero"`

	// MaxUses is the max number of uses or 0 for unlimited. between 0 and 100.
	MaxUses int `json:"max_uses,omitempty"`

	// Temporary is whether this invite only grants temporary membership.
	Temporary bool `json:"temporary,omitempty"`

	// Unique if true, don't try to reuse a similar invite (useful for creating many unique one time use invites)
	Unique bool `json:"unique,omitzero"`

	// TargetType is the type of target for this voice channel invite.
	TargetType InviteTargetType `json:"target_type,omitzero"`

	// TargetUserID is the id of the user whose stream to display for this invite,
	// required if TargetType is InviteTargetTypeStream, the user must be streaming in the channel.
	TargetUserID Snowflake `json:"target_user_id,omitempty"`

	// TargetApplicationID is the id of the embedded application to open for this invite,
	// required if TargetType is InviteTargetTypeEmbeddedApplication, the application must have the EMBEDDED flag.
	TargetApplicationID Snowflake `json:"target_application_id,omitempty"`

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// CreateChannelInvite creates a new invite object for the channel.
//
// Note:
//   - Only usable for guild channels.
//
// Requires the PermissionCreateInstantInvite permission.
func (r *requester) CreateChannelInvite(channelID Snowflake, opts CreateChannelInviteOptions) result.Result[Invite] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/channels/" + channelID.String() + "/invites",
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[Invite](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var invite Invite
	if err := json.NewDecoder(body).Decode(&invite); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "POST",
			"url":    "/channels/{id}/invites",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Invite](err)
	}
	return result.Ok(invite)
}

// DeleteChannelPermissionOptions contains parameters for deleting a channel permission.
type DeleteChannelPermissionOptions struct {
	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// DeleteChannelPermission deletes a channel permission overwrite for a user or role in a channel.
//
// Note:
//   - Only usable for guild channels.
//
// Requires the PermissionManageRoles permission.
func (r *requester) DeleteChannelPermission(channelID Snowflake, overwriteID Snowflake, opts DeleteChannelPermissionOptions) result.Void {
	res := r.DoRequest(Request{
		Method: "DELETE",
		URL:    "/channels/" + channelID.String() + "/permissions/" + overwriteID.String(),
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// FollowAnnouncementChannelOptions contains parameters for following a Announcement channel.
type FollowAnnouncementChannelOptions struct {
	// WebhookChannelID is the id of target channel.
	WebhookChannelID Snowflake `json:"webhook_channel_id"`

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// FollowedChannel represents a channel that is followed.
type FollowedChannel struct {
	// ChannelID is the source channel id.
	ChannelID Snowflake `json:"channel_id"`

	// WebhookID is the created target webhook id.
	WebhookID Snowflake `json:"webhook_id"`
}

func (c *FollowedChannel) CreatedAt() time.Time {
	return c.ChannelID.Timestamp()
}

func (c *FollowedChannel) FollowedAt() time.Time {
	return c.WebhookID.Timestamp()
}

// Mention returns a Discord mention string for the channel.
//
// Example output: "<#123456789012345678>"
func (c *FollowedChannel) Mention() string {
	return "<#" + c.ChannelID.String() + ">"
}

// String implements the fmt.Stringer interface.
func (c *FollowedChannel) String() string {
	return c.Mention()
}

// FollowAnnouncementChannel follows an Announcement Channel to send messages to a target channel.
//
// Note:
//   - Only usable for guild channels.
//
// Requires the PermissionManageWebhooks permission in the target channel.
func (r *requester) FollowAnnouncementChannel(channelID Snowflake, opts FollowAnnouncementChannelOptions) result.Result[FollowedChannel] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/channels/" + channelID.String() + "/followers",
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[FollowedChannel](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var followedChannel FollowedChannel
	if err := json.NewDecoder(body).Decode(&followedChannel); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "POST",
			"url":    "/channels/{id}/followers",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[FollowedChannel](err)
	}
	return result.Ok(followedChannel)
}

// TriggerTypingIndicator posts a typing indicator for the specified channel.
//
// Note:
//   - The typing indicator expires after 10 seconds.
func (r *requester) TriggerTypingIndicator(channelID Snowflake) result.Void {
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/channels/" + channelID.String() + "/typing",
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// GroupDMAddRecipientOptions contains parameters for adding a recipient to a group dm channel.
type GroupDMAddRecipientOptions struct {
	// AccessToken is the access token of a user that has granted your app the 'gdm.join' scope.
	AccessToken string `json:"access_token"`

	// Nick is the nickname of the user being added.
	Nick string `json:"nick,omitempty"`
}

// GroupDMAddRecipient adds a recipient to a Group DM using their access token.
func (r *requester) GroupDMAddRecipient(channelID Snowflake, userID Snowflake, opts GroupDMAddRecipientOptions) result.Void {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "PUT",
		URL:    "/channels/" + channelID.String() + "/recipients/" + userID.String(),
		Body:   reqBody,
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// GroupDMRemoveRecipient removes a recipient from a Group DM.
func (r *requester) GroupDMRemoveRecipient(channelID Snowflake, userID Snowflake) result.Void {
	res := r.DoRequest(Request{
		Method: "DELETE",
		URL:    "/channels/" + channelID.String() + "/recipients/" + userID.String(),
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// StartThreadFromMessageOptions contains parameters for starting a thread from a message.
type StartThreadFromMessageOptions struct {
	// Name is a 1-100 character channel name
	Name string `json:"name"`

	// AutoArchiveDuration is the number of minutes of inactivity after which the
	// thread will be automatically archived and stop showing in the channel list.
	// Valid values are: 60, 1440, 4320, and 10080.
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration,omitempty"`

	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message (0-21600).
	RateLimitPerUser int `json:"rate_limit_per_user"`

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// StartThreadFromMessage creates a new thread from an existing message.
//
// Note:
//
//   - When called on a 'GuildText' channel, creates a 'PublicThread'.
//     When called on a 'GuildAnnouncement' channel, creates a 'AnnouncementThread'.
//     Does not work on a 'GuildForum' or a 'GuildMedia' channel.
//
//   - The id of the created thread will be the same as the id of the source message,
//     and as such a message can only have a single thread created from it.
//
// Requires the PermissionManageWebhooks permission in the target channel.
func (r *requester) StartThreadFromMessage(channelID Snowflake, messageID Snowflake, opts StartThreadFromMessageOptions) result.Result[GuildChannel] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/channels/" + channelID.String() + "/messages/" + messageID.String() + "/threads",
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[GuildChannel](res.Err())
	}
	body := res.Value()
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return result.Err[GuildChannel](err)
	}

	channel, err := UnmarshalChannel(bodyBytes)
	if err != nil {
		r.logger.WithFields(map[string]any{
			"method": "POST",
			"url":    "/channels/{id}/messages/{msg_id}/threads",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildChannel](err)
	}
	return result.Ok(channel.(GuildChannel))
}

// StartThreadWithoutMessageOptions contains parameters for starting a thread without a message.
type StartThreadWithoutMessageOptions struct {
	// Type the type of thread to create
	//
	// Can be one of: ChannelTypeAnnouncementThread, ChannelTypePublicThread or ChannelTypePrivateThread.
	Type ChannelType `json:"type,omitzero"`

	// Name is a 1-100 character channel name
	Name string `json:"name"`

	// AutoArchiveDuration is the number of minutes of inactivity after which the
	// thread will be automatically archived and stop showing in the channel list.
	// Valid values are: 60, 1440, 4320, and 10080.
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration,omitempty"`

	// RateLimitPerUser is the amount of seconds a user has to wait before sending another message (0-21600).
	RateLimitPerUser int `json:"rate_limit_per_user"`

	// Invitable is whether non-moderators can add other non-moderators to a thread
	//
	// Note:
	//   - only available when creating a private thread.
	Invitable optional.Option[bool] `json:"invitable,omitzero"`

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// StartThreadWithoutMessage creates a new thread.
func (r *requester) StartThreadWithoutMessage(channelID Snowflake, opts StartThreadWithoutMessageOptions) result.Result[GuildChannel] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/channels/" + channelID.String() + "/threads",
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[GuildChannel](res.Err())
	}
	body := res.Value()
	defer body.Close()

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return result.Err[GuildChannel](err)
	}

	channel, err := UnmarshalChannel(bodyBytes)
	if err != nil {
		r.logger.WithFields(map[string]any{
			"method": "POST",
			"url":    "/channels/{id}/threads",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildChannel](err)
	}
	return result.Ok(channel.(GuildChannel))
}

// JoinThread adds the current user to a thread.
func (r *requester) JoinThread(channelID Snowflake) result.Void {
	res := r.DoRequest(Request{
		Method: "PUT",
		URL:    "/channels/" + channelID.String() + "/thread-members/@me",
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// AddThreadMember adds a member to a thread.
func (r *requester) AddThreadMember(channelID Snowflake, userID Snowflake) result.Void {
	res := r.DoRequest(Request{
		Method: "PUT",
		URL:    "/channels/" + channelID.String() + "/thread-members/" + userID.String(),
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// LeaveThread removes the current user from a thread.
func (r *requester) LeaveThread(channelID Snowflake) result.Void {
	res := r.DoRequest(Request{
		Method: "DELETE",
		URL:    "/channels/" + channelID.String() + "/thread-members/@me",
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// RemoveThreadMember removes a member from a thread.
//
// Requires the PermissionManageThreads permission, or the creator of the thread if it is a 'PrivateThread'.
func (r *requester) RemoveThreadMember(channelID Snowflake, userID Snowflake) result.Void {
	res := r.DoRequest(Request{
		Method: "DELETE",
		URL:    "/channels/" + channelID.String() + "/thread-members/" + userID.String(),
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// FetchThreadMemberOptions contains parameters for fetching a thread member.
type FetchThreadMemberOptions struct {
	// WithMember is whether to include a guild member object for the thread member.
	WithMember bool `json:"member,omitempty"`
}

// FetchThreadMember retrieves a thread member.
func (r *requester) FetchThreadMember(channelID Snowflake, userID Snowflake, opts FetchThreadMemberOptions) result.Result[ThreadMember] {
	endpoint := "/channels/" + channelID.String() + "/thread-members/" + userID.String() + "?with_member=" + strconv.FormatBool(opts.WithMember)
	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[ThreadMember](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var member ThreadMember
	if err := json.NewDecoder(body).Decode(&member); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/channels/{id}/thread-members/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[ThreadMember](err)
	}
	return result.Ok(member)
}

// ListThreadMembersOptions contains parameters for listing a thread members.
type ListThreadMembersOptions struct {
	// WithMember is whether to include a guild member object for the thread member.
	WithMember bool `json:"member,omitempty"`

	// Limit is the maximum number of members to return (1-100).
	//
	//  Note:
	//   - Defaults to 100 if not specified.
	Limit int `json:"limit,omitempty"`

	// After is the member ID to start after for pagination.
	// Used to get the next page of results.
	After Snowflake `json:"after,omitempty"`
}

// ListThreadMembers retrieves a list of thread members.
func (r *requester) ListThreadMembers(channelID Snowflake, opts ListThreadMembersOptions) result.Result[[]ThreadMember] {
	params := url.Values{}
	params.Set("with_member", strconv.FormatBool(opts.WithMember))

	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if !opts.After.UnSet() {
		params.Set("after", opts.After.String())
	}

	endpoint := "/channels/" + channelID.String() + "/thread-members?" + params.Encode()

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[[]ThreadMember](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var member []ThreadMember
	if err := json.NewDecoder(body).Decode(&member); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/channels/{id}/thread-members",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]ThreadMember](err)
	}
	return result.Ok(member)
}

// ListArchivedThreadsOptions contains parameters for listing archived threads.
type ListArchivedThreadsOptions struct {
	// Limit is the maximum number of members to return (1-100).
	//
	//  Note:
	//   - Defaults to 100 if not specified.
	Limit int `json:"limit,omitempty"`

	// Returns threads archived before this timestamp.
	Before time.Time `json:"before,omitzero"`
}

type ListArchivedThreadsResponse struct {
	// Threads are the archived threads.
	Threads []ThreadChannel `json:"threads"`

	// A thread member object for each returned thread the current user has joined.
	Members []ThreadMember `json:"members"`

	// HasMore is whether there are potentially additional threads that could be returned on a subsequent call.
	HasMore bool `json:"has_more"`
}

func (r *requester) listArchivedThreads(channelID Snowflake, subEndPoint string, params url.Values) result.Result[ListArchivedThreadsResponse] {
	endpoint := "/channels/" + channelID.String() + subEndPoint
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[ListArchivedThreadsResponse](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var response ListArchivedThreadsResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/channels/{id}" + subEndPoint,
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[ListArchivedThreadsResponse](err)
	}
	return result.Ok(response)
}

// ListPublicArchivedThreads retrieves a list of public archived threads.
func (r *requester) ListPublicArchivedThreads(channelID Snowflake, opts ListArchivedThreadsOptions) result.Result[ListArchivedThreadsResponse] {
	params := url.Values{}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if !opts.Before.IsZero() {
		params.Set("before", opts.Before.Format(time.RFC3339))
	}
	return r.listArchivedThreads(channelID, "/threads/archived/public", params)
}

// ListPrivateArchivedThreads retrieves a list of private archived threads.
func (r *requester) ListPrivateArchivedThreads(channelID Snowflake, opts ListArchivedThreadsOptions) result.Result[ListArchivedThreadsResponse] {
	params := url.Values{}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if !opts.Before.IsZero() {
		params.Set("before", opts.Before.Format(time.RFC3339))
	}
	return r.listArchivedThreads(channelID, "/threads/archived/private", params)
}

// ListJoinedPrivateArchivedThreadsOptions contains parameters for listing joined private archived threads.
type ListJoinedPrivateArchivedThreadsOptions struct {
	// Limit is the maximum number of members to return (1-100).
	//
	//  Note:
	//   - Defaults to 100 if not specified.
	Limit int `json:"limit,omitempty"`

	// Returns threads archived before this id.
	Before Snowflake `json:"before,omitzero"`
}

// ListJoinedPrivateArchivedThreads retrieves a list of private archived threads that the current user has joined.
func (r *requester) ListJoinedPrivateArchivedThreads(channelID Snowflake, opts ListJoinedPrivateArchivedThreadsOptions) result.Result[ListArchivedThreadsResponse] {
	params := url.Values{}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if !opts.Before.UnSet() {
		params.Set("before", opts.Before.String())
	}
	return r.listArchivedThreads(channelID, "/users/@me/threads/archived/private", params)
}
