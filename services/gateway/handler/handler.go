package handler

import (
	"context"
	"github.com/Inoi-K/Find-Me/pkg/config"
	"github.com/Inoi-K/Find-Me/services/gateway/command"
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
		command.SignUpCommand: &command.SignUp{},
		//command.StartCommand: &command.Start{},
		//command.HelpCommand:     &command.Help{},
		//command.LanguageCommand: &command.Language{},
		//command.LanguageButton:  &command.LanguageButton{},
		//command.PingCommand:     &command.Ping{},
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
			handleUpdate(ctx, update)
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
		// This is equivalent to forwarding, without the sender's name
		copyMsg := tgbotapi.NewCopyMessage(message.Chat.ID, message.Chat.ID, message.MessageID)
		_, err = bot.CopyMessage(copyMsg)
	}

	if err != nil {
		log.Printf("couldn't process the message: %s", err.Error())
	}
}

// handleCommand handles commands specifically
func handleCommand(ctx context.Context, update tgbotapi.Update) error {
	curCommand := update.Message.Command()
	if cmd, ok := commands[curCommand]; ok {
		return cmd.Execute(ctx, bot, update, update.Message.CommandArguments())
	}

	return command.UnknownCommandError
}

// handleButton handles buttons callback specifically
func handleButton(ctx context.Context, update tgbotapi.Update) {
	query := update.CallbackQuery
	command, args, _ := strings.Cut(query.Data, config.C.ArgumentsSeparator)

	err := commands[command].Execute(ctx, bot, update, args)
	if err != nil {
		log.Printf("couldn't process button callback: %v", err)
	}

	// close the query
	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	_, err = bot.Request(callbackCfg)
	if err != nil {
		log.Printf("callback config error: %v", err)
	}
}