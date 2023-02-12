package command

import (
	"context"
	loc "github.com/Inoi-K/Find-Me/services/gateway/localization"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ICommand provides an interface for all commands and buttons callbacks
type ICommand interface {
	Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error
}

// Start command begins an interaction with the chat and creates the record in database
type Start struct{}

func (c *Start) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	usr := upd.SentFrom()

	// send request to profile services to register a new user

	ok := loc.ChangeLanguage(usr.LanguageCode)
	// if user's language is not supported then set default language to english
	if !ok {
		loc.ChangeLanguage("en")
	}

	defer Reply(bot, chat, loc.Message(loc.Help))
	return Reply(bot, chat, loc.Message(loc.Start))
}
