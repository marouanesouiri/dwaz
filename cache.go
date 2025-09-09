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

import "sync"

type CacheFlags int

const (
	CacheFlagUsers CacheFlags = 1 << iota
	CacheFlagGuilds
	CacheFlagMembers
	CacheFlagThreadMembers
	CacheFlagMessages
	CacheFlagChannels
	CacheFlagRoles
	CacheFlagVoiceStates

	CacheFlagsNone CacheFlags = 0

	CacheFlagsAll = CacheFlagUsers | CacheFlagGuilds | CacheFlagMembers | CacheFlagThreadMembers |
		CacheFlagMessages | CacheFlagChannels | CacheFlagRoles | CacheFlagVoiceStates
)

func (f CacheFlags) Has(bits ...CacheFlags) bool {
	return BitFieldHas(f, bits...)
}

type SnowflakePairKey struct {
	A Snowflake
	B Snowflake
}

type CacheManager interface {
	Flags() CacheFlags
	SetFlags(flags ...CacheFlags)

	GetUser(userID Snowflake) (User, bool)
	GetGuild(guildID Snowflake) (Guild, bool)
	GetMember(guildID, userID Snowflake) (Member, bool)
	GetChannel(channelID Snowflake) (Channel, bool)
	GetMessage(messageID Snowflake) (Message, bool)
	GetVoiceState(guildID, userID Snowflake) (VoiceState, bool)
	GetGuildChannels(guildID Snowflake) (map[Snowflake]Channel, bool)
	GetGuildMembers(guildID Snowflake) (map[Snowflake]Member, bool)
	GetGuildVoiceStates(guildID Snowflake) (map[Snowflake]VoiceState, bool)
	GetGuildRoles(guildID Snowflake) (map[Snowflake]Role, bool)

	CountUsers() int
	CountGuilds() int
	CountMembers() int
	CountChannels() int
	CountMessages() int
	CountVoiceStates() int
	CountRoles() int
	CountGuildChannels(guildID Snowflake) int
	CountGuildMembers(guildID Snowflake) int
	CountGuildRoles(guildID Snowflake) int

	PutUser(user User)
	PutGuild(guild Guild)
	PutMember(member Member)
	PutChannel(channel Channel)
	PutMessage(message Message)
	PutVoiceState(voiceState VoiceState)
	PutRole(role Role)

	DelUser(userID Snowflake) bool
	DelGuild(guildID Snowflake) bool
	DelMember(guildID, userID Snowflake) bool
	DelChannel(channelID Snowflake) bool
	DelMessage(messageID Snowflake) bool
	DelVoiceState(guildID, userID Snowflake) bool
	DelGuildChannels(guildID Snowflake) bool
	DelGuildMembers(guildID Snowflake) bool
	DelRole(guildID, roleID Snowflake) bool
}

type DefaultCache struct {
	flags CacheFlags

	usersCache   map[Snowflake]User
	usersCacheMu sync.RWMutex

	guildsCache   map[Snowflake]Guild
	guildsCacheMu sync.RWMutex

	membersCache   map[SnowflakePairKey]Member
	membersCacheMu sync.RWMutex

	channelsCache   map[Snowflake]Channel
	channelsCacheMu sync.RWMutex

	messagesCache   map[Snowflake]Message
	messagesCacheMu sync.RWMutex

	voiceStatesCache   map[SnowflakePairKey]VoiceState
	voiceStatesCacheMu sync.RWMutex

	rolesCache   map[Snowflake]Role
	rolesCacheMu sync.RWMutex

	// Index: guildID -> set[userID]
	guildToMemberIDs   map[Snowflake]map[Snowflake]struct{}
	guildToMemberIDsMu sync.RWMutex

	// Index: guildID -> map[channelID]
	guildToChannelIDs   map[Snowflake]map[Snowflake]struct{}
	guildToChannelIDsMu sync.RWMutex

	// Index: guildID -> map[userID]
	guildToVoiceStateUserIDs   map[Snowflake]map[Snowflake]struct{}
	guildToVoiceStateUserIDsMu sync.RWMutex

	// Index: guildID -> map[roleID]
	guildToRoleIDs   map[Snowflake]map[Snowflake]struct{}
	guildToRoleIDsMu sync.RWMutex
}

