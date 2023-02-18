package command

import (
	"errors"
	"github.com/Inoi-K/Find-Me/pkg/config"
	loc "github.com/Inoi-K/Find-Me/services/gateway/localization"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StartCommand       = "start"
	PingCommand        = "ping"
	EditProfileCommand = "edit"
	EditFieldButton    = "editButton"
	EditFieldCommand   = "editField"
	SignUpCommand      = "signup"

	HelpCommand     = "help"
	LanguageCommand = "lang"
	LanguageButton  = "langButton"

	// FIELDS
	Name        = "name"
	Gender      = "gender"
	Photo       = "photo"
	Description = "description"
)

var (
	EditProfileMarkup tgbotapi.InlineKeyboardMarkup

	UnknownCommandError = errors.New("unknown command")
	ContextDoneError    = errors.New("context is done")
)

func UpdateIndex() {
	// format is '<EditFieldButton><argumentSeparator><field>
	EditProfileMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Photo), EditFieldButton+config.C.Separator+Photo),
			tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Description), EditFieldButton+config.C.Separator+Description),
		),
	)
}
