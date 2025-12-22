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
	"runtime/debug"

	"github.com/marouanesouiri/stdx/xlog"
)

/*****************************
 *   EventhandlersManager
 *****************************/

// eventhandlersManager defines the interface for managing event handlers of a specific event type.
//
// Implementations must support adding handlers and dispatching raw JSON event data to those handlers.
type eventhandlersManager interface {
	// handleEvent unmarshals the raw JSON data and calls all registered handlers.
	handleEvent(cache CacheManager, runAsync bool, shardID int, buf []byte)
	// addHandler adds a new handler function for the event type.
	addHandler(handler any)
}

/*****************************
 *        dispatcher
 *****************************/

// HandlerExecutionMode defines how event handlers are executed.
type HandlerExecutionMode int

const (
	// HandlerExecutionSync runs all handlers for an event sequentially in the dispatch goroutine.
	HandlerExecutionSync HandlerExecutionMode = iota
	// HandlerExecutionAsync runs each handler for an event in a separate goroutine/task using the executor.
	HandlerExecutionAsync
)

// dispatcher manages registration of event handlers and dispatching of events.
//
// It stores handlers by event name string and invokes the correct handlers for incoming events.
//
// WARNING:
//   - This implementation is not fully thread-safe for handler registration. You must register
//     all handlers sequentially before starting event dispatching (usually at startup).
//   - Dispatching handlers is done asynchronously in separate goroutines for each event.
type dispatcher struct {
	logger               xlog.Logger
	cacheManager         CacheManager
	handlersManagers     map[string]eventhandlersManager
	handlerExecutionMode HandlerExecutionMode
}

// newDispatcher creates a new dispatcher instance.
//
// If logger is nil, it creates a default logger that writes to os.Stdout with debug-level logging.
func newDispatcher(logger xlog.Logger, cacheManager CacheManager, mode HandlerExecutionMode) *dispatcher {
	if logger == nil {
		logger = xlog.NewTextLogger(nil, xlog.LogLevelInfoLevel)
	}
	d := &dispatcher{
		logger:               logger,
		cacheManager:         cacheManager,
		handlerExecutionMode: mode,
		handlersManagers:     make(map[string]eventhandlersManager, 20),
	}

	// Register some necessary events for caching
	d.handlersManagers["READY"] = &readyHandlers{logger: logger}
	d.handlersManagers["GUILD_CREATE"] = &guildCreateHandlers{logger: logger}

	return d
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
	d.logger.WithFields(map[string]any{
		"shard_id": shardID,
		"event":    eventName,
	}).Debug("event dispatched")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				d.logger.WithField("event", eventName).
					WithField("shard_id", shardID).
					WithField("panic", r).
					WithField("stack", string(debug.Stack())).
					Error("Recovered from panic while handling event")
			}
		}()

		if hm, ok := d.handlersManagers[eventName]; ok {
			hm.handleEvent(d.cacheManager, d.handlerExecutionMode == HandlerExecutionAsync, shardID, data)
		}
	}()
}

/*****************************
 *      Register Handlers
 *****************************/

