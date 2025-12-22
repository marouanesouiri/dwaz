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
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/marouanesouiri/stdx/optional"
	"github.com/marouanesouiri/stdx/result"
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

// String implements the fmt.Stringer interface.
func (c *GuildWelcomeChannel) String() string {
	return c.Mention()
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

// GuildIncidentsData represent incidents data of a Discord guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#incidents-data-object
type GuildIncidentsData struct {
	// InvitesDisabledUntil is when invites get enabled again,
	InvitesDisabledUntil optional.Option[time.Time] `json:"invites_disabled_until,omitzero"`

	// DMsDisabledUntil is when direct messages get enabled again.
	DMsDisabledUntil optional.Option[time.Time] `json:"dms_disabled_until,omitzero"`

	// DMSpamDetectedAt is when the dm spam was detected.
	DMSpamDetectedAt optional.Option[time.Time] `json:"dm_spam_detected_at,omitzero"`

	// RaidDetectedAt is when the raid was detected.
	RaidDetectedAt optional.Option[time.Time] `json:"raid_detected_at,omitzero"`
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
	//  - Always not present, apart from the largest of guilds.
	MaxPresences optional.Option[int] `json:"max_presences"`

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
		return GuildDiscoverySplashURL(g.ID, g.DiscoverySplash, format, size)
	}
	return ""
}

// Ban represents a guild ban.
//
// Reference: https://discord.com/developers/docs/resources/guild#ban-object
type Ban struct {
	// Reason is the reason for the ban.
	Reason string `json:"reason,omitempty"`

	// User is the banned user.
	User User `json:"user"`
}

// GuildPreview represents a preview of a guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-preview-object
type GuildPreview struct {
	// ID is the guild id.
	ID Snowflake `json:"id"`

	// Name is the guild name (2-100 characters).
	Name string `json:"name"`

	// Icon is the icon hash.
	Icon string `json:"icon,omitempty"`

	// Splash is the splash hash.
	Splash string `json:"splash,omitempty"`

	// DiscoverySplash is the discovery splash hash.
	DiscoverySplash string `json:"discovery_splash,omitempty"`

	// Emojis are the custom guild emojis.
	Emojis []Emoji `json:"emojis"`

	// Features are the enabled guild features.
	Features []GuildFeature `json:"features"`

	// ApproximateMemberCount is the approximate number of members in this guild.
	ApproximateMemberCount int `json:"approximate_member_count"`

	// ApproximatePresenceCount is the approximate number of online members in this guild.
	ApproximatePresenceCount int `json:"approximate_presence_count"`

	// Description is the description for the guild.
	Description string `json:"description,omitempty"`

	// Stickers are the custom guild stickers.
	Stickers []Sticker `json:"stickers"`
}

// GuildWidgetSettings represents a guild widget settings.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-widget-settings-object
type GuildWidgetSettings struct {
	// Enabled is whether the widget is enabled.
	Enabled bool `json:"enabled"`

	// ChannelID is the widget channel id.
	ChannelID optional.Option[Snowflake] `json:"channel_id,omitzero"`
}

// GuildWidget represents a guild widget.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-widget-object
type GuildWidget struct {
	// ID is the guild id.
	ID Snowflake `json:"id"`

	// Name is the guild name.
	Name string `json:"name"`

	// InstantInvite is the instant invite for the guilds specified widget invite channel.
	InstantInvite string `json:"instant_invite,omitempty"`

	// Channels are the voice and stage channels which are accessible by @everyone.
	Channels []Channel `json:"channels"`

	// Members are the special widget user objects.
	Members []User `json:"members"`

	// PresenceCount is the number of online members in this guild.
	PresenceCount int `json:"presence_count"`
}

// OnboardingMode defines the criteria used to satisfy Onboarding constraints that are required for enabling.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-onboarding-object-onboarding-mode
type OnboardingMode int

const (
	// OnboardingModeDefault counts only Default Channels towards constraints.
	OnboardingModeDefault OnboardingMode = 0
	// OnboardingModeAdvanced counts Default Channels and Questions towards constraints.
	OnboardingModeAdvanced OnboardingMode = 1
)

// Is checks if the onboarding mode matches the provided mode.
func (m OnboardingMode) Is(mode OnboardingMode) bool {
	return m == mode
}

// GuildOnboarding represents guild onboarding configuration.

// Reference: https://discord.com/developers/docs/resources/guild#guild-onboarding-object
type GuildOnboarding struct {
	// GuildID is the ID of the guild this onboarding is part of.
	GuildID Snowflake `json:"guild_id"`

	// Prompts are the prompts shown during onboarding.
	Prompts []OnboardingPrompt `json:"prompts"`

	// DefaultChannelIDs are the channel IDs that members get opted into automatically.
	DefaultChannelIDs []Snowflake `json:"default_channel_ids"`

	// Enabled is whether onboarding is enabled in the guild.
	Enabled bool `json:"enabled"`

	// Mode is the current mode of onboarding.
	Mode OnboardingMode `json:"mode"`
}

// PromptType represents the type of onboarding prompt.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-onboarding-object-prompt-type
type PromptType int

const (
	// PromptTypeMultipleChoice represents a multiple choice prompt.
	PromptTypeMultipleChoice PromptType = 0
	// PromptTypeDropdown represents a dropdown prompt.
	PromptTypeDropdown PromptType = 1
)

// Is checks if the prompt type matches the provided type.
func (t PromptType) Is(typ PromptType) bool {
	return t == typ
}

// OnboardingPrompt represents an onboarding prompt.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-onboarding-object-onboarding-prompt-structure
type OnboardingPrompt struct {
	// ID is the ID of the prompt.
	ID Snowflake `json:"id,omitempty"`

	// Type is the type of prompt.
	Type PromptType `json:"type"`

	// Options are the options available within the prompt.
	Options []OnboardingPromptOption `json:"options"`

	// Title is the title of the prompt.
	Title string `json:"title"`

	// SingleSelect indicates whether users are limited to selecting one option.
	SingleSelect bool `json:"single_select"`

	// Required indicates whether the prompt is required before completing onboarding.
	Required bool `json:"required"`

	// InOnboarding indicates whether the prompt is present in the onboarding flow.
	InOnboarding bool `json:"in_onboarding"`
}

// OnboardingPromptOption represents an option within an onboarding prompt.
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-onboarding-object-prompt-option-structure
type OnboardingPromptOption struct {
	// ID is the ID of the prompt option.
	ID Snowflake `json:"id,omitempty"`

	// ChannelIDs are the IDs for channels a member is added to.
	ChannelIDs []Snowflake `json:"channel_ids"`

	// RoleIDs are the IDs for roles assigned to a member.
	RoleIDs []Snowflake `json:"role_ids"`

	// EmojiID is the emoji ID of the option.
	EmojiID Snowflake `json:"emoji_id,omitempty"`

	// EmojiName is the emoji name of the option.
	EmojiName string `json:"emoji_name,omitempty"`

	// EmojiAnimated is whether the emoji is animated.
	EmojiAnimated bool `json:"emoji_animated"`

	// Title is the title of the option.
	Title string `json:"title"`

	// Description is the description of the option.
	Description string `json:"description,omitempty"`
}

// IntegrationExpireBehavior represents the behavior of expiring subscribers.
type IntegrationExpireBehavior int

const (
	IntegrationExpireBehaviorRemoveRole IntegrationExpireBehavior = 0
	IntegrationExpireBehaviorKick       IntegrationExpireBehavior = 1
)

// IntegrationAccount represents an integration account.
//
// Reference: https://discord.com/developers/docs/resources/guild#integration-account-object
type IntegrationAccount struct {
	// ID is the id of the account.
	ID string `json:"id"`

	// Name is the name of the account.
	Name string `json:"name"`
}

