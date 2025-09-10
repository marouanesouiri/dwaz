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
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/bytedance/sonic"
)

// VerificationLevel represents the verification level required on a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-object-verification-level
type VerificationLevel int

const (
	// Unrestricted.
	VerificationLevelNone VerificationLevel = iota
	// Must have verified email on account.
	VerificationLevelLow
	// Must be registered on Discord for longer than 5 minutes.
	VerificationLevelMedium
	// Must be a member of the server for longer than 10 minutes.
	VerificationLevelHigh
	// Must have a verified phone number
	VerificationLevelVeryHigh
)

// Is returns true if the verification level matches the provided one.
func (l VerificationLevel) Is(verifLevel VerificationLevel) bool {
	return l == verifLevel
}

// MessageNotificationLevel represents the default notification level on a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-object-default-message-notification-level
type MessageNotificationsLevel int

const (
	// Members will receive notifications for all messages by default.
	MessageNotificationsLevelAllMessages MessageNotificationsLevel = iota
	// Members will receive notifications only for messages that @mention them by default.
	MessageNotificationsLevelOnlyMentions
)

// Is returns true if the message notifaction level matches the provided one.
func (l MessageNotificationsLevel) Is(messageNotificationLevel MessageNotificationsLevel) bool {
	return l == messageNotificationLevel
}

// ExplicitContentFilterLevel represents the explicit content filter level on a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-object-explicit-content-filter-level
type ExplicitContentFilterLevel int

const (
	// Media content will not be scanned.
	ExplicitContentFilterLevelDisabled ExplicitContentFilterLevel = iota
	// Media content sent by members without roles will be scanned.
	ExplicitContentFilterLevelMembersWithoutRoles
	// Media content sent by all members will be scanned
	ExplicitContentFilterLevelAllMembers
)

// Is returns true if the explicit content level matches the provided one.
func (l ExplicitContentFilterLevel) Is(level ExplicitContentFilterLevel) bool {
	return l == level
}

// ExplicitContentFilterLevel represents the mfa level on a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-object-mfa-level
type MFALevel int

const (
	// Guild has no MFA/2FA requirement for moderation actions.
	MFALevelNone MFALevel = iota
	// Guild has a 2FA requirement for moderation actions.
	MFALevelElevated
)

// Is returns true if the MFA level matches the provided one.
func (l MFALevel) Is(level MFALevel) bool {
	return l == level
}

// GuildFeature represents the features of a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-object-guild-features
type GuildFeature string

