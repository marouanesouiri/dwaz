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

// Locale represents a Discord supported locale code.
//
// Reference: https://discord.com/developers/docs/reference#locales
type Locale string

const (
	LocaleIndonesian        Locale = "id"     // Bahasa Indonesia
	LocaleDanish            Locale = "da"     // Dansk
	LocaleGerman            Locale = "de"     // Deutsch
	LocaleEnglishUK         Locale = "en-GB"  // English, UK
	LocaleEnglishUS         Locale = "en-US"  // English, US
	LocaleSpanishSpain      Locale = "es-ES"  // Español
	LocaleSpanishLatam      Locale = "es-419" // Español, LATAM
	LocaleFrench            Locale = "fr"     // Français
	LocaleCroatian          Locale = "hr"     // Hrvatski
	LocaleItalian           Locale = "it"     // Italiano
	LocaleLithuanian        Locale = "lt"     // Lietuviškai
	LocaleHungarian         Locale = "hu"     // Magyar
	LocaleDutch             Locale = "nl"     // Nederlands
	LocaleNorwegian         Locale = "no"     // Norsk
	LocalePolish            Locale = "pl"     // Polski
	LocalePortugueseBrazil  Locale = "pt-BR"  // Português do Brasil
	LocaleRomanian          Locale = "ro"     // Română
	LocaleFinnish           Locale = "fi"     // Suomi
	LocaleSwedish           Locale = "sv-SE"  // Svenska
	LocaleVietnamese        Locale = "vi"     // Tiếng Việt
	LocaleTurkish           Locale = "tr"     // Türkçe
	LocaleCzech             Locale = "cs"     // Čeština
	LocaleGreek             Locale = "el"     // Ελληνικά
	LocaleBulgarian         Locale = "bg"     // български
	LocaleRussian           Locale = "ru"     // Pусский
	LocaleUkrainian         Locale = "uk"     // Українська
	LocaleHindi             Locale = "hi"     // हिन्दी
	LocaleChineseChina      Locale = "zh-CN"  // 中文
	LocaleJapanese          Locale = "ja"     // 日本語
	LocaleChineseTaiwan     Locale = "zh-TW"  // 繁體中文
	LocaleKorean            Locale = "ko"     // 한국어
)
