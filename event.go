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

// ReadyCreateEvent Shard is ready
type ReadyEvent struct {
	ShardsID int // shard that dispatched this event
}

// MessageCreateEvent Message was created
type MessageCreateEvent struct {
	ShardsID int // shard that dispatched this event
}

// MessageDeleteEvent Message was deleted
type MessageDeleteEvent struct {
	ShardsID int // shard that dispatched this event
}

// TODO: add other events
