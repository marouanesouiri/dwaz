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
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
)

// InteractionType represents the type of an interaction in Discord.
// It indicates how the interaction was triggered.
//
// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-type
type InteractionType int

const (
	InteractionTypePing InteractionType = iota + 1
	// InteractionTypeApplicationCommand is triggered when a user invokes a slash command.
	InteractionTypeApplicationCommand

	// InteractionTypeComponent is triggered when a user interacts with a message component, like a button or select menu.
	InteractionTypeComponent

	// InteractionTypeAutocomplete is triggered when a user is typing in a command option that supports autocomplete.
	InteractionTypeAutocomplete

	// InteractionTypeModalSubmit is triggered when a user submits a modal dialog.
	InteractionTypeModalSubmit
)

// InteractionContextType is the context in Discord where an interaction can be used, or where it was triggered from.
// Details about using interaction contexts for application commands is in the commands context [documentation].
//
// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-interaction-context-types
//
// [documentation]: https://discord.com/developers/docs/interactions/application-commands#interaction-contexts
type InteractionContextType int

const (
	// Interaction can be used within servers
	InteractionContextTypeGuild InteractionContextType = iota

	// Interaction can be used within DMs with the app's bot user
	InteractionContextTypeBotDM

	// Interaction can be used within Group DMs and DMs other than the app's bot user
	InteractionContextTypePrivateChannel
)

// Is returns true if the interaction context's Type matches the provided one.
func (t InteractionContextType) Is(interactionType InteractionContextType) bool {
	return t == interactionType
}

// ApplicationCommandInteractionDataFields holds fields common to all application command interaction data.
type ApplicationCommandInteractionDataFields struct {
	// ID is the unique ID of the invoked command.
	ID Snowflake `json:"id"`

	// Type is the type of the invoked command.
	Type ApplicationCommandType `json:"type"`

	// Name is the name of the invoked command.
	Name string `json:"name"`

	// GuildID is the ID of the guild the command is registered to.
	//
	// Optional:
	//   - Will be 0 for global commands.
	GuildID Snowflake `json:"guild_id"`
}

// ChatInputCommandResolvedInteractionData represents the resolved data inside
// Interaction.Data.Resolved for chat input command interactions. This includes
// users, members, roles, channels, and attachments referenced by the command.
//
// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-resolved-data-structure
type ChatInputCommandResolvedInteractionData struct {
	// Users is a map of user IDs to User objects referenced by the command.
	Users map[Snowflake]User `json:"users"`

	// Members is a map of user IDs to partial Member objects for the guild.
	Members map[Snowflake]ResolvedMember `json:"members"`

	// Roles is a map of role IDs to Role objects referenced by the command.
	Roles map[Snowflake]Role `json:"roles"`

	// Channels is a map of channel IDs to partial Channel objects referenced by the command.
	Channels map[Snowflake]ResolvedChannel `json:"channels"`

	// Attachments is a map of attachment IDs to Attachment objects referenced by the command.
	Attachments map[Snowflake]Attachment `json:"attachments"`
}

// ChatInputInteractionCommandOption represents a single option provided
// by a user when invoking a chat input command (slash command).
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#interaction-object-application-command-interaction-data-option-object
type ChatInputInteractionCommandOption struct {
	// Name is the name of the option as defined in the command.
	Name string `json:"name"`

	// Type is the type of the option (string, integer, boolean, etc.).
	Type ApplicationCommandOptionType `json:"type"`

	// Value is the raw JSON value provided by the user for this option.
	//
	// Note: It's recommended to use the helper methods provided
	// (String, Int, Float, Bool, Snowflake) to extract the value.
	// These methods will panic if called on a type that doesn't match the
	// expected option type, so you must be sure of the type or check the
	// Type field before calling. For example, if Type is
	// ApplicationCommandOptionTypeString, calling String() is safe.
	Value json.RawMessage `json:"value"`
}

// String returns the option's value as a string.
//
// Panics if the underlying JSON value cannot be unmarshaled into a string.
// Make sure the option's Type is ApplicationCommandOptionTypeString before calling.
func (o *ChatInputInteractionCommandOption) String() string {
	var v string
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal option %q as string: %v", o.Name, err))
	}
	return v
}

// Int returns the option's value as an int.
//
// Panics if the underlying JSON value cannot be unmarshaled into an int.
// Make sure the option's Type is ApplicationCommandOptionTypeInteger before calling.
func (o *ChatInputInteractionCommandOption) Int() int {
	var v int
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal option %q as int: %v", o.Name, err))
	}
	return v
}

