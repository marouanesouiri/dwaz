package yada

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	"github.com/bytedance/sonic"
)

// ApplicationCommandOptionType represents the type of an application command option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-type
type ApplicationCommandOptionType int

const (
	// ApplicationCommandOptionTypeSubCommand is a sub-command option, representing a nested command.
	ApplicationCommandOptionTypeSubCommand ApplicationCommandOptionType = iota + 1

	// ApplicationCommandOptionTypeSubCommandGroup is a group of sub-commands.
	ApplicationCommandOptionTypeSubCommandGroup

	// ApplicationCommandOptionTypeString is a string option.
	ApplicationCommandOptionTypeString

	// ApplicationCommandOptionTypeInteger is an integer option, supporting values between -2^53 and 2^53.
	ApplicationCommandOptionTypeInteger

	// ApplicationCommandOptionTypeBool is a boolean option.
	ApplicationCommandOptionTypeBool

	// ApplicationCommandOptionTypeUser is a user option, referencing a Discord user.
	ApplicationCommandOptionTypeUser

	// ApplicationCommandOptionTypeChannel is a channel option, including all channel types and categories.
	ApplicationCommandOptionTypeChannel

	// ApplicationCommandOptionTypeRole is a role option, referencing a Discord role.
	ApplicationCommandOptionTypeRole

	// ApplicationCommandOptionTypeMentionable is a mentionable option, including users and roles.
	ApplicationCommandOptionTypeMentionable

	// ApplicationCommandOptionTypeFloat is a float option, supporting any double between -2^53 and 2^53.
	ApplicationCommandOptionTypeFloat

	// ApplicationCommandOptionTypeAttachment is an attachment option, referencing an uploaded file.
	ApplicationCommandOptionTypeAttachment
)

// Is returns true if the option's Type matches the provided one.
func (t ApplicationCommandOptionType) Is(optionType ApplicationCommandOptionType) bool {
	return t == optionType
}

// ApplicationCommandOption is the interface representing a Discord application command option.
//
// This interface can represent any type of option returned by Discord, including sub-commands,
// sub-command groups, strings, integers, booleans, users, channels, roles, mentionables, floats,
// and attachments.
//
// Use this interface when you want to handle options generically without knowing the specific
// concrete type in advance.
//
// You can convert (assert) it to a specific option type using a type assertion or a type switch,
// as described in the official Go documentation:
//   - https://go.dev/ref/spec#Type_assertions
//   - https://go.dev/doc/effective_go#type_switch
//
// Example usage:
//
//	var myOption ApplicationCommandOption
//
//	switch opt := myOption.(type) {
//	case *ApplicationCommandOptionString:
//	    fmt.Println("String option:", opt.Name)
//	case *ApplicationCommandOptionInteger:
//	    fmt.Println("Integer option:", opt.Name)
//	case *ApplicationCommandOptionSubCommand:
//	    fmt.Println("Sub-command options:", opt.Options)
//	default:
//	    fmt.Println("Other option type:", opt.GetType())
//	}
type ApplicationCommandOption interface {
	GetType() ApplicationCommandOptionType
	GetName() string
	GetDescription() string
	json.Marshaler
}

// OptionBase contains fields common to all application command option types.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type OptionBase struct {
	// Type is the type of the option.
	Type ApplicationCommandOptionType `json:"type"`

	// Name is the name of the option.
	//
	// Info:
	//  - Must be 1-32 characters.
	//  - Must be unique within an array of application command options.
	Name string `json:"name"`

	// Description is the description of the option.
	//
	// Info:
	//  - Must be 1-100 characters.
	Description string `json:"description"`

	// NameLocalizations is a localization dictionary for the name field.
	//
	// Info:
	//  - Keys are available locales.
	//  - Values follow the same restrictions as Name (1-32 characters).
	NameLocalizations map[Locale]string `json:"name_localizations,omitempty"`

	// DescriptionLocalizations is a localization dictionary for the description field.
	//
	// Info:
	//  - Keys are available locales.
	//  - Values follow the same restrictions as Description (1-100 characters).
	DescriptionLocalizations map[Locale]string `json:"description_localizations,omitempty"`
}

