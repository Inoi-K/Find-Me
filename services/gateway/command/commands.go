package command

import (
	"context"
	"fmt"
	pb "github.com/Inoi-K/Find-Me/pkg/api"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/model"
	loc "github.com/Inoi-K/Find-Me/services/gateway/localization"
	"github.com/Inoi-K/Find-Me/services/gateway/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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
	usr := upd.SentFrom()

	ok := loc.ChangeLanguage(usr.LanguageCode)
	// if user's language is not supported then set default language to english
	if !ok {
		loc.ChangeLanguage("en")
	}

	// TODO validate terms & agreement

	// TODO validate user with corporate email

	// main information
	// name
	// create session with user
	session.Users[usr.ID] = &model.User{}
	return reply(bot, chat, loc.Message(loc.EnterName))
}

// SignUp sends a request to the profile service to register a new user
type SignUp struct{}

func (c *SignUp) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	usr := upd.SentFrom()

	// Set up a connection to the server.
	address := fmt.Sprintf("%s:%s", config.C.ProfileHost, config.C.ProfilePort)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewProfileClient(conn)

	// Contact the server and print out its response.
	ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = client.SignUp(ctx2, &pb.SignUpRequest{
		UserID: usr.ID,
		Name:   args,
	})
	if err != nil {
		log.Fatalf("couldn't sign up: %v", err)
	}

	return nil
}

// Help command shows information about all commands
//type Help struct{}
//
//func (c *Help) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
//	chat := upd.FromChat()
//
//	return reply(bot, chat, loc.Message(loc.Help))
//}
//
//type Language struct{}
//
//func (c *Language) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
//	chat := upd.FromChat()
//
//	return replyKeyboard(bot, chat, loc.Message(loc.Lang), makeInlineKeyboard(loc.SupportedLanguages, consts.LanguageButton))
//}
//
//type Ping struct{}
//
//func (c *Ping) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
//	chat := upd.FromChat()
//	return reply(bot, chat, loc.Message(loc.Pong))
//}