func NewDefaultCache(flags CacheFlags) CacheManager {
	return &DefaultCache{
		flags:                    flags,
		usersCache:               make(map[Snowflake]User),
		guildsCache:              make(map[Snowflake]Guild),
		membersCache:             make(map[SnowflakePairKey]Member),
		channelsCache:            make(map[Snowflake]Channel),
		messagesCache:            make(map[Snowflake]Message),
		voiceStatesCache:         make(map[SnowflakePairKey]VoiceState),
		rolesCache:               make(map[Snowflake]Role),
		guildToMemberIDs:         make(map[Snowflake]map[Snowflake]struct{}),
		guildToChannelIDs:        make(map[Snowflake]map[Snowflake]struct{}),
		guildToVoiceStateUserIDs: make(map[Snowflake]map[Snowflake]struct{}),
		guildToRoleIDs:           make(map[Snowflake]map[Snowflake]struct{}),
	}
}

func (c *DefaultCache) Flags() CacheFlags {
	return c.flags
}

func (c *DefaultCache) SetFlags(flags ...CacheFlags) {
	c.flags = CacheFlagsNone
	for _, f := range flags {
		c.flags |= f
	}
}

func (c *DefaultCache) GetUser(userID Snowflake) (user User, ok bool) {
	c.usersCacheMu.RLock()
	user, ok = c.usersCache[userID]
	c.usersCacheMu.RUnlock()
	return
}

func (c *DefaultCache) GetGuild(guildID Snowflake) (guild Guild, ok bool) {
	c.guildsCacheMu.RLock()
	guild, ok = c.guildsCache[guildID]
	c.guildsCacheMu.RUnlock()
	return
}

func (c *DefaultCache) GetMember(guildID, userID Snowflake) (member Member, ok bool) {
	c.membersCacheMu.RLock()
	member, ok = c.membersCache[SnowflakePairKey{A: userID, B: guildID}]
	c.membersCacheMu.RUnlock()
	return
}

func (c *DefaultCache) GetChannel(channelID Snowflake) (channel Channel, ok bool) {
	c.channelsCacheMu.RLock()
	channel, ok = c.channelsCache[channelID]
	c.channelsCacheMu.RUnlock()
	return
}

func (c *DefaultCache) GetMessage(messageID Snowflake) (message Message, ok bool) {
	c.messagesCacheMu.RLock()
	message, ok = c.messagesCache[messageID]
	c.messagesCacheMu.RUnlock()
	return
}

func (c *DefaultCache) GetVoiceState(guildID, userID Snowflake) (voiceState VoiceState, ok bool) {
	c.voiceStatesCacheMu.RLock()
	voiceState, ok = c.voiceStatesCache[SnowflakePairKey{A: guildID, B: userID}]
	c.voiceStatesCacheMu.RUnlock()
	return
}

func (c *DefaultCache) GetGuildChannels(guildID Snowflake) (map[Snowflake]Channel, bool) {
	c.guildToChannelIDsMu.RLock()
	set, ok := c.guildToChannelIDs[guildID]
	c.guildToChannelIDsMu.RUnlock()
	if !ok {
		return nil, false
	}
	c.channelsCacheMu.RLock()
	defer c.channelsCacheMu.RUnlock()
	res := make(map[Snowflake]Channel, len(set))
	for channelID := range set {
		if channel, exists := c.channelsCache[channelID]; exists {
			res[channelID] = channel
		}
	}
	return res, true
}

func (c *DefaultCache) GetGuildMembers(guildID Snowflake) (map[Snowflake]Member, bool) {
	c.guildToMemberIDsMu.RLock()
	set, ok := c.guildToMemberIDs[guildID]
	c.guildToMemberIDsMu.RUnlock()
	if !ok {
		return nil, false
	}
	c.membersCacheMu.RLock()
	defer c.membersCacheMu.RUnlock()
	res := make(map[Snowflake]Member, len(set))
	for userID := range set {
		key := SnowflakePairKey{A: guildID, B: userID}
		if member, exists := c.membersCache[key]; exists {
			res[userID] = member
		}
	}
	return res, true
}