const (
	// Guild has access to set an animated guild banner image.
	GuildFeatureAnimatedBanner GuildFeature = "ANIMATED_BANNER"
	// Guild has access to set an animated guild icon.
	GuildFeatureAnimatedIcon GuildFeature = "ANIMATED_ICON"
	// Guild is using the old permissions configuration behavior.
	//
	// Reference: https://discord.com/developers/docs/change-log#upcoming-application-command-permission-changes
	GuildFeatureAPPLICATION_COMMAND_PERMISSIONS_V2 GuildFeature = "APPLICATION_COMMAND_PERMISSIONS_V2"
	// guild has set up auto moderation rules
	GuildFeatureAutoModeration GuildFeature = "AUTO_MODERATION"
	// Guild has access to set a guild banner image.
	GuildFeatureBanner GuildFeature = "BANNER"
	// Guild can enable welcome screen, Membership Screening, stage channels and discovery, and receives community updates.
	GuildFeatureCommunity GuildFeature = "COMMUNITY"
	// Guild has enabled monetization
	GuildFeatureCreatorMonetizableProvisional GuildFeature = "CREATOR_MONETIZABLE_PROVISIONAL"
	// Guild has enabled the role subscription promo page.
	GuildFeatureCreatorStorePage GuildFeature = "CREATOR_STORE_PAGE"
	// Guild has been set as a support server on the App Directory.
	GuildFeatureDeveloperSupportServer GuildFeature = "DEVELOPER_SUPPORT_SERVER"
	// Guild is able to be discovered in the directory.
	GuildFeatureDiscoverable GuildFeature = "DISCOVERABLE"
	// Guild is able to be featured in the directory.
	GuildFeatureFeaturable GuildFeature = "FEATURABLE"
	// Guild has paused invites, preventing new users from joining.
	GuildFeatureInvitesDisabled GuildFeature = "INVITES_DISABLED"
	// Guild has access to set an invite splash background.
	GuildFeatureInviteSplash GuildFeature = "INVITE_SPLASH"
	// Guild has enabled Membership Screening.
	//
	// Reference: https://discord.com/developers/docs/resources/guild#membership-screening-object
	GuildFeatureMemberVerificationGateEnabled GuildFeature = "MEMBER_VERIFICATION_GATE_ENABLED"
	// Guild has increased custom soundboard sound slots.
	GuildFeatureMoreSoundboard GuildFeature = "MORE_SOUNDBOARD"
	// Guild has increased custom sticker slots.
	GuildFeatureMoreStickers GuildFeature = "MORE_STICKERS"
	// Guild has access to create announcement channels.
	GuildFeatureNews GuildFeature = "NEWS"
	// Guild is partnered.
	GuildFeaturePartnered GuildFeature = "PARTNERED"
	// Guild can be previewed before joining via Membership Screening or the directory.
	GuildFeaturePreviewEnabled GuildFeature = "PREVIEW_ENABLED"
	// Guild has disabled alerts for join raids in the configured safety alerts channel
	GuildFeatureRaidAlertsDisabled GuildFeature = "RAID_ALERTS_DISABLED"
	// Guild is able to set role icons.
	GuildFeatureRoleIcons GuildFeature = "ROLE_ICONS"
	// Guild has role subscriptions that can be purchased.
	GuildFeatureRoleSubscriptionsAvailableForPurchase GuildFeature = "ROLE_SUBSCRIPTIONS_AVAILABLE_FOR_PURCHASE"
	// Guild has enabled role subscriptions.
	GuildFeatureRoleSubscriptionsEnabled GuildFeature = "ROLE_SUBSCRIPTIONS_ENABLED"
	// Guild has created soundboard sounds.
	GuildFeatureSoundboard GuildFeature = "SOUNDBOARD"
	// Guild has enabled ticketed events.
	GuildFeatureTicketedEventsEnabled GuildFeature = "TICKETED_EVENTS_ENABLED"
	// Guild has access to set a vanity URL.
	GuildFeatureVanityURL GuildFeature = "VANITY_URL"
	// Guild is verified.
	GuildFeatureVerified GuildFeature = "VERIFIED"
	// Guild has access to set 384kbps bitrate in voice (previously VIP voice servers).
	GuildFeatureVipRegions GuildFeature = "VIP_REGIONS"
	// Guild has enabled the welcome screen.
	GuildFeatureWelcomeScreenEnabled GuildFeature = "WELCOME_SCREEN_ENABLED"
	// Guild has access to guest invites.
	GuildFeatureGuestsEnabled GuildFeature = "GUESTS_ENABLED"
	// Guild has access to set guild tags.
	GuildFeatureGuildTags GuildFeature = "GUILD_TAGS"
	// Guild is able to set gradient colors to roles.
	GuildFeatureEnhancedRoleColors GuildFeature = "ENHANCED_ROLE_COLORS"
)

// SystemChannelFlags contains the settings for the Guild(s) system channel
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-object-system-channel-flags
type SystemChannelFlags int

const (
	// Suppress member join notifications.
	SystemChannelFlagSuppressJoinNotifications SystemChannelFlags = 1 << iota
	// Suppress server boost notifications.
	SystemChannelFlagSuppressPremiumSubscriptions
	// Suppress server setup tips.
	SystemChannelFlagSuppressGuildReminderNotifications
	// Hide member join sticker reply buttons.
	SystemChannelFlagSuppressJoinNotificationReplies
	// Suppress role subscription purchase and renewal notifications.
	SystemChannelFlagSuppressRoleSubscriptionPurchaseNotifications
	// Hide role subscription sticker reply buttons
	SystemChannelFlagSuppressRoleSubscriptionPurchaseNotificationReplies
)

// Has returns true if all provided flags are set.
func (f SystemChannelFlags) Has(flags ...SystemChannelFlags) bool {
	return BitFieldHas(f, flags...)
}

// PremiumTier represents the boost level of a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-object-premium-tier
type PremiumTier int

