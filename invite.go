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
	"strconv"
	"time"

	"github.com/marouanesouiri/stdx/optional"
	"github.com/marouanesouiri/stdx/result"
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
	Guild *PartialGuild `json:"guild"`

	// Channel is the channel this invite is for.
	Channel *PartialChannel `json:"channel"`

	// Inviter is the user who created the invite.
	Inviter *User `json:"inviter"`

	// TargetType is the type of target for this voice channel invite.
	TargetType optional.Option[InviteTargetType] `json:"target_type"`

	// TargetUser is the user whose stream to display for this voice channel stream invite.
	TargetUser *User `json:"target_user"`

	// TargetApplication is the embedded application to open for this voice channel embedded application invite.
	TargetApplication *PartialApplication `json:"target_application"`

	// ApproximatePresenceCount is the approximate count of online members.
	ApproximatePresenceCount optional.Option[int] `json:"approximate_presence_count"`

	// ApproximateMemberCount is the approximate count of total members.
	ApproximateMemberCount optional.Option[int] `json:"approximate_member_count"`

	// ExpiresAt is the expiration date of this invite.
	ExpiresAt optional.Option[time.Time] `json:"expires_at"`

	// GuildScheduledSvent guild scheduled event data, only included if GuildScheduledEventID contains a valid guild scheduled event id.
	GuildScheduledSvent *GuildScheduledEvent `json:"guild_scheduled_event"`

	// Flags is the guild invite flags for guild invites.
	Flags optional.Option[GuildInviteFlags] `json:"flags"`
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

// FetchInviteOptions contains parameters for fetching a invite.
type FetchInviteOptions struct {
	// WithCounts is whether the invite should contain approximate member counts.
	WithCounts bool `json:"with_counts"`

	// GuildScheduledEventID is the guild scheduled event to include with the invite.
	GuildScheduledEventID Snowflake `json:"guild_scheduled_event_id,omitempty"`
}

// FetchInvite retrieves a invite by its code.
func (r *requester) FetchInvite(code string, opts FetchInviteOptions) result.Result[Invite] {
	endpoint := "/invites/" + code + "?with_counts=" + strconv.FormatBool(opts.WithCounts)
	if opts.GuildScheduledEventID > 0 {
		endpoint += "&guild_scheduled_event_id=" + opts.GuildScheduledEventID.String()
	}
	res := r.DoRequest(Request{
		Method: "GET",
		URL:    endpoint,
	})
	if res.IsErr() {
		return result.Err[Invite](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var invite Invite
	if err := json.NewDecoder(body).Decode(&invite); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "GET",
			"url":    "/invites/{code}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Invite](err)
	}
	return result.Ok(invite)
}

// DeleteInviteOptions contains parameters for deleting invite.
type DeleteInviteOptions struct {
	// Reason is the reason shown in the audit log for this action.
	Reason string `json:"-"`
}

// DeleteInvite deletes an invite.
//
// Requires the PermissionManageChannels permission on the channel this invite belongs to, or PermissionManageGuild to remove any invite across the guild.
func (r *requester) DeleteInvite(code string, opts DeleteInviteOptions) result.Result[Invite] {
	res := r.DoRequest(Request{
		Method: "DELETE",
		URL:    "/invites/" + code,
		Reason: opts.Reason,
	})
	if res.IsErr() {
		return result.Err[Invite](res.Err())
	}
	body := res.Value()
	defer body.Close()

	var invite Invite
	if err := json.NewDecoder(body).Decode(&invite); err != nil {
		r.logger.WithFields(map[string]any{
			"method": "DELETE",
			"url":    "/invites/{code}",
			"error":  err.Error(),
		}).Error("failed parsing response")
		return result.Err[Invite](err)
	}
	return result.Ok(invite)
}
