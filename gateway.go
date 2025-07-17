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

import "github.com/bytedance/sonic"

// gateway holds the Discord Gateway URL.
type gateway struct {
	// WSS URL that can be used for connecting to the Gateway
	URL string `json:"url"`
}

func (o *gateway) fillFromJson(json []byte) error {
	return sonic.Unmarshal(json, o)
}

// gatewayBot is Discord Gateway Bot.
type gatewayBot struct {
	// WSS URL that can be used for connecting to the Gateway
	URL string `json:"url"`
	// Recommended number of shards to use when connecting
	Shards int `json:"shards"`
	// Information on the current session start limit
	SessionStartLimit int `json:"session_start_limit"`
}

func (o *gatewayBot) fillFromJson(json []byte) error {
	return sonic.Unmarshal(json, o)
}
