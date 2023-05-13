package command

import (
	"github.com/Inoi-K/Find-Me/pkg/config"
	loc "github.com/Inoi-K/Find-Me/services/gateway/localization"
	"github.com/Inoi-K/Find-Me/services/gateway/model"
	"github.com/Inoi-K/Find-Me/services/gateway/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// send builds message and sends it directly to the user
func send(bot *tgbotapi.BotAPI, userID int64, text string) error {
	msg := newMessage(userID, text, tgbotapi.InlineKeyboardMarkup{})
	_, err := bot.Send(msg)
	return err
}

// reply builds message and sends it to the chat
func reply(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := newMessage(chat.ID, text, tgbotapi.InlineKeyboardMarkup{})
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
func newMessage(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = config.C.ParseMode
	if len(keyboard.InlineKeyboard) != 0 {
		msg.ReplyMarkup = keyboard
	}
	return msg
}

// makeInlineKeyboard builds inline keyboard from the content
func makeInlineKeyboard(content []model.Content, commandButton string, columnCount int) tgbotapi.InlineKeyboardMarkup {
	// count the rows number taking in mind divisible/indivisible numbers
	rowCount := len(content) / columnCount
	if len(content)%columnCount != 0 {
		rowCount++
	}
	keyboard := make([][]tgbotapi.InlineKeyboardButton, rowCount)

	// specify command in the first part of the button data if needed
	commandPart := commandButton
	if len(commandPart) > 0 {
		commandPart += config.C.Separator
	}

	// build the keyboard
	for i := 0; i < rowCount; i++ {
		buttonPlaced := i * columnCount
		currentColumnCount := columnCount
		buttonLeft := len(content) - buttonPlaced
		if buttonLeft < columnCount {
			currentColumnCount = buttonLeft
		}
		columns := make([]tgbotapi.InlineKeyboardButton, currentColumnCount)
		for j := 0; j < len(columns); j++ {
			columns[j] = tgbotapi.NewInlineKeyboardButtonData(content[buttonPlaced+j].Text, commandPart+content[buttonPlaced+j].Data)
		}
		keyboard[i] = columns
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
	case Email:
		state = session.EnterEmail
		message = loc.EnterEmail
	}
	return
}
