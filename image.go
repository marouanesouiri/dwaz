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
	"strconv"
)

// NOTE:
// Image format enums are duplicated per endpoint section because
// not all Discord image endpoints support the same set of image formats.
// This duplication ensures type safety and clear intent per endpoint.
// Also, the CDN base URLs are defined as constants for maintainability,
// since some media like Sticker GIFs do NOT use the CDN base URL and
// instead use a different base URL for serving content.

const (
	ImageBaseUrl = "https://cdn.discordapp.com/"
	MediaBaseUrl = "https://media.discordapp.net/"
)

type ImageSize int

const (
	ImageSize_16   ImageSize = 16
	ImageSize_32   ImageSize = 32
	ImageSize_64   ImageSize = 64
	ImageSize_128  ImageSize = 128
	ImageSize_256  ImageSize = 256
	ImageSize_512  ImageSize = 512
	ImageSize_1024 ImageSize = 1024
	ImageSize_2048 ImageSize = 2048
	ImageSize_4096 ImageSize = 4096
)

/***********************
 *    Emoji Endpoints   *
 ***********************/

type EmojiFormat string

const (
	EmojiFormat_PNG  EmojiFormat = ".png"
	EmojiFormat_JPEG EmojiFormat = ".jpeg"
	EmojiFormat_WebP EmojiFormat = ".webp"
	EmojiFormat_GIF  EmojiFormat = ".gif"
	EmojiFormat_AVIF EmojiFormat = ".avif"
)

func EmojiURL(EmojiID Snowflake, format EmojiFormat, size ImageSize) string {
	return ImageBaseUrl + "emojis/" + EmojiID.String() + string(format) + "?size=" + strconv.Itoa(int(size))
}

/***********************
 *   Guild Endpoints    *
 ***********************/

type GuildIconFormat string

const (
	GuildIconFormat_PNG  GuildIconFormat = ".png"
	GuildIconFormat_JPEG GuildIconFormat = ".jpeg"
	GuildIconFormat_WebP GuildIconFormat = ".webp"
	GuildIconFormat_GIF  GuildIconFormat = ".gif"
)