// Float returns the option's value as a float64.
//
// Panics if the underlying JSON value cannot be unmarshaled into a float64.
// Make sure the option's Type is ApplicationCommandOptionTypeNumber before calling.
func (o *ChatInputInteractionCommandOption) Float() float64 {
	var v float64
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal option %q as float64: %v", o.Name, err))
	}
	return v
}

// Bool returns the option's value as a boolean.
//
// Panics if the underlying JSON value cannot be unmarshaled into a bool.
// Make sure the option's Type is ApplicationCommandOptionTypeBoolean before calling.
func (o *ChatInputInteractionCommandOption) Bool() bool {
	var v bool
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal option %q as bool: %v", o.Name, err))
	}
	return v
}

// Snowflake returns the option's value as a Snowflake.
//
// Panics if the underlying JSON value cannot be unmarshaled into a Snowflake.
// Make sure the option's Type corresponds to an ID-based option before calling.
func (o *ChatInputInteractionCommandOption) Snowflake() Snowflake {
	var v Snowflake
	if err := json.Unmarshal(o.Value, &v); err != nil {
		panic(fmt.Sprintf("failed to unmarshal option %q as Snowflake: %v", o.Name, err))
	}
	return v
}

// ChatInputCommandInteractionData represents the data payload for a chat input
// command interaction.
//
// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-application-command-data-structure
type ChatInputCommandInteractionData struct {
	ApplicationCommandInteractionDataFields

	// Resolved contains the resolved objects referenced by this command.
	Resolved ChatInputCommandResolvedInteractionData `json:"resolved"`

	// Options contains the parameters and values provided by the user for this command.
	Options []ChatInputInteractionCommandOption `json:"options"`
}

// MessageCommandInteractionDataResolved represents the resolved data for a
// message command interaction.
//
// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-application-command-data-structure
type MessageCommandInteractionDataResolved struct {
	// Messages is a map of message IDs to partial Message objects
	// referenced by the command.
	Messages map[Snowflake]Message `json:"messages"`
}

// MessageCommandInteractionData represents the data for a message command interaction.
//
// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-application-command-data-structure
type MessageCommandInteractionData struct {
	ApplicationCommandInteractionDataFields

	// Resolved contains the resolved objects referenced by the command, messages in this case.
	Resolved MessageCommandInteractionDataResolved `json:"resolved"`

	// TargetID is the id of the message targeted by the command.
	TargetID Snowflake `json:"target_id"`
}

// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-application-command-data-structure
type UserCommandInteractionDataResolved struct {
	// Users is a map of user IDs to User objects referenced by the command.
	Users map[Snowflake]User `json:"users"`

	// Members is a map of user IDs to partial Member objects for the guild.
	Members map[Snowflake]ResolvedMember `json:"members"`
}

// UserCommandInteractionData represents the data for a user command interaction.
//
// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-application-command-data-structure
type UserCommandInteractionData struct {
	ApplicationCommandInteractionDataFields

	// Resolved contains the resolved objects referenced by the command, members and users in this case.
	Resolved UserCommandInteractionDataResolved `json:"resolved"`

	// TargetID is the id of the user targeted by the command.
	TargetID Snowflake `json:"target_id"`
}

