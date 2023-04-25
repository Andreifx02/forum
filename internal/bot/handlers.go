package bot

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	Handler func(msg tgbotapi.Message, dialog chan tgbotapi.Message)
	DialogSize int
}

func (b *Bot) signupHandler(msg tgbotapi.Message, dialog chan tgbotapi.Message) {
	err := b.SignUp(msg.From.UserName)
	if err != nil {
		b.answerWithError(msg, err)
		return 
	}
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "User is created"))
}

func (b *Bot) createPostHandler(msg tgbotapi.Message, dialog chan tgbotapi.Message) {
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Please, enter the topic"))
	
	reply := <-dialog
	topic := reply.Text

	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Enter the text"))

	reply = <-dialog
	text := reply.Text
	
	err := b.CreatePost(msg.From.UserName, topic, text, time.Unix(int64(msg.Date), 0))
	if err != nil {
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, err.Error()))
		return
	}
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Succesful:)"))
}

func (b *Bot) helpHandler(msg tgbotapi.Message, dialog chan tgbotapi.Message) {
	text := 
`
/signup : If you use this bot at first time
/post : Creates post
`
	b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
}

func (b *Bot) answerWithError(msg tgbotapi.Message, err error) {
	var text string
	if err != nil {
		text = err.Error()
		if text == "" {
			text = "Something went wrong. Please try again;("
		}
		b.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, text))
	} 
}