func GuildIconURL(guildID Snowflake, iconHash string, format GuildIconFormat, size ImageSize) string {
	if format == GuildIconFormat_GIF && (len(iconHash) < 2 || iconHash[:2] != "a_") {
		format = GuildIconFormat_PNG
	}

	url := ImageBaseUrl + "icons/" + guildID.String() + "/" + iconHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == GuildIconFormat_WebP && len(iconHash) >= 2 && iconHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

type GuildSplashFormat string

const (
	GuildSplashFormat_PNG  GuildSplashFormat = ".png"
	GuildSplashFormat_JPEG GuildSplashFormat = ".jpeg"
	GuildSplashFormat_WebP GuildSplashFormat = ".webp"
)

func GuildSplashURL(guildID Snowflake, splashHash string, format GuildSplashFormat, size ImageSize) string {
	return ImageBaseUrl + "splashes/" + guildID.String() + "/" + splashHash + string(format) + "?size=" + strconv.Itoa(int(size))
}

type GuildBannerFormat string

const (
	GuildBannerFormat_PNG  GuildBannerFormat = ".png"
	GuildBannerFormat_JPEG GuildBannerFormat = ".jpeg"
	GuildBannerFormat_WebP GuildBannerFormat = ".webp"
	GuildBannerFormat_GIF  GuildBannerFormat = ".gif"
)

func GuildBannerURL(guildID Snowflake, bannerHash string, format GuildBannerFormat, size ImageSize) string {
	if format == GuildBannerFormat_GIF && (len(bannerHash) < 2 || bannerHash[:2] != "a_") {
		format = GuildBannerFormat_PNG
	}

	url := ImageBaseUrl + "banners/" + guildID.String() + "/" + bannerHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == GuildBannerFormat_WebP && len(bannerHash) >= 2 && bannerHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

/***********************
 *    User Endpoints    *
 ***********************/

type UserAvatarFormat string

const (
	UserAvatarFormat_PNG  UserAvatarFormat = ".png"
	UserAvatarFormat_JPEG UserAvatarFormat = ".jpeg"
	UserAvatarFormat_WebP UserAvatarFormat = ".webp"
	UserAvatarFormat_GIF  UserAvatarFormat = ".gif"
)

func UserAvatarURL(userID Snowflake, avatarHash string, format UserAvatarFormat, size ImageSize) string {
	if format == UserAvatarFormat_GIF && (len(avatarHash) < 2 || avatarHash[:2] != "a_") {
		format = UserAvatarFormat_PNG
	}

	url := ImageBaseUrl + "avatars/" + userID.String() + "/" + avatarHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == UserAvatarFormat_WebP && len(avatarHash) >= 2 && avatarHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

type UserBannerFormat string

const (
	UserBannerFormat_PNG  UserBannerFormat = ".png"
	UserBannerFormat_JPEG UserBannerFormat = ".jpeg"
	UserBannerFormat_WebP UserBannerFormat = ".webp"
	UserBannerFormat_GIF  UserBannerFormat = ".gif"
)

func UserBannerURL(userID Snowflake, bannerHash string, format UserBannerFormat, size ImageSize) string {
	if format == UserBannerFormat_GIF && (len(bannerHash) < 2 || bannerHash[:2] != "a_") {
		format = UserBannerFormat_PNG
	}

	url := ImageBaseUrl + "banners/" + userID.String() + "/" + bannerHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == UserBannerFormat_WebP && len(bannerHash) >= 2 && bannerHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

/***********************
 * Application Endpoints *
 ***********************/

type ApplicationIconFormat string

const (
	ApplicationIconFormat_PNG  ApplicationIconFormat = ".png"
	ApplicationIconFormat_JPEG ApplicationIconFormat = ".jpeg"
	ApplicationIconFormat_WebP ApplicationIconFormat = ".webp"
)

func ApplicationIconURL(appID Snowflake, iconHash string, format ApplicationIconFormat, size ImageSize) string {
	return ImageBaseUrl + "app-icons/" + appID.String() + "/" + iconHash + string(format) + "?size=" + strconv.Itoa(int(size))
}

type ApplicationCoverFormat string

const (
	ApplicationCoverFormat_PNG  ApplicationCoverFormat = ".png"
	ApplicationCoverFormat_JPEG ApplicationCoverFormat = ".jpeg"
	ApplicationCoverFormat_WebP ApplicationCoverFormat = ".webp"
)

func ApplicationCoverURL(appID Snowflake, coverHash string, format ApplicationCoverFormat, size ImageSize) string {
	return ImageBaseUrl + "app-icons/" + appID.String() + "/" + coverHash + string(format) + "?size=" + strconv.Itoa(int(size))
}

/***********************
 *     Sticker Endpoints *
 ***********************/

type StickerFormat string

const (
	StickerFormat_PNG    StickerFormat = ".png"
	StickerFormat_GIF    StickerFormat = ".gif"
	StickerFormat_Lottie StickerFormat = ".json"
)

// Stickers with GIF format are served from MediaBaseUrl, not CDN base.
func StickerURL(stickerID Snowflake, format StickerFormat) string {
	base := ImageBaseUrl + "stickers/" + stickerID.String()
	if format == StickerFormat_GIF {
		base = MediaBaseUrl + "stickers/" + stickerID.String()
	}
	return base + string(format)
}

/***********************
 *      Other Endpoints  *
 ***********************/

// Default User Avatar: embed/avatars/index.png
// Size param ignored, fixed size only
func DefaultUserAvatarURL(index int) string {
	return ImageBaseUrl + "embed/avatars/" + strconv.Itoa(index) + ".png"
}

type GuildMemberAvatarFormat string

const (
	GuildMemberAvatarFormat_PNG  GuildMemberAvatarFormat = ".png"
	GuildMemberAvatarFormat_JPEG GuildMemberAvatarFormat = ".jpeg"
	GuildMemberAvatarFormat_WebP GuildMemberAvatarFormat = ".webp"
	GuildMemberAvatarFormat_GIF  GuildMemberAvatarFormat = ".gif"
)

func GuildMemberAvatarURL(guildID, userID Snowflake, avatarHash string, format GuildMemberAvatarFormat, size ImageSize) string {
	if format == GuildMemberAvatarFormat_GIF && (len(avatarHash) < 2 || avatarHash[:2] != "a_") {
		format = GuildMemberAvatarFormat_PNG
	}

	url := ImageBaseUrl + "guilds/" + guildID.String() + "/users/" + userID.String() + "/avatars/" + avatarHash + string(format) + "?size=" + strconv.Itoa(int(size))

	if format == GuildMemberAvatarFormat_WebP && len(avatarHash) >= 2 && avatarHash[:2] == "a_" {
		url += "&animated=true"
	}

	return url
}

type RoleIconFormat string

const (
	RoleIconFormat_PNG  RoleIconFormat = ".png"
	RoleIconFormat_JPEG RoleIconFormat = ".jpeg"
	RoleIconFormat_WebP RoleIconFormat = ".webp"
)

func RoleIconURL(roleID Snowflake, iconHash string, format RoleIconFormat, size ImageSize) string {
	return ImageBaseUrl + "role-icons/" + roleID.String() + "/" + iconHash + string(format) + "?size=" + strconv.Itoa(int(size))
}