func (c *DefaultCache) GetGuildVoiceStates(guildID Snowflake) (map[Snowflake]VoiceState, bool) {
	c.guildToVoiceStateUserIDsMu.RLock()
	set, ok := c.guildToVoiceStateUserIDs[guildID]
	c.guildToVoiceStateUserIDsMu.RUnlock()
	if !ok {
		return nil, false
	}
	c.voiceStatesCacheMu.RLock()
	defer c.voiceStatesCacheMu.RUnlock()
	res := make(map[Snowflake]VoiceState, len(set))
	for userID := range set {
		key := SnowflakePairKey{A: guildID, B: userID}
		if voiceState, exists := c.voiceStatesCache[key]; exists {
			res[userID] = voiceState
		}
	}
	return res, true
}

func (c *DefaultCache) GetGuildRoles(guildID Snowflake) (map[Snowflake]Role, bool) {
	c.guildToRoleIDsMu.RLock()
	set, ok := c.guildToRoleIDs[guildID]
	c.guildToRoleIDsMu.RUnlock()
	if !ok {
		return nil, false
	}
	c.rolesCacheMu.RLock()
	defer c.rolesCacheMu.RUnlock()
	res := make(map[Snowflake]Role, len(set))
	for roleID := range set {
		if role, exists := c.rolesCache[roleID]; exists {
			res[roleID] = role
		}
	}
	return res, true
}

func (c *DefaultCache) CountUsers() int {
	c.usersCacheMu.RLock()
	count := len(c.usersCache)
	c.usersCacheMu.RUnlock()
	return count
}

func (c *DefaultCache) CountGuilds() int {
	c.guildsCacheMu.RLock()
	count := len(c.guildsCache)
	c.guildsCacheMu.RUnlock()
	return count
}

func (c *DefaultCache) CountMembers() int {
	c.membersCacheMu.RLock()
	count := len(c.membersCache)
	c.membersCacheMu.RUnlock()
	return count
}

func (c *DefaultCache) CountChannels() int {
	c.channelsCacheMu.RLock()
	count := len(c.channelsCache)
	c.channelsCacheMu.RUnlock()
	return count
}

func (c *DefaultCache) CountMessages() int {
	c.messagesCacheMu.RLock()
	count := len(c.messagesCache)
	c.messagesCacheMu.RUnlock()
	return count
}

func (c *DefaultCache) CountVoiceStates() int {
	c.voiceStatesCacheMu.RLock()
	count := len(c.voiceStatesCache)
	c.voiceStatesCacheMu.RUnlock()
	return count
}

func (c *DefaultCache) CountRoles() int {
	c.rolesCacheMu.RLock()
	count := len(c.rolesCache)
	c.rolesCacheMu.RUnlock()
	return count
}

func (c *DefaultCache) CountGuildChannels(guildID Snowflake) int {
	c.guildToChannelIDsMu.RLock()
	set, exists := c.guildToChannelIDs[guildID]
	c.guildToChannelIDsMu.RUnlock()
	if !exists {
		return 0
	}
	return len(set)
}

func (c *DefaultCache) CountGuildMembers(guildID Snowflake) int {
	c.guildToMemberIDsMu.RLock()
	set, exists := c.guildToMemberIDs[guildID]
	c.guildToMemberIDsMu.RUnlock()
	if !exists {
		return 0
	}
	return len(set)
}

func (c *DefaultCache) CountGuildRoles(guildID Snowflake) int {
	c.guildToRoleIDsMu.RLock()
	set, exists := c.guildToRoleIDs[guildID]
	c.guildToRoleIDsMu.RUnlock()
	if !exists {
		return 0
	}
	return len(set)
}

func (c *DefaultCache) PutUser(user User) {
	if !c.flags.Has(CacheFlagUsers) {
		return
	}
	c.usersCacheMu.Lock()
	c.usersCache[user.ID] = user
	c.usersCacheMu.Unlock()
}

