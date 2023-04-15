package bot

import (
	"strings"
	"time"

	"github.com/Andreifx02/forum/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	serverAddress string 
}

func New(cfg *config.Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, err
	}
	return &Bot{
		bot: bot,
		serverAddress: cfg.Bot.Server,
	}, nil
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if cmdName := update.Message.Command(); cmdName != "" {
				if cmdName == "signup" {
					text := "User is created"
					err := b.SignUp(update.Message.From.UserName)
					if err != nil {
						text = err.Error()
					} 
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
					b.bot.Send(msg)
					continue
				} 
				if cmdName == "create_post" {
					args := strings.Split(update.Message.CommandArguments(), "/")
					if len(args) != 2 {
						b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid format"))
						continue
					}
					err := b.CreatePost(update.Message.From.UserName, args[0], args[1], time.Unix(int64(update.Message.Date), 0))
					if err != nil {
						b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
						continue
					}
					b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Succesful:)"))
					continue
				}
				b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unckown command"))
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			
			b.bot.Send(msg)
		}
	}
}

