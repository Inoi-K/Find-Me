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
)

var (
	UnknownCommandError = errors.New("unknown command")
)
