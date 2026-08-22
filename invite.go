/************************************************************************************
 *
 * wumpkit, A Lightweight Go library for Discord API
 *
 * SPDX-License-Identifier: BSD-3-Clause
 *
 * Copyright 2025 Marouane Souiri
 *
 * Licensed under the BSD 3-Clause License.
 * See the LICENSE file for details.
 *
 ************************************************************************************/

package wumpkit

import (
	"time"
)

// InviteType represents invite types.
//
// Reference: https://discord.com/developers/docs/resources/invite#invite-object-invite-types
type InviteType int

const (
	InviteTypeGuild InviteType = iota
	InviteTypeGroupDM
	InviteTypeFriend
)

// Is returns true if the invite's Type matches the provided one.
func (t InviteType) Is(inviteType InviteType) bool {
	return t == inviteType
}

// InviteTargetType represents invite target types.
//
// Reference: https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
type InviteTargetType int

const (
	InviteTargetTypeStream              InviteTargetType = 1
	InviteTargetTypeEmbeddedApplication InviteTargetType = 2
)

// Is returns true if the invite target's Type matches the provided one.
func (t InviteTargetType) Is(inviteType InviteTargetType) bool {
	return t == inviteType
}

// GuildInviteFlags represents invite flags.
//
// Reference: https://discord.com/developers/docs/resources/invite#invite-object-invite-target-types
type GuildInviteFlags int

const (
	// GuildInviteFlagsIsGuestInvite indicates that this invite is a guest invite for a voice channel
	GuildInviteFlagsIsGuestInvite GuildInviteFlags = 1 << 0
)

// Has returns true if all provided flags are set.
func (f GuildInviteFlags) Has(flags ...GuildInviteFlags) bool {
	return BitFieldHas(f, flags...)
}

// Invite represents a invite object.
//
// Reference: https://discord.com/developers/docs/resources/invite#invite-object
type Invite struct {
	// Type is the type of invite.
	Type InviteType `json:"type"`

	// Code is the invite code (unique ID).
	Code string `json:"code"`

	// Guild is the guild this invite is for.
	Guild *PartialGuild `json:"guild,omitempty"`

	// Channel is the channel this invite is for.
	Channel *PartialChannel `json:"channel,omitempty"`

	// Inviter is the user who created the invite.
	Inviter *User `json:"inviter,omitempty"`

	// TargetType is the type of target for this voice channel invite.
	TargetType *InviteTargetType `json:"target_type,omitempty"`

	// TargetUser is the user whose stream to display for this voice channel stream invite.
	TargetUser *User `json:"target_user,omitempty"`

	// TargetApplication is the embedded application to open for this voice channel embedded application invite.
	TargetApplication *PartialApplication `json:"target_application,omitempty"`

	// ApproximatePresenceCount is the approximate count of online members.
	ApproximatePresenceCount *int `json:"approximate_presence_count,omitempty"`

	// ApproximateMemberCount is the approximate count of total members.
	ApproximateMemberCount *int `json:"approximate_member_count,omitempty"`

	// ExpiresAt is the expiration date of this invite.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// GuildScheduledEvent guild scheduled event data, only included if GuildScheduledEventID contains a valid guild scheduled event id.
	GuildScheduledEvent *GuildScheduledEvent `json:"guild_scheduled_event,omitempty"`

	// Flags is the guild invite flags for guild invites.
	Flags *GuildInviteFlags `json:"flags,omitempty"`
}

// InviteMetadata represents extra information about an invite, will extend the invite object.
//
// Reference: https://discord.com/developers/docs/resources/invite#invite-metadata-object
type InviteMetadata struct {
	// Uses is the number of times this invite has been used.
	Uses int `json:"uses"`

	// MaxUses is the max number of times this invite can be used.
	MaxUses int `json:"max_uses"`

	// MaxAge is the duration (in seconds) after which the invite expires.
	MaxAge int `json:"max_age"`

	// Temporary is whether this invite only grants temporary membership.
	Temporary bool `json:"temporary"`

	// CreatedAt is when this invite was created.
	CreatedAt time.Time `json:"created_at"`
}

type FullInvite struct {
	Invite
	InviteMetadata
}
