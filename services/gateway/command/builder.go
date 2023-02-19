package command

import (
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/model"
	loc "github.com/Inoi-K/Find-Me/services/gateway/localization"
	"github.com/Inoi-K/Find-Me/services/gateway/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// reply builds message and sends it to the chat
func reply(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := newMessage(chat.ID, text, nil)
	_, err := bot.Send(msg)
	return err
}

// replyKeyboard builds message with keyboard and sends it to the chat
func replyKeyboard(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := newMessage(chat.ID, text, keyboard)
	_, err := bot.Send(msg)
	return err
}

// newMessage builds message with all needed parameters
func newMessage(chatID int64, text string, keyboard interface{}) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = config.C.ParseMode
	msg.ReplyMarkup = keyboard
	return msg
}

// makeInlineKeyboard builds inline keyboard from the content
func makeInlineKeyboard(content []model.Content, commandButton string) tgbotapi.InlineKeyboardMarkup {
	keyboard := make([][]tgbotapi.InlineKeyboardButton, len(content))

	for i := 0; i < len(content); i++ {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(content[i].Text, commandButton+config.C.Separator+content[i].Data),
		)
		keyboard[i] = row
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

// handleFieldSpecifics returns specific to a field state and a text message to ask
func handleFieldSpecifics(field string) (state session.State, message string, keyboard tgbotapi.InlineKeyboardMarkup) {
	switch field {
	case Name:
		state = session.EnterName
		message = loc.EnterName
	case Gender:
		state = session.EnterGender
		message = loc.EnterGender
		keyboard = EditGenderMarkup
	case Age:
		state = session.EnterAge
		message = loc.EnterAge
	case Faculty:
		state = session.EnterFaculty
		message = loc.EnterFaculty
		keyboard = EditFacultyMarkup
	case Photo:
		state = session.EnterPhoto
		message = loc.EnterPhoto
	case Description:
		state = session.EnterDescription
		message = loc.EnterDescription
	case Tags:
		state = session.EnterTags
		message = loc.EnterTags
		keyboard = EditTagsMarkup
	}
	return
}