func (o *OptionBase) GetType() ApplicationCommandOptionType { return o.Type }
func (o *OptionBase) GetName() string                       { return o.Name }
func (o *OptionBase) GetDescription() string                { return o.Description }

// RequiredBase contains the required field for value-based options.
type RequiredBase struct {
	// Required indicates whether the parameter is required or optional.
	//
	// Info:
	//  - Defaults to false.
	//  - Required options must be listed before optional options in an array of options.
	Required bool `json:"required,omitempty"`
}

// ChoiceBase contains the autocomplete field for choice-based options.
type ChoiceBase struct {
	// Autocomplete indicates whether autocomplete interactions are enabled for this option.
	//
	// Info:
	//  - May not be set to true if choices are present.
	//  - Options using autocomplete are not confined to only use choices given by the application.
	Autocomplete bool `json:"autocomplete,omitempty"`
}

// ChoiceOptionBase contains fields common to all choice option types.
type ChoiceOptionBase struct {
	// Name is the name of the choice.
	//
	// Info:
	//  - Must be 1-100 characters.
	Name string `json:"name"`

	// NameLocalizations is a localization dictionary for the choice name.
	//
	// Info:
	//  - Keys are available locales.
	//  - Values follow the same restrictions as Name (1-100 characters).
	NameLocalizations map[Locale]string `json:"name_localizations,omitempty"`
}

// ApplicationCommandOptionChoiceString represents a choice for string options.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoiceString struct {
	ChoiceOptionBase
	// Value is the string value of the choice.
	Value string `json:"value"`
}

// ChoiceFieldsString contains fields for string choice options.
type ChoiceFieldsString struct {
	// Choices is an array of choices for the user to pick from.
	//
	// Info:
	//  - Maximum of 25 choices.
	Choices []ApplicationCommandOptionChoiceString `json:"choices,omitempty"`
}

// StringConstraints contains constraints for string options.
type StringConstraints struct {
	// MinLength is the minimum allowed length for the string.
	//
	// Info:
	//  - Minimum of 0, maximum of 6000.
	//
	// Optional:
	//  - May be nil if no minimum length is specified.
	MinLength *int `json:"min_length,omitempty"`

	// MaxLength is the maximum allowed length for the string.
	//
	// Info:
	//  - Minimum of 1, maximum of 6000.
	//
	// Optional:
	//  - May be nil if no maximum length is specified.
	MaxLength *int `json:"max_length,omitempty"`
}

// ApplicationCommandOptionString represents a string option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionString struct {
	OptionBase
	RequiredBase
	ChoiceBase
	ChoiceFieldsString
	StringConstraints
}

func (o *ApplicationCommandOptionString) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ApplicationCommandOptionChoiceInteger represents a choice for integer options.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoiceInteger struct {
	ChoiceOptionBase
	// Value is the integer value of the choice.
	//
	// Info:
	//  - Must be between -2^53 and 2^53.
	Value int `json:"value"`
}

// ChoiceFieldsInteger contains fields for integer choice options.
type ChoiceFieldsInteger struct {
	// Choices is an array of choices for the user to pick from.
	//
	// Info:
	//  - Maximum of 25 choices.
	Choices []ApplicationCommandOptionChoiceInteger `json:"choices,omitempty"`
}

// IntegerConstraints contains constraints for integer options.
type IntegerConstraints struct {
	// MinValue is the minimum value permitted for the integer.
	//
	// Info:
	//  - Must be between -2^53 and 2^53.
	//
	// Optional:
	//  - May be nil if no minimum value is specified.
	MinValue *int `json:"min_value,omitempty"`

	// MaxValue is the maximum value permitted for the integer.
	//
	// Info:
	//  - Must be between -2^53 and 2^53.
	//
	// Optional:
	//  - May be nil if no maximum value is specified.
	MaxValue *int `json:"max_value,omitempty"`
}

