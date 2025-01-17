package bot

import (
	"github.com/TimIrwing/nyashka-butler/internal/message"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Wrapper struct {
	_b *tgbotapi.BotAPI
}

func New(token string) *Wrapper {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Can't get botapi from telegram: %s", err)
	}
	return &Wrapper{_b: bot}
}

func (bot Wrapper) Start() {
	log.Printf("Authorized on account %s", bot._b.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	ch := bot._b.GetUpdatesChan(updateConfig)

	for {
		u := <-ch
		go bot.HandleUpdate(u)
	}
}

func (bot Wrapper) HandleUpdate(u tgbotapi.Update) {
	switch {
	case u.Message != nil:
		m := message.From(u.Message, bot)
		m.Handle()
	}
}

func (bot Wrapper) Send(m *message.Message) {
	res := tgbotapi.NewMessage(m.ChatID, m.Text)
	res.ReplyToMessageID = m.ReplyTo
	_, err := bot._b.Send(res)
	if err != nil {
		log.Printf("Couldn't send message: %s", err)
	}
}