const (
	// Guild has not unlocked any Server Boost perks.
	PremiumTierNone PremiumTier = iota
	// Guild has unlocked Server Boost level 1 perks.
	PremiumTierOne
	// Guild has unlocked Server Boost level 2 perks.
	PremiumTierTwo
	// Guild has unlocked Server Boost level 3 perks.
	PremiumTierThree
)

// Is returns true if the guild's premium tier matches the provided premium tier.
func (p PremiumTier) Is(premiumTier PremiumTier) bool {
	return p == premiumTier
}

// GuildWelcomeChannel is one of the channels in a GuildWelcomeScreen
//
// Reference: https://discord.com/developers/docs/resources/guild#welcome-screen-object-welcome-screen-channel-structure
type GuildWelcomeChannel struct {
	// ChannelID is the channel's id.
	ChannelID Snowflake `json:"channel_id"`

	// Description is the description shown for the channel.
	Description string `json:"description"`

	// EmojiID is the emoji id, if the emoji is custom
	//
	// Optional:
	//  - May be equal to 0 if no emoji is set.
	//  - May be equal to 0 if the emoji is set but its a unicode emoji.
	EmojiID Snowflake `json:"emoji_id,omitempty"`

	// EmojiID is the emoji name if custom, the unicode character if standard, or empty string if no emoji is set
	//
	// Optional:
	//  - May be empty string if no emoji is set.
	EmojiName string `json:"emoji_name,omitempty"`
}

// Mention returns a Discord mention string for the channel.
//
// Example output: "<#123456789012345678>"
func (c *GuildWelcomeChannel) Mention() string {
	return "<#" + c.ChannelID.String() + ">"
}

// GuildWelcomeScreen is the Welcome Screen of a Guild
//
// Reference: https://discord.com/developers/docs/resources/guild#welcome-screen-object
type GuildWelcomeScreen struct {
	// Description is the server description shown in the welcome screen.
	Description string `json:"description,omitempty"`

	// WelcomeChannels is the channels shown in the welcome screen,
	//
	// Note:
	//  - Can be up to 5 channels.
	WelcomeChannels []GuildWelcomeChannel `json:"welcome_channels"`
}

// NSFWLevel represent the NSFW level on a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-object-guild-nsfw-level
type NSFWLevel int

const (
	NSFWLevelDefault NSFWLevel = iota
	NSFWLevelExplicit
	NSFWLevelSafe
	NSFWLevelAgeRestricted
)

// Is returns true if the guild's NSFW level matches the provided NSFW level.
func (l NSFWLevel) Is(level NSFWLevel) bool {
	return l == level
}

// GuildWelcomeScreen represent incidents data of a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#incidents-data-object
type GuildIncidentsData struct {
	// InvitesDisabledUntil is when invites get enabled again,
	InvitesDisabledUntil *time.Time `json:"invites_disabled_until"`
	// DMsDisabledUntil is when direct messages get enabled again.
	DMsDisabledUntil *time.Time `json:"dms_disabled_until"`
	// DMSpamDetectedAt is when the dm spam was detected.
	DMSpamDetectedAt *time.Time `json:"dm_spam_detected_at"`
	// RaidDetectedAt is when the raid was detected.
	RaidDetectedAt *time.Time `json:"raid_detected_at"`
}

