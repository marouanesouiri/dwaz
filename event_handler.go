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

	"github.com/marouanesouiri/stdx/xlog"
)

/*****************************
 *   READY Handler
 *****************************/

// readyHandlers manages all registered handlers for MESSAGE_CREATE events.
type readyHandlers struct {
	logger   xlog.Logger
	handlers []func(ReadyEvent)
}

// handleEvent parses the READY event data and calls each registered handler.
func (h *readyHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	evt := ReadyEvent{Client: client, ShardID: shardID}
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("readyHandlers: Failed parsing event data")
		return
	}

	for i := range len(evt.Guilds) {
		client.PutGuild(evt.Guilds[i])
	}

	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
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
	logger   xlog.Logger
	handlers []func(GuildCreateEvent)
}

// handleEvent parses the GUILD_CREATE event data and calls each registered handler.
func (h *guildCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	evt := GuildCreateEvent{Client: client, ShardID: shardID}

	if err := json.Unmarshal(data, &evt.Guild); err != nil {
		h.logger.Error("guildCreateHandlers: Failed parsing event data")
		return
	}

	flags := client.Flags()

	if flags.Has(CacheFlagGuilds) {
		client.PutGuild(evt.Guild.Guild)
	}
	if flags.Has(CacheFlagChannels) {
		for i := range len(evt.Guild.Channels) {
			client.PutChannel(evt.Guild.Channels[i])
		}
	}
	if flags.Has(CacheFlagRoles) {
		for i := range len(evt.Guild.Roles) {
			client.PutRole(evt.Guild.Roles[i])
		}
	}
	if flags.Has(CacheFlagVoiceStates) {
		for i := range len(evt.Guild.VoiceStates) {
			client.PutVoiceState(evt.Guild.VoiceStates[i])
		}
	}
	if flags.Has(CacheFlagMembers) {
		for i := range len(evt.Guild.Members) {
			client.PutMember(evt.Guild.Members[i].Member)
		}
	}
	if flags.Has(CacheFlagUsers) {
		for i := range len(evt.Guild.Members) {
			client.PutUser(evt.Guild.Members[i].User)
		}
	}

	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
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
	logger   xlog.Logger
	handlers []func(MessageCreateEvent)
}

// handleEvent parses the MESSAGE_CREATE event data and calls each registered handler.
func (h *messageCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	evt := MessageCreateEvent{Client: client, ShardID: shardID}

	if err := json.Unmarshal(data, &evt.Message); err != nil {
		h.logger.Error("messageCreateHandlers: Failed parsing event data")
		return
	}

	if client.Flags().Has(CacheFlagMessages) {
		client.PutMessage(evt.Message)
	}

	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
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
	logger   xlog.Logger
	handlers []func(MessageDeleteEvent)
}

// handleEvent parses the MESSAGE_DELETE event data and calls each registered handler.
func (h *messageDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	evt := MessageDeleteEvent{Client: client, ShardID: shardID}
	if err := json.Unmarshal(data, &evt.Message); err != nil {
		h.logger.Error("messageDeleteHandlers: Failed parsing event data")
		return
	}

	if msgOpt := client.GetMessage(evt.Message.ID); msgOpt.IsPresent() {
		evt.Message = msgOpt.Get()
	}
	client.DelMessage(evt.Message.ID)

	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
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
	logger   xlog.Logger
	handlers []func(MessageUpdateEvent)
}

// handleEvent parses the MESSAGE_UPDATE event data and calls each registered handler.
func (h *messageUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	evt := MessageUpdateEvent{Client: client, ShardID: shardID}
	if err := json.Unmarshal(data, &evt.NewMessage); err != nil {
		h.logger.Error("messageUpdateHandlers: Failed parsing event data")
		return
	}

	if oldMsgOpt := client.GetMessage(evt.NewMessage.ID); oldMsgOpt.IsPresent() {
		evt.OldMessage = oldMsgOpt.Get()
	} else {
		evt.OldMessage.ID = evt.NewMessage.ID
		evt.OldMessage.ChannelID = evt.NewMessage.ChannelID
		evt.OldMessage.GuildID = evt.NewMessage.GuildID
		evt.OldMessage.Author = evt.NewMessage.Author
		evt.OldMessage.Timestamp = evt.NewMessage.Timestamp
		evt.OldMessage.ApplicationID = evt.NewMessage.ApplicationID
	}

	if client.Flags().Has(CacheFlagMessages) {
		client.PutMessage(evt.NewMessage)
	}

	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
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
	logger   xlog.Logger
	handlers []func(InteractionCreateEvent)
}

