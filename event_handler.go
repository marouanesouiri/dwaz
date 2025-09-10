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

import "github.com/bytedance/sonic"

/*****************************
 *   READY Handler
 *****************************/

// readyHandlers manages all registered handlers for MESSAGE_CREATE events.
type readyHandlers struct {
	logger   Logger
	handlers []func(ReadyEvent)
}

// handleEvent parses the READY event data and calls each registered handler.
func (h *readyHandlers) handleEvent(cache CacheManager, shardID int, data []byte) {
	evt := ReadyEvent{ShardsID: shardID}
	if err := sonic.Unmarshal(data, &evt); err != nil {
		h.logger.Error("readyHandlers: Failed parsing event data")
		return
	}

	for i := range len(evt.Guilds) {
		cache.PutGuild(evt.Guilds[i])
	}

	for _, handler := range h.handlers {
		handler(evt)
	}
}

// addHandler registers a new READY handler function.
//
// This method is not thread-safe.
func (h *readyHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ReadyEvent)))
}

/*****************************
 *   GUILD_CREATE Handler
 *****************************/

// guildCreateHandlers manages all registered handlers for GUILD_CREATE events.
type guildCreateHandlers struct {
	logger   Logger
	handlers []func(GuildCreateEvent)
}

// handleEvent parses the GUILD_CREATE event data and calls each registered handler.
func (h *guildCreateHandlers) handleEvent(cache CacheManager, shardID int, data []byte) {
	evt := GuildCreateEvent{ShardsID: shardID}

	if err := sonic.Unmarshal(data, &evt.Guild); err != nil {
		h.logger.Error("guildCreateHandlers: Failed parsing event data")
		return
	}

	flags := cache.Flags()

	if flags.Has(CacheFlagGuilds) {
		cache.PutGuild(evt.Guild.Guild)
	}
	if flags.Has(CacheFlagMembers) {
		for i := range len(evt.Guild.Members) {
			cache.PutMember(evt.Guild.Members[i])
		}
	}
	if flags.Has(CacheFlagChannels) {
		for i := range len(evt.Guild.Channels) {
			cache.PutChannel(evt.Guild.Channels[i])
		}
	}
	if flags.Has(CacheFlagRoles) {
		for i := range len(evt.Guild.Roles) {
			cache.PutRole(evt.Guild.Roles[i])
		}
	}
	if flags.Has(CacheFlagVoiceStates) {
		for i := range len(evt.Guild.VoiceStates) {
			cache.PutVoiceState(evt.Guild.VoiceStates[i])
		}
	}

	for _, handler := range h.handlers {
		handler(evt)
	}
}

// addHandler registers a new GUILD_CREATE handler function.
//
// This method is not thread-safe.
func (h *guildCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildCreateEvent)))
}

/*****************************
 *   MESSAGE_CREATE Handler
 *****************************/

// messageCreateHandlers manages all registered handlers for MESSAGE_CREATE events.
type messageCreateHandlers struct {
	logger   Logger
	handlers []func(MessageCreateEvent)
}

// handleEvent parses the MESSAGE_CREATE event data and calls each registered handler.
func (h *messageCreateHandlers) handleEvent(cache CacheManager, shardID int, data []byte) {
	evt := MessageCreateEvent{ShardsID: shardID}

	if err := sonic.Unmarshal(data, &evt.Message); err != nil {
		h.logger.Error("messageCreateHandlers: Failed parsing event data")
		return
	}

	if cache.Flags().Has(CacheFlagMessages) {
		cache.PutMessage(evt.Message)
	}

	for _, handler := range h.handlers {
		handler(evt)
	}
}

// addHandler registers a new MESSAGE_CREATE handler function.
//
// This method is not thread-safe.
func (h *messageCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessageCreateEvent)))
}

/*****************************
 *   MESSAGE_DELETE Handler
 *****************************/

// messageDeleteHandlers manages all registered handlers for MESSAGE_DELETE events.
type messageDeleteHandlers struct {
	logger   Logger
	handlers []func(MessageDeleteEvent)
}

// handleEvent parses the MESSAGE_DELETE event data and calls each registered handler.
func (h *messageDeleteHandlers) handleEvent(cache CacheManager, shardID int, data []byte) {
	evt := MessageDeleteEvent{ShardsID: shardID}
	if err := sonic.Unmarshal(data, &evt.Message); err != nil {
		h.logger.Error("messageDeleteHandlers: Failed parsing event data")
		return
	}

	if message, ok := cache.GetMessage(evt.Message.ID); ok {
		evt.Message = message
	}
	cache.DelMessage(evt.Message.ID)

	for _, handler := range h.handlers {
		handler(evt)
	}
}

