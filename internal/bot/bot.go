package bot

import (
	"sync"

	"github.com/Andreifx02/forum/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	serverAddress string 
	handlers map[string]Handler
	updatesChan chan tgbotapi.Update
	dialogs map[string]Dialog
}

type Dialog struct {
	C chan tgbotapi.Message
	Counter int
}


func New(cfg *config.Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		return nil, err
	}

	b := Bot{
		bot: bot,
		serverAddress: cfg.Bot.Server,
		updatesChan: make(chan tgbotapi.Update),
		dialogs: make(map[string]Dialog),
	}

	b.addHandlers()

	return &b, nil
}

func (b *Bot) addHandlers() {
	b.handlers = make(map[string]Handler)
	b.handlers["help"] = Handler{
		Handler: b.helpHandler,
		DialogSize: 0,
	}
	b.handlers["post"] = Handler{
		Handler: b.createPostHandler,
		DialogSize: 2,
	}
	b.handlers["signup"] = Handler{
		Handler: b.signupHandler,
		DialogSize: 0,
	}
}

func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		for update := range updates {
			b.updatesChan <- update
		}
		close(b.updatesChan)
		wg.Done()
	}()
	

	for update := range b.updatesChan {
		if update.Message != nil { // If we got a message
			dialog, ok := b.dialogs[update.Message.From.UserName]
			if ok {
				if dialog.Counter != 0 {
					dialog.C <- *update.Message
					dialog.Counter--
					b.dialogs[update.Message.From.UserName] = dialog
					continue
				} else {
					close(dialog.C)
					delete(b.dialogs, update.Message.From.UserName)
				}
			}
			if cmdName := update.Message.Command(); cmdName != "" {
				handler, ok := b.handlers[cmdName]
				if !ok {
					b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command"))
				} else {
					wg.Add(1)
					b.dialogs[update.Message.From.UserName] = Dialog {
						C: make(chan tgbotapi.Message, handler.DialogSize),
						Counter: handler.DialogSize,
					}
					go handler.Handler(*update.Message, b.dialogs[update.Message.From.UserName].C)
					wg.Done()
				}
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			
			b.bot.Send(msg)
		}
	}

	wg.Wait()
}