// handleEvent parses the INTERACTION_CREATE event data and calls each registered handler.
func (h *interactionCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	evt := InteractionCreateEvent{Client: client, ShardID: shardID}
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("interactionCreateHandlers: Failed parsing event data")
		return
	}

	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
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
	logger   xlog.Logger
	handlers []func(VoiceStateUpdateEvent)
}

// handleEvent parses the VOICE_STATE_UPDATE event data and calls each registered handler.
func (h *voiceStateUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	evt := VoiceStateUpdateEvent{Client: client, ShardID: shardID}
	if err := json.Unmarshal(data, &evt.NewState); err != nil {
		h.logger.Error("voiceStateCreateHandlers: Failed parsing event data")
		return
	}

	if oldVoiceStateOpt := client.GetVoiceState(evt.NewState.GuildID, evt.NewState.UserID); oldVoiceStateOpt.IsPresent() {
		evt.OldState.VoiceState = oldVoiceStateOpt.Get()
	} else {
		evt.OldState = evt.NewState
		evt.OldState.ChannelID = 0
	}

	evt.OldState.Member = evt.NewState.Member

	if client.Flags().Has(CacheFlagVoiceStates) {
		client.PutVoiceState(evt.NewState.VoiceState)
	}

	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

// addHandler registers a new VOICE_STATE_UPDATE handler function.
//
// This method is not thread-safe.
func (h *voiceStateUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(VoiceStateUpdateEvent)))
}

/*********************************
 *   ApplicationCommandPermissionsUpdate Handler
 *********************************/

type applicationCommandPermissionsUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(ApplicationCommandPermissionsUpdateEvent)
}

func (h *applicationCommandPermissionsUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ApplicationCommandPermissionsUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("applicationCommandPermissionsUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *applicationCommandPermissionsUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ApplicationCommandPermissionsUpdateEvent)))
}

/*********************************
 *   AutoModeration Handlers
 *********************************/

type autoModerationRuleCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(AutoModerationRuleCreateEvent)
}

func (h *autoModerationRuleCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt AutoModerationRuleCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("autoModerationRuleCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *autoModerationRuleCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(AutoModerationRuleCreateEvent)))
}

type autoModerationRuleUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(AutoModerationRuleUpdateEvent)
}

func (h *autoModerationRuleUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt AutoModerationRuleUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("autoModerationRuleUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *autoModerationRuleUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(AutoModerationRuleUpdateEvent)))
}

type autoModerationRuleDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(AutoModerationRuleDeleteEvent)
}

func (h *autoModerationRuleDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt AutoModerationRuleDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("autoModerationRuleDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *autoModerationRuleDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(AutoModerationRuleDeleteEvent)))
}

type autoModerationActionExecutionHandlers struct {
	logger   xlog.Logger
	handlers []func(AutoModerationActionExecutionEvent)
}

func (h *autoModerationActionExecutionHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt AutoModerationActionExecutionEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("autoModerationActionExecutionHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *autoModerationActionExecutionHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(AutoModerationActionExecutionEvent)))
}

/*********************************
 *   Channel Handlers
 *********************************/

type channelCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(ChannelCreateEvent)
}

func (h *channelCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ChannelCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("channelCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *channelCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ChannelCreateEvent)))
}

type channelUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(ChannelUpdateEvent)
}

func (h *channelUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ChannelUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("channelUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *channelUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ChannelUpdateEvent)))
}

type channelDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(ChannelDeleteEvent)
}

func (h *channelDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ChannelDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("channelDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *channelDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ChannelDeleteEvent)))
}

type channelPinsUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(ChannelPinsUpdateEvent)
}

func (h *channelPinsUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ChannelPinsUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("channelPinsUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *channelPinsUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ChannelPinsUpdateEvent)))
}

/*********************************
 *   Thread Handlers
 *********************************/

type threadCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(ThreadCreateEvent)
}

func (h *threadCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ThreadCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("threadCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *threadCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ThreadCreateEvent)))
}

type threadUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(ThreadUpdateEvent)
}

func (h *threadUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ThreadUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("threadUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *threadUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ThreadUpdateEvent)))
}

type threadDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(ThreadDeleteEvent)
}

func (h *threadDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ThreadDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("threadDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *threadDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ThreadDeleteEvent)))
}

type threadListSyncHandlers struct {
	logger   xlog.Logger
	handlers []func(ThreadListSyncEvent)
}

func (h *threadListSyncHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ThreadListSyncEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("threadListSyncHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *threadListSyncHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ThreadListSyncEvent)))
}

type threadMemberUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(ThreadMemberUpdateEvent)
}

func (h *threadMemberUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ThreadMemberUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("threadMemberUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *threadMemberUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ThreadMemberUpdateEvent)))
}

type threadMembersUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(ThreadMembersUpdateEvent)
}

func (h *threadMembersUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt ThreadMembersUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("threadMembersUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *threadMembersUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(ThreadMembersUpdateEvent)))
}

/*********************************
 *   Entitlement Handlers
 *********************************/

type entitlementCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(EntitlementCreateEvent)
}

func (h *entitlementCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt EntitlementCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("entitlementCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *entitlementCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(EntitlementCreateEvent)))
}

type entitlementUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(EntitlementUpdateEvent)
}

func (h *entitlementUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt EntitlementUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("entitlementUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *entitlementUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(EntitlementUpdateEvent)))
}

type entitlementDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(EntitlementDeleteEvent)
}

func (h *entitlementDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt EntitlementDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("entitlementDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *entitlementDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(EntitlementDeleteEvent)))
}

/*********************************
 *   Guild Handlers
 *********************************/

type guildUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildUpdateEvent)
}

func (h *guildUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildUpdateEvent)))
}

type guildDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildDeleteEvent)
}

func (h *guildDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildDeleteEvent)))
}

type guildAuditLogEntryCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildAuditLogEntryCreateEvent)
}

func (h *guildAuditLogEntryCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildAuditLogEntryCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildAuditLogEntryCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildAuditLogEntryCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildAuditLogEntryCreateEvent)))
}

type guildBanAddHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildBanAddEvent)
}

func (h *guildBanAddHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildBanAddEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildBanAddHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildBanAddHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildBanAddEvent)))
}

type guildBanRemoveHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildBanRemoveEvent)
}

func (h *guildBanRemoveHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildBanRemoveEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildBanRemoveHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildBanRemoveHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildBanRemoveEvent)))
}

type guildEmojisUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildEmojisUpdateEvent)
}

func (h *guildEmojisUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildEmojisUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildEmojisUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildEmojisUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildEmojisUpdateEvent)))
}

type guildStickersUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildStickersUpdateEvent)
}

func (h *guildStickersUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildStickersUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildStickersUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildStickersUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildStickersUpdateEvent)))
}

type guildIntegrationsUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildIntegrationsUpdateEvent)
}

func (h *guildIntegrationsUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildIntegrationsUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildIntegrationsUpdateHandlers: Failed parsing event data")
		return
	}
	for _, handler := range h.handlers {
		if runAsync {
			go handler(evt)
		} else {
			handler(evt)
		}
	}
}

func (h *guildIntegrationsUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildIntegrationsUpdateEvent)))
}

type guildMemberAddHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildMemberAddEvent)
}

func (h *guildMemberAddHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildMemberAddEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildMemberAddHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildMemberAddHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildMemberAddEvent)))
}

type guildMemberRemoveHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildMemberRemoveEvent)
}

func (h *guildMemberRemoveHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildMemberRemoveEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildMemberRemoveHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildMemberRemoveHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildMemberRemoveEvent)))
}

type guildMemberUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildMemberUpdateEvent)
}

func (h *guildMemberUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildMemberUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildMemberUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildMemberUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildMemberUpdateEvent)))
}

type guildMembersChunkHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildMembersChunkEvent)
}

func (h *guildMembersChunkHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildMembersChunkEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildMembersChunkHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildMembersChunkHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildMembersChunkEvent)))
}

type guildRoleCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildRoleCreateEvent)
}

func (h *guildRoleCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildRoleCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildRoleCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildRoleCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildRoleCreateEvent)))
}

type guildRoleUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildRoleUpdateEvent)
}

func (h *guildRoleUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildRoleUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildRoleUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildRoleUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildRoleUpdateEvent)))
}

type guildRoleDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildRoleDeleteEvent)
}

func (h *guildRoleDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildRoleDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildRoleDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildRoleDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildRoleDeleteEvent)))
}

type guildScheduledEventCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildScheduledEventCreateEvent)
}

func (h *guildScheduledEventCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildScheduledEventCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildScheduledEventCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildScheduledEventCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildScheduledEventCreateEvent)))
}

type guildScheduledEventUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildScheduledEventUpdateEvent)
}

func (h *guildScheduledEventUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildScheduledEventUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildScheduledEventUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildScheduledEventUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildScheduledEventUpdateEvent)))
}

type guildScheduledEventDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildScheduledEventDeleteEvent)
}