// Guild represent a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild
type Guild struct {
	// ID is the guild's unique Discord snowflake ID.
	ID Snowflake `json:"id"`

	// Unavailable is whether this guild is available or not.
	Unavailable bool `json:"unavailable"`

	// Name is the guild's name.
	Name string `json:"name"`

	// Description is the description of a guild.
	//
	// Optional:
	//  - May be empty string if no description is set.
	Description string `json:"description"`

	// Icon is the guild's icon hash.
	//
	// Optional:
	//  - May be empty string if no icon.
	Icon string `json:"icon"`

	// Splash is the guild's splash hash.
	//
	// Optional:
	//  - May be empty string if no splash.
	Splash string `json:"splash"`

	// DiscoverySplash is the guild's discovery splash hash.
	//
	// Optional:
	//  - May be empty string if no discovery splash.
	DiscoverySplash string `json:"discovery_splash"`

	// OwnerID is the guild's owner id.
	OwnerID Snowflake `json:"owner_id"`

	// AfkChannelID is the guild's afk channel id.
	//
	// Optional:
	//  - May be equal to 0 if no Afk channel is set.
	AfkChannelID Snowflake `json:"afk_channel_id"`

	// AfkTimeout is the afk timeout in seconds.
	AfkTimeout int `json:"afk_timeout"`

	// WidgetEnabled is whether the server widget is enabled.
	WidgetEnabled bool `json:"widget_enabled"`

	// WidgetChannelID is the channel id that the widget will generate an invite to, or 0 if set to no invite.
	//
	// Optional:
	//  - May be equal to 0 if no widget channel is set.
	WidgetChannelID Snowflake `json:"widget_channel_id"`

	// VerificationLevel is the verification level required for the guild.
	VerificationLevel VerificationLevel `json:"verification_level"`

	// DefaultMessageNotifications is the default message notifications level.
	DefaultMessageNotifications MessageNotificationsLevel `json:"default_message_notifications"`

	// ExplicitContentFilter is the explicit content filter level.
	ExplicitContentFilter ExplicitContentFilterLevel `json:"explicit_content_filter"`

	// Features is the enabled guild features.
	Features []GuildFeature `json:"features"`

	// MFALevel is the required MFA level for the guild
	MFALevel MFALevel `json:"mfa_level"`

	// SystemChannelID is the guild's system channel id.
	//
	// Optional:
	//  - May be equal to 0 if no system channel is set.
	SystemChannelID Snowflake `json:"system_channel_id"`

	// SystemChannelFlags is the system channel flags on this guild.
	SystemChannelFlags SystemChannelFlags `json:"system_channel_flags"`

	// RulesChannelID is the guild's rules channel id.
	//
	// Optional:
	//  - May be equal to 0 if no rules channel is set.
	RulesChannelID Snowflake `json:"rules_channel_id"`

	// MaxPresences is the maximum number of presences for the guild.
	//
	// Optional:
	//  - ALways nil, apart from the largest of guilds.
	MaxPresences *int `json:"max_presences"`

	// MaxMembers is the maximum number of members for the guild.
	MaxMembers int `json:"max_members"`

	// VanityURLCode is the vanity url code for the guild
	//
	// Optional:
	//  - May be empty string if no vanity url code is set.
	VanityURLCode string `json:"vanity_url_code"`

	// Banner is the guild's banner hash.
	//
	// Optional:
	//  - May be empty string if no banner is set.
	Banner string `json:"banner"`

	// PremiumTier is premium tier of this guild (Server Boost level).
	PremiumTier PremiumTier `json:"premium_tier"`

	// PremiumSubscriptionCount is the number of boosts this guild currently has.
	PremiumSubscriptionCount int `json:"premium_subscription_count"`

	// PreferredLocale is the preferred locale of a Community guild;
	// used in server discovery and notices from Discord, and sent in interactions; defaults to "en-US"
	PreferredLocale Locale `json:"preferred_locale"`

	// PublicUpdatesChannelID is the id of the channel where admins and moderators
	// of Community guilds receive notices from Discord
	//
	// Optional:
	//  - May be equal to 0 if no public updates channel is set.
	PublicUpdatesChannelID Snowflake `json:"public_updates_channel_id"`

	// MaxVideoChannelUsers is the maximum amount of users in a video channel.
	MaxVideoChannelUsers int `json:"max_video_channel_users"`

	// MaxStageVideoChannelUsers is the maximum amount of users in a stage video channel.
	MaxStageVideoChannelUsers int `json:"max_stage_video_channel_users"`

	// WelcomeScreen is the welcome screen of a Community guild, shown to new members.
	WelcomeScreen GuildWelcomeScreen `json:"welcome_screen"`

	// NSFWLevel is the guild NSFW level.
	NSFWLevel NSFWLevel `json:"nsfw_level"`

	// PremiumProgressBarEnabled is whether the guild has the boost progress bar enabled.
	PremiumProgressBarEnabled bool `json:"premium_progress_bar_enabled"`

	// SafetyAlertsChannelID is the id of the channel where admins and moderators
	// of Community guilds receive safety alerts from Discord.
	//
	// Optional:
	//  - May be equal to 0 if no safety alerts channel is set.
	SafetyAlertsChannelID Snowflake `json:"safety_alerts_channel_id"`

	// IncidentsData is the incidents data for this guild.
	//
	// Optional:
	//  - May be nil if guild has no incidents data.
	IncidentsData *GuildIncidentsData `json:"incidents_data"`
}

