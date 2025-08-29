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
	"os"
	"runtime/debug"
	"sync"
)

/*****************************
 *   EventhandlersManager
 *****************************/

// eventhandlersManager defines the interface for managing event handlers of a specific event type.
//
// Implementations must support adding handlers and dispatching raw JSON event data to those handlers.
type eventhandlersManager interface {
	// handleEvent unmarshals the raw JSON data and calls all registered handlers.
	handleEvent(shardID int, data []byte)
	// addHandler adds a new handler function for the event type.
	addHandler(h any)
}

/*****************************
 *        dispatcher
 *****************************/

// dispatcher manages registration of event handlers and dispatching of events.
//
// It stores handlers by event name string and invokes the correct handlers for incoming events.
//
// WARNING:
//   - This implementation is not fully thread-safe for handler registration. You must register
//     all handlers sequentially before starting event dispatching (usually at startup).
//   - Dispatching handlers is done asynchronously in separate goroutines for each event.
type dispatcher struct {
	logger           Logger
	workerPool       WorkerPool
	handlersManagers map[string]eventhandlersManager
	mu               sync.RWMutex
}

// newDispatcher creates a new dispatcher instance.
//
// If logger is nil, it creates a default logger that writes to os.Stdout with debug-level logging.
func newDispatcher(logger Logger, workerPool WorkerPool) *dispatcher {
	if logger == nil {
		logger = NewDefaultLogger(os.Stdout, LogLevel_DebugLevel)
	}
	return &dispatcher{
		logger:           logger,
		workerPool:       workerPool,
		handlersManagers: make(map[string]eventhandlersManager, 20),
	}
}

/*****************************
 *     Dispatch Event
 *****************************/

// dispatch sends raw event JSON data to all registered handlers for that event name.
//
// The eventName must exactly match the Discord event string (e.g., "MESSAGE_CREATE").
//
// This method spawns a new goroutine for each dispatch to avoid blocking the main event loop.
func (d *dispatcher) dispatch(shardID int, eventName string, data []byte) {
	d.logger.Debug("Event '" + eventName + "' dispatched")
	if !d.workerPool.Submit(func() {
		defer func() {
			if r := recover(); r != nil {
				d.logger.WithField("event", eventName).
					WithField("shard_id", shardID).
					WithField("panic", r).
					WithField("stack", string(debug.Stack())).
					Error("Recovered from panic while handling event")
			}
		}()

		d.mu.RLock()
		hm, ok := d.handlersManagers[eventName]
		d.mu.RUnlock()

		if ok {
			hm.handleEvent(shardID, data)
		}
	}) {
		d.logger.Warn("Dispatcher: dropped event '" + eventName + "' due to full queue")
	}
}

/*****************************
 *      Register Handlers
 *****************************/

// OnMessageCreate registers a handler function for 'MESSAGE_CREATE' events.
//
// Note:
//   - This method is thread-safe via internal locking.
//   - However, it is strongly recommended to register all event handlers sequentially during startup,
//     before starting event dispatching, to avoid runtime mutations and ensure stable configuration.
//   - Handlers are called sequentially when dispatching in the order they were added.
func (d *dispatcher) OnMessageCreate(h func(MessageCreateEvent)) {
	const key = "MESSAGE_CREATE" // event name
	d.logger.Debug(key + " event handler registered")

	d.mu.Lock()
	defer d.mu.Unlock()

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessageDelete registers a handler function for 'MESSAGE_DELETE' events.
//
// Note:
//   - This method is thread-safe via internal locking.
//   - However, it is strongly recommended to register all event handlers sequentially during startup,
//     before starting event dispatching, to avoid runtime mutations and ensure stable configuration.
//   - Handlers are called sequentially when dispatching in the order they were added.
func (d *dispatcher) OnMessageDelete(h func(MessageDeleteEvent)) {
	const key = "MESSAGE_DELETE" // event name
	d.logger.Debug(key + " event handler registered")

	d.mu.Lock()
	defer d.mu.Unlock()

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// TODO: Add other OnXXX methods to register handlers for additional Discord events.
