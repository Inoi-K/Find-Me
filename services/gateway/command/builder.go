package command

import (
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/pkg/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// Reply builds message and sends it to the chat
func Reply(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := NewMessage(chat.ID, text, nil)
	_, err := bot.Send(msg)
	return err
}

// ReplyKeyboard builds message with keyboard and sends it to the chat
func ReplyKeyboard(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string, keyboard interface{}) error {
	msg := NewMessage(chat.ID, text, keyboard)
	_, err := bot.Send(msg)
	return err
}

// NewMessage builds message with all needed parameters
func NewMessage(chatID int64, text string, keyboard interface{}) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = config.C.ParseMode
	msg.ReplyMarkup = keyboard
	return msg
}

// MakeInlineKeyboard builds inline keyboard from the content
func MakeInlineKeyboard(content []model.Content, commandButton string) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, c := range content {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(c.Text, strings.Join([]string{commandButton, c.Data}, config.C.ArgumentsSeparator)),
		)
		keyboard = append(keyboard, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