// IntegrationApplication represents an integration application.
//
// Reference: https://discord.com/developers/docs/resources/guild#integration-application-object
type IntegrationApplication struct {
	// ID is the id of the app.
	ID Snowflake `json:"id"`

	// Name is the name of the app.
	Name string `json:"name"`

	// Icon is the icon hash of the app.
	Icon string `json:"icon,omitempty"`

	// Description is the description of the app.
	Description string `json:"description"`

	// Bot is the bot associated with this application.
	Bot *User `json:"bot,omitempty"`
}

// Integration represents a guild integration.
//
// Reference: https://discord.com/developers/docs/resources/guild#integration-object
type Integration struct {
	// ID is the integration id.
	ID Snowflake `json:"id"`

	// Name is the integration name.
	Name string `json:"name"`

	// Type is the integration type (twitch, youtube, discord, or guild_subscription).
	Type string `json:"type"`

	// Enabled is whether this integration is enabled.
	Enabled bool `json:"enabled"`

	// Syncing is whether this integration is syncing.
	Syncing bool `json:"syncing,omitempty"`

	// RoleID is the id that this integration uses for "subscribers".
	RoleID Snowflake `json:"role_id,omitempty"`

	// EnableEmoticons is whether emoticons should be synced for this integration.
	EnableEmoticons bool `json:"enable_emoticons,omitempty"`

	// ExpireBehavior is the behavior of expiring subscribers.
	ExpireBehavior *IntegrationExpireBehavior `json:"expire_behavior,omitempty"`

	// ExpireGracePeriod is the grace period (in days) before expiring subscribers.
	ExpireGracePeriod int `json:"expire_grace_period,omitempty"`

	// User is the user for this integration.
	User *User `json:"user,omitempty"`

	// Account is the integration account information.
	Account IntegrationAccount `json:"account"`

	// SyncedAt is when this integration was last synced.
	SyncedAt optional.Option[time.Time] `json:"synced_at,omitzero"`

	// SubscriberCount is how many subscribers this integration has.
	SubscriberCount int `json:"subscriber_count,omitempty"`

	// Revoked is whether this integration has been revoked.
	Revoked bool `json:"revoked,omitempty"`

	// Application is the bot/OAuth2 application for discord integrations.
	Application *IntegrationApplication `json:"application,omitempty"`
	// Scopes are the scopes the application has been authorized for.
	Scopes []string `json:"scopes,omitempty"`
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
	Members []FullMember `json:"members"`

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
		Members          []FullMember      `json:"members"`
		Channels         []json.RawMessage `json:"channels"`
		Threads          []ThreadChannel   `json:"threads"`
		StageInstances   []StageInstance   `json:"stage_instances"`
		SoundboardSounds []SoundBoardSound `json:"soundboard_sounds"`
	}

	var temp tempGuild
	if err := json.Unmarshal(buf, &temp); err != nil {
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
	for i := range len(g.VoiceStates) {
		g.VoiceStates[i].GuildID = g.ID
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

	// Name is the guild's name.
	Name string `json:"name"`

	// Icon is the guild's icon hash.
	//
	// Optional:
	//  - May be empty string if no icon.
	Icon string `json:"icon"`

	// Banner is the guild's banner hash.
	//
	// Optional:
	//  - May be empty string if no banner is set.
	Banner string `json:"banner"`

	// Locale is the preferred locale of the guild;
	Locale Locale `json:"locale"`

	// Features is the enabled guild features.
	Features []GuildFeature `json:"features"`
}

// IconURL returns the URL to the guild's icon image.
//
// If the guild has a custom icon set, it returns the URL to that icon, otherwise empty string.
// By default, it uses GIF format if the icon is animated, otherwise PNG.
//
// Example usage:
//
//	url := guild.IconURL()
func (g *PartialGuild) IconURL() string {
	if g.Icon != "" {
		return GuildIconURL(g.ID, g.Icon, ImageFormatDefault, ImageSizeDefault)
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
func (g *PartialGuild) BannerURL() string {
	if g.Icon != "" {
		return GuildBannerURL(g.ID, g.Icon, ImageFormatDefault, ImageSizeDefault)
	}
	return ""
}

// FetchGuildOptions contains parameters for fetching a guild.
type FetchGuildOptions struct {
	// When 'true', will return approximate member and presence counts for the guild
	WithCounts bool `json:"with_counts,omitempty"`
}

// FetchGuild retrieves a guild by its ID.
//
// Reference: https://discord.com/developers/docs/resources/guild#get-guild
func (r *requester) FetchGuild(guildID Snowflake, opts FetchGuildOptions) result.Result[RestGuild] {
	endpoint := "/guilds/" + guildID.String() + "?with_counts=" + strconv.FormatBool(opts.WithCounts)

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[RestGuild](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var guild RestGuild
	if err := json.NewDecoder(body).Decode(&guild); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[RestGuild](err)
	}
	return result.Ok(guild)
}

// FetchGuildPreview retrieves a guild preview by its ID.
//
// Reference: https://discord.com/developers/docs/resources/guild#get-guild-preview
func (r *requester) FetchGuildPreview(guildID Snowflake) result.Result[GuildPreview] {
	endpoint := "/guilds/" + guildID.String() + "/preview"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[GuildPreview](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var preview GuildPreview
	if err := json.NewDecoder(body).Decode(&preview); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/preview",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildPreview](err)
	}
	return result.Ok(preview)
}

type AfkTimeout int

const (
	AfkTimeout60   AfkTimeout = 60
	AfkTimeout300  AfkTimeout = 300
	AfkTimeout900  AfkTimeout = 900
	AfkTimeout1800 AfkTimeout = 1800
	AfkTimeout3600 AfkTimeout = 3600
)

// ModifyGuildOptions contains parameters for modifying a guild's settings.
//
// Reference: https://discord.com/developers/docs/resources/guild#modify-guild
type ModifyGuildOptions struct {
	Name                        string                                      `json:"name,omitempty"`
	VerificationLevel           optional.Option[VerificationLevel]          `json:"verification_level,omitzero"`
	DefaultMessageNotifications optional.Option[MessageNotificationsLevel]  `json:"default_message_notifications,omitzero"`
	ExplicitContentFilter       optional.Option[ExplicitContentFilterLevel] `json:"explicit_content_filter,omitzero"`
	AfkChannelID                optional.Option[Snowflake]                  `json:"afk_channel_id,omitzero"`
	AfkTimeout                  AfkTimeout                                  `json:"afk_timeout,omitempty"`
	Icon                        optional.Option[Base64Image]                `json:"icon,omitzero"`
	Splash                      optional.Option[Base64Image]                `json:"splash,omitzero"`
	DiscoverySplash             optional.Option[Base64Image]                `json:"discovery_splash,omitzero"`
	Banner                      optional.Option[Base64Image]                `json:"banner,omitzero"`
	SystemChannelID             optional.Option[Snowflake]                  `json:"system_channel_id,omitzero"`
	SystemChannelFlags          optional.Option[SystemChannelFlags]         `json:"system_channel_flags,omitzero"`
	RulesChannelID              optional.Option[Snowflake]                  `json:"rules_channel_id,omitzero"`
	PublicUpdatesChannelID      optional.Option[Snowflake]                  `json:"public_updates_channel_id,omitzero"`
	PreferredLocale             Locale                                      `json:"preferred_locale,omitempty"`
	Features                    []GuildFeature                              `json:"features"`
	Description                 string                                      `json:"description"`
	PremiumProgressBarEnabled   optional.Option[bool]                       `json:"premium_progress_bar_enabled,omitzero"`
	SafetyAlertsChannelID       optional.Option[Snowflake]                  `json:"safety_alerts_channel_id,omitzero"`

	Reason string `json:"-"`
}

// ModifyGuild modifies a guild's settings and returns the updated guild object.
//
// Requires the PermissionManageGuild permission.
//
// Note:
//   - Attempting to add or remove the GuildFeatureCommunity guild feature requires the PermissionAdministrator permission
//
// Reference: https://discord.com/developers/docs/resources/guild#modify-guild
func (r *requester) ModifyGuild(guildID Snowflake, opts ModifyGuildOptions) result.Result[Guild] {
	reqBody, _ := json.Marshal(opts)
	body := r.DoRequest(Request{
		Method: "PATCH",
		URL:    "/guilds/" + guildID.String(),
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if body.IsErr() {
		return result.Err[Guild](body.Err())
	}
	defer body.Value().Close()

	var guild Guild
	if err := json.NewDecoder(body.Value()).Decode(&guild); err != nil {
		return result.Err[Guild](err)
	}
	return result.Ok(guild)
}

// FetchGuildChannels returns a list of guild channel objects.
//
// Note:
//   - Does not include threads.
//
// Reference: https://discord.com/developers/docs/resources/guild#get-guild-channels
func (r *requester) FetchGuildChannels(guildID Snowflake) result.Result[[]GuildChannel] {
	endpoint := "/guilds/" + guildID.String() + "/channels"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[[]GuildChannel](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var rawChannels []json.RawMessage
	if err := json.NewDecoder(body).Decode(&rawChannels); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/channels",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]GuildChannel](err)
	}

	channels := make([]GuildChannel, 0, len(rawChannels))
	for _, raw := range rawChannels {
		ch, err := UnmarshalChannel(raw)
		if err != nil {
			return result.Err[[]GuildChannel](err)
		}
		if guildCh, ok := ch.(GuildChannel); ok {
			channels = append(channels, guildCh)
		}
	}

	return result.Ok(channels)
}

// CreateChannelOptions defines the configuration for creating a new Discord guild channel.
//
// Note:
//   - This struct configures properties for a new channel, such as text, voice, or forum.
//   - Only set fields applicable to the channel type to avoid errors.
//
// Reference: https://discord.com/developers/docs/resources/guild#create-guild-channel-json-params
type CreateChannelOptions struct {
	// Name is the channel's name (1-100 characters).
	//
	// Applies to All Channels.
	Name string `json:"name"`

	// Type specifies the type of channel to create.
	//
	// Note:
	//  - Defaults to ChannelTypeGuildText if unset.
	//
	// Applies to All Channels.
	Type ChannelType `json:"type"`

	// Topic is a description of the channel (0-1024 characters).
	//
	// Note:
	//  - This field is optional.
	//
	// Applies to Channels of Type: Text, Announcement, Forum, Media.
	Topic string `json:"topic,omitempty"`

	// Bitrate sets the audio quality for voice or stage channels (in bits, minimum 8000).
	//
	// Note:
	//  - This field is ignored for non-voice channels.
	//
	// Applies to Channels of Type: Voice, Stage.
	Bitrate Bitrate `json:"bitrate,omitempty"`

	// UserLimit caps the number of users in a voice or stage channel (0 for unlimited, 1-99 for a limit).
	//
	// Note:
	//  - Set to 0 to allow unlimited users.
	//
	// Applies to Channels of Type: Voice, Stage.
	UserLimit int `json:"user_limit,omitempty"`

	// RateLimitPerUser sets the seconds a user must wait before sending another message (0-21600).
	//
	// Note:
	//  - Bots and users with manage_messages or manage_channel permissions are unaffected.
	//
	// Applies to Channels of Type: Text, Voice, Stage, Forum, Media.
	RateLimitPerUser int `json:"rate_limit_per_user,omitempty"`

	// Position determines the channel’s position in the server’s channel list (lower numbers appear higher).
	//
	// Note:
	//  - Channels with the same position are sorted by their internal ID.
	//
	// Applies to All Channels.
	Position optional.Option[int] `json:"position,omitzero"`

	// PermissionOverwrites defines custom permissions for specific roles or users.
	//
	// Applies to All Channels.
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitempty"`

	// ParentID is the ID of the category to nest the channel under.
	//
	// Applies to Channels of Type: Text, Voice, Announcement, Stage, Forum, Media.
	ParentID Snowflake `json:"parent_id,omitempty"`

	// Nsfw marks the channel as Not Safe For Work, restricting it to 18+ users.
	//
	// Note:
	//  - Set to true to enable the age restriction.
	//
	// Applies to Channels of Type: Text, Voice, Announcement, Stage, Forum.
	Nsfw bool `json:"nsfw,omitempty"`

	// RTCRegion sets the channel voice region id of the voice or stage channel.
	//
	// Applies to Channels of Type: Voice, Stage.
	RTCRegion string `json:"rtc_region,omitempty"`

	// VideoQualityMode sets the camera video quality for voice or stage channels.
	//
	// Applies to Channels of Type: Voice, Stage.
	VideoQualityMode VideoQualityModes `json:"video_quality_mode,omitempty"`

	// DefaultAutoArchiveDuration sets the default time (in minutes) before threads are archived.
	//
	// Applies to Channels of Type: Text, Announcement, Forum, Media.
	DefaultAutoArchiveDuration AutoArchiveDuration `json:"default_auto_archive_duration,omitempty"`

	// DefaultReactionEmoji is the default emoji for the add reaction button on threads.
	//
	// Applies to Channels of Type: Forum, Media.
	DefaultReactionEmoji optional.Option[DefaultReactionEmoji] `json:"default_reaction_emoji,omitzero"`

	// AvailableTags lists tags that can be applied to threads for organization.
	//
	// Note:
	//  - This field defines tags users can select for threads.
	//
	// Applies to Channels of Type: Forum, Media.
	AvailableTags []ForumTag `json:"available_tags,omitempty"`

	// DefaultSortOrder sets how threads are sorted by default.
	//
	// Note:
	//  - Valid options are defined in ForumPostsSortOrder.
	//
	// Applies to Channels of Type: Forum, Media.
	DefaultSortOrder optional.Option[ForumPostsSortOrder] `json:"default_sort_order,omitzero"`

	// DefaultForumLayout sets the default view for forum posts.
	//
	// Applies to Channels of Type: Forum.
	DefaultForumLayout ForumLayout `json:"default_forum_layout,omitempty"`

	// DefaultThreadRateLimitPerUser sets the default slow mode for messages in new threads.
	//
	// Note:
	//  - This value is copied to new threads at creation and does not update live.
	//
	// Applies to Channels of Type: Text, Announcement, Forum, Media.
	DefaultThreadRateLimitPerUser int `json:"default_thread_rate_limit_per_user,omitzero"`

	// Reason specifies the audit log reason for creating the channel.
	Reason string `json:"-"`
}

// CreateChannel creates a new channel for the guild.
//
// Requires the PermissionManageChannels permission.
//
// Reference: https://discord.com/developers/docs/resources/guild#create-guild-channel
func (r *requester) CreateChannel(guildID Snowflake, opts CreateChannelOptions) result.Result[GuildChannel] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/guilds/" + guildID.String() + "/channels",
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
			"url":    "/guilds/{id}/channels",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildChannel](err)
	}
	if gc, ok := channel.(GuildChannel); ok {
		return result.Ok(gc)
	}
	return result.Err[GuildChannel](errors.New("created channel is not a guild channel"))
}

type ChannelPosition struct {
	// Channel id
	ID Snowflake `json:"id"`

	// Position is the sorting position of the channel (channels with the same position are sorted by id).
	Position optional.Option[int] `json:"position,omitzero"`

	// Syncs the permission overwrites with the new parent, if moving to a new category.
	LockPermissions bool `json:"lock_permissions,omitzero"`

	// ParentID is the new parent ID for the channel that is moved.
	ParentID Snowflake `json:"parent_id,omitempty"`
}

// ModifyChannelPositionOptions contains parameters for modifying channel positions.
//
// Reference: https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions
type ModifyChannelPositionOptions struct {
	Channels []ChannelPosition
}

// ModifyChannelPositions modifies the positions of a set of channel objects for the guild.
//
// Requires the PermissionManageChannels permission.
//
// Reference: https://discord.com/developers/docs/resources/guild#modify-guild-channel-positions
func (r *requester) ModifyChannelPositions(guildID Snowflake, opts ModifyChannelPositionOptions) result.Void {
	reqBody, _ := json.Marshal(opts.Channels)
	res := r.DoRequest(Request{
		Method: "PATCH",
		URL:    "/guilds/" + guildID.String() + "/channels",
		Body:   reqBody,
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// ActiveThreadsResponse is the response for listing active threads.
//
// Reference: https://discord.com/developers/docs/resources/guild#list-active-guild-threads
type ActiveThreadsResponse struct {
	// Threads is a list of active threads.
	//
	//  Note:
	//   - Threads are ordered by their id, in descending order.
	Threads []ThreadChannel `json:"threads"`

	// Members is list of a thread member object for each returned thread the current user has joined.
	Members []ThreadMember `json:"members"`
}

// ListActiveGuildThreads returns all active threads in the guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#list-active-guild-threads
func (r *requester) ListActiveGuildThreads(guildID Snowflake) result.Result[ActiveThreadsResponse] {
	endpoint := "/guilds/" + guildID.String() + "/threads/active"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[ActiveThreadsResponse](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var response ActiveThreadsResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/threads/active",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[ActiveThreadsResponse](err)
	}
	return result.Ok(response)
}

// FetchMember retrieves a guild member object for the specified user.
func (r *requester) FetchMember(guildID, userID Snowflake) result.Result[FullMember] {
	endpoint := "/guilds/" + guildID.String() + "/members/" + userID.String()

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[FullMember](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var member FullMember
	if err := json.NewDecoder(body).Decode(&member); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/members/{user_id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[FullMember](err)
	}
	member.GuildID = guildID
	return result.Ok(member)
}

// ListMembersOptions contains parameters for paginating through guild members.
type ListMembersOptions struct {
	// Limit is the maximum number of members to return (1-1000).
	//
	//  Note:
	//   - Defaults to 1 if not specified.
	Limit int `json:"limit,omitempty"`

	// After is the user ID to start after for pagination.
	// Used to get the next page of results.
	After Snowflake `json:"after,omitempty"`
}

// ListMembers retrieves a paginated list of members in a guild.
//
//	Note:
//	 - This endpoint is restricted according to whether the GUILD_MEMBERS Privileged Intent is enabled for your application.
func (r *requester) ListMembers(guildID Snowflake) result.Result[[]FullMember] {
	return r.ListMembersWithOptions(guildID, ListMembersOptions{})
}

// ListMembers retrieves a paginated list of members in a guild.
//
//	Note:
//	 - This endpoint is restricted according to whether the GUILD_MEMBERS Privileged Intent is enabled for your application.
func (r *requester) ListMembersWithOptions(guildID Snowflake, opts ListMembersOptions) result.Result[[]FullMember] {
	endpoint := "/guilds/" + guildID.String() + "/members"

	params := url.Values{}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.After != 0 {
		params.Set("after", opts.After.String())
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	body := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if body.IsErr() {
		return result.Err[[]FullMember](body.Err())
	}
	defer body.Value().Close()

	var members []FullMember
	if err := json.NewDecoder(body.Value()).Decode(&members); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/members",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]FullMember](err)
	}
	for i := range members {
		members[i].GuildID = guildID
	}
	return result.Ok(members)
}

// SearchMembersOptions contains parameters for searching members by name.
type SearchMembersOptions struct {
	// Query is the text to search for in usernames and nicknames.
	//
	//  Note:
	//   - Query is required to be set
	Query string `json:"query"`

	// Limit is the maximum number of members to return (1-1000).
	// Defaults to 1 if not specified.
	Limit int `json:"limit,omitempty"`
}

// SearchMembers returns a list of guild member objects whose username or nickname starts with a provided string.
func (r *requester) SearchMembers(guildID Snowflake, opts SearchMembersOptions) result.Result[[]FullMember] {
	endpoint := "/guilds/" + guildID.String() + "/members/search"

	params := url.Values{}
	params.Set("query", opts.Query)
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	body := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if body.IsErr() {
		return result.Err[[]FullMember](body.Err())
	}
	defer body.Value().Close()

	var members []FullMember
	if err := json.NewDecoder(body.Value()).Decode(&members); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/members/search",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]FullMember](err)
	}
	for i := range members {
		members[i].GuildID = guildID
	}
	return result.Ok(members)
}

// AddMemberOptions contains parameters for adding a user to a guild.
//
// Requires a valid OAuth2 access token with the guilds.join scope.
type AddMemberOptions struct {
	// AccessToken is the OAuth2 access token for the user you want to add.
	// This must have the guilds.join scope.
	AccessToken string `json:"access_token"`

	// Nick sets the member's initial nickname in the guild.
	//
	// Requires the PermissionManageNicknames permission.
	Nick string `json:"nick,omitempty"`

	// Roles is a list of role IDs to assign to the member initially.
	//
	// Requires the PermissionManageRoles permission.
	Roles []Snowflake `json:"roles,omitempty"`

	// Mute sets whether the user is muted in voice channels.
	//
	// Requires the PermissionMuteMember permission.
	Mute bool `json:"mute,omitempty"`

	// Deaf sets whether the user is deafened in voice channels.
	//
	// Requires the PermissionDeafenMembers permission.
	Deaf bool `json:"deaf,omitempty"`
}

// AddMember adds a user to the guild using their OAuth2 access token.
//
// The access token must have the guilds.join scope.
// If the user is already in the guild, this has no effect.
//
//	Note:
//	 - The bot must be a member of the guild with PermissionCreateInstantInvite permission.
//	 - For guilds with Membership Screening enabled, this endpoint will default to adding new members as
//	   pending in the guild member object. Members that are pending will have to complete membership screening
//	   before they become full members that can talk.
func (r *requester) AddMember(guildID, userID Snowflake, opts AddMemberOptions) result.Result[FullMember] {
	reqBody, _ := json.Marshal(opts)
	endpoint := "/guilds/" + guildID.String() + "/members/" + userID.String()

	res := r.DoRequest(Request{
		Method: "PUT",
		URL:    endpoint,
		Body:   reqBody,
	})
	if res.IsErr() {
		return result.Err[FullMember](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var member FullMember
	if err := json.NewDecoder(body).Decode(&member); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PUT",
			"url":    "/guilds/{id}/members/{user_id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[FullMember](err)
	}
	member.GuildID = guildID
	return result.Ok(member)
}

// ModifyMemberOptions contains parameters for modifying a guild member.
//
// All fields are optional. Only provide the fields you want to change.
type ModifyMemberOptions struct {
	// Nickname sets the member's guild nickname.
	// Set to empty string to remove the nickname.
	//
	// Requires the PermissionManageNicknames permission.
	Nickname optional.Option[string] `json:"nick,omitzero"`

	// Roles sets the complete list of role IDs for the member.
	// This replaces all existing roles with the provided list.
	//
	// Requires the PermissionManageRoles permission.
	Roles optional.Option[[]Snowflake] `json:"roles,omitzero"`

	// Mute sets whether the member is muted in voice channels.
	// The member must be in a voice channel for this to work.
	//
	// Requires the PermissionMuteMembers permission.
	Mute optional.Option[bool] `json:"mute,omitzero"`

	// Deaf sets whether the member is deafened in voice channels.
	// The member must be in a voice channel for this to work.
	//
	// Requires the PermissionDeafenMembers permission.
	Deaf optional.Option[bool] `json:"deaf,omitzero"`

	// ChannelID moves the member to a different voice channel.
	// The member must be connected to voice for this to work.
	//
	// Requires the PermissionMoveMembers permission.
	ChannelID Snowflake `json:"channel_id,omitempty"`

	// CommunicationDisabledUntil sets when the member's timeout expires.
	// Can be up to 28 days in the future.
	//
	// Note: Supplying 'optional.Nil[time.Time]()' disables the action.
	//
	// Requires the PermissionModerateMembers permission.
	CommunicationDisabledUntil optional.Option[time.Time] `json:"communication_disabled_until,omitzero"`

	// Flags sets the member's guild-specific flags.
	//
	// Requires the PermissionManageGuild or PermissionManageRoles or (PermissionModerateMembers and PermissionKickMembers and PermissionBanMembers).
	Flags MemberFlags `json:"flags,omitempty"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// ModifyMember updates a member's properties in a guild.
func (r *requester) ModifyMember(guildID, userID Snowflake, opts ModifyMemberOptions) result.Result[FullMember] {
	reqBody, _ := json.Marshal(opts)
	endpoint := "/guilds/" + guildID.String() + "/members/" + userID.String()

	body := r.DoRequest(Request{
		Method: "PATCH",
		URL:    endpoint,
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if body.IsErr() {
		return result.Err[FullMember](body.Err())
	}
	defer body.Value().Close()

	var member FullMember
	if err := json.NewDecoder(body.Value()).Decode(&member); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PATCH",
			"url":    "/guilds/{id}/members/{user_id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[FullMember](err)
	}
	member.GuildID = guildID
	return result.Ok(member)
}

// ModifyCurrentMemberOptions contains parameters for modifying your bot's member profile.
//
// This allows you to update your bot's nickname, avatar, banner, and bio in a specific guild.
type ModifyCurrentMemberOptions struct {
	// Nick sets your bot's nickname in this guild.
	//
	// Requires the PermissionManageNicknames permission.
	Nick optional.Option[string] `json:"nick,omitzero"`

	// Banner sets your bot's guild-specific banner image.
	Banner optional.Option[Base64Image] `json:"banner,omitzero"`

	// Avatar sets your bot's guild-specific avatar image.
	Avatar optional.Option[Base64Image] `json:"avatar,omitzero"`

	// Bio sets your bot's bio text for this guild.
	Bio optional.Option[string] `json:"bio,omitzero"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// ModifyCurrentMember updates your bot's member properties in a guild.
func (r *requester) ModifyCurrentMember(guildID Snowflake, opts ModifyCurrentMemberOptions) result.Result[FullMember] {
	reqBody, _ := json.Marshal(opts)
	endpoint := "/guilds/" + guildID.String() + "/members/@me"

	body := r.DoRequest(Request{
		Method: "PATCH",
		URL:    endpoint,
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if body.IsErr() {
		return result.Err[FullMember](body.Err())
	}
	defer body.Value().Close()

	var member FullMember
	if err := json.NewDecoder(body.Value()).Decode(&member); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PATCH",
			"url":    "/guilds/{id}/members/@me",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[FullMember](err)
	}
	member.GuildID = guildID
	return result.Ok(member)
}

// AddMemberRoleOptions contains parameters for adding a role to a member.
type AddMemberRoleOptions struct {
	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// AddMemberRole assigns a role to a member in a guild.
//
// Requires the PermissionManageRoles permission.
func (r *requester) AddMemberRole(guildID, userID, roleID Snowflake, opts AddMemberRoleOptions) result.Void {
	endpoint := "/guilds/" + guildID.String() + "/members/" + userID.String() + "/roles/" + roleID.String()
	res := r.DoRequest(Request{
		Method: "PUT",
		URL:    endpoint,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// RemoveMemberRoleOptions contains parameters for removing a role from a member.
type RemoveMemberRoleOptions struct {
	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// RemoveMemberRole unassigns a role from a member in a guild.
//
// Requires the PermissionManageRoles permission.
func (r *requester) RemoveMemberRole(guildID, userID, roleID Snowflake, opts RemoveMemberRoleOptions) result.Void {
	endpoint := "/guilds/" + guildID.String() + "/members/" + userID.String() + "/roles/" + roleID.String()
	res := r.DoRequest(Request{Method: "DELETE", URL: endpoint, Reason: opts.Reason})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	return result.OkVoid()
}

// KickMemberOptions contains parameters for kicking a member from a guild.
type KickMemberOptions struct {
	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// KickMember kicks a member from a guild.
//
// Requires the PermissionKickMembers permission.
func (r *requester) KickMember(guildID, userID Snowflake, opts KickMemberOptions) result.Void {
	endpoint := "/guilds/" + guildID.String() + "/members/" + userID.String()
	res := r.DoRequest(Request{Method: "DELETE", URL: endpoint, Reason: opts.Reason})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	return result.OkVoid()
}

// FetchGuildBansOptions contains parameters for fetching guild bans.
//
// Reference: https://discord.com/developers/docs/resources/guild#get-guild-bans
type FetchGuildBansOptions struct {
	// Limit is the number of users to return (up to maximum 1000)
	//
	// Default to 1000 if not spesified
	Limit int `json:"limit,omitempty"`

	// Before consider only users before given user id.
	Before Snowflake `json:"before,omitempty"`

	// After consider only users after given user id.
	After Snowflake `json:"after,omitempty"`
}

// FetchGuildBans returns a list of ban objects for the users banned from this guild.
//
//	Note:
//	 - Provide a user id to before and after for pagination. Users will always be returned in ascending order by 'user.id'.
//	   If both before and after are provided, only before is respected.
//
// Requires the PermissionBanMembers permission.
func (r *requester) FetchGuildBans(guildID Snowflake, opts FetchGuildBansOptions) result.Result[[]Ban] {
	endpoint := "/guilds/" + guildID.String() + "/bans"

	params := url.Values{}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Before != 0 {
		params.Set("before", opts.Before.String())
	}
	if opts.After != 0 {
		params.Set("after", opts.After.String())
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[[]Ban](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var bans []Ban
	if err := json.NewDecoder(body).Decode(&bans); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/bans",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]Ban](err)
	}
	return result.Ok(bans)
}

// FetchGuildBan returns a ban object for the given user.
//
//	Note:
//	 - When the method returns Ok(Some(Ban)), the request succeeded and the
//	   user is banned. When it returns Ok(None), the request also succeeded but there is no
//	   ban for the given user ID. If the method returns Err(err), the request did not succeed.
//
// Requires the PermissionBanMembers permission.
func (r *requester) FetchGuildBan(guildID, userID Snowflake) result.Result[optional.Option[Ban]] {
	endpoint := "/guilds/" + guildID.String() + "/bans/" + userID.String()

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[optional.Option[Ban]](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var ban Ban
	if err := json.NewDecoder(body).Decode(&ban); err != nil {
		if errors.Is(err, io.EOF) {
			return result.Ok(optional.None[Ban]())
		}

		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/bans/{user_id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[optional.Option[Ban]](err)
	}
	return result.Ok(optional.Some(ban))
}

// BanMemberOptions contains parameters for banning a guild member.
type BanMemberOptions struct {
	// DeleteMessageSeconds is the number of seconds to delete messages for, between 0 and 604800 (7 days)
	DeleteMessageSeconds int `json:"delete_message_seconds,omitempty"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// BanMember ban's a guild member.
//
// Requires the PermissionBanMembers permission.
func (r *requester) BanMember(guildID, userID Snowflake, opts BanMemberOptions) result.Void {
	reqBody, _ := json.Marshal(opts)
	endpoint := "/guilds/" + guildID.String() + "/bans/" + userID.String()
	res := r.DoRequest(Request{
		Method: "PUT",
		URL:    endpoint,
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// UnbanMemberOptions contains parameters for unbanning a guild member.
type UnbanMemberOptions struct {
	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// UnbanMember removes the ban for a user.
//
// Requires the PermissionBanMembers permission.
func (r *requester) UnbanMember(guildID, userID Snowflake, opts UnbanMemberOptions) result.Void {
	endpoint := "/guilds/" + guildID.String() + "/bans/" + userID.String()
	res := r.DoRequest(Request{Method: "DELETE", URL: endpoint, Reason: opts.Reason})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	return result.OkVoid()
}

// BulkBanMembersOptions contains parameters for bulk banning guild members.
//
// Reference: https://discord.com/developers/docs/resources/guild#bulk-guild-ban
type BulkBanMembersOptions struct {
	// UserIDs is a list of user ids to ban (max 200).
	UserIDs []Snowflake `json:"user_ids"`

	// DeleteMessageSeconds is the number of seconds to delete messages for, between 0 and 604800 (7 days).
	DeleteMessageSeconds int `json:"delete_message_seconds,omitempty"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// BulkBanResponse is the response for bulk banning members.
//
// Reference: https://discord.com/developers/docs/resources/guild#bulk-guild-ban
type BulkBanResponse struct {
	// BannedUsers is a list of user ids, that were successfully banned.
	BannedUsers []Snowflake `json:"banned_users"`

	// FailedUsers is a list of user ids, that were not banned.
	FailedUsers []Snowflake `json:"failed_users"`
}

// BulkBanMembers bans up to 200 users from a guild.
//
// Reference: https://discord.com/developers/docs/resources/guild#bulk-guild-ban
func (r *requester) BulkBanMembers(guildID Snowflake, opts BulkBanMembersOptions) result.Result[BulkBanResponse] {
	if len(opts.UserIDs) > 200 {
		panic("BulkBanMembers: UserIDs exceeds Discord limit of 200 users")
	}
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/guilds/" + guildID.String() + "/bulk-ban",
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[BulkBanResponse](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var response BulkBanResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "POST",
			"url":    "/guilds/{id}/bulk-ban",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[BulkBanResponse](err)
	}
	return result.Ok(response)
}

// FetchRoles returns a list of role objects for the guild.
func (r *requester) FetchRoles(guildID Snowflake) result.Result[[]Role] {
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    "/guilds/" + guildID.String() + "/roles",
	})
	if res.IsErr() {
		return result.Err[[]Role](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var roles []Role
	if err := json.NewDecoder(body).Decode(&roles); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/roles",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]Role](err)
	}
	for i := range len(roles) {
		roles[i].GuildID = guildID
	}
	return result.Ok(roles)
}

// FetchRole returns a role object for the specified role id.
func (r *requester) FetchRole(guildID, roleID Snowflake) result.Result[Role] {
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    "/guilds/" + guildID.String() + "/roles/" + roleID.String(),
	})
	if res.IsErr() {
		return result.Err[Role](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var role Role
	if err := json.NewDecoder(body).Decode(&role); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/roles/{id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Role](err)
	}
	role.GuildID = guildID
	return result.Ok(role)
}

// FetchRolesMemberCount returns a map of role IDs to the number of members with the role.
//
//	Note:
//	 - Does not include the @everyone role.
func (r *requester) FetchRolesMemberCount(guildID Snowflake) result.Result[map[Snowflake]int] {
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    "/guilds/" + guildID.String() + "/roles/member-counts",
	})
	if res.IsErr() {
		return result.Err[map[Snowflake]int](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var rolesMemberCount map[Snowflake]int
	if err := json.NewDecoder(body).Decode(&rolesMemberCount); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/roles/member-counts",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[map[Snowflake]int](err)
	}
	return result.Ok(rolesMemberCount)
}

// CreateRoleOptions defines the configuration for creating a new Discord guild role.
type CreateRoleOptions struct {
	// Name is the channel's name (max 100 characters).
	//
	// Default to "new role" if unspesified.
	Name string `json:"name,omitempty"`

	// Permissions is the permission to set for the role.
	Permissions Permissions `json:"permissions,omitempty"`

	// Colors are the colors to set for the role.
	Colors RoleColors `json:"colors"`

	// Hoist is whether the role should be displayed separately in the sidebar.
	Hoist bool `json:"hoist,omitempty"`

	// Mentionable is whether the role should be displayed separately in the sidebar.
	Mentionable bool `json:"mentionable,omitempty"`

	// Icon is the role's icon image (if the guild has the GuildFeatureRoleIcons feature).
	Icon Base64Image `json:"icon,omitempty"`

	// UnicodeEmoji is the role's unicode emoji as a standard emoji (if the guild has the GuildFeatureRoleIcons feature).
	UnicodeEmoji string `json:"unicode_emoji,omitempty"`

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// CreateRole creates a new role for the guild.
//
// Requires the PermissionManageRoles permission.
func (r *requester) CreateRole(guildID Snowflake, opts CreateRoleOptions) result.Result[Role] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    "/guilds/" + guildID.String() + "/roles",
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[Role](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var role Role
	if err := json.NewDecoder(body).Decode(&role); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "POST",
			"url":    "/guilds/{id}/roles",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Role](err)
	}
	return result.Ok(role)
}

type RolePosition struct {
	// Channel id
	ID Snowflake `json:"id"`

	// Position is the sorting position of the channel (channels with the same position are sorted by id).
	Position optional.Option[int] `json:"position,omitzero"`
}

// ModifyRolePositionsOptions defines the configuration for modifying roles positions.
type ModifyRolePositionsOptions struct {
	Roles []RolePosition

	// Reason specifies the audit log reason for this action.
	Reason string `json:"-"`
}

// ModifyRolePositions modifies the positions of a set of roles.
//
// Requires the PermissionManageRoles permission.
func (r *requester) ModifyRolePositions(guildID Snowflake, opts ModifyRolePositionsOptions) result.Result[[]Role] {
	reqBody, _ := json.Marshal(opts.Roles)
	res := r.DoRequest(Request{
		Method: "PATCH",
		URL:    "/guilds/" + guildID.String() + "/roles",
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[[]Role](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var roles []Role
	if err := json.NewDecoder(body).Decode(&roles); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PATCH",
			"url":    "/guilds/{id}/roles",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]Role](err)
	}
	return result.Ok(roles)
}

// ModifyRoleOptions contains parameters for modifying a guild role.
type ModifyRoleOptions struct {
	// Name is the channel's name (max 100 characters).
	Name string `json:"name,omitempty"`

	// Permissions is the permission to set for the role.
	Permissions optional.Option[Permissions] `json:"permissions,omitzero"`

	// Colors are the colors to set for the role.
	Colors optional.Option[RoleColors] `json:"colors,omitzero"`

	// Hoist is whether the role should be displayed separately in the sidebar.
	Hoist optional.Option[bool] `json:"hoist,omitzero"`

	// Mentionable is whether the role should be displayed separately in the sidebar.
	Mentionable optional.Option[bool] `json:"mentionable,omitzero"`

	// Icon is the role's icon image (if the guild has the GuildFeatureRoleIcons feature).
	Icon optional.Option[Base64Image] `json:"icon,omitzero"`

	// UnicodeEmoji is the role's unicode emoji as a standard emoji (if the guild has the GuildFeatureRoleIcons feature).
	UnicodeEmoji optional.Option[string] `json:"unicode_emoji,omitzero"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// ModifyRole updates a member's properties in a guild.
//
// Requires the PermissionManageRoles permission.
func (r *requester) ModifyRole(guildID, roleID Snowflake, opts ModifyRoleOptions) result.Result[Role] {
	reqBody, _ := json.Marshal(opts)
	endpoint := "/guilds/" + guildID.String() + "/roles/" + roleID.String()

	body := r.DoRequest(Request{
		Method: "PATCH",
		URL:    endpoint,
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if body.IsErr() {
		return result.Err[Role](body.Err())
	}
	defer body.Value().Close()

	var role Role
	if err := json.NewDecoder(body.Value()).Decode(&role); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PATCH",
			"url":    "/guilds/{id}/roles/{role_id}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Role](err)
	}
	role.GuildID = guildID
	return result.Ok(role)
}

// DeleteRoleOptions contains parameters for deleting a guild role.
type DeleteRoleOptions struct {
	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// DeleteRole delete's a guild role.
//
// Requires the PermissionManageRoles permission.
func (r *requester) DeleteRole(guildID, roleID Snowflake, opts DeleteRoleOptions) result.Void {
	endpoint := "/guilds/" + guildID.String() + "/roles/" + roleID.String()
	body := r.DoRequest(Request{Method: "DELETE", URL: endpoint, Reason: opts.Reason})
	if body.IsErr() {
		return result.ErrVoid(body.Err())
	}
	body.Value().Close()
	return result.OkVoid()
}

// PruneCount contains parameters for fetching guild prune count.
type FetchGuildPruneCountOptions struct {
	// Pruned is the number of days to count prune for (1-30).
	Days int `json:"days,omitempty"`

	// By default, prune will not remove users with roles. You can optionally include specific roles in your prune by them in this field.
	// Any inactive user that has a subset of the provided role(s) will be counted in the prune and users with additional roles will not.
	IncludeRoles []Snowflake `json:"include_roles,omitempty"`
}

// PruneCount represents the result of a prune count.
type PruneCount struct {
	// Pruned is the number of members that would be/were removed in a prune operation.
	Pruned int `json:"pruned"`
}

// FetchGuildPruneCount returns an object with one pruned key indicating the number
// of members that would be removed in a prune operation.
//
// Requires the PermissionManageGuild and PermissionKickMembers permission.
func (r *requester) FetchGuildPruneCount(guildID Snowflake, opts FetchGuildPruneCountOptions) result.Result[PruneCount] {
	endpoint := "/guilds/" + guildID.String() + "/prune"

	params := url.Values{}
	if opts.Days > 0 {
		params.Set("days", strconv.Itoa(opts.Days))
	}
	if len(opts.IncludeRoles) > 0 {
		var b strings.Builder
		for i, snowflake := range opts.IncludeRoles {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(snowflake.String())
		}
		params.Set("include_roles", b.String())
	}
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[PruneCount](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var count PruneCount
	if err := json.NewDecoder(body).Decode(&count); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/prune",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[PruneCount](err)
	}
	return result.Ok(count)
}

// BeginGuildPruneOptions contains parameters for begining guild prune.
type BeginGuildPruneOptions struct {
	// Pruned is the number of days to count prune for (1-30).
	//
	// Default 7 if not set.
	Days int `json:"days,omitempty"`

	// ComputePruneCount is whether pruned is returned, discouraged for large guilds.
	ComputePruneCount optional.Option[bool] `json:"compute_prune_count,omitzero"`

	// IncludeRoles are the roles to include.
	IncludeRoles []Snowflake `json:"include_roles,omitempty"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// BeginGuildPrune Begin a prune operation.
//
// Requires the PermissionManageGuild and PermissionKickMembers permission.
func (r *requester) BeginGuildPrune(guildID Snowflake, opts BeginGuildPruneOptions) result.Result[PruneCount] {
	reqBody, _ := json.Marshal(opts)
	endpoint := "/guilds/" + guildID.String() + "/prune"
	res := r.DoRequest(Request{
		Method: "POST",
		URL:    endpoint,
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[PruneCount](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var count PruneCount
	if err := json.NewDecoder(body).Decode(&count); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "POST",
			"url":    "/guilds/{id}/prune",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[PruneCount](err)
	}
	return result.Ok(count)
}

// FetchGuildVoiceRegions returns a list of voice region objects for the guild.
// Unlike the similar /voice route, this returns VIP servers when the guild is VIP-enabled.
func (r *requester) FetchGuildVoiceRegions(guildID Snowflake) result.Result[[]VoiceRegion] {
	endpoint := "/guilds/" + guildID.String() + "/regions"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[[]VoiceRegion](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var regions []VoiceRegion
	if err := json.NewDecoder(body).Decode(&regions); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/regions",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]VoiceRegion](err)
	}
	return result.Ok(regions)
}

// FetchGuildVoiceRegions returns a list of voice region objects for the guild.
// Unlike the similar /voice route, this returns VIP servers when the guild is VIP-enabled.
//
// Requires the PermissionManageGuild and PermissionKickMembers permission.
func (r *requester) FetchGuildInvites(guildID Snowflake) result.Result[[]FullInvite] {
	endpoint := "/guilds/" + guildID.String() + "/invites"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[[]FullInvite](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var invites []FullInvite
	if err := json.NewDecoder(body).Decode(&invites); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/invites",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]FullInvite](err)
	}
	return result.Ok(invites)
}

// FetchGuildIntegrations returns a list of integration objects for the guild.
//
// Note:
//   - This endpoint returns a maximum of 50 integrations. If a guild has more integrations, they cannot be accessed.
//
// Requires the PermissionManageGuild permission.
func (r *requester) FetchGuildIntegrations(guildID Snowflake) result.Result[[]Integration] {
	endpoint := "/guilds/" + guildID.String() + "/integrations"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[[]Integration](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var integrations []Integration
	if err := json.NewDecoder(body).Decode(&integrations); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/integrations",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[[]Integration](err)
	}
	return result.Ok(integrations)
}

// DeleteGuildIntegrationOptions contains parameters for deleting a guild integration.
type DeleteGuildIntegrationOptions struct {
	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// DeleteGuildIntegration deletes the attached integration object for the guild.
// Deletes any associated webhooks and kicks the associated bot if there is one.
//
// Requires the PermissionManageGuild permission.
func (r *requester) DeleteGuildIntegration(guildID, integrationID Snowflake, opts DeleteGuildIntegrationOptions) result.Void {
	endpoint := "/guilds/" + guildID.String() + "/integrations/" + integrationID.String()
	res := r.DoRequest(Request{Method: "DELETE", URL: endpoint, Reason: opts.Reason})
	if res.IsErr() {
		return result.ErrVoid(res.Err())
	}
	res.Value().Close()
	return result.OkVoid()
}

// FetchGuildWidgetSettings returns a guild widget settings object.
//
// Requires the PermissionManageGuild permission.
func (r *requester) FetchGuildWidgetSettings(guildID Snowflake) result.Result[GuildWidgetSettings] {
	endpoint := "/guilds/" + guildID.String() + "/widget"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[GuildWidgetSettings](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var settings GuildWidgetSettings
	if err := json.NewDecoder(body).Decode(&settings); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/widget",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildWidgetSettings](err)
	}
	return result.Ok(settings)
}

// ModifyGuildWidgetOptions contains parameters for modifying a guild widget.
type ModifyGuildWidgetOptions struct {
	Enabled   optional.Option[bool] `json:"enabled,omitzero"`
	ChannelID Snowflake             `json:"channel_id,omitempty"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// ModifyGuildWidget modifies a guild widget settings object for the guild.
func (r *requester) ModifyGuildWidget(guildID Snowflake, opts ModifyGuildWidgetOptions) result.Result[GuildWidgetSettings] {
	reqBody, _ := json.Marshal(opts)
	res := r.DoRequest(Request{
		Method: "PATCH",
		URL:    "/guilds/" + guildID.String() + "/widget",
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[GuildWidgetSettings](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var settings GuildWidgetSettings
	if err := json.NewDecoder(body).Decode(&settings); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PATCH",
			"url":    "/guilds/{id}/widget",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildWidgetSettings](err)
	}
	return result.Ok(settings)
}

// FetchGuildWidget returns the widget for the guild.
func (r *requester) FetchGuildWidget(guildID Snowflake) result.Result[GuildWidget] {
	endpoint := "/guilds/" + guildID.String() + "/widget.json"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[GuildWidget](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var widget GuildWidget
	if err := json.NewDecoder(body).Decode(&widget); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/widget.json",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildWidget](err)
	}
	return result.Ok(widget)
}

// PartialInvite represents a guild's partial object.
//
// Reference: https://discord.com/developers/docs/resources/guild#get-guild-vanity-url
type PartialInvite struct {
	// Code is the vanity url code.
	Code string `json:"code,omitempty"`

	// Uses is the number of times this invite has been used.
	Uses int `json:"uses"`
}

// FetchGuildVanityURL returns a partial invite object for guilds with the feature enabled.
//
// Requires the PermissionManageGuild permission.
func (r *requester) FetchGuildVanityURL(guildID Snowflake) result.Result[PartialInvite] {
	endpoint := "/guilds/" + guildID.String() + "/vanity-url"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[PartialInvite](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var partialInvite PartialInvite
	if err := json.NewDecoder(body).Decode(&partialInvite); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/vanity-url",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[PartialInvite](err)
	}
	return result.Ok(partialInvite)
}

// GuildFeature represents the style of a Discord guild widget.
//
// Reference: https://discord.com/developers/docs/resources/guild#get-guild-widget-image-widget-style-options
type GuildWidgetStyle string

const (
	// Shield style widget with Discord icon and guild members online count
	GuildWidgetStyleShield GuildFeature = "shield"
	// Large image with guild icon, name and online count. "POWERED BY DISCORD" as the footer of the widget
	GuildWidgetStyleBanner1 GuildFeature = "banner1"
	// Smaller widget style with guild icon, name and online count. Split on the right with Discord logo
	GuildWidgetStyleBanner2 GuildFeature = "banner2"
	// Large image with guild icon, name and online count. In the footer, Discord logo on the left and "Chat Now" on the right
	GuildWidgetStyleBanner3 GuildFeature = "banner3"
	// large Discord logo at the top of the widget. Guild icon, name and online count in the middle portion of the widget
	// and a "JOIN MY SERVER" button at the bottom
	GuildWidgetStyleBanner4 GuildFeature = "banner4"
)

// FetchGuildWidgetImageOptions contains parameters for fetching guild widget image.
//
// Reference: https://discord.com/developers/docs/resources/guild#get-guild-widget-image
type FetchGuildWidgetImageOptions struct {
	// Style is the style of the widget image returned.
	Style GuildWidgetStyle `json:"style,omitempty"`
}

// FetchGuildWidgetImage returns a URL for a PNG image widget for the guild.
func (r *requester) FetchGuildWidgetImage(guildID Snowflake, opts FetchGuildWidgetImageOptions) string {
	url := "https://discord.com/api/v10/guilds/" + guildID.String() + "/widget.png"
	if opts.Style != "" {
		url += "?style=" + string(opts.Style)
	}
	return url
}

// FetchGuildWelcomeScreen returns the Welcome Screen object for the guild.
//
// Requires the PermissionManageGuild permission.
func (r *requester) FetchGuildWelcomeScreen(guildID Snowflake) result.Result[GuildWelcomeScreen] {
	endpoint := "/guilds/" + guildID.String() + "/welcome-screen"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[GuildWelcomeScreen](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var screen GuildWelcomeScreen
	if err := json.NewDecoder(body).Decode(&screen); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/welcome-screen",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildWelcomeScreen](err)
	}
	return result.Ok(screen)
}

// ModifyGuildWelcomeScreenOptions contains parameters for modifying a guild welcome screen.
//
// Reference: https://discord.com/developers/docs/resources/guild#modify-guild-welcome-screen
type ModifyGuildWelcomeScreenOptions struct {
	// Enabled is whether the welcome screen is enabled.
	Enabled optional.Option[bool] `json:"enabled,omitzero"`

	// WelcomeChannels are the channels linked in the welcome screen and their display options.
	WelcomeChannels optional.Option[[]GuildWelcomeChannel] `json:"welcome_channels,omitzero"`

	// Description is the the server description to show in the welcome screen.
	Description optional.Option[string] `json:"description,omitzero"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// ModifyGuildWelcomeScreen modifies the guild's Welcome Screen.
//
// Requires the PermissionManageGuild permission.
func (r *requester) ModifyGuildWelcomeScreen(guildID Snowflake, opts ModifyGuildWelcomeScreenOptions) result.Result[GuildWelcomeScreen] {
	reqBody, _ := json.Marshal(opts)
	endpoint := "/guilds/" + guildID.String() + "/welcome-screen"
	res := r.DoRequest(Request{
		Method: "PATCH",
		URL:    endpoint,
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[GuildWelcomeScreen](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var screen GuildWelcomeScreen
	if err := json.NewDecoder(body).Decode(&screen); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PATCH",
			"url":    "/guilds/{id}/welcome-screen",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildWelcomeScreen](err)
	}
	return result.Ok(screen)
}

// FetchGuildOnboarding returns the Onboarding object for the guild.
//
// Requires the PermissionManageGuild and PermissionManageRoles permissions.
func (r *requester) FetchGuildOnboarding(guildID Snowflake) result.Result[GuildOnboarding] {
	endpoint := "/guilds/" + guildID.String() + "/onboarding"

	res := r.DoRequest(Request{Method: "GET", URL: endpoint})
	if res.IsErr() {
		return result.Err[GuildOnboarding](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var onboarding GuildOnboarding
	if err := json.NewDecoder(body).Decode(&onboarding); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/guilds/{id}/onboarding",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildOnboarding](err)
	}
	return result.Ok(onboarding)
}

// ModifyGuildOnboardingOptions contains parameters for modifying guild onboarding.
//
// Reference: https://discord.com/developers/docs/resources/guild#modify-guild-onboarding
type ModifyGuildOnboardingOptions struct {
	// Prompts are the prompts shown during onboarding and in customize community.
	Prompts optional.Option[[]OnboardingPrompt] `json:"prompts,omitzero"`

	// DefaultChannelIDs are the channel IDs that members get opted into automatically.
	DefaultChannelIDs optional.Option[[]Snowflake] `json:"default_channel_ids,omitzero"`

	// Enabled is whether onboarding is enabled in the guild.
	Enabled optional.Option[bool] `json:"enabled,omitzero"`

	// Mode is the current mode of onboarding.
	Mode OnboardingMode `json:"mode,omitempty"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// ModifyGuildOnboarding modifies the onboarding configuration of the guild.
//
// Note:
//   - Onboarding enforces constraints when enabled. These constraints are that there must be at
//     least 7 Default Channels and at least 5 of them must allow sending messages to the @everyone role.
//     The mode field modifies what is considered when enforcing these constraints.
//
// Requires the PermissionManageGuild and PermissionManageRoles permissions.
func (r *requester) ModifyGuildOnboarding(guildID Snowflake, opts ModifyGuildOnboardingOptions) result.Result[GuildOnboarding] {
	reqBody, _ := json.Marshal(opts)
	endpoint := "/guilds/" + guildID.String() + "/onboarding"
	res := r.DoRequest(Request{
		Method: "PUT",
		URL:    endpoint,
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[GuildOnboarding](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var onboarding GuildOnboarding
	if err := json.NewDecoder(body).Decode(&onboarding); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PUT",
			"url":    "/guilds/{id}/onboarding",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildOnboarding](err)
	}
	return result.Ok(onboarding)
}

// ModifyGuildIncidentActionsOptions contains parameters for modifying guild incident actions.
//
// Reference: https://discord.com/developers/docs/resources/guild#modify-guild-incident-actions
type ModifyGuildIncidentActionsOptions struct {
	// InvitesDisabledUntil is when invites will be enabled again.
	//
	// Note: Supplying 'optional.Nil[time.Time]()' disables the action.
	InvitesDisabledUntil optional.Option[time.Time] `json:"invites_disabled_until,omitzero"`

	// DMsDisabledUntil is when direct messages will be enabled again.
	//
	// Note: Supplying 'optional.Nil[time.Time]()' disables the action.
	DMsDisabledUntil optional.Option[time.Time] `json:"dms_disabled_until,omitzero"`

	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// ModifyGuildIncidentActions modifies the incident actions of the guild.
//
// Requires the PermissionManageGuild permission.
func (r *requester) ModifyGuildIncidentActions(guildID Snowflake, opts ModifyGuildIncidentActionsOptions) result.Result[GuildIncidentsData] {
	reqBody, _ := json.Marshal(opts)
	endpoint := "/guilds/" + guildID.String() + "/incident-actions"
	res := r.DoRequest(Request{
		Method: "PUT",
		URL:    endpoint,
		Body:   reqBody,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[GuildIncidentsData](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var incidentsData GuildIncidentsData
	if err := json.NewDecoder(body).Decode(&incidentsData); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "PUT",
			"url":    "/guilds/{id}/incident-actions",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[GuildIncidentsData](err)
	}
	return result.Ok(incidentsData)
}