// addHandler registers a new MESSAGE_DELETE handler function.
//
// This method is not thread-safe.
func (h *messageDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessageDeleteEvent)))
}

/*****************************
 *   MESSAGE_UPDATE Handler
 *****************************/

// messageUpdateHandlers manages all registered handlers for MESSAGE_UPDATE events.
type messageUpdateHandlers struct {
	logger   Logger
	handlers []func(MessageUpdateEvent)
}

// handleEvent parses the MESSAGE_UPDATE event data and calls each registered handler.
func (h *messageUpdateHandlers) handleEvent(cache CacheManager, shardID int, data []byte) {
	evt := MessageUpdateEvent{ShardsID: shardID}
	if err := sonic.Unmarshal(data, &evt.NewMessage); err != nil {
		h.logger.Error("messageUpdateHandlers: Failed parsing event data")
		return
	}

	if oldMessage, ok := cache.GetMessage(evt.NewMessage.ID); ok {
		evt.OldMessage = oldMessage
	} else {
		evt.OldMessage.ID = evt.NewMessage.ID
		evt.OldMessage.ChannelID = evt.NewMessage.ChannelID
		evt.OldMessage.GuildID = evt.NewMessage.GuildID
		evt.OldMessage.Author = evt.NewMessage.Author
		evt.OldMessage.Timestamp = evt.NewMessage.Timestamp
		evt.OldMessage.ApplicationID = evt.NewMessage.ApplicationID
	}

	if cache.Flags().Has(CacheFlagMessages) {
		cache.PutMessage(evt.NewMessage)
	}

	for _, handler := range h.handlers {
		handler(evt)
	}
}

// addHandler registers a new MESSAGE_UPDATE handler function.
//
// This method is not thread-safe.
func (h *messageUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessageUpdateEvent)))
}

/*****************************
 * INTERACTION_CREATE Handler
 *****************************/

// interactionCreateHandlers manages all registered handlers for INTERACTION_CREATE events.
type interactionCreateHandlers struct {
	logger   Logger
	handlers []func(InteractionCreateEvent)
}

// handleEvent parses the INTERACTION_CREATE event data and calls each registered handler.
func (h *interactionCreateHandlers) handleEvent(cache CacheManager, shardID int, data []byte) {
	evt := InteractionCreateEvent{ShardsID: shardID}
	if err := sonic.Unmarshal(data, &evt); err != nil {
		h.logger.Error("interactionCreateHandlers: Failed parsing event data")
		return
	}

	for _, handler := range h.handlers {
		handler(evt)
	}
}

// addHandler registers a new INTERACTION_CREATE handler function.
//
// This method is not thread-safe.
func (h *interactionCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(InteractionCreateEvent)))
}

/*****************************
 * VOICE_STATE_UPDATE Handler
 *****************************/

// voiceStateUpdateHandlers manages all registered handlers for VOICE_STATE_UPDATE events.
type voiceStateUpdateHandlers struct {
	logger   Logger
	handlers []func(VoiceStateUpdateEvent)
}

// handleEvent parses the VOICE_STATE_UPDATE event data and calls each registered handler.
func (h *voiceStateUpdateHandlers) handleEvent(cache CacheManager, shardID int, data []byte) {
	evt := VoiceStateUpdateEvent{ShardsID: shardID}
	if err := sonic.Unmarshal(data, &evt.NewState); err != nil {
		h.logger.Error("voiceStateCreateHandlers: Failed parsing event data")
		return
	}

	if oldVoiceState, ok := cache.GetVoiceState(evt.NewState.GuildID, evt.NewState.UserID); ok {
		evt.OldState = oldVoiceState
	} else {
		evt.OldState = evt.NewState
		evt.OldState.ChannelID = 0
	}

	if cache.Flags().Has(CacheFlagVoiceStates) {
		cache.PutVoiceState(evt.NewState)
	}

	for _, handler := range h.handlers {
		handler(evt)
	}
}

// addHandler registers a new VOICE_STATE_UPDATE handler function.
//
// This method is not thread-safe.
func (h *voiceStateUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(VoiceStateUpdateEvent)))
}
