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
	"github.com/bytedance/sonic"
)

/*****************************
 *   READY Handler
 *****************************/

// readyHandlers manages all registered handlers for MESSAGE_CREATE events.
type readyHandlers struct {
	logger   Logger
	handlers []func(ReadyEvent)
}

// handleEvent parses the READY event data and calls each registered handler.
func (h *readyHandlers) handleEvent(shardID int, data []byte) {
	evt := ReadyEvent{ShardsID: shardID}
	if err := sonic.Unmarshal(data, &evt); err != nil {
		h.logger.Error("readyHandlers: Failed parsing event data")
		return
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
 *   MESSAGE_CREATE Handler
 *****************************/

// messageCreateHandlers manages all registered handlers for MESSAGE_CREATE events.
type messageCreateHandlers struct {
	logger   Logger
	handlers []func(MessageCreateEvent)
}

// handleEvent parses the MESSAGE_CREATE event data and calls each registered handler.
func (h *messageCreateHandlers) handleEvent(shardID int, data []byte) {
	evt := MessageCreateEvent{ShardsID: shardID}

	if err := sonic.Unmarshal(data, &evt); err != nil {
		h.logger.Error("messageCreateHandlers: Failed parsing event data")
		return
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
func (h *messageDeleteHandlers) handleEvent(shardID int, data []byte) {
	evt := MessageDeleteEvent{ShardsID: shardID}
	if err := sonic.Unmarshal(data, &evt); err != nil {
		h.logger.Error("messageDeleteHandlers: Failed parsing event data")
		return
	}

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