func (h *guildScheduledEventDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildScheduledEventDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildScheduledEventDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildScheduledEventDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildScheduledEventDeleteEvent)))
}

type guildScheduledEventUserAddHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildScheduledEventUserAddEvent)
}

func (h *guildScheduledEventUserAddHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildScheduledEventUserAddEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildScheduledEventUserAddHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildScheduledEventUserAddHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildScheduledEventUserAddEvent)))
}

type guildScheduledEventUserRemoveHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildScheduledEventUserRemoveEvent)
}

func (h *guildScheduledEventUserRemoveHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildScheduledEventUserRemoveEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildScheduledEventUserRemoveHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildScheduledEventUserRemoveHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildScheduledEventUserRemoveEvent)))
}

type guildSoundboardSoundCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildSoundboardSoundCreateEvent)
}

func (h *guildSoundboardSoundCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildSoundboardSoundCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildSoundboardSoundCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildSoundboardSoundCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildSoundboardSoundCreateEvent)))
}

type guildSoundboardSoundUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildSoundboardSoundUpdateEvent)
}

func (h *guildSoundboardSoundUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildSoundboardSoundUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildSoundboardSoundUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildSoundboardSoundUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildSoundboardSoundUpdateEvent)))
}

type guildSoundboardSoundDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildSoundboardSoundDeleteEvent)
}

func (h *guildSoundboardSoundDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildSoundboardSoundDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildSoundboardSoundDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildSoundboardSoundDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildSoundboardSoundDeleteEvent)))
}

type guildSoundboardSoundsUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(GuildSoundboardSoundsUpdateEvent)
}