// CreatedAt returns the time when this guild is created.
func (g *Guild) CreatedAt() time.Time {
	return g.ID.Timestamp()
}

// IconURL returns the URL to the guild's icon image.
//
// If the guild has a custom icon set, it returns the URL to that icon, otherwise empty string.
// By default, it uses GIF format if the icon is animated, otherwise PNG.
//
// Example usage:
//
//	url := guild.IconURL()
func (g *Guild) IconURL() string {
	if g.Icon != "" {
		return GuildIconURL(g.ID, g.Icon, ImageFormatDefault, ImageSizeDefault)
	}
	return ""
}

// IconURLWith returns the URL to the guild's icon image,
// allowing explicit specification of image format and size.
//
// If the guild has a custom icon set, it returns the URL to that icon (otherwise empty string)
// using the provided format and size.
//
// Example usage:
//
//	url := guild.IconURLWith(ImageFormatWebP, ImageSize512)
func (g *Guild) IconURLWith(format ImageFormat, size ImageSize) string {
	if g.Icon != "" {
		return GuildIconURL(g.ID, g.Icon, format, size)
	}
	return ""
}

// BannerURL returns the URL to the guild's banner image.
//
// If the guild has a custom banner set, it returns the URL to that banner, otherwise empty string.
// By default, it uses GIF format if the banner is animated, otherwise PNG.
//
// Example usage:
//
//	url := guild.BannerURL()
func (g *Guild) BannerURL() string {
	if g.Icon != "" {
		return GuildBannerURL(g.ID, g.Icon, ImageFormatDefault, ImageSizeDefault)
	}
	return ""
}

// BannerURLWith returns the URL to the guild's banner image,
// allowing explicit specification of image format and size.
//
// If the guild has a custom banner set, it returns the URL to that banner (otherwise empty string)
// using the provided format and size.
//
// Example usage:
//
//	url := guild.BannerURLWith(ImageFormatWebP, ImageSize512)
func (g *Guild) BannerURLWith(format ImageFormat, size ImageSize) string {
	if g.Icon != "" {
		return GuildBannerURL(g.ID, g.Icon, format, size)
	}
	return ""
}

// SplashURL returns the URL to the guild's splash image.
//
// If the guild has a splash image set, it returns the URL to that image,
// Otherwise empty string, By default it uses PNG.
//
// Example usage:
//
//	url := guild.SplashURL()
func (g *Guild) SplashURL() string {
	if g.Splash != "" {
		return GuildSplashURL(g.ID, g.Splash, ImageFormatDefault, ImageSizeDefault)
	}
	return ""
}

// SplashURLWith returns the URL to the guild's splash image,
// allowing explicit specification of image format and size.
//
// If the guild has a splash image set, it returns the URL to that image (otherwise empty string).
// using the provided format and size.
//
// Example usage:
//
//	url := guild.SplashURLWith(ImageFormatWebP, ImageSize512)
func (g *Guild) SplashURLWith(format ImageFormat, size ImageSize) string {
	if g.Splash != "" {
		return GuildSplashURL(g.ID, g.Icon, format, size)
	}
	return ""
}

// DiscoverySplashURL returns the URL to the guild's discovery splash image.
//
// If the guild has a discovery splash image set, it returns the URL to that image,
// Otherwise empty string, By default it uses PNG.
//
// Example usage:
//
//	url := guild.DiscoverySplashURL()
func (g *Guild) DiscoverySplashURL() string {
	if g.DiscoverySplash != "" {
		return GuildDiscoverySplashURL(g.ID, g.Splash, ImageFormatDefault, ImageSizeDefault)
	}
	return ""
}

// DiscoverySplashURLWith returns the URL to the guild's discovery splash image,
// allowing explicit specification of image format and size.
//
// If the guild has a discovery splash image set, it returns the URL to that image (otherwise empty string).
// using the provided format and size.
//
// Example usage:
//
//	url := guild.DiscoverySplashURLWith(ImageFormatWebP, ImageSize512)
func (g *Guild) DiscoverySplashURLWith(format ImageFormat, size ImageSize) string {
	if g.DiscoverySplash != "" {
		return GuildDiscoverySplashURL(g.ID, g.Icon, format, size)
	}
	return ""
}

