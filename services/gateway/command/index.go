package command

import "errors"

const (
	StartCommand    = "start"
	HelpCommand     = "help"
	LanguageCommand = "lang"
	PingCommand     = "ping"

	SignUpCommand = "signup"

	// CALLBACKS
	LanguageButton = "language"

	// STATES
	State        = "state"
	DefaultState = "default"
	NameState    = "name"
)

var (
	UnknownCommandError = errors.New("unknown command")
	UnknownStateError   = errors.New("unknown state")
)
