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
	"sync"

	"github.com/marouanesouiri/stdx/optional"
)

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

	GetUser(userID Snowflake) optional.Option[User]
	GetGuild(guildID Snowflake) optional.Option[Guild]
	GetMember(guildID, userID Snowflake) optional.Option[Member]
	GetChannel(channelID Snowflake) optional.Option[Channel]
	GetMessage(messageID Snowflake) optional.Option[Message]
	GetVoiceState(guildID, userID Snowflake) optional.Option[VoiceState]
	GetGuildChannels(guildID Snowflake) optional.Option[map[Snowflake]GuildChannel]
	GetGuildMembers(guildID Snowflake) optional.Option[map[Snowflake]Member]
	GetGuildVoiceStates(guildID Snowflake) optional.Option[map[Snowflake]VoiceState]
	GetGuildRoles(guildID Snowflake) optional.Option[map[Snowflake]Role]
	GetRoles(rolesIDs ...Snowflake) map[Snowflake]Role

	HasUser(userID Snowflake) bool
	HasGuild(guildID Snowflake) bool
	HasMember(guildID, userID Snowflake) bool
	HasChannel(channelID Snowflake) bool
	HasMessage(messageID Snowflake) bool
	HasVoiceState(guildID, userID Snowflake) bool
	HasGuildChannels(guildID Snowflake) bool
	HasGuildMembers(guildID Snowflake) bool
	HasGuildVoiceStates(guildID Snowflake) bool
	HasGuildRoles(guildID Snowflake) bool
	HasRoles(rolesIDs ...Snowflake) bool

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
	PutRoles(roles ...Role)

	DelUser(userID Snowflake) bool
	DelGuild(guildID Snowflake) bool
	DelMember(guildID, userID Snowflake) bool
	DelChannel(channelID Snowflake) bool
	DelMessage(messageID Snowflake) bool
	DelVoiceState(guildID, userID Snowflake) bool
	DelGuildChannels(guildID Snowflake) bool
	DelGuildMembers(guildID Snowflake) bool
	DelRole(guildID, roleID Snowflake) bool
	DelRoles(roleIDs ...Snowflake) bool
}

