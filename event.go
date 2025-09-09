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

	"github.com/bytedance/sonic"
)

// ReadyCreateEvent Shard is ready
type ReadyEvent struct {
	ShardsID int // shard that dispatched this event
}

// MessageCreateEvent Message was created
type MessageCreateEvent struct {
	ShardsID int // shard that dispatched this event
	Message  Message
}

var _ json.Unmarshaler = (*MessageCreateEvent)(nil)

// UnmarshalJSON implements json.Unmarshaler for MessageCreateEvent.
func (e *MessageCreateEvent) UnmarshalJSON(buf []byte) error {
	return sonic.Unmarshal(buf, &e.Message)
}

// MessageCreateEvent Message was created
type MessageUpdateEvent struct {
	ShardsID   int // shard that dispatched this event
	OldMessage Message
	NewMessage Message
}

var _ json.Unmarshaler = (*MessageUpdateEvent)(nil)

// UnmarshalJSON implements json.Unmarshaler for MessageCreateEvent.
func (e *MessageUpdateEvent) UnmarshalJSON(buf []byte) error {
	return sonic.Unmarshal(buf, &e.NewMessage)
}

// MessageDeleteEvent Message was deleted
type MessageDeleteEvent struct {
	ShardsID int // shard that dispatched this event
	Message  Message
}

var _ json.Unmarshaler = (*MessageDeleteEvent)(nil)

// UnmarshalJSON implements json.Unmarshaler for MessageDeleteEvent.
func (e *MessageDeleteEvent) UnmarshalJSON(buf []byte) error {
	return sonic.Unmarshal(buf, &e.Message)
}

// InteractionCreateEvent Interaction created
type InteractionCreateEvent struct {
	ShardsID    int // shard that dispatched this event
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

// VoiceStateUpdateEvent VoiceState was updated
type VoiceStateUpdateEvent struct {
	ShardsID int // shard that dispatched this event
	OldState VoiceState
	NewState VoiceState
}

var _ json.Unmarshaler = (*VoiceStateUpdateEvent)(nil)

// UnmarshalJSON implements json.Unmarshaler for InteractionCreateEvent.
func (c *VoiceStateUpdateEvent) UnmarshalJSON(buf []byte) error {
	return sonic.Unmarshal(buf, c.NewState)
}

// TODO: add other events