func (h *guildSoundboardSoundsUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt GuildSoundboardSoundsUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("guildSoundboardSoundsUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *guildSoundboardSoundsUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(GuildSoundboardSoundsUpdateEvent)))
}

type soundboardSoundsHandlers struct {
	logger   xlog.Logger
	handlers []func(SoundboardSoundsEvent)
}

func (h *soundboardSoundsHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt SoundboardSoundsEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("soundboardSoundsHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *soundboardSoundsHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(SoundboardSoundsEvent)))
}

/*********************************
 *   Integration Handlers
 *********************************/

type integrationCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(IntegrationCreateEvent)
}

func (h *integrationCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt IntegrationCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("integrationCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *integrationCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(IntegrationCreateEvent)))
}

type integrationUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(IntegrationUpdateEvent)
}

func (h *integrationUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt IntegrationUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("integrationUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *integrationUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(IntegrationUpdateEvent)))
}

type integrationDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(IntegrationDeleteEvent)
}

func (h *integrationDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt IntegrationDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("integrationDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *integrationDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(IntegrationDeleteEvent)))
}

/*********************************
 *   Invite Handlers
 *********************************/

type inviteCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(InviteCreateEvent)
}

func (h *inviteCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt InviteCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("inviteCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *inviteCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(InviteCreateEvent)))
}

type inviteDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(InviteDeleteEvent)
}

func (h *inviteDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt InviteDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("inviteDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *inviteDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(InviteDeleteEvent)))
}

/*********************************
 *   Message Misc Handlers
 *********************************/

type messageDeleteBulkHandlers struct {
	logger   xlog.Logger
	handlers []func(MessageDeleteBulkEvent)
}

func (h *messageDeleteBulkHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt MessageDeleteBulkEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("messageDeleteBulkHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *messageDeleteBulkHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessageDeleteBulkEvent)))
}

type messageReactionAddHandlers struct {
	logger   xlog.Logger
	handlers []func(MessageReactionAddEvent)
}

func (h *messageReactionAddHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt MessageReactionAddEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("messageReactionAddHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *messageReactionAddHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessageReactionAddEvent)))
}

type messageReactionRemoveHandlers struct {
	logger   xlog.Logger
	handlers []func(MessageReactionRemoveEvent)
}

func (h *messageReactionRemoveHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt MessageReactionRemoveEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("messageReactionRemoveHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *messageReactionRemoveHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessageReactionRemoveEvent)))
}

type messageReactionRemoveAllHandlers struct {
	logger   xlog.Logger
	handlers []func(MessageReactionRemoveAllEvent)
}

func (h *messageReactionRemoveAllHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt MessageReactionRemoveAllEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("messageReactionRemoveAllHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *messageReactionRemoveAllHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessageReactionRemoveAllEvent)))
}

type messageReactionRemoveEmojiHandlers struct {
	logger   xlog.Logger
	handlers []func(MessageReactionRemoveEmojiEvent)
}

func (h *messageReactionRemoveEmojiHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt MessageReactionRemoveEmojiEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("messageReactionRemoveEmojiHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *messageReactionRemoveEmojiHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessageReactionRemoveEmojiEvent)))
}

type messagePollVoteAddHandlers struct {
	logger   xlog.Logger
	handlers []func(MessagePollVoteAddEvent)
}

func (h *messagePollVoteAddHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt MessagePollVoteAddEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("messagePollVoteAddHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *messagePollVoteAddHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessagePollVoteAddEvent)))
}

type messagePollVoteRemoveHandlers struct {
	logger   xlog.Logger
	handlers []func(MessagePollVoteRemoveEvent)
}

func (h *messagePollVoteRemoveHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt MessagePollVoteRemoveEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("messagePollVoteRemoveHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *messagePollVoteRemoveHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(MessagePollVoteRemoveEvent)))
}

/*********************************
 *   Presence Handlers
 *********************************/

type presenceUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(PresenceUpdateEvent)
}

func (h *presenceUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt PresenceUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("presenceUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *presenceUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(PresenceUpdateEvent)))
}

/*********************************
 *   Stage Instance Handlers
 *********************************/

type stageInstanceCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(StageInstanceCreateEvent)
}

func (h *stageInstanceCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt StageInstanceCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("stageInstanceCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *stageInstanceCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(StageInstanceCreateEvent)))
}

type stageInstanceUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(StageInstanceUpdateEvent)
}

func (h *stageInstanceUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt StageInstanceUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("stageInstanceUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *stageInstanceUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(StageInstanceUpdateEvent)))
}

type stageInstanceDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(StageInstanceDeleteEvent)
}

func (h *stageInstanceDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt StageInstanceDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("stageInstanceDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *stageInstanceDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(StageInstanceDeleteEvent)))
}

/*********************************
 *   Subscription Handlers
 *********************************/

type subscriptionCreateHandlers struct {
	logger   xlog.Logger
	handlers []func(SubscriptionCreateEvent)
}

func (h *subscriptionCreateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt SubscriptionCreateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("subscriptionCreateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *subscriptionCreateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(SubscriptionCreateEvent)))
}

type subscriptionUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(SubscriptionUpdateEvent)
}

func (h *subscriptionUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt SubscriptionUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("subscriptionUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *subscriptionUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(SubscriptionUpdateEvent)))
}

type subscriptionDeleteHandlers struct {
	logger   xlog.Logger
	handlers []func(SubscriptionDeleteEvent)
}

func (h *subscriptionDeleteHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt SubscriptionDeleteEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("subscriptionDeleteHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *subscriptionDeleteHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(SubscriptionDeleteEvent)))
}

/*********************************
 *   Typing Handlers
 *********************************/

type typingStartHandlers struct {
	logger   xlog.Logger
	handlers []func(TypingStartEvent)
}

func (h *typingStartHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt TypingStartEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("typingStartHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *typingStartHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(TypingStartEvent)))
}

/*********************************
 *   User Handlers
 *********************************/

type userUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(UserUpdateEvent)
}

func (h *userUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt UserUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("userUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *userUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(UserUpdateEvent)))
}

/*********************************
 *   Voice Handlers
 *********************************/

type voiceChannelEffectSendHandlers struct {
	logger   xlog.Logger
	handlers []func(VoiceChannelEffectSendEvent)
}

func (h *voiceChannelEffectSendHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt VoiceChannelEffectSendEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("voiceChannelEffectSendHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *voiceChannelEffectSendHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(VoiceChannelEffectSendEvent)))
}

type voiceServerUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(VoiceServerUpdateEvent)
}

func (h *voiceServerUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt VoiceServerUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("voiceServerUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *voiceServerUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(VoiceServerUpdateEvent)))
}

/*********************************
 *   Webhooks Handlers
 *********************************/

type webhooksUpdateHandlers struct {
	logger   xlog.Logger
	handlers []func(WebhooksUpdateEvent)
}

func (h *webhooksUpdateHandlers) handleEvent(client *Client, runAsync bool, shardID int, data []byte) {
	var evt WebhooksUpdateEvent
	if err := json.Unmarshal(data, &evt); err != nil {
		h.logger.Error("webhooksUpdateHandlers: Failed parsing event data")
		return
	}
	if runAsync {
		for _, handler := range h.handlers {
			go handler(evt)
		}
	} else {
		for _, handler := range h.handlers {
			handler(evt)
		}
	}
}

func (h *webhooksUpdateHandlers) addHandler(handler any) {
	h.handlers = append(h.handlers, handler.(func(WebhooksUpdateEvent)))
}
