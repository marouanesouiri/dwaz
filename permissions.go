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

import "strconv"

// Permissions flags for roles and channels permissions.
//
// Reference: https://discord.com/developers/docs/topics/permissions
type Permissions uint64

const (
	// CreateInstantInvite allows creation of instant invites.
	//
	// Channel types: Text, Voice, Stage
	PermissionCreateInstantInvite Permissions = 1 << 0

	// KickMembers allows kicking members.
	PermissionKickMembers Permissions = 1 << 1

	// BanMembers allows banning members.
	PermissionBanMembers Permissions = 1 << 2

	// Administrator allows all permissions and bypasses channel permission overwrites.
	PermissionAdministrator Permissions = 1 << 3

	// ManageChannels allows management and editing of channels.
	//
	// Channel types: Text, Voice, Stage
	PermissionManageChannels Permissions = 1 << 4

	// ManageGuild allows management and editing of the guild.
	PermissionManageGuild Permissions = 1 << 5

	// AddReactions allows adding new reactions to messages.
	// Does not apply to reacting with an existing reaction.
	//
	// Channel types: Text, Voice, Stage
	PermissionAddReactions Permissions = 1 << 6

	// ViewAuditLog allows viewing audit logs.
	PermissionViewAuditLog Permissions = 1 << 7

	// PrioritySpeaker allows using priority speaker in a voice channel.
	//
	// Channel types: Voice
	PermissionPrioritySpeaker Permissions = 1 << 8

	// Stream allows the user to go live.
	//
	// Channel types: Voice, Stage
	PermissionStream Permissions = 1 << 9

	// ViewChannel allows viewing a channel, reading messages, or joining voice channels.
	//
	// Channel types: Text, Voice, Stage
	PermissionViewChannel Permissions = 1 << 10

	// SendMessages allows sending messages and creating threads in forums.
	// Does not allow sending messages in threads.
	//
	// Channel types: Text, Voice, Stage
	PermissionSendMessages Permissions = 1 << 11

	// SendTTSMessages allows sending /tts messages.
	//
	// Channel types: Text, Voice, Stage
	PermissionSendTTSMessages Permissions = 1 << 12

	// ManageMessages allows deletion of other users' messages.
	//
	// Channel types: Text, Voice, Stage
	PermissionManageMessages Permissions = 1 << 13

	// EmbedLinks allows links to be auto-embedded.
	//
	// Channel types: Text, Voice, Stage
	PermissionEmbedLinks Permissions = 1 << 14

	// AttachFiles allows uploading images and files.
	//
	// Channel types: Text, Voice, Stage
	PermissionAttachFiles Permissions = 1 << 15

	// ReadMessageHistory allows reading message history.
	//
	// Channel types: Text, Voice, Stage
	PermissionReadMessageHistory Permissions = 1 << 16

	// MentionEveryone allows using @everyone and @here tags.
	//
	// Channel types: Text, Voice, Stage
	PermissionMentionEveryone Permissions = 1 << 17

	// UseExternalEmojis allows using custom emojis from other servers.
	//
	// Channel types: Text, Voice, Stage
	PermissionUseExternalEmojis Permissions = 1 << 18

	// ViewGuildInsights allows viewing guild insights.
	PermissionViewGuildInsights Permissions = 1 << 19

	// Connect allows joining a voice channel.
	//
	// Channel types: Voice, Stage
	PermissionConnect Permissions = 1 << 20

	// Speak allows speaking in a voice channel.
	//
	// Channel types: Voice
	PermissionSpeak Permissions = 1 << 21

	// MuteMembers allows muting members in a voice channel.
	//
	// Channel types: Voice, Stage
	PermissionMuteMembers Permissions = 1 << 22

	// DeafenMembers allows deafening members in a voice channel.
	//
	// Channel types: Voice
	PermissionDeafenMembers Permissions = 1 << 23

	// MoveMembers allows moving members between voice channels.
	//
	// Channel types: Voice, Stage
	PermissionMoveMembers Permissions = 1 << 24

	// UseVAD allows using voice activity detection in a voice channel.
	//
	// Channel types: Voice
	PermissionUseVAD Permissions = 1 << 25

	// ChangeNickname allows modification of own nickname.
	PermissionChangeNickname Permissions = 1 << 26

	// ManageNicknames allows modification of other users' nicknames.
	PermissionManageNicknames Permissions = 1 << 27

	// ManageRoles allows management and editing of roles.
	//
	// Channel types: Text, Voice, Stage
	PermissionManageRoles Permissions = 1 << 28

	// ManageWebhooks allows management and editing of webhooks.
	//
	// Channel types: Text, Voice, Stage
	PermissionManageWebhooks Permissions = 1 << 29

	// ManageGuildExpressions allows editing/deleting emojis, stickers, and soundboard sounds created by all users.
	PermissionManageGuildExpressions Permissions = 1 << 30

	// UseApplicationCommands allows using application (slash) commands.
	//
	// Channel types: Text, Voice, Stage
	PermissionUseApplicationCommands Permissions = 1 << 31

	// RequestToSpeak allows requesting to speak in stage channels.
	//
	// Channel types: Stage
	PermissionRequestToSpeak Permissions = 1 << 32
	// ManageEvents allows editing and deleting scheduled events created by all users.
	//
	// Channel types: Voice, Stage
	PermissionManageEvents Permissions = 1 << 33

	// ManageThreads allows deleting, archiving, and viewing all private threads.
	//
	// Channel types: Text
	PermissionManageThreads Permissions = 1 << 34

	// CreatePublicThreads allows creating public and announcement threads.
	//
	// Channel types: Text
	PermissionCreatePublicThreads Permissions = 1 << 35

	// CreatePrivateThreads allows creating private threads.
	//
	// Channel types: Text
	PermissionCreatePrivateThreads Permissions = 1 << 36

	// UseExternalStickers allows using custom stickers from other servers.
	//
	// Channel types: Text, Voice, Stage
	PermissionUseExternalStickers Permissions = 1 << 37

	// SendMessagesInThreads allows sending messages in threads.
	//
	// Channel types: Text
	PermissionSendMessagesInThreads Permissions = 1 << 38

	// UseEmbeddedActivities allows using Activities (applications with the EMBEDDED flag).
	//
	// Channel types: Text, Voice
	PermissionUseEmbeddedActivities Permissions = 1 << 39

	// ModerateMembers allows timing out users to prevent sending/reacting to messages or speaking.
	PermissionModerateMembers Permissions = 1 << 40

	// ViewCreatorMonetizationAnalytics allows viewing role subscription insights.
	PermissionViewCreatorMonetizationAnalytics Permissions = 1 << 41

	// UseSoundboard allows using soundboard in a voice channel.
	//
	// Channel types: Voice
	PermissionUseSoundboard Permissions = 1 << 42

	// CreateGuildExpressions allows creating emojis, stickers, and soundboard sounds, and editing/deleting those created by self.
	PermissionCreateGuildExpressions Permissions = 1 << 43

	// CreateEvents allows creating scheduled events, editing, and deleting those created by self.
	//
	// Channel types: Voice, Stage
	PermissionCreateEvents Permissions = 1 << 44

	// UseExternalSounds allows using custom soundboard sounds from other servers.
	//
	// Channel types: Voice
	PermissionUseExternalSounds Permissions = 1 << 45

	// SendVoiceMessages allows sending voice messages.
	//
	// Channel types: Text, Voice, Stage
	PermissionSendVoiceMessages Permissions = 1 << 46

	// SendPolls allows sending polls.
	//
	// Channel types: Text, Voice, Stage
	PermissionSendPolls Permissions = 1 << 49

	// UseExternalApps allows user-installed apps to send public responses.
	//
	// Channel types: Text, Voice, Stage
	PermissionUseExternalApps Permissions = 1 << 50
)