func (c *DefaultCache) PutGuild(guild Guild) {
	if !c.flags.Has(CacheFlagGuilds) {
		return
	}
	c.guildsCacheMu.Lock()
	c.guildsCache[guild.ID] = guild
	c.guildsCacheMu.Unlock()
}

func (c *DefaultCache) PutMember(member Member) {
	if !c.flags.Has(CacheFlagMembers) {
		return
	}
	userID := member.User.ID
	guildID := member.GuildID
	key := SnowflakePairKey{A: guildID, B: userID}
	c.membersCacheMu.Lock()
	c.membersCache[key] = member
	c.membersCacheMu.Unlock()
	c.guildToMemberIDsMu.Lock()
	if _, exists := c.guildToMemberIDs[guildID]; !exists {
		c.guildToMemberIDs[guildID] = make(map[Snowflake]struct{})
	}
	c.guildToMemberIDs[guildID][userID] = struct{}{}
	c.guildToMemberIDsMu.Unlock()
}

func (c *DefaultCache) PutChannel(channel Channel) {
	if !c.flags.Has(CacheFlagChannels) {
		return
	}
	channelID := channel.GetID()
	c.channelsCacheMu.Lock()
	c.channelsCache[channelID] = channel
	c.channelsCacheMu.Unlock()
	if guildChannel, ok := channel.(GuildChannel); ok {
		guildID := guildChannel.GetGuildID()
		c.guildToChannelIDsMu.Lock()
		if _, exists := c.guildToChannelIDs[guildID]; !exists {
			c.guildToChannelIDs[guildID] = make(map[Snowflake]struct{})
		}
		c.guildToChannelIDs[guildID][channelID] = struct{}{}
		c.guildToChannelIDsMu.Unlock()
	}
}

func (c *DefaultCache) PutMessage(message Message) {
	if !c.flags.Has(CacheFlagMessages) {
		return
	}
	c.messagesCacheMu.Lock()
	c.messagesCache[message.ID] = message
	c.messagesCacheMu.Unlock()
}

func (c *DefaultCache) PutVoiceState(voiceState VoiceState) {
	if !c.flags.Has(CacheFlagVoiceStates) {
		return
	}
	guildID := voiceState.GuildID
	userID := voiceState.UserID
	key := SnowflakePairKey{A: guildID, B: userID}
	c.voiceStatesCacheMu.Lock()
	c.voiceStatesCache[key] = voiceState
	c.voiceStatesCacheMu.Unlock()
	c.guildToVoiceStateUserIDsMu.Lock()
	if _, exists := c.guildToVoiceStateUserIDs[guildID]; !exists {
		c.guildToVoiceStateUserIDs[guildID] = make(map[Snowflake]struct{})
	}
	c.guildToVoiceStateUserIDs[guildID][userID] = struct{}{}
	c.guildToVoiceStateUserIDsMu.Unlock()
}

func (c *DefaultCache) PutRole(role Role) {
	if !c.flags.Has(CacheFlagRoles) {
		return
	}
	guildID := role.GuildID
	roleID := role.ID
	c.rolesCacheMu.Lock()
	c.rolesCache[roleID] = role
	c.rolesCacheMu.Unlock()
	c.guildToRoleIDsMu.Lock()
	if _, exists := c.guildToRoleIDs[guildID]; !exists {
		c.guildToRoleIDs[guildID] = make(map[Snowflake]struct{})
	}
	c.guildToRoleIDs[guildID][roleID] = struct{}{}
	c.guildToRoleIDsMu.Unlock()
}

func (c *DefaultCache) DelUser(userID Snowflake) bool {
	c.usersCacheMu.Lock()
	_, ok := c.usersCache[userID]
	if ok {
		delete(c.usersCache, userID)
	}
	c.usersCacheMu.Unlock()
	return ok
}

func (c *DefaultCache) DelGuild(guildID Snowflake) bool {
	c.guildsCacheMu.Lock()
	_, ok := c.guildsCache[guildID]
	if ok {
		delete(c.guildsCache, guildID)
	}
	c.guildsCacheMu.Unlock()
	return ok
}

