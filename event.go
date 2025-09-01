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

import "encoding/json"

// ReadyCreateEvent Shard is ready
type ReadyEvent struct {
	ShardsID int // shard that dispatched this event
}

// MessageCreateEvent Message was created
type MessageCreateEvent struct {
	ShardsID int // shard that dispatched this event
	Message Message
}

// MessageCreateEvent Message was created
type MessageUpdateEvent struct {
	ShardsID   int // shard that dispatched this event
	OldMessage Message
	NewMessage Message
}

// MessageDeleteEvent Message was deleted
type MessageDeleteEvent struct {
	ShardsID int // shard that dispatched this event
	Message Message
}

// InteractionCreateEvent Interaction created
type InteractionCreateEvent struct {
	ShardsID int // shard that dispatched this event
	Interaction Interaction
}

var _ json.Unmarshaler = (*InteractionCreateEvent)(nil)

// UnmarshalJSON implements json.Unmarshaler for InteractionCreateEvent.
func (c *InteractionCreateEvent) UnmarshalJSON(buf []byte) error {
	interaction, err := UnmarshalInteraction(buf)
	if err == nil {
		c.Interaction = interaction
	}
	return err
}

// TODO: add other events