// PermissionName is a human-readable name for a Discord permission.
type PermissionName = string

const (
	PermissionNameCreateInstantInvite              PermissionName = "CreateInstantInvite"
	PermissionNameKickMembers                      PermissionName = "KickMembers"
	PermissionNameBanMembers                       PermissionName = "BanMembers"
	PermissionNameAdministrator                    PermissionName = "Administrator"
	PermissionNameManageChannels                   PermissionName = "ManageChannels"
	PermissionNameManageGuild                      PermissionName = "ManageGuild"
	PermissionNameAddReactions                     PermissionName = "AddReactions"
	PermissionNameViewAuditLog                     PermissionName = "ViewAuditLog"
	PermissionNamePrioritySpeaker                  PermissionName = "PrioritySpeaker"
	PermissionNameStream                           PermissionName = "Stream"
	PermissionNameViewChannel                      PermissionName = "ViewChannel"
	PermissionNameSendMessages                     PermissionName = "SendMessages"
	PermissionNameSendTTSMessages                  PermissionName = "SendTTSMessages"
	PermissionNameManageMessages                   PermissionName = "ManageMessages"
	PermissionNameEmbedLinks                       PermissionName = "EmbedLinks"
	PermissionNameAttachFiles                      PermissionName = "AttachFiles"
	PermissionNameReadMessageHistory               PermissionName = "ReadMessageHistory"
	PermissionNameMentionEveryone                  PermissionName = "MentionEveryone"
	PermissionNameUseExternalEmojis                PermissionName = "UseExternalEmojis"
	PermissionNameViewGuildInsights                PermissionName = "ViewGuildInsights"
	PermissionNameConnect                          PermissionName = "Connect"
	PermissionNameSpeak                            PermissionName = "Speak"
	PermissionNameMuteMembers                      PermissionName = "MuteMembers"
	PermissionNameDeafenMembers                    PermissionName = "DeafenMembers"
	PermissionNameMoveMembers                      PermissionName = "MoveMembers"
	PermissionNameUseVAD                           PermissionName = "UseVAD"
	PermissionNameChangeNickname                   PermissionName = "ChangeNickname"
	PermissionNameManageNicknames                  PermissionName = "ManageNicknames"
	PermissionNameManageRoles                      PermissionName = "ManageRoles"
	PermissionNameManageWebhooks                   PermissionName = "ManageWebhooks"
	PermissionNameManageGuildExpressions           PermissionName = "ManageGuildExpressions"
	PermissionNameUseApplicationCommands           PermissionName = "UseApplicationCommands"
	PermissionNameRequestToSpeak                   PermissionName = "RequestToSpeak"
	PermissionNameManageEvents                     PermissionName = "ManageEvents"
	PermissionNameManageThreads                    PermissionName = "ManageThreads"
	PermissionNameCreatePublicThreads              PermissionName = "CreatePublicThreads"
	PermissionNameCreatePrivateThreads             PermissionName = "CreatePrivateThreads"
	PermissionNameUseExternalStickers              PermissionName = "UseExternalStickers"
	PermissionNameSendMessagesInThreads            PermissionName = "SendMessagesInThreads"
	PermissionNameUseEmbeddedActivities            PermissionName = "UseEmbeddedActivities"
	PermissionNameModerateMembers                  PermissionName = "ModerateMembers"
	PermissionNameViewCreatorMonetizationAnalytics PermissionName = "ViewCreatorMonetizationAnalytics"
	PermissionNameUseSoundboard                    PermissionName = "UseSoundboard"
	PermissionNameCreateGuildExpressions           PermissionName = "CreateGuildExpressions"
	PermissionNameCreateEvents                     PermissionName = "CreateEvents"
	PermissionNameUseExternalSounds                PermissionName = "UseExternalSounds"
	PermissionNameSendVoiceMessages                PermissionName = "SendVoiceMessages"
	PermissionNameSendPolls                        PermissionName = "SendPolls"
	PermissionNameUseExternalApps                  PermissionName = "UseExternalApps"
)