func (c *DefaultCache) DelMember(guildID, userID Snowflake) bool {
	key := SnowflakePairKey{A: guildID, B: userID}
	c.membersCacheMu.Lock()
	_, ok := c.membersCache[key]
	if ok {
		delete(c.membersCache, key)
	}
	c.membersCacheMu.Unlock()
	if ok {
		c.guildToMemberIDsMu.Lock()
		if m, has := c.guildToMemberIDs[guildID]; has {
			delete(m, userID)
			if len(m) == 0 {
				delete(c.guildToMemberIDs, guildID)
			}
		}
		c.guildToMemberIDsMu.Unlock()
	}
	return ok
}

func (c *DefaultCache) DelChannel(channelID Snowflake) bool {
	c.channelsCacheMu.Lock()
	channel, ok := c.channelsCache[channelID]
	if ok {
		delete(c.channelsCache, channelID)
	}
	c.channelsCacheMu.Unlock()
	if ok {
		if guildChannel, ok := channel.(GuildChannel); ok {
			c.guildToChannelIDsMu.Lock()
			if m, has := c.guildToChannelIDs[guildChannel.GetGuildID()]; has {
				delete(m, channelID)
				if len(m) == 0 {
					delete(c.guildToChannelIDs, guildChannel.GetGuildID())
				}
			}
			c.guildToChannelIDsMu.Unlock()
		}
	}
	return ok
}

func (c *DefaultCache) DelMessage(messageID Snowflake) bool {
	c.messagesCacheMu.Lock()
	_, ok := c.messagesCache[messageID]
	if ok {
		delete(c.messagesCache, messageID)
	}
	c.messagesCacheMu.Unlock()
	return ok
}

func (c *DefaultCache) DelVoiceState(guildID, userID Snowflake) bool {
	key := SnowflakePairKey{A: guildID, B: userID}
	c.voiceStatesCacheMu.Lock()
	_, ok := c.voiceStatesCache[key]
	if ok {
		delete(c.voiceStatesCache, key)
	}
	c.voiceStatesCacheMu.Unlock()
	if ok {
		c.guildToVoiceStateUserIDsMu.Lock()
		if m, has := c.guildToVoiceStateUserIDs[guildID]; has {
			delete(m, userID)
			if len(m) == 0 {
				delete(c.guildToVoiceStateUserIDs, guildID)
			}
		}
		c.guildToVoiceStateUserIDsMu.Unlock()
	}
	return ok
}

func (c *DefaultCache) DelRole(guildID, roleID Snowflake) bool {
	c.rolesCacheMu.Lock()
	_, ok := c.rolesCache[roleID]
	if ok {
		delete(c.rolesCache, roleID)
	}
	c.rolesCacheMu.Unlock()
	if ok {
		c.guildToRoleIDsMu.Lock()
		if m, has := c.guildToRoleIDs[guildID]; has {
			delete(m, roleID)
			if len(m) == 0 {
				delete(c.guildToRoleIDs, guildID)
			}
		}
		c.guildToRoleIDsMu.Unlock()
	}
	return ok
}

func (c *DefaultCache) DelGuildChannels(guildID Snowflake) bool {
	c.guildToChannelIDsMu.Lock()
	set, ok := c.guildToChannelIDs[guildID]
	if ok {
		delete(c.guildToChannelIDs, guildID)
	}
	c.guildToChannelIDsMu.Unlock()
	if ok {
		c.channelsCacheMu.Lock()
		for channelID := range set {
			delete(c.channelsCache, channelID)
		}
		c.channelsCacheMu.Unlock()
	}
	return ok
}

func (c *DefaultCache) DelGuildMembers(guildID Snowflake) bool {
	c.guildToMemberIDsMu.Lock()
	set, ok := c.guildToMemberIDs[guildID]
	if ok {
		delete(c.guildToMemberIDs, guildID)
	}
	c.guildToMemberIDsMu.Unlock()
	if ok {
		c.membersCacheMu.Lock()
		for userID := range set {
			key := SnowflakePairKey{A: guildID, B: userID}
			delete(c.membersCache, key)
		}
		c.membersCacheMu.Unlock()
	}
	return ok
}