// ApplicationCommandOptionInteger represents an integer option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionInteger struct {
	OptionBase
	RequiredBase
	ChoiceBase
	ChoiceFieldsInteger
	IntegerConstraints
}

func (o *ApplicationCommandOptionInteger) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ApplicationCommandOptionChoiceFloat represents a choice for float options.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-choice-structure
type ApplicationCommandOptionChoiceFloat struct {
	ChoiceOptionBase
	// Value is the float value of the choice.
	//
	// Info:
	//  - Must be between -2^53 and 2^53.
	Value float64 `json:"value"`
}

// ChoiceFieldsFloat contains fields for float choice options.
type ChoiceFieldsFloat struct {
	// Choices is an array of choices for the user to pick from.
	//
	// Info:
	//  - Maximum of 25 choices.
	Choices []ApplicationCommandOptionChoiceFloat `json:"choices,omitempty"`
}

// FloatConstraints contains constraints for float options.
type FloatConstraints struct {
	// MinValue is the minimum value permitted for the float.
	//
	// Info:
	//  - Must be between -2^53 and 2^53.
	//
	// Optional:
	//  - May be nil if no minimum value is specified.
	MinValue *float64 `json:"min_value,omitempty"`

	// MaxValue is the maximum value permitted for the float.
	//
	// Info:
	//  - Must be between -2^53 and 2^53.
	//
	// Optional:
	//  - May be nil if no maximum value is specified.
	MaxValue *float64 `json:"max_value,omitempty"`
}

// ApplicationCommandOptionFloat represents a float option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionFloat struct {
	OptionBase
	RequiredBase
	ChoiceBase
	ChoiceFieldsFloat
	FloatConstraints
}

func (o *ApplicationCommandOptionFloat) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ChannelConstraints contains constraints for channel options.
type ChannelConstraints struct {
	// ChannelTypes is an array of channel types that the option is restricted to.
	//
	// Info:
	//  - If not specified, includes all channel types and categories.
	ChannelTypes []ChannelType `json:"channel_types,omitempty"`
}

// ApplicationCommandOptionChannel represents a channel option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionChannel struct {
	OptionBase
	RequiredBase
	ChannelConstraints
}

func (o *ApplicationCommandOptionChannel) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ApplicationCommandOptionSubCommand represents a sub-command option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionSubCommand struct {
	OptionBase
	// Options is an array of nested options for the sub-command.
	//
	// Info:
	//  - Up to 25 options.
	//  - These are the parameters of the sub-command.
	Options []ApplicationCommandOption `json:"options,omitempty"`
}

func (o *ApplicationCommandOptionSubCommand) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ApplicationCommandOptionSubCommandGroup represents a sub-command group option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionSubCommandGroup struct {
	OptionBase
	// Options is an array of sub-commands for the group.
	//
	// Info:
	//  - Up to 25 sub-commands.
	Options []ApplicationCommandOptionSubCommand `json:"options,omitempty"`
}

func (o *ApplicationCommandOptionSubCommandGroup) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ApplicationCommandOptionBool represents a boolean option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionBool struct {
	OptionBase
	RequiredBase
}

func (o *ApplicationCommandOptionBool) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ApplicationCommandOptionUser represents a user option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionUser struct {
	OptionBase
	RequiredBase
}

func (o *ApplicationCommandOptionUser) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ApplicationCommandOptionRole represents a role option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionRole struct {
	OptionBase
	RequiredBase
}

func (o *ApplicationCommandOptionRole) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ApplicationCommandOptionMentionable represents a mentionable option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionMentionable struct {
	OptionBase
	RequiredBase
}

func (o *ApplicationCommandOptionMentionable) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

// ApplicationCommandOptionAttachment represents an attachment option.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-structure
type ApplicationCommandOptionAttachment struct {
	OptionBase
	RequiredBase
}

func (o *ApplicationCommandOptionAttachment) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(o)
}

