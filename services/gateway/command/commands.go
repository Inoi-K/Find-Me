package command

import (
	"context"
	"fmt"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/services/gateway/client"
	loc "github.com/Inoi-K/Find-Me/services/gateway/localization"
	"github.com/Inoi-K/Find-Me/services/gateway/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

// ICommand provides an interface for all commands and buttons callbacks
type ICommand interface {
	Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error
}

// Ping replies with 'pong' message
type Ping struct{}

func (c *Ping) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()
	return reply(bot, chat, loc.Message(loc.Pong))
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
	ctx2, cancel := context.WithTimeout(context.Background(), config.C.Timeout)
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
	var signUpArgs string
	for _, field := range []string{Name, Gender, Age, Faculty} {
		newArg, err := askStateField(ctx, bot, chat, user.ID, field)
		if err != nil {
			return err
		}
		signUpArgs += newArg + config.C.Separator
	}
	// sign up
	err = (&SignUp{}).Execute(ctx, bot, upd, signUpArgs)
	if err != nil {
		return err
	}

	// additional information
	editFieldButton := &EditFieldCallback{}
	for _, field := range []string{Photo, Description, Tags} {
		err = editFieldButton.Execute(ctx, bot, upd, field)
		if err != nil {
			return err
		}
	}

	return reply(bot, chat, loc.Message(loc.Rubicon))
}

// askStateField handles getting a new value for the field
func askStateField(ctx context.Context, bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, userID int64, field string) (string, error) {
	state, message, keyboard := handleFieldSpecifics(field)

	// create user state
	session.UserState[userID] = state

	// get arguments
	var newArg string
	var err error
	if len(keyboard.InlineKeyboard) == 0 {
		newArg, err = askArg(ctx, bot, chat, userID, message)
	} else {
		newArg, err = askArgKeyboard(ctx, bot, chat, userID, message, keyboard)
	}
	if err != nil {
		return "", err
	}

	// clear user state
	delete(session.UserState, userID)

	return newArg, nil
}

// askArg asks a user for a new value of the field and waits for it - response
func askArg(ctx context.Context, bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, userID int64, messageKey string) (string, error) {
	return askArgKeyboard(ctx, bot, chat, userID, messageKey, tgbotapi.InlineKeyboardMarkup{})
}

// askArgKeyboard asks with a keyboard a user for a new value of the field and waits for it - response
func askArgKeyboard(ctx context.Context, bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, userID int64, messageKey string, keyboard tgbotapi.InlineKeyboardMarkup) (string, error) {
	err := replyKeyboard(bot, chat, loc.Message(messageKey), keyboard)
	if err != nil {
		return "", err
	}

	session.UserStateArg[userID] = make(chan string)
	arg := ""
	switch session.UserState[userID] {
	case session.EnterName, session.EnterGender, session.EnterAge, session.EnterFaculty, session.EnterPhoto, session.EnterDescription:
		select {
		case <-ctx.Done():
			return "", ContextDoneError
		case arg = <-session.UserStateArg[userID]:
		}

	case session.EnterTags:
		tags := make(map[string]struct{}, config.C.TagsLimit)

		for len(tags) < config.C.TagsLimit {
			select {
			case <-ctx.Done():
				return "", ContextDoneError
			case newArg := <-session.UserStateArg[userID]:
				if _, picked := tags[newArg]; picked {
					delete(tags, newArg)
				} else {
					tags[newArg] = struct{}{}
				}
			}
		}

		for tag := range tags {
			arg += tag + config.C.Separator
		}
		arg = strings.TrimSpace(arg)
	}
	close(session.UserStateArg[userID])

	return arg, nil
}

// SignUp sends a request to the profile service to register a new user
type SignUp struct{}

func (c *SignUp) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	user := upd.SentFrom()
	chat := upd.FromChat()

	info := strings.Split(strings.TrimSpace(args), config.C.Separator)

	// Contact the server and print out its response.
	ctx2, cancel := context.WithTimeout(context.Background(), config.C.Timeout)
	defer cancel()
	_, err := client.Profile.SignUp(ctx2, &pb.SignUpRequest{
		UserID:   user.ID,
		SphereID: config.C.SphereID,
		Name:     info[0],
		Gender:   info[1],
		Age:      info[2],
		Faculty:  info[3],
	})
	if err != nil {
		log.Printf("couldn't sign up: %v", err)
		_ = reply(bot, chat, loc.Message(loc.SignUpFail))
		return err
	}

	return reply(bot, chat, loc.Message(loc.SignUpSuccess))
}

// EditProfile sends inline keyboard with fields available for editing
type EditProfile struct{}

func (c *EditProfile) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	return replyKeyboard(bot, chat, loc.Message(loc.EditMenu), EditProfileMarkup)
}

// ShowProfile shows a profile with its image, main and additional info
type ShowProfile struct{}

func (c *ShowProfile) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	userID := upd.SentFrom().ID
	// get user
	main, err := client.Profile.GetUserMain(ctx, &pb.GetUserMainRequest{
		UserID: userID,
	})
	if err != nil {
		return err
	}
	add, err := client.Profile.GetUserAdditional(ctx, &pb.GetUserAdditionalRequest{
		UserID:   userID,
		SphereID: config.C.SphereID,
	})
	if err != nil {
		return err
	}

	file := tgbotapi.FileID(add.PhotoID)
	photoMsg := tgbotapi.NewPhoto(upd.FromChat().ID, file)

	tagsPart := "#" + strings.Join(add.Tags, " #")

	photoMsg.Caption = fmt.Sprintf("%s, %d y.o.\n%s\n%s\n\n%s", main.Name, main.Age, main.Faculty, add.Description, tagsPart)

	_, err = bot.Send(photoMsg)
	return err
}

// Help command shows information about all commands
type Help struct{}

func (c *Help) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	return reply(bot, chat, loc.Message(loc.Help))
}

// Language replies with buttons of languages available for change
type Language struct{}

func (c *Language) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	chat := upd.FromChat()

	return replyKeyboard(bot, chat, loc.Message(loc.Lang), makeInlineKeyboard(loc.SupportedLanguages, LanguageButton, 1))
}
