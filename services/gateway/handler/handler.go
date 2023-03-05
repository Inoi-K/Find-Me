package handler

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/services/gateway/command"
	"github.com/Inoi-K/Find-Me/services/gateway/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

var (
	bot      *tgbotapi.BotAPI
	commands map[string]command.ICommand
)

func Start(ctx context.Context) error {
	var err error
	// Connect to the bot
	bot, err = tgbotapi.NewBotAPI(config.C.Token)
	if err != nil {
		return err
	}
	// Set this to true to log all interactions with telegram servers
	bot.Debug = true

	// Generate structs for commands
	commands = makeCommands()
	command.UpdateIndex(ctx)

	// Set update rate
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)
	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	return nil
}

// makeCommands creates all bot commands and buttons
func makeCommands() map[string]command.ICommand {
	return map[string]command.ICommand{
		// General commands
		command.StartCommand:    &command.Start{},
		command.PingCommand:     &command.Ping{},
		command.HelpCommand:     &command.Help{},
		command.LanguageCommand: &command.Language{},
		command.LanguageButton:  &command.LanguageCallback{},

		// Specific commands
		command.ShowProfileCommand: &command.ShowProfile{},
		command.EditProfileCommand: &command.EditProfile{},
		command.EditFieldButton:    &command.EditFieldCallback{},

		// Shortcut commands for testing
		command.FindCommand:  &command.Find{},
		command.MatchCommand: &command.Match{},
		//command.SignUpCommand: &command.SignUp{},
		//command.EditFieldCommand: &command.EditField{},
	}
}

// receiveUpdates handles updates and context cancel
func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		case update := <-updates:
			go handleUpdate(ctx, update)
		}
	}
}

// handleUpdate distributes incoming update
func handleUpdate(ctx context.Context, update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(ctx, update)

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(ctx, update)
	}
}

// handleMessage defines the type of the message (command or other - replies as echo in the latter case)
func handleMessage(ctx context.Context, update tgbotapi.Update) {
	message := update.Message
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	var err error
	if message.IsCommand() {
		err = handleCommand(ctx, update)
	} else {
		err = handleText(ctx, update)
	}

	if err != nil {
		log.Printf("couldn't process the message: %s", err.Error())
	}
}

// handleCommand handles commands specifically
func handleCommand(ctx context.Context, upd tgbotapi.Update) error {
	curCommand := upd.Message.Command()
	args := upd.Message.CommandArguments()
	return executeCommand(ctx, upd, curCommand, args)
}

func executeCommand(ctx context.Context, upd tgbotapi.Update, cmd string, args string) error {
	if c, ok := commands[cmd]; ok {
		return c.Execute(ctx, bot, upd, args)
	}

	return command.UnknownCommandError
}

// handleText handles text specifically
func handleText(ctx context.Context, upd tgbotapi.Update) error {
	message := upd.Message
	user := upd.SentFrom()

	// validated user (but might be not fully registered)
	if state, ok := session.UserState[user.ID]; ok {
		switch state {
		case session.EnterName, session.EnterDescription:
			session.UserStateArg[user.ID] <- message.Text
		case session.EnterAge:
			// TODO check for decimal
			session.UserStateArg[user.ID] <- message.Text
		case session.EnterPhoto:
			session.UserStateArg[user.ID] <- message.Photo[0].FileID
		case session.Matching:
			return executeCommand(ctx, upd, command.MatchCommand, message.Text)
		default:
			return command.UnknownStateError
		}
		return nil
	}

	// This is equivalent to forwarding, without the sender's name
	copyMsg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
	_, err := bot.CopyMessage(copyMsg)
	return err
}

// handleButton handles buttons callback specifically
func handleButton(ctx context.Context, upd tgbotapi.Update) {
	query := upd.CallbackQuery
	mainPart, args, isCommand := strings.Cut(query.Data, config.C.Separator)

	if isCommand {
		err := executeCommand(ctx, upd, mainPart, args)
		if err != nil {
			log.Printf("couldn't process button callback: %v", err)
		}
	} else {
		user := upd.SentFrom()

		if state, ok := session.UserState[user.ID]; ok {
			switch state {
			case session.EnterGender, session.EnterFaculty, session.EnterTags:
				_, err := updateInlineKeyboard(upd, query.Data)
				if err != nil {
					log.Printf(command.KeyboardUpdateError.Error())
					return
				}
				session.UserStateArg[user.ID] <- mainPart
			default:
				log.Printf(command.UnknownStateError.Error())
				return
			}
		}
	}

	// close the query
	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	_, err := bot.Request(callbackCfg)
	if err != nil {
		log.Printf("callback config error: %v", err)
	}
}

func updateInlineKeyboard(upd tgbotapi.Update, queryData string) (bool, error) {
	var isNew bool

	message := upd.CallbackQuery.Message
	newKeyboard := message.ReplyMarkup
	for rowIndex, row := range newKeyboard.InlineKeyboard {
		for columnIndex, button := range row {
			if *button.CallbackData == queryData {
				// TODO check a possible error with language changes
				if li := strings.LastIndex(button.Text, config.C.MarkSuccess); li != -1 {
					isNew = false
					newKeyboard.InlineKeyboard[rowIndex][columnIndex].Text = button.Text[:(li - 1)]
				} else {
					isNew = true
					newKeyboard.InlineKeyboard[rowIndex][columnIndex].Text += " " + config.C.MarkSuccess
				}
				break
			}
		}
	}

	msg := tgbotapi.NewEditMessageReplyMarkup(upd.FromChat().ID, message.MessageID, *newKeyboard)
	_, err := bot.Send(msg)
	if err != nil {
		return false, err
	}

	return isNew, nil
}