func UnmarshalApplicationCommandOption(buf []byte) (ApplicationCommandOption, error) {
	var meta struct {
		Type ApplicationCommandOptionType `json:"type"`
	}
	if err := sonic.Unmarshal(buf, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case ApplicationCommandOptionTypeSubCommand:
		var o ApplicationCommandOptionSubCommand
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeSubCommandGroup:
		var o ApplicationCommandOptionSubCommand
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeString:
		var o ApplicationCommandOptionString
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeInteger:
		var o ApplicationCommandOptionInteger
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeBool:
		var o ApplicationCommandOptionBool
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeUser:
		var o ApplicationCommandOptionUser
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeChannel:
		var o ApplicationCommandOptionChannel
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeRole:
		var o ApplicationCommandOptionRole
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeMentionable:
		var o ApplicationCommandOptionMentionable
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeFloat:
		var o ApplicationCommandOptionFloat
		return &o, sonic.Unmarshal(buf, &o)
	case ApplicationCommandOptionTypeAttachment:
		var o ApplicationCommandOptionAttachment
		return &o, sonic.Unmarshal(buf, &o)
	default:
		return nil, errors.New("unknown application command option type")
	}
}

/*****************************
 *   Application Commands
 *****************************/

// ApplicationCommandType represents the type of an application command.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-types
type ApplicationCommandType int

const (
	// ApplicationCommandTypeChatInput is a text-based slash command that shows up when a user types "/".
	ApplicationCommandTypeChatInput ApplicationCommandType = 1 + iota

	// ApplicationCommandTypeUser is a UI-based command that shows up when you right-click or tap on a user.
	ApplicationCommandTypeUser

	// ApplicationCommandTypeMessage is a UI-based command that shows up when you right-click or tap on a message.
	ApplicationCommandTypeMessage

	// ApplicationCommandTypePrimaryEntryPoint is a UI-based command that represents the primary way to invoke an app's Activity.
	ApplicationCommandTypePrimaryEntryPoint
)

// Is returns true if the command's Type matches the provided one.
func (t ApplicationCommandType) Is(commandType ApplicationCommandType) bool {
	return t == commandType
}

// ApplicationCommandHandlerType represents the handler type for PRIMARY_ENTRY_POINT commands.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object-entry-point-command-handler-types
type ApplicationCommandHandlerType int

const (
	// ApplicationCommandHandlerTypeApp indicates the app handles the interaction using an interaction token.
	ApplicationCommandHandlerTypeApp ApplicationCommandHandlerType = 1 + iota

	// ApplicationCommandHandlerTypeDiscord indicates Discord handles the interaction by launching an Activity and sending a follow-up message.
	ApplicationCommandHandlerTypeDiscord
)

// Is returns true if the command handler's Type matches the provided one.
func (t ApplicationCommandHandlerType) Is(handlerType ApplicationCommandHandlerType) bool {
	return t == handlerType
}

// ApplicationCommand is the interface representing a Discord application command.
//
// This interface can represent any type of command returned by Discord, including slash commands,
// user commands, message commands, and primary entry point commands.
//
// Use this interface when you want to handle commands generically without knowing the specific
// concrete type in advance.
//
// You can convert (assert) it to a specific command type using a type assertion or a type switch,
// as described in the official Go documentation:
//   - https://go.dev/ref/spec#Type_assertions
//   - https://go.dev/doc/effective_go#type_switch
//
// Example usage:
//
//	var myCommand ApplicationCommand
//
//	switch cmd := myCommand.(type) {
//	case *ChatInputCommand:
//	    fmt.Println("Slash command:", cmd.Name)
//	case *ApplicationUserCommand:
//	    fmt.Println("User command:", cmd.Name)
//	case *ApplicationMessageCommand:
//	    fmt.Println("Message command:", cmd.Name)
//	default:
//	    fmt.Println("Other command type:", cmd.GetType())
//	}
type ApplicationCommand interface {
	GetType() ApplicationCommandType
	GetID() Snowflake
	GetApplicationID() Snowflake
	GetGuildID() Snowflake
	GetName() string
	GetNameLocalizations() map[Locale]string
	GetDefaultMemberPermissions() Permissions
	GetVersion() Snowflake
	CreatedAt() time.Time
	IsNSFW() bool
	GetIntegrationTypes() []ApplicationIntegrationType
	GetContexts() []InteractionContextType
	json.Marshaler
}