// OnMessageCreate registers a handler function for 'MESSAGE_CREATE' events.
func (d *dispatcher) OnMessageCreate(h func(MessageCreateEvent)) {
	const key = "MESSAGE_CREATE" // event name
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessageDelete registers a handler function for 'MESSAGE_DELETE' events.
func (d *dispatcher) OnMessageDelete(h func(MessageDeleteEvent)) {
	const key = "MESSAGE_DELETE" // event name
	d.logger.Debug(key + " event handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessageUpdate registers a handler function for 'MESSAGE_UPDATE' events.
func (d *dispatcher) OnMessageUpdate(h func(MessageDeleteEvent)) {
	const key = "MESSAGE_UPDATE" // event name
	d.logger.Debug(key + " event handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnInteractionCreate registers a handler function for 'INTERACTION_CREATE' events.
func (d *dispatcher) OnInteractionCreate(h func(InteractionCreateEvent)) {
	const key = "INTERACTION_CREATE" // event name
	d.logger.Debug(key + " event handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &interactionCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnVoiceStateUpdate registers a handler function for 'VOICE_STATE_UPDATE' events.
func (d *dispatcher) OnVoiceStateUpdate(h func(VoiceStateUpdateEvent)) {
	const key = "VOICE_STATE_UPDATE" // event name
	d.logger.Debug(key + " event handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &voiceStateUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnApplicationCommandPermissionsUpdate registers a handler for 'APPLICATION_COMMAND_PERMISSIONS_UPDATE' events.
func (d *dispatcher) OnApplicationCommandPermissionsUpdate(h func(ApplicationCommandPermissionsUpdateEvent)) {
	const key = "APPLICATION_COMMAND_PERMISSIONS_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &applicationCommandPermissionsUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnAutoModerationRuleCreate registers a handler for 'AUTO_MODERATION_RULE_CREATE' events.
func (d *dispatcher) OnAutoModerationRuleCreate(h func(AutoModerationRuleCreateEvent)) {
	const key = "AUTO_MODERATION_RULE_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &autoModerationRuleCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnAutoModerationRuleUpdate registers a handler for 'AUTO_MODERATION_RULE_UPDATE' events.
func (d *dispatcher) OnAutoModerationRuleUpdate(h func(AutoModerationRuleUpdateEvent)) {
	const key = "AUTO_MODERATION_RULE_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &autoModerationRuleUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnAutoModerationRuleDelete registers a handler for 'AUTO_MODERATION_RULE_DELETE' events.
func (d *dispatcher) OnAutoModerationRuleDelete(h func(AutoModerationRuleDeleteEvent)) {
	const key = "AUTO_MODERATION_RULE_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &autoModerationRuleDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnAutoModerationActionExecution registers a handler for 'AUTO_MODERATION_ACTION_EXECUTION' events.
func (d *dispatcher) OnAutoModerationActionExecution(h func(AutoModerationActionExecutionEvent)) {
	const key = "AUTO_MODERATION_ACTION_EXECUTION"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &autoModerationActionExecutionHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnChannelCreate registers a handler for 'CHANNEL_CREATE' events.
func (d *dispatcher) OnChannelCreate(h func(ChannelCreateEvent)) {
	const key = "CHANNEL_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &channelCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnChannelUpdate registers a handler for 'CHANNEL_UPDATE' events.
func (d *dispatcher) OnChannelUpdate(h func(ChannelUpdateEvent)) {
	const key = "CHANNEL_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &channelUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnChannelDelete registers a handler for 'CHANNEL_DELETE' events.
func (d *dispatcher) OnChannelDelete(h func(ChannelDeleteEvent)) {
	const key = "CHANNEL_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &channelDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnChannelPinsUpdate registers a handler for 'CHANNEL_PINS_UPDATE' events.
func (d *dispatcher) OnChannelPinsUpdate(h func(ChannelPinsUpdateEvent)) {
	const key = "CHANNEL_PINS_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &channelPinsUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnThreadCreate registers a handler for 'THREAD_CREATE' events.
func (d *dispatcher) OnThreadCreate(h func(ThreadCreateEvent)) {
	const key = "THREAD_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &threadCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnThreadUpdate registers a handler for 'THREAD_UPDATE' events.
func (d *dispatcher) OnThreadUpdate(h func(ThreadUpdateEvent)) {
	const key = "THREAD_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &threadUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnThreadDelete registers a handler for 'THREAD_DELETE' events.
func (d *dispatcher) OnThreadDelete(h func(ThreadDeleteEvent)) {
	const key = "THREAD_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &threadDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnThreadListSync registers a handler for 'THREAD_LIST_SYNC' events.
func (d *dispatcher) OnThreadListSync(h func(ThreadListSyncEvent)) {
	const key = "THREAD_LIST_SYNC"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &threadListSyncHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnThreadMemberUpdate registers a handler for 'THREAD_MEMBER_UPDATE' events.
func (d *dispatcher) OnThreadMemberUpdate(h func(ThreadMemberUpdateEvent)) {
	const key = "THREAD_MEMBER_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &threadMemberUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnThreadMembersUpdate registers a handler for 'THREAD_MEMBERS_UPDATE' events.
func (d *dispatcher) OnThreadMembersUpdate(h func(ThreadMembersUpdateEvent)) {
	const key = "THREAD_MEMBERS_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &threadMembersUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnEntitlementCreate registers a handler for 'ENTITLEMENT_CREATE' events.
func (d *dispatcher) OnEntitlementCreate(h func(EntitlementCreateEvent)) {
	const key = "ENTITLEMENT_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &entitlementCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnEntitlementUpdate registers a handler for 'ENTITLEMENT_UPDATE' events.
func (d *dispatcher) OnEntitlementUpdate(h func(EntitlementUpdateEvent)) {
	const key = "ENTITLEMENT_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &entitlementUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnEntitlementDelete registers a handler for 'ENTITLEMENT_DELETE' events.
func (d *dispatcher) OnEntitlementDelete(h func(EntitlementDeleteEvent)) {
	const key = "ENTITLEMENT_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &entitlementDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildUpdate registers a handler for 'GUILD_UPDATE' events.
func (d *dispatcher) OnGuildUpdate(h func(GuildUpdateEvent)) {
	const key = "GUILD_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildDelete registers a handler for 'GUILD_DELETE' events.
func (d *dispatcher) OnGuildDelete(h func(GuildDeleteEvent)) {
	const key = "GUILD_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildAuditLogEntryCreate registers a handler for 'GUILD_AUDIT_LOG_ENTRY_CREATE' events.
func (d *dispatcher) OnGuildAuditLogEntryCreate(h func(GuildAuditLogEntryCreateEvent)) {
	const key = "GUILD_AUDIT_LOG_ENTRY_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildAuditLogEntryCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildBanAdd registers a handler for 'GUILD_BAN_ADD' events.
func (d *dispatcher) OnGuildBanAdd(h func(GuildBanAddEvent)) {
	const key = "GUILD_BAN_ADD"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildBanAddHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildBanRemove registers a handler for 'GUILD_BAN_REMOVE' events.
func (d *dispatcher) OnGuildBanRemove(h func(GuildBanRemoveEvent)) {
	const key = "GUILD_BAN_REMOVE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildBanRemoveHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildEmojisUpdate registers a handler for 'GUILD_EMOJIS_UPDATE' events.
func (d *dispatcher) OnGuildEmojisUpdate(h func(GuildEmojisUpdateEvent)) {
	const key = "GUILD_EMOJIS_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildEmojisUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildStickersUpdate registers a handler for 'GUILD_STICKERS_UPDATE' events.
func (d *dispatcher) OnGuildStickersUpdate(h func(GuildStickersUpdateEvent)) {
	const key = "GUILD_STICKERS_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildStickersUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildIntegrationsUpdate registers a handler for 'GUILD_INTEGRATIONS_UPDATE' events.
func (d *dispatcher) OnGuildIntegrationsUpdate(h func(GuildIntegrationsUpdateEvent)) {
	const key = "GUILD_INTEGRATIONS_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildIntegrationsUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildMemberAdd registers a handler for 'GUILD_MEMBER_ADD' events.
func (d *dispatcher) OnGuildMemberAdd(h func(GuildMemberAddEvent)) {
	const key = "GUILD_MEMBER_ADD"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildMemberAddHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildMemberRemove registers a handler for 'GUILD_MEMBER_REMOVE' events.
func (d *dispatcher) OnGuildMemberRemove(h func(GuildMemberRemoveEvent)) {
	const key = "GUILD_MEMBER_REMOVE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildMemberRemoveHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildMemberUpdate registers a handler for 'GUILD_MEMBER_UPDATE' events.
func (d *dispatcher) OnGuildMemberUpdate(h func(GuildMemberUpdateEvent)) {
	const key = "GUILD_MEMBER_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildMemberUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildMembersChunk registers a handler for 'GUILD_MEMBERS_CHUNK' events.
func (d *dispatcher) OnGuildMembersChunk(h func(GuildMembersChunkEvent)) {
	const key = "GUILD_MEMBERS_CHUNK"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildMembersChunkHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildRoleCreate registers a handler for 'GUILD_ROLE_CREATE' events.
func (d *dispatcher) OnGuildRoleCreate(h func(GuildRoleCreateEvent)) {
	const key = "GUILD_ROLE_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildRoleCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildRoleUpdate registers a handler for 'GUILD_ROLE_UPDATE' events.
func (d *dispatcher) OnGuildRoleUpdate(h func(GuildRoleUpdateEvent)) {
	const key = "GUILD_ROLE_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildRoleUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildRoleDelete registers a handler for 'GUILD_ROLE_DELETE' events.
func (d *dispatcher) OnGuildRoleDelete(h func(GuildRoleDeleteEvent)) {
	const key = "GUILD_ROLE_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildRoleDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildScheduledEventCreate registers a handler for 'GUILD_SCHEDULED_EVENT_CREATE' events.
func (d *dispatcher) OnGuildScheduledEventCreate(h func(GuildScheduledEventCreateEvent)) {
	const key = "GUILD_SCHEDULED_EVENT_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildScheduledEventCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildScheduledEventUpdate registers a handler for 'GUILD_SCHEDULED_EVENT_UPDATE' events.
func (d *dispatcher) OnGuildScheduledEventUpdate(h func(GuildScheduledEventUpdateEvent)) {
	const key = "GUILD_SCHEDULED_EVENT_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildScheduledEventUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildScheduledEventDelete registers a handler for 'GUILD_SCHEDULED_EVENT_DELETE' events.
func (d *dispatcher) OnGuildScheduledEventDelete(h func(GuildScheduledEventDeleteEvent)) {
	const key = "GUILD_SCHEDULED_EVENT_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildScheduledEventDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildScheduledEventUserAdd registers a handler for 'GUILD_SCHEDULED_EVENT_USER_ADD' events.
func (d *dispatcher) OnGuildScheduledEventUserAdd(h func(GuildScheduledEventUserAddEvent)) {
	const key = "GUILD_SCHEDULED_EVENT_USER_ADD"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildScheduledEventUserAddHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildScheduledEventUserRemove registers a handler for 'GUILD_SCHEDULED_EVENT_USER_REMOVE' events.
func (d *dispatcher) OnGuildScheduledEventUserRemove(h func(GuildScheduledEventUserRemoveEvent)) {
	const key = "GUILD_SCHEDULED_EVENT_USER_REMOVE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildScheduledEventUserRemoveHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildSoundboardSoundCreate registers a handler for 'GUILD_SOUNDBOARD_SOUND_CREATE' events.
func (d *dispatcher) OnGuildSoundboardSoundCreate(h func(GuildSoundboardSoundCreateEvent)) {
	const key = "GUILD_SOUNDBOARD_SOUND_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildSoundboardSoundCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildSoundboardSoundUpdate registers a handler for 'GUILD_SOUNDBOARD_SOUND_UPDATE' events.
func (d *dispatcher) OnGuildSoundboardSoundUpdate(h func(GuildSoundboardSoundUpdateEvent)) {
	const key = "GUILD_SOUNDBOARD_SOUND_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildSoundboardSoundUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildSoundboardSoundDelete registers a handler for 'GUILD_SOUNDBOARD_SOUND_DELETE' events.
func (d *dispatcher) OnGuildSoundboardSoundDelete(h func(GuildSoundboardSoundDeleteEvent)) {
	const key = "GUILD_SOUNDBOARD_SOUND_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildSoundboardSoundDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnGuildSoundboardSoundsUpdate registers a handler for 'GUILD_SOUNDBOARD_SOUNDS_UPDATE' events.
func (d *dispatcher) OnGuildSoundboardSoundsUpdate(h func(GuildSoundboardSoundsUpdateEvent)) {
	const key = "GUILD_SOUNDBOARD_SOUNDS_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &guildSoundboardSoundsUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnSoundboardSounds registers a handler for 'SOUNDBOARD_SOUNDS' events.
func (d *dispatcher) OnSoundboardSounds(h func(SoundboardSoundsEvent)) {
	const key = "SOUNDBOARD_SOUNDS"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &soundboardSoundsHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnIntegrationCreate registers a handler for 'INTEGRATION_CREATE' events.
func (d *dispatcher) OnIntegrationCreate(h func(IntegrationCreateEvent)) {
	const key = "INTEGRATION_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &integrationCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnIntegrationUpdate registers a handler for 'INTEGRATION_UPDATE' events.
func (d *dispatcher) OnIntegrationUpdate(h func(IntegrationUpdateEvent)) {
	const key = "INTEGRATION_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &integrationUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnIntegrationDelete registers a handler for 'INTEGRATION_DELETE' events.
func (d *dispatcher) OnIntegrationDelete(h func(IntegrationDeleteEvent)) {
	const key = "INTEGRATION_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &integrationDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnInviteCreate registers a handler for 'INVITE_CREATE' events.
func (d *dispatcher) OnInviteCreate(h func(InviteCreateEvent)) {
	const key = "INVITE_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &inviteCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnInviteDelete registers a handler for 'INVITE_DELETE' events.
func (d *dispatcher) OnInviteDelete(h func(InviteDeleteEvent)) {
	const key = "INVITE_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &inviteDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessageDeleteBulk registers a handler for 'MESSAGE_DELETE_BULK' events.
func (d *dispatcher) OnMessageDeleteBulk(h func(MessageDeleteBulkEvent)) {
	const key = "MESSAGE_DELETE_BULK"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageDeleteBulkHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessageReactionAdd registers a handler for 'MESSAGE_REACTION_ADD' events.
func (d *dispatcher) OnMessageReactionAdd(h func(MessageReactionAddEvent)) {
	const key = "MESSAGE_REACTION_ADD"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageReactionAddHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessageReactionRemove registers a handler for 'MESSAGE_REACTION_REMOVE' events.
func (d *dispatcher) OnMessageReactionRemove(h func(MessageReactionRemoveEvent)) {
	const key = "MESSAGE_REACTION_REMOVE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageReactionRemoveHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessageReactionRemoveAll registers a handler for 'MESSAGE_REACTION_REMOVE_ALL' events.
func (d *dispatcher) OnMessageReactionRemoveAll(h func(MessageReactionRemoveAllEvent)) {
	const key = "MESSAGE_REACTION_REMOVE_ALL"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageReactionRemoveAllHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessageReactionRemoveEmoji registers a handler for 'MESSAGE_REACTION_REMOVE_EMOJI' events.
func (d *dispatcher) OnMessageReactionRemoveEmoji(h func(MessageReactionRemoveEmojiEvent)) {
	const key = "MESSAGE_REACTION_REMOVE_EMOJI"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messageReactionRemoveEmojiHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnPresenceUpdate registers a handler for 'PRESENCE_UPDATE' events.
func (d *dispatcher) OnPresenceUpdate(h func(PresenceUpdateEvent)) {
	const key = "PRESENCE_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &presenceUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnStageInstanceCreate registers a handler for 'STAGE_INSTANCE_CREATE' events.
func (d *dispatcher) OnStageInstanceCreate(h func(StageInstanceCreateEvent)) {
	const key = "STAGE_INSTANCE_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &stageInstanceCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnStageInstanceUpdate registers a handler for 'STAGE_INSTANCE_UPDATE' events.
func (d *dispatcher) OnStageInstanceUpdate(h func(StageInstanceUpdateEvent)) {
	const key = "STAGE_INSTANCE_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &stageInstanceUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnStageInstanceDelete registers a handler for 'STAGE_INSTANCE_DELETE' events.
func (d *dispatcher) OnStageInstanceDelete(h func(StageInstanceDeleteEvent)) {
	const key = "STAGE_INSTANCE_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &stageInstanceDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnSubscriptionCreate registers a handler for 'SUBSCRIPTION_CREATE' events.
func (d *dispatcher) OnSubscriptionCreate(h func(SubscriptionCreateEvent)) {
	const key = "SUBSCRIPTION_CREATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &subscriptionCreateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnSubscriptionUpdate registers a handler for 'SUBSCRIPTION_UPDATE' events.
func (d *dispatcher) OnSubscriptionUpdate(h func(SubscriptionUpdateEvent)) {
	const key = "SUBSCRIPTION_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &subscriptionUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnSubscriptionDelete registers a handler for 'SUBSCRIPTION_DELETE' events.
func (d *dispatcher) OnSubscriptionDelete(h func(SubscriptionDeleteEvent)) {
	const key = "SUBSCRIPTION_DELETE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &subscriptionDeleteHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnTypingStart registers a handler for 'TYPING_START' events.
func (d *dispatcher) OnTypingStart(h func(TypingStartEvent)) {
	const key = "TYPING_START"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &typingStartHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnUserUpdate registers a handler for 'USER_UPDATE' events.
func (d *dispatcher) OnUserUpdate(h func(UserUpdateEvent)) {
	const key = "USER_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &userUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnVoiceChannelEffectSend registers a handler for 'VOICE_CHANNEL_EFFECT_SEND' events.
func (d *dispatcher) OnVoiceChannelEffectSend(h func(VoiceChannelEffectSendEvent)) {
	const key = "VOICE_CHANNEL_EFFECT_SEND"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &voiceChannelEffectSendHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnVoiceServerUpdate registers a handler for 'VOICE_SERVER_UPDATE' events.
func (d *dispatcher) OnVoiceServerUpdate(h func(VoiceServerUpdateEvent)) {
	const key = "VOICE_SERVER_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &voiceServerUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnWebhooksUpdate registers a handler for 'WEBHOOKS_UPDATE' events.
func (d *dispatcher) OnWebhooksUpdate(h func(WebhooksUpdateEvent)) {
	const key = "WEBHOOKS_UPDATE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &webhooksUpdateHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessagePollVoteAdd registers a handler for 'MESSAGE_POLL_VOTE_ADD' events.
func (d *dispatcher) OnMessagePollVoteAdd(h func(MessagePollVoteAddEvent)) {
	const key = "MESSAGE_POLL_VOTE_ADD"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messagePollVoteAddHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}

// OnMessagePollVoteRemove registers a handler for 'MESSAGE_POLL_VOTE_REMOVE' events.
func (d *dispatcher) OnMessagePollVoteRemove(h func(MessagePollVoteRemoveEvent)) {
	const key = "MESSAGE_POLL_VOTE_REMOVE"
	d.logger.WithField("event", key).Debug("handler registered")

	hm, ok := d.handlersManagers[key]
	if !ok {
		hm = &messagePollVoteRemoveHandlers{logger: d.logger}
		d.handlersManagers[key] = hm
	}
	hm.addHandler(h)
}