// Has returns true if all given permissions are set.
func (p Permissions) Has(perms ...Permissions) bool {
	for _, perm := range perms {
		if p&perm != perm {
			return false
		}
	}
	return true
}

// Missing returns a Permissions bitmask containing the permissions
// that are present in the input perms but missing from p.
//
// Example:
//
//	p := PermissionSendMessages
//	missing := p.Missing(PermissionSendMessages, PermissionManageChannels)
//	// missing will contain PermissionManageChannels
func (p Permissions) Missing(perms ...Permissions) Permissions {
	var missing Permissions
	for _, perm := range perms {
		if p&perm == 0 {
			missing |= perm
		}
	}
	return missing
}

// Add sets all given permissions.
func (p *Permissions) Add(perms ...Permissions) {
	for _, perm := range perms {
		*p |= perm
	}
}

// Remove clears all given permissions.
func (p *Permissions) Remove(perms ...Permissions) {
	for _, perm := range perms {
		*p &^= perm
	}
}

// PermissionsToNames returns a slice of PermissionName for all permissions set in the mask.
func (p Permissions) Names() []PermissionName {
	var names []PermissionName
	for bit := Permissions(1); bit != 0; bit <<= 1 {
		if p&bit != 0 {
			switch bit {
			case PermissionCreateInstantInvite:
				names = append(names, PermissionNameCreateInstantInvite)
			case PermissionKickMembers:
				names = append(names, PermissionNameKickMembers)
			case PermissionBanMembers:
				names = append(names, PermissionNameBanMembers)
			case PermissionAdministrator:
				names = append(names, PermissionNameAdministrator)
			case PermissionManageChannels:
				names = append(names, PermissionNameManageChannels)
			case PermissionManageGuild:
				names = append(names, PermissionNameManageGuild)
			case PermissionAddReactions:
				names = append(names, PermissionNameAddReactions)
			case PermissionViewAuditLog:
				names = append(names, PermissionNameViewAuditLog)
			case PermissionPrioritySpeaker:
				names = append(names, PermissionNamePrioritySpeaker)
			case PermissionStream:
				names = append(names, PermissionNameStream)
			case PermissionViewChannel:
				names = append(names, PermissionNameViewChannel)
			case PermissionSendMessages:
				names = append(names, PermissionNameSendMessages)
			case PermissionSendTTSMessages:
				names = append(names, PermissionNameSendTTSMessages)
			case PermissionManageMessages:
				names = append(names, PermissionNameManageMessages)
			case PermissionEmbedLinks:
				names = append(names, PermissionNameEmbedLinks)
			case PermissionAttachFiles:
				names = append(names, PermissionNameAttachFiles)
			case PermissionReadMessageHistory:
				names = append(names, PermissionNameReadMessageHistory)
			case PermissionMentionEveryone:
				names = append(names, PermissionNameMentionEveryone)
			case PermissionUseExternalEmojis:
				names = append(names, PermissionNameUseExternalEmojis)
			case PermissionViewGuildInsights:
				names = append(names, PermissionNameViewGuildInsights)
			case PermissionConnect:
				names = append(names, PermissionNameConnect)
			case PermissionSpeak:
				names = append(names, PermissionNameSpeak)
			case PermissionMuteMembers:
				names = append(names, PermissionNameMuteMembers)
			case PermissionDeafenMembers:
				names = append(names, PermissionNameDeafenMembers)
			case PermissionMoveMembers:
				names = append(names, PermissionNameMoveMembers)
			case PermissionUseVAD:
				names = append(names, PermissionNameUseVAD)
			case PermissionChangeNickname:
				names = append(names, PermissionNameChangeNickname)
			case PermissionManageNicknames:
				names = append(names, PermissionNameManageNicknames)
			case PermissionManageRoles:
				names = append(names, PermissionNameManageRoles)
			case PermissionManageWebhooks:
				names = append(names, PermissionNameManageWebhooks)
			case PermissionManageGuildExpressions:
				names = append(names, PermissionNameManageGuildExpressions)
			case PermissionUseApplicationCommands:
				names = append(names, PermissionNameUseApplicationCommands)
			case PermissionRequestToSpeak:
				names = append(names, PermissionNameRequestToSpeak)
			case PermissionManageEvents:
				names = append(names, PermissionNameManageEvents)
			case PermissionManageThreads:
				names = append(names, PermissionNameManageThreads)
			case PermissionCreatePublicThreads:
				names = append(names, PermissionNameCreatePublicThreads)
			case PermissionCreatePrivateThreads:
				names = append(names, PermissionNameCreatePrivateThreads)
			case PermissionUseExternalStickers:
				names = append(names, PermissionNameUseExternalStickers)
			case PermissionSendMessagesInThreads:
				names = append(names, PermissionNameSendMessagesInThreads)
			case PermissionUseEmbeddedActivities:
				names = append(names, PermissionNameUseEmbeddedActivities)
			case PermissionModerateMembers:
				names = append(names, PermissionNameModerateMembers)
			case PermissionViewCreatorMonetizationAnalytics:
				names = append(names, PermissionNameViewCreatorMonetizationAnalytics)
			case PermissionUseSoundboard:
				names = append(names, PermissionNameUseSoundboard)
			case PermissionCreateGuildExpressions:
				names = append(names, PermissionNameCreateGuildExpressions)
			case PermissionCreateEvents:
				names = append(names, PermissionNameCreateEvents)
			case PermissionUseExternalSounds:
				names = append(names, PermissionNameUseExternalSounds)
			case PermissionSendVoiceMessages:
				names = append(names, PermissionNameSendVoiceMessages)
			case PermissionSendPolls:
				names = append(names, PermissionNameSendPolls)
			case PermissionUseExternalApps:
				names = append(names, PermissionNameUseExternalApps)
			}
		}
	}
	return names
}

// Method used internally by the library.
func (p *Permissions) UnmarshalJSON(data []byte) error {
	str, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}

	*p = Permissions(id)
	return nil
}

// Method used internally by the library.
func (p Permissions) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(p), 10) + `"`), nil
}