// ApplicationCommandBase contains fields common to all application command types.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object
type ApplicationCommandBase struct {
	// Type is the type of the command.
	Type ApplicationCommandType `json:"type"`

	// ID is the unique ID of the command.
	ID Snowflake `json:"id"`

	// ApplicationID is the ID of the parent application.
	ApplicationID Snowflake `json:"application_id"`

	// GuildID is the guild ID of the command, if not global.
	GuildID Snowflake `json:"guild_id"`

	// Name is the name of the command.
	//
	// Info:
	//  - Must be 1-32 characters.
	//  - For CHAT_INPUT commands, must match the regex ^[-_'\p{L}\p{N}\p{sc=Deva}\p{sc=Thai}]{1,32}$ with unicode flag.
	//  - For USER and MESSAGE commands, may be mixed case and include spaces.
	//  - Must use lowercase variants of letters when available.
	Name string `json:"name"`

	// NameLocalizations is a localization dictionary for the name field.
	//
	// Info:
	//  - Keys are available locales.
	//  - Values follow the same restrictions as Name.
	NameLocalizations map[Locale]string `json:"name_localizations"`

	// DefaultMemberPermissions is the set of permissions required to use the command.
	//
	// Info:
	//  - Represented as a bit set.
	//  - Set to "0" to disable the command for everyone except admins by default.
	DefaultMemberPermissions Permissions `json:"default_member_permissions"`

	// DefaultPermission indicates whether the command is enabled by default when the app is added to a guild.
	//
	// Info:
	//  - Defaults to true.
	//  - Deprecated; use DefaultMemberPermissions or Contexts instead.
	DefaultPermission bool `json:"default_permission"`

	// NSFW indicates whether the command is age-restricted.
	//
	// Info:
	//  - Defaults to false.
	NSFW bool `json:"nsfw"`

	// IntegrationTypes is the list of installation contexts where the command is available.
	//
	// Info:
	//  - Only for globally-scoped commands.
	//  - Defaults to the app's configured contexts.
	IntegrationTypes []ApplicationIntegrationType `json:"integration_types"`

	// Contexts is the list of interaction contexts where the command can be used.
	//
	// Info:
	//  - Only for globally-scoped commands.
	Contexts []InteractionContextType `json:"contexts"`

	// Version is the autoincrementing version identifier updated during substantial record changes.
	Version Snowflake `json:"version"`
}

func (a *ApplicationCommandBase) GetType() ApplicationCommandType {
	return a.Type
}

func (a *ApplicationCommandBase) GetID() Snowflake {
	return a.ID
}

func (a *ApplicationCommandBase) GetApplicationID() Snowflake {
	return a.ApplicationID
}

// GetGuildID returns the guild ID of the command, if it is guild-specific.
//
// Returns:
//  - The Snowflake ID of the guild if the command is associated with a guild.
//  - 0 if the command is global (not tied to a specific guild).
func (a *ApplicationCommandBase) GetGuildID() Snowflake {
	return a.GuildID
}

func (a *ApplicationCommandBase) GetName() string {
	return a.Name
}

func (a *ApplicationCommandBase) GetNameLocalizations() map[Locale]string {
	return a.NameLocalizations
}

func (a *ApplicationCommandBase) GetDefaultMemberPermissions() Permissions {
	return a.DefaultMemberPermissions
}

func (a *ApplicationCommandBase) GetVersion() Snowflake {
	return a.Version
}

func (a *ApplicationCommandBase) CreatedAt() time.Time {
	return a.ID.Timestamp()
}

func (a *ApplicationCommandBase) IsNSFW() bool {
	return a.NSFW
}

func (a *ApplicationCommandBase) GetIntegrationTypes() []ApplicationIntegrationType {
	return a.IntegrationTypes
}

