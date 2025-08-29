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
	"errors"

	"github.com/bytedance/sonic"
)

type InteractionType int

const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
	InteractionTypeComponent
	InteractionTypeAutocomplete
	InteractionTypeModalSubmit
)

type InteractionContextType int

const (
	InteractionContextTypeGuild InteractionContextType = iota
	InteractionContextTypeBotDM
	InteractionContextTypePrivateChannel
)

type InteractionGuild struct {
	ID       Snowflake      `json:"id"`
	Locale   Locale         `json:"locale"`
	Features []GuildFeature `json:"features"`
}

type Interaction interface {
	GetID() Snowflake
	GetType() InteractionType
	GetApplicationID() Snowflake
	GetToken() string
	GetVersion() int
}

type InteractionFields struct {
	ID            Snowflake       `json:"id"`
	Type          InteractionType `json:"type"`
	ApplicationID Snowflake       `json:"application_id"`
	Token         string          `json:"token"`
	Version       int             `json:"version"`
}

func (i *InteractionFields) GetID() Snowflake {
	return i.ID
}

func (i *InteractionFields) GetType() InteractionType {
	return i.Type
}

func (i *InteractionFields) GetApplicationID() Snowflake {
	return i.ApplicationID
}

func (i *InteractionFields) GetToken() string {
	return i.Token
}

func (i *InteractionFields) GetVersion() int {
	return i.Version
}

type base2 struct {
	Entitlements                 []Entitlement                            `json:"entitlements"`
	AuthorizingIntegrationOwners map[ApplicationIntegrationType]Snowflake `json:"authorizing_integration_owners"`
	AttachmentSizeLimit          int                                      `json:"attachment_size_limit"`
}

type PingInteraction struct {
	InteractionFields
}

type ApplicationCommandInteraction struct {
	InteractionFields
	Guild                        *InteractionGuild
	Channel                      ResolvedChannel
	Locale                       Locale
	Member                       *ResolvedMember
	User                         *User
	AppPermissions               *Permissions
	Entitlements                 []Entitlement
	AuthorizingIntegrationOwners map[ApplicationIntegrationType]Snowflake
	Context                      InteractionContextType
	AttachmentSizeLimit          int
}

type ComponentInteraction struct {
	InteractionFields
}

type AutoCompleteInteraction struct {
	InteractionFields
}

type ModalSubmitInteraction struct {
	InteractionFields
}

// Helper func to Unmarshal any channel type to a Channel interface.
func UnmarshalInteraction(buf []byte) (Interaction, error) {
	var meta struct {
		Type InteractionType `json:"type"`
	}
	if err := sonic.Unmarshal(buf, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case InteractionTypePing:
		var i PingInteraction
		return &i, sonic.Unmarshal(buf, &i)
	case InteractionTypeApplicationCommand:
		var i ApplicationCommandInteraction
		return &i, sonic.Unmarshal(buf, &i)
	case InteractionTypeComponent:
		var i ComponentInteraction
		return &i, sonic.Unmarshal(buf, &i)
	case InteractionTypeAutocomplete:
		var i AutoCompleteInteraction
		return &i, sonic.Unmarshal(buf, &i)
	case InteractionTypeModalSubmit:
		var i ModalSubmitInteraction
		return &i, sonic.Unmarshal(buf, &i)
	default:
		return nil, errors.New("unknown interaction type")
	}
}
