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
	"strconv"
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
	ctx2, cancel := context.WithTimeout(ctx, config.C.Timeout)
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
	// add university and username
	signUpArgs += config.C.University + config.C.Separator + "@" + user.UserName
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

	// Contact the profile server
	ctx2, cancel := context.WithTimeout(ctx, config.C.Timeout)
	defer cancel()
	_, err := client.Profile.SignUp(ctx2, &pb.SignUpRequest{
		UserID:     user.ID,
		SphereID:   config.C.SphereID,
		Name:       info[0],
		Gender:     info[1],
		Age:        info[2],
		Faculty:    info[3],
		University: info[4],
		Username:   info[5],
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
	photoMsg, err := prepareProfile(ctx, upd.SentFrom().ID, upd.FromChat().ID)
	if err != nil {
		return err
	}

	_, err = bot.Send(photoMsg)
	return err
}

func prepareProfile(ctx context.Context, userID, chatID int64) (*tgbotapi.PhotoConfig, error) {
	// get user
	main, err := client.Profile.GetUserMain(ctx, &pb.GetUserMainRequest{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	ctx2, cancel := context.WithTimeout(ctx, config.C.Timeout)
	defer cancel()
	add, err := client.Profile.GetUserAdditional(ctx2, &pb.GetUserAdditionalRequest{
		UserID:   userID,
		SphereID: config.C.SphereID,
	})
	if err != nil {
		return nil, err
	}

	// build a message with profile photo and other info
	file := tgbotapi.FileID(add.PhotoID)
	photoMsg := tgbotapi.NewPhoto(chatID, file)

	tagsPart := "#" + strings.Join(add.Tags, " #")

	photoMsg.Caption = fmt.Sprintf("%s, %d y.o.\n%s\n%s\n\n%s", main.Name, main.Age, main.Faculty, add.Description, tagsPart)

	return &photoMsg, nil
}

type Find struct{}

func (c *Find) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	user := upd.SentFrom()
	chat := upd.FromChat()

	session.UserState[user.ID] = session.Matching
	// get next user id
	ctx2, cancel := context.WithTimeout(ctx, config.C.Timeout)
	defer cancel()
	next, err := client.Match.Next(ctx2, &pb.NextRequest{
		UserID:   user.ID,
		SphereID: config.C.SphereID,
	})
	if err != nil {
		// recommendations failed
		log.Printf("couldn't get recommendations: %v", err)
		_ = reply(bot, chat, loc.Message(loc.FindFail))
		return err
	}
	// no more users handling
	if next != nil && next.NextUserID == -1 {
		return reply(bot, chat, loc.Message(loc.FindEnd))
	}

	// prepare profile with like/dislike buttons
	nextProfile, err := prepareProfile(ctx, next.NextUserID, chat.ID)
	if err != nil {
		log.Printf("couldn't get profile %d: %v", next.NextUserID, err)
		_ = reply(bot, chat, loc.Message(loc.FindFail))
		return err
	}
	nextProfile.ReplyMarkup = MatchMarkup
	go func() {
		// TODO handle possible data loss
		if _, ok := session.UserStateArg[user.ID]; !ok {
			session.UserStateArg[user.ID] = make(chan string)
		}
		session.UserStateArg[user.ID] <- strconv.FormatInt(next.NextUserID, 10)
	}()

	_, err = bot.Send(nextProfile)
	return err
}

type Match struct{}

func (c *Match) Execute(ctx context.Context, bot *tgbotapi.BotAPI, upd tgbotapi.Update, args string) error {
	user := upd.SentFrom()
	likedUserID, err := strconv.ParseInt(<-session.UserStateArg[user.ID], 10, 0)
	if err != nil {
		return err
	}

	// no special handling for dislike
	if args == config.C.DislikeButton {
		return (&Find{}).Execute(ctx, bot, upd, args)
	}

	// like handling
	rep, err := client.Match.Like(ctx, &pb.LikeRequest{
		LikerID: user.ID,
		LikedID: likedUserID,
	})
	if err != nil {
		return err
	}

	// send contacts each user
	// or notify the liked one
	if rep.IsReciprocated {
		// get usernames
		// TODO parallel
		ctx2, cancel := context.WithTimeout(ctx, config.C.Timeout)
		defer cancel()
		likerMain, err := client.Profile.GetUserMain(ctx2, &pb.GetUserMainRequest{
			UserID: user.ID,
		})
		if err != nil {
			return err
		}
		ctx3, cancel := context.WithTimeout(ctx, config.C.Timeout)
		defer cancel()
		likedMain, err := client.Profile.GetUserMain(ctx3, &pb.GetUserMainRequest{
			UserID: likedUserID,
		})
		if err != nil {
			return err
		}

		// send contacts
		err = send(bot, user.ID, loc.Message(loc.Match)+likedMain.Username)
		if err != nil {
			return err
		}

		err = send(bot, likedUserID, loc.Message(loc.Match)+likerMain.Username)
		if err != nil {
			return err
		}
	} else {
		// notify the liked one that he/she has received a like
		err = send(bot, likedUserID, loc.Message(loc.LikeReceived))
		if err != nil {
			return err
		}
	}

	return (&Find{}).Execute(ctx, bot, upd, args)
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