func (a *ApplicationCommandBase) GetContexts() []InteractionContextType {
	return a.Contexts
}

// DescriptionConstraints contains description fields for application commands.
type DescriptionConstraints struct {
	// Description is the description of the command.
	//
	// Info:
	//  - Must be 1-100 characters.
	Description string `json:"description"`

	// DescriptionLocalizations is a localization dictionary for the description field.
	//
	// Info:
	//  - Keys are available locales.
	//  - Values follow the same restrictions as Description (1-100 characters).
	DescriptionLocalizations map[Locale]string `json:"description_localizations"`
}

// ChatInputCommand represents a slash command.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object
type ChatInputCommand struct {
	ApplicationCommandBase
	DescriptionConstraints
	// Options is an array of parameters for the command.
	//
	// Info:
	//  - Maximum of 25 options.
	Options []ApplicationCommandOption `json:"options"`
}

var _ json.Unmarshaler = (*ChatInputCommand)(nil)

func (c *ChatInputCommand) UnmarshalJSON(buf []byte) error {
	type tempChatInputCommand struct {
		ApplicationCommandBase
		DescriptionConstraints
		Options []json.RawMessage `json:"options"`
	}

	var temp tempChatInputCommand
	if err := sonic.Unmarshal(buf, &temp); err != nil {
		return err
	}

	c.ApplicationCommandBase = temp.ApplicationCommandBase
	c.DescriptionConstraints = temp.DescriptionConstraints


	if temp.Options != nil {
		c.Options = make([]ApplicationCommandOption, 0, len(temp.Options))
		for i := range len(temp.Options) {
			if len(temp.Options[i]) == 0 || bytes.Equal(temp.Options[i], []byte("null")) {
				continue
			}
			option, err := UnmarshalApplicationCommandOption(temp.Options[i])
			if err != nil {
				return err
			}
			c.Options[i] = option.(ApplicationCommandOption)
		}
	}

	return nil
}

func (c *ChatInputCommand) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(c)
}

// ApplicationUserCommand represents a UI-based command that appears when right-clicking or tapping on a user.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object
type ApplicationUserCommand struct {
	ApplicationCommandBase
}

func (c *ApplicationUserCommand) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(c)
}

// ApplicationMessageCommand represents a UI-based command that appears when right-clicking or tapping on a message.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object
type ApplicationMessageCommand struct {
	ApplicationCommandBase
}

func (c *ApplicationMessageCommand) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(c)
}

// ApplicationEntryPointCommand represents a UI-based command that is the primary way to invoke an app's Activity.
//
// Reference: https://discord.com/developers/docs/interactions/application-commands#application-command-object
type ApplicationEntryPointCommand struct {
	ApplicationCommandBase
	DescriptionConstraints
	// Handler determines whether the interaction is handled by the app's interactions handler or by Discord.
	//
	// Info:
	//  - Only applicable for commands with the EMBEDDED flag (i.e., applications with an Activity).
	Handler ApplicationCommandHandlerType `json:"handler"`
}

func (c *ApplicationEntryPointCommand) MarshalJSON() ([]byte, error) {
	return sonic.Marshal(c)
}

func UnmarshalApplicationCommand(buf []byte) (ApplicationCommand, error) {
	var meta struct {
		Type ApplicationCommandType `json:"type"`
	}
	if err := sonic.Unmarshal(buf, &meta); err != nil {
		return nil, err
	}

	switch meta.Type {
	case ApplicationCommandTypeChatInput:
		var c ChatInputCommand
		return &c, sonic.Unmarshal(buf, &c)
	case ApplicationCommandTypeUser:
		var c ApplicationUserCommand
		return &c, sonic.Unmarshal(buf, &c)
	case ApplicationCommandTypeMessage:
		var c ApplicationMessageCommand
		return &c, sonic.Unmarshal(buf, &c)
	case ApplicationCommandTypePrimaryEntryPoint:
		var c ApplicationEntryPointCommand
		return &c, sonic.Unmarshal(buf, &c)
	default:
		return nil, errors.New("unknown application command type")
	}
}
