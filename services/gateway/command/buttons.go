package command

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/services/gateway/client"
	loc "github.com/Inoi-K/Find-Me/services/gateway/localization"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

// EditFieldCallback asks a user for a new value of a field and edits it
// callbacks come from EditMenu
type EditFieldCallback struct{}

func (c *EditFieldCallback) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	user := upd.SentFrom()
	chat := upd.FromChat()

	// get an argument
	newArg, err := askStateField(ctx, bot, chat, user.ID, args)
	if err != nil {
		return err
	}
	// format arguments
	editArgs := args + config.C.Separator + newArg
	// run edit field command
	return (&EditField{}).Execute(ctx, bot, upd, editArgs)
}

// EditField edits a field in additional information of user.
// args format is `<field><argumentSeparator><value>`
type EditField struct{}

func (c *EditField) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	user := upd.SentFrom()
	chat := upd.FromChat()

	// parse the formatted arguments
	field, value, _ := strings.Cut(args, config.C.Separator)
	editRequest := &pb.EditRequest{
		UserID:   user.ID,
		SphereID: config.C.SphereID,
		Field:    field,
		// FIXME description cuts down to a first word only
		Value: strings.Split(value, config.C.Separator),
	}

	// contact the profile server for editing
	ctx2, cancel := context.WithTimeout(ctx, config.C.Timeout)
	defer cancel()
	_, err := client.Profile.Edit(ctx2, editRequest)
	if err != nil {
		log.Printf("couldn't edit field %s: %v", field, err)
		_ = reply(bot, chat, loc.Message(loc.EditFail))
		return err
	}

	return reply(bot, chat, loc.Message(loc.EditSuccess))
}

// LanguageCallback changes language
type LanguageCallback struct{}

func (c *LanguageCallback) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	result := loc.Message(loc.LangFail) // fail in current language
	if loc.ChangeLanguage(args) {
		result = loc.Message(loc.LangSuccess) // success in new language
	}

	return reply(bot, chat, result)
}