type InMemoryCacheManager struct {
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

func NewInMemoryCacheManager(flags CacheFlags) *InMemoryCacheManager {
	return &InMemoryCacheManager{
		flags: flags,

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

func (c *InMemoryCacheManager) Flags() CacheFlags {
	return c.flags
}

func (c *InMemoryCacheManager) SetFlags(flags ...CacheFlags) {
	c.flags = CacheFlagsNone
	for _, f := range flags {
		c.flags |= f
	}
}

func (c *InMemoryCacheManager) GetUser(userID Snowflake) optional.Option[User] {
	c.usersCacheMu.RLock()
	val, ok := c.usersCache[userID]
	c.usersCacheMu.RUnlock()
	return optional.FromPair(val, ok)
}

func (c *InMemoryCacheManager) GetGuild(guildID Snowflake) optional.Option[Guild] {
	c.guildsCacheMu.RLock()
	val, ok := c.guildsCache[guildID]
	c.guildsCacheMu.RUnlock()
	return optional.FromPair(val, ok)
}

func (c *InMemoryCacheManager) GetMember(guildID, userID Snowflake) optional.Option[Member] {
	c.membersCacheMu.RLock()
	val, ok := c.membersCache[SnowflakePairKey{A: guildID, B: userID}]
	c.membersCacheMu.RUnlock()
	return optional.FromPair(val, ok)
}

func (c *InMemoryCacheManager) GetChannel(channelID Snowflake) optional.Option[Channel] {
	c.channelsCacheMu.RLock()
	val, ok := c.channelsCache[channelID]
	c.channelsCacheMu.RUnlock()
	return optional.FromPair(val, ok)
}

func (c *InMemoryCacheManager) GetMessage(messageID Snowflake) optional.Option[Message] {
	c.messagesCacheMu.RLock()
	val, ok := c.messagesCache[messageID]
	c.messagesCacheMu.RUnlock()
	return optional.FromPair(val, ok)
}

func (c *InMemoryCacheManager) GetVoiceState(guildID, userID Snowflake) optional.Option[VoiceState] {
	c.voiceStatesCacheMu.RLock()
	val, ok := c.voiceStatesCache[SnowflakePairKey{A: guildID, B: userID}]
	c.voiceStatesCacheMu.RUnlock()
	return optional.FromPair(val, ok)
}

func (c *InMemoryCacheManager) GetGuildChannels(guildID Snowflake) optional.Option[map[Snowflake]GuildChannel] {
	c.guildToChannelIDsMu.RLock()
	set, ok := c.guildToChannelIDs[guildID]
	c.guildToChannelIDsMu.RUnlock()
	if !ok {
		return optional.None[map[Snowflake]GuildChannel]()
	}
	c.channelsCacheMu.RLock()
	defer c.channelsCacheMu.RUnlock()
	res := make(map[Snowflake]GuildChannel, len(set))
	for channelID := range set {
		if channel, exists := c.channelsCache[channelID]; exists {
			res[channelID] = channel.(GuildChannel)
		}
	}
	return optional.Some(res)
}

func (c *InMemoryCacheManager) GetGuildMembers(guildID Snowflake) optional.Option[map[Snowflake]Member] {
	c.guildToMemberIDsMu.RLock()
	set, ok := c.guildToMemberIDs[guildID]
	c.guildToMemberIDsMu.RUnlock()
	if !ok {
		return optional.None[map[Snowflake]Member]()
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
	return optional.Some(res)
}

func (c *InMemoryCacheManager) GetGuildVoiceStates(guildID Snowflake) optional.Option[map[Snowflake]VoiceState] {
	c.guildToVoiceStateUserIDsMu.RLock()
	set, ok := c.guildToVoiceStateUserIDs[guildID]
	c.guildToVoiceStateUserIDsMu.RUnlock()
	if !ok {
		return optional.None[map[Snowflake]VoiceState]()
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
	return optional.Some(res)
}

func (c *InMemoryCacheManager) GetGuildRoles(guildID Snowflake) optional.Option[map[Snowflake]Role] {
	c.guildToRoleIDsMu.RLock()
	set, ok := c.guildToRoleIDs[guildID]
	c.guildToRoleIDsMu.RUnlock()
	if !ok {
		return optional.None[map[Snowflake]Role]()
	}
	c.rolesCacheMu.RLock()
	defer c.rolesCacheMu.RUnlock()
	res := make(map[Snowflake]Role, len(set))
	for roleID := range set {
		if role, exists := c.rolesCache[roleID]; exists {
			res[roleID] = role
		}
	}
	return optional.Some(res)
}

func (c *InMemoryCacheManager) HasUser(userID Snowflake) bool {
	if !c.flags.Has(CacheFlagUsers) {
		return false
	}
	c.usersCacheMu.RLock()
	_, exists := c.usersCache[userID]
	c.usersCacheMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) HasGuild(guildID Snowflake) bool {
	if !c.flags.Has(CacheFlagGuilds) {
		return false
	}
	c.guildsCacheMu.RLock()
	_, exists := c.guildsCache[guildID]
	c.guildsCacheMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) HasMember(guildID, userID Snowflake) bool {
	if !c.flags.Has(CacheFlagMembers) {
		return false
	}
	c.membersCacheMu.RLock()
	_, exists := c.membersCache[SnowflakePairKey{A: guildID, B: userID}]
	c.membersCacheMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) HasChannel(channelID Snowflake) bool {
	if !c.flags.Has(CacheFlagChannels) {
		return false
	}
	c.channelsCacheMu.RLock()
	_, exists := c.channelsCache[channelID]
	c.channelsCacheMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) HasMessage(messageID Snowflake) bool {
	if !c.flags.Has(CacheFlagMessages) {
		return false
	}
	c.messagesCacheMu.RLock()
	_, exists := c.messagesCache[messageID]
	c.messagesCacheMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) HasVoiceState(guildID, userID Snowflake) bool {
	if !c.flags.Has(CacheFlagVoiceStates) {
		return false
	}
	c.voiceStatesCacheMu.RLock()
	_, exists := c.voiceStatesCache[SnowflakePairKey{A: guildID, B: userID}]
	c.voiceStatesCacheMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) HasGuildChannels(guildID Snowflake) bool {
	if !c.flags.Has(CacheFlagChannels) {
		return false
	}
	c.guildToChannelIDsMu.RLock()
	_, exists := c.guildToChannelIDs[guildID]
	c.guildToChannelIDsMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) HasGuildMembers(guildID Snowflake) bool {
	if !c.flags.Has(CacheFlagMembers) {
		return false
	}
	c.guildToMemberIDsMu.RLock()
	_, exists := c.guildToMemberIDs[guildID]
	c.guildToMemberIDsMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) HasGuildVoiceStates(guildID Snowflake) bool {
	if !c.flags.Has(CacheFlagVoiceStates) {
		return false
	}
	c.guildToVoiceStateUserIDsMu.RLock()
	_, exists := c.guildToVoiceStateUserIDs[guildID]
	c.guildToVoiceStateUserIDsMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) HasGuildRoles(guildID Snowflake) bool {
	if !c.flags.Has(CacheFlagRoles) {
		return false
	}
	c.guildToRoleIDsMu.RLock()
	_, exists := c.guildToRoleIDs[guildID]
	c.guildToRoleIDsMu.RUnlock()
	return exists
}

func (c *InMemoryCacheManager) CountUsers() int {
	c.usersCacheMu.RLock()
	count := len(c.usersCache)
	c.usersCacheMu.RUnlock()
	return count
}

func (c *InMemoryCacheManager) CountGuilds() int {
	c.guildsCacheMu.RLock()
	count := len(c.guildsCache)
	c.guildsCacheMu.RUnlock()
	return count
}

func (c *InMemoryCacheManager) CountMembers() int {
	c.membersCacheMu.RLock()
	count := len(c.membersCache)
	c.membersCacheMu.RUnlock()
	return count
}

func (c *InMemoryCacheManager) CountChannels() int {
	c.channelsCacheMu.RLock()
	count := len(c.channelsCache)
	c.channelsCacheMu.RUnlock()
	return count
}

func (c *InMemoryCacheManager) CountMessages() int {
	c.messagesCacheMu.RLock()
	count := len(c.messagesCache)
	c.messagesCacheMu.RUnlock()
	return count
}

func (c *InMemoryCacheManager) CountVoiceStates() int {
	c.voiceStatesCacheMu.RLock()
	count := len(c.voiceStatesCache)
	c.voiceStatesCacheMu.RUnlock()
	return count
}

func (c *InMemoryCacheManager) CountRoles() int {
	c.rolesCacheMu.RLock()
	count := len(c.rolesCache)
	c.rolesCacheMu.RUnlock()
	return count
}

func (c *InMemoryCacheManager) CountGuildChannels(guildID Snowflake) int {
	c.guildToChannelIDsMu.RLock()
	set, exists := c.guildToChannelIDs[guildID]
	c.guildToChannelIDsMu.RUnlock()
	if !exists {
		return 0
	}
	return len(set)
}

func (c *InMemoryCacheManager) CountGuildMembers(guildID Snowflake) int {
	c.guildToMemberIDsMu.RLock()
	set, exists := c.guildToMemberIDs[guildID]
	c.guildToMemberIDsMu.RUnlock()
	if !exists {
		return 0
	}
	return len(set)
}

func (c *InMemoryCacheManager) CountGuildRoles(guildID Snowflake) int {
	c.guildToRoleIDsMu.RLock()
	set, exists := c.guildToRoleIDs[guildID]
	c.guildToRoleIDsMu.RUnlock()
	if !exists {
		return 0
	}
	return len(set)
}

func (c *InMemoryCacheManager) PutUser(user User) {
	if !c.flags.Has(CacheFlagUsers) {
		return
	}
	c.usersCacheMu.Lock()
	c.usersCache[user.ID] = user
	c.usersCacheMu.Unlock()
}

func (c *InMemoryCacheManager) PutGuild(guild Guild) {
	if !c.flags.Has(CacheFlagGuilds) {
		return
	}
	c.guildsCacheMu.Lock()
	c.guildsCache[guild.ID] = guild
	c.guildsCacheMu.Unlock()
}

func (c *InMemoryCacheManager) PutMember(member Member) {
	if !c.flags.Has(CacheFlagMembers) {
		return
	}
	userID := member.ID
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

func (c *InMemoryCacheManager) PutChannel(channel Channel) {
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

func (c *InMemoryCacheManager) PutMessage(message Message) {
	if !c.flags.Has(CacheFlagMessages) {
		return
	}
	c.messagesCacheMu.Lock()
	c.messagesCache[message.ID] = message
	c.messagesCacheMu.Unlock()
}

func (c *InMemoryCacheManager) PutVoiceState(voiceState VoiceState) {
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

func (c *InMemoryCacheManager) PutRole(role Role) {
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

func (c *InMemoryCacheManager) DelUser(userID Snowflake) bool {
	c.usersCacheMu.Lock()
	_, ok := c.usersCache[userID]
	if ok {
		delete(c.usersCache, userID)
	}
	c.usersCacheMu.Unlock()
	return ok
}

func (c *InMemoryCacheManager) DelGuild(guildID Snowflake) bool {
	c.guildsCacheMu.Lock()
	_, ok := c.guildsCache[guildID]
	if ok {
		delete(c.guildsCache, guildID)
	}
	c.guildsCacheMu.Unlock()
	return ok
}

func (c *InMemoryCacheManager) DelMember(guildID, userID Snowflake) bool {
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

func (c *InMemoryCacheManager) DelChannel(channelID Snowflake) bool {
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

func (c *InMemoryCacheManager) DelMessage(messageID Snowflake) bool {
	c.messagesCacheMu.Lock()
	_, ok := c.messagesCache[messageID]
	if ok {
		delete(c.messagesCache, messageID)
	}
	c.messagesCacheMu.Unlock()
	return ok
}

func (c *InMemoryCacheManager) DelVoiceState(guildID, userID Snowflake) bool {
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

func (c *InMemoryCacheManager) DelRole(guildID, roleID Snowflake) bool {
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

func (c *InMemoryCacheManager) DelGuildChannels(guildID Snowflake) bool {
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

func (c *InMemoryCacheManager) DelGuildMembers(guildID Snowflake) bool {
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

func (c *InMemoryCacheManager) GetRoles(rolesIDs ...Snowflake) map[Snowflake]Role {
	c.rolesCacheMu.RLock()
	defer c.rolesCacheMu.RUnlock()
	res := make(map[Snowflake]Role, len(rolesIDs))
	for _, id := range rolesIDs {
		if role, ok := c.rolesCache[id]; ok {
			res[id] = role
		}
	}
	return res
}

func (c *InMemoryCacheManager) HasRoles(rolesIDs ...Snowflake) bool {
	if !c.flags.Has(CacheFlagRoles) {
		return false
	}
	c.rolesCacheMu.RLock()
	defer c.rolesCacheMu.RUnlock()
	for _, id := range rolesIDs {
		if _, ok := c.rolesCache[id]; !ok {
			return false
		}
	}
	return true
}

func (c *InMemoryCacheManager) PutRoles(roles ...Role) {
	for _, role := range roles {
		c.PutRole(role)
	}
}

func (c *InMemoryCacheManager) DelRoles(roleIDs ...Snowflake) bool {
	anyDeleted := false
	for _, id := range roleIDs {
		c.rolesCacheMu.RLock()
		role, ok := c.rolesCache[id]
		c.rolesCacheMu.RUnlock()

		if ok {
			if c.DelRole(role.GuildID, id) {
				anyDeleted = true
			}
		}
	}
	return anyDeleted
}
