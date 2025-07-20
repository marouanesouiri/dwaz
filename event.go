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

// ReadyCreateEvent Shard is ready
type ReadyEvent struct {
	ShardsID int // shard that dispatched this event
}

func (e *ReadyEvent) fillFromJson(json []byte) error {
	return sonic.Unmarshal(json, e)
}

// MessageCreateEvent Message was created
type MessageCreateEvent struct {
	ShardsID int // shard that dispatched this event
}

func (e *MessageCreateEvent) fillFromJson(json []byte) error {
	return sonic.Unmarshal(json, e)
}

// MessageDeleteEvent Message was deleted
type MessageDeleteEvent struct {
	ShardsID int // shard that dispatched this event
}

func (e *MessageDeleteEvent) fillFromJson(json []byte) error {
	return sonic.Unmarshal(json, e)
}

// TODO: add other events
