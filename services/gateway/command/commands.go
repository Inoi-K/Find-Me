package command

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/services/gateway/client"
	loc "github.com/Inoi-K/Find-Me/services/gateway/localization"
	"github.com/Inoi-K/Find-Me/services/gateway/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"time"
)

// ICommand provides an interface for all commands and buttons callbacks
type ICommand interface {
	Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error
}

// Start command begins an interaction with the chat and creates the record in database
type Start struct{}

func (c *Start) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	user := upd.SentFrom()

	ok := loc.ChangeLanguage(user.LanguageCode)
	// if user's language is not supported then set default language to english
	if !ok {
		loc.ChangeLanguage("en")
	}

	// check user existence
	// contact the profile service
	ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rep, err := client.Profile.Exists(ctx2, &pb.ExistsRequest{
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("couldn't check existance of user with id = %d : %v", user.ID, err)
		return reply(bot, chat, loc.Message(loc.TryAgain))
	}
	// break execution if user already exists
	if rep.Exists {
		return reply(bot, chat, loc.Message(loc.AlreadyRegistered))
	}

	// TODO validate terms & agreement

	// TODO validate user with corporate email

	// main information
	var signUpArgs, newArg string
	session.UserStateArg[user.ID] = make(chan string)
	// name
	session.UserState[user.ID] = session.EnterName
	newArg, err = askNewArg(ctx, bot, chat, user.ID, loc.EnterName)
	if err != nil {
		return err
	}
	signUpArgs += newArg
	// gender
	session.UserState[user.ID]++
	newArg, err = askNewArg(ctx, bot, chat, user.ID, loc.EnterGender)
	if err != nil {
		return err
	}
	signUpArgs += config.C.ArgumentsSeparator + newArg

	// clear user state
	close(session.UserStateArg[user.ID])
	delete(session.UserState, user.ID)

	signUp := &SignUp{}
	return signUp.Execute(ctx, bot, upd, signUpArgs)
}

func askNewArg(ctx context.Context, bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, userID int64, messageKey string) (string, error) {
	err := reply(bot, chat, loc.Message(messageKey))
	if err != nil {
		return "", err
	}
	newArg := ""
	select {
	case <-ctx.Done():
		return "", ContextDoneError
	case newArg = <-session.UserStateArg[userID]:
	}

	return newArg, nil
}

// SignUp sends a request to the profile service to register a new user
type SignUp struct{}

func (c *SignUp) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	user := upd.SentFrom()
	chat := upd.FromChat()

	info := strings.Split(strings.TrimSpace(args), config.C.ArgumentsSeparator)

	// Contact the server and print out its response.
	ctx2, cancel := context.WithTimeout(context.Background(), config.C.Timeout)
	defer cancel()
	_, err := client.Profile.SignUp(ctx2, &pb.SignUpRequest{
		UserID: user.ID,
		Name:   info[0],
	})
	if err != nil {
		log.Printf("couldn't sign up: %v", err)
		return reply(bot, chat, loc.Message(loc.SignUpFail))
	}

	return reply(bot, chat, loc.Message(loc.SignUpSuccess))
}

// Help command shows information about all commands
// type Help struct{}
//
//	func (c *Help) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
//		chat := upd.FromChat()
//
//		return reply(bot, chat, loc.Message(loc.Help))
//	}
//
// type Language struct{}
//
//	func (c *Language) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
//		chat := upd.FromChat()
//
//		return replyKeyboard(bot, chat, loc.Message(loc.Lang), makeInlineKeyboard(loc.SupportedLanguages, consts.LanguageButton))
//	}

// Ping replies with 'pong' message
type Ping struct{}

func (c *Ping) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	return reply(bot, chat, loc.Message(loc.Pong))
}