// InteractionFields holds fields common to all Discord interactions.
type InteractionFields struct {
	// ID is the unique ID of the interaction.
	ID Snowflake `json:"id"`

	// Type is the type of the interaction.
	Type InteractionType `json:"type"`

	// ApplicationID is the ID of the application/bot this interaction belongs to.
	ApplicationID Snowflake `json:"application_id"`

	// Token is a token used to respond to the interaction.
	Token string `json:"token"`
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

// ApplicationCommandInteractionFields holds fields common to all application command interactions.
//
// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding
type ApplicationCommandInteractionFields struct {
	InteractionFields

	// Guild is the guild the interaction happened in.
	//
	// Optional:
	//   - Will be nil if the interaction occurred in a DM.
	Guild *PartialGuild `json:"guild"`

	// Channel is the channel the interaction happened in.
	Channel ResolvedMessageChannel `json:"channel"`

	// Locale is the selected language of the invoking user.
	Locale Locale `json:"locale"`

	// Member is the guild member data for the invoking user.
	//
	// Optional:
	//   - Present when the interaction is invoked in a guild.
	Member *ResolvedMember `json:"member"`

	// User is the user object for the invoking user, if invoked in a DM.
	//
	// Optional:
	//   - Present only when the interaction is invoked in a DM.
	User *User `json:"user"`

	// AppPermissions is a bitwise set of permissions the app has in the source location of the interaction.
	AppPermissions Permissions `json:"app_permissions"`

	// Entitlements is a list of entitlements for the invoking user.
	Entitlements []Entitlement `json:"entitlement"`

	// AuthorizingIntegrationOwners maps installation contexts that the interaction was authorized for
	// to related user or guild IDs.
	AuthorizingIntegrationOwners map[ApplicationIntegrationType]Snowflake `json:"authorizing_integration_owners"`

	// Context indicates where the interaction was triggered from.
	Context InteractionContextType `json:"context"`

	// AttachmentSizeLimit is the maximum size of attachments in bytes for this interaction.
	AttachmentSizeLimit int `json:"attachment_size_limit"`
}

// PingInteraction represents a Discord Ping interaction.
//
// Reference: https://discord.com/developers/docs/interactions/receiving-and-responding
type PingInteraction struct {
	InteractionFields
}

// ChatInputCommandInteraction represents an interaction triggered by a chat input (slash) command.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands
type ChatInputCommandInteraction struct {
	ApplicationCommandInteractionFields

	// Data contains the payload of the interaction specific to chat input commands.
	Data ChatInputCommandInteractionData `json:"data"`
}

// UserCommandInteraction represents an interaction triggered by a user context menu command.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands
type UserCommandInteraction struct {
	ApplicationCommandInteractionFields

	// Data contains the payload of the interaction specific to user commands.
	Data UserCommandInteractionData `json:"data"`
}

// MessageCommandInteraction represents an interaction triggered by a message context menu command.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands
type MessageCommandInteraction struct {
	ApplicationCommandInteractionFields

	// Data contains the payload of the interaction specific to message commands.
	Data MessageCommandInteractionData `json:"data"`
}

// TODO: continue the three interactions under this comment.
type ComponentInteraction struct {
	InteractionFields
}
type AutoCompleteInteraction struct {
	InteractionFields
}
type ModalSubmitInteraction struct {
	InteractionFields
}

// Interaction is the interface representing a Discord interaction.
//
// This interface can represent any type of interaction returned by Discord,
// including ping interactions, chat input commands, user commands, message commands,
// component interactions, autocomplete interactions, and modal submits.
//
// Use this interface when you want to handle interactions generically without knowing
// the specific concrete type in advance.
//
// You can convert (assert) it to a specific interaction type using a type assertion or
// a type switch, as described in the official Go documentation:
//   - https://go.dev/ref/spec#Type_assertions
//   - https://go.dev/doc/effective_go#type_switch
//
// Example usage:
//
//	var i Interaction
//
//	switch in := i.(type) {
//	case *PingInteraction:
//	    fmt.Println("Ping interaction ID:", in.ID)
//	case *ChatInputCommandInteraction:
//	    fmt.Println("Chat input command name:", in.Data.Name)
//	case *UserCommandInteraction:
//	    fmt.Println("User command target ID:", in.Data.TargetID)
//	default:
//	    fmt.Println("Other interaction type:", in.GetType())
//	}
//
// You can also use an if-condition to check a specific type:
//
//	if chatInputIn, ok := i.(*ChatInputCommandInteraction); ok {
//	    fmt.Println("Chat input command:", chatInputIntr.Data.Name)
//	}
type Interaction interface {
	GetID() Snowflake
	GetType() InteractionType
	GetApplicationID() Snowflake
	GetToken() string
}

var (
	_ Interaction = (*PingInteraction)(nil)
	_ Interaction = (*ChatInputCommandInteraction)(nil)
	_ Interaction = (*UserCommandInteraction)(nil)
	_ Interaction = (*MessageCommandInteraction)(nil)
	_ Interaction = (*ComponentInteraction)(nil)
	_ Interaction = (*AutoCompleteInteraction)(nil)
	_ Interaction = (*ModalSubmitInteraction)(nil)
)

// Helper func to Unmarshal any interaction type to a Interaction interface.
func UnmarshalInteraction(buf []byte) (Interaction, error) {
	var meta struct {
		Type InteractionType `json:"type"`
		Data struct {
			Type ApplicationCommandType `json:"type"`
		} `json:"data"`
	}
	if err := sonic.Unmarshal(buf, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case InteractionTypePing:
		var i PingInteraction
		return &i, sonic.Unmarshal(buf, &i)

	case InteractionTypeApplicationCommand:
		switch meta.Data.Type {
		case ApplicationCommandTypeChatInput:
			var i ChatInputCommandInteraction
			return &i, sonic.Unmarshal(buf, &i)
		case ApplicationCommandTypeUser:
			var i UserCommandInteraction
			return &i, sonic.Unmarshal(buf, &i)
		case ApplicationCommandTypeMessage:
			var i MessageCommandInteraction
			return &i, sonic.Unmarshal(buf, &i)
		default:
			return nil, errors.New("unknown application interacton type")
		}

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
