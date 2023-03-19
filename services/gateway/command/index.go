package command

import (
	"context"
	"errors"
	"github.com/Inoi-K/Find-Me/pkg/api/pb"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/model"
	"github.com/Inoi-K/Find-Me/services/gateway/client"
	loc "github.com/Inoi-K/Find-Me/services/gateway/localization"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const (
	StartCommand       = "start"
	PingCommand        = "ping"
	ShowProfileCommand = "profile"
	EditProfileCommand = "edit"
	EditFieldButton    = "editButton"
	EditFieldCommand   = "editField"
	SignUpCommand      = "signup"
	FindCommand        = "find"
	MatchCommand       = "match"

	HelpCommand     = "help"
	LanguageCommand = "lang"
	LanguageButton  = "langButton"

	// FIELDS
	Name        = "name"
	Gender      = "gender"
	Age         = "age"
	Faculty     = "faculty"
	Photo       = "photo"
	Description = "description"
	Tags        = "tags"
	Email       = "email"

	Male   = "M"
	Female = "F"

	DeleteAccount = "deleteAccount"
)

var (
	EditProfileMarkup tgbotapi.InlineKeyboardMarkup
	EditGenderMarkup  tgbotapi.InlineKeyboardMarkup
	EditFacultyMarkup tgbotapi.InlineKeyboardMarkup
	EditTagsMarkup    tgbotapi.InlineKeyboardMarkup
	MatchMarkup       tgbotapi.ReplyKeyboardMarkup

	UnknownCommandError = errors.New("unknown command")
	ContextDoneError    = errors.New("context is done")
	UnknownStateError   = errors.New("unknown state")
	KeyboardUpdateError = errors.New("couldn't update inline keyboard")
)

// UpdateIndex creates keyboard markups in format '<EditFieldButton><argumentSeparator><field>
func UpdateIndex(ctx context.Context) {
	// EditMenu
	EditProfileMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Photo), EditFieldButton+config.C.Separator+Photo),
			tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Description), EditFieldButton+config.C.Separator+Description),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Age), EditFieldButton+config.C.Separator+Age),
			tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Faculty), EditFieldButton+config.C.Separator+Faculty),
		),
		tgbotapi.NewInlineKeyboardRow(
			//tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Age), EditFieldButton+config.C.Separator+Age),
			tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Tags), EditFieldButton+config.C.Separator+Tags),
		),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.DeleteAccount), EditFieldButton+config.C.Separator+DeleteAccount),
		//),
	)

	// gender
	EditGenderMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Male), Male),
			tgbotapi.NewInlineKeyboardButtonData(loc.Message(loc.Female), Female),
		),
	)

	// faculties
	faculties := make([]model.Content, len(config.C.Faculties))
	for i := 0; i < len(config.C.Faculties); i++ {
		faculties[i] = model.Content{
			Text: config.C.Faculties[i],
			Data: config.C.Faculties[i],
		}
	}
	EditFacultyMarkup = makeInlineKeyboard(faculties, "", 2)

	// tags
	rep, err := client.Profile.GetTags(ctx, &pb.GetTagsRequest{SphereID: config.C.SphereID})
	if err != nil {
		log.Fatalf("couldn't get tags: %v", err)
	}
	tags := make([]model.Content, len(rep.TagIDs))
	for i := 0; i < len(rep.TagIDs); i++ {
		tags[i] = model.Content{
			Text: rep.TagNames[i],
			Data: rep.TagIDs[i],
		}
	}
	EditTagsMarkup = makeInlineKeyboard(tags, "", 4)

	// match
	MatchMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(config.C.DislikeButton),
			tgbotapi.NewKeyboardButton(config.C.LikeButton),
		),
	)
}