// RestGuild represents a guild object returned by the Discord API.
// It embeds Guild and adds additional fields provided by the REST endpoint.
//
// Reference: https://discord.com/developers/docs/resources/guild
type RestGuild struct {
	Guild

	// Stickers contains the custom stickers available in the guild.
	Stickers []Sticker `json:"stickers"`

	// Roles contains all roles defined in the guild.
	Roles []Role `json:"roles"`

	// Emojis contains the custom emojis available in the guild.
	Emojis []Emoji `json:"emojis"`
}

// RestGuild represents a guild object returned by the Discord gateway.
// It embeds RestGuild and adds additional fields provided in the gateway.
//
// Reference: https://discord.com/developers/docs/events/gateway-events#guild-create
type GatewayGuild struct {
	RestGuild

	// Large if true this is considered a large guild.
	Large bool `json:"large"`

	// MemberCount is the total number of members in this guild.
	MemberCount int `json:"member_count"`

	// VoiceStates is the states of members currently in voice channels; lacks the GuildID key.
	VoiceStates []VoiceState `json:"voice_states"`

	// Members is a slice of the Users in the guild.
	Members []Member `json:"members"`

	// Channels is a slice of the Channels in the guild.
	Channels []GuildChannel `json:"channels"`

	// Threads are all active threads in the guild that current user has permission to view.
	Threads []ThreadChannel `json:"threads"`

	// StageInstances is a slice of the Stage instances in the guild.
	StageInstances []StageInstance `json:"stage_instances"`

	// SoundboardSounds is a slice of the Soundboard sounds in the guild.
	SoundboardSounds []SoundBoardSound `json:"soundboard_sounds"`
}

var _ json.Unmarshaler = (*GatewayGuild)(nil)

// UnmarshalJSON implements json.Unmarshaler for GatewayGuild.
func (g *GatewayGuild) UnmarshalJSON(buf []byte) error {
	type tempGuild struct {
		RestGuild
		Large            bool              `json:"large"`
		MemberCount      int               `json:"member_count"`
		VoiceStates      []VoiceState      `json:"voice_states"`
		Members          []Member          `json:"members"`
		Channels         []json.RawMessage `json:"channels"`
		Threads          []ThreadChannel   `json:"threads"`
		StageInstances   []StageInstance   `json:"stage_instances"`
		SoundboardSounds []SoundBoardSound `json:"soundboard_sounds"`
	}

	var temp tempGuild
	if err := sonic.Unmarshal(buf, &temp); err != nil {
		return err
	}

	g.RestGuild = temp.RestGuild
	g.Large = temp.Large
	g.MemberCount = temp.MemberCount
	g.VoiceStates = temp.VoiceStates
	g.Members = temp.Members
	g.Threads = temp.Threads
	g.StageInstances = temp.StageInstances
	g.SoundboardSounds = temp.SoundboardSounds

	for i := range len(g.Roles) {
		g.Roles[i].GuildID = g.ID
	}
	for i := range len(g.Members) {
		g.Members[i].GuildID = g.ID
	}

	if temp.Channels != nil {
		g.Channels = make([]GuildChannel, 0, len(temp.Channels))
		for i := range len(temp.Channels) {
			if len(temp.Channels[i]) == 0 || bytes.Equal(temp.Channels[i], []byte("null")) {
				continue
			}
			channel, err := UnmarshalChannel(temp.Channels[i])
			if err != nil {
				return err
			}
			if guildCh, ok := channel.(GuildChannel); ok {
				g.Channels = append(g.Channels, guildCh)
			} else {
				return errors.New("cannot unmarshal non-GuildChannel into GuildChannel")
			}
		}
	}

	return nil
}

// PartialGuild represents a partial struct of a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild
type PartialGuild struct {
	// ID is the guild's unique Discord snowflake ID.
	ID Snowflake `json:"id"`

	// Locale is the preferred locale of the guild;
	Locale Locale `json:"locale"`

	// Features is the enabled guild features.
	Features []GuildFeature `json:"features"`
}
