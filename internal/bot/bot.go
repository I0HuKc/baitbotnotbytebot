package baitbot

import (
	"fmt"
	"log"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type baitbot struct {
	botApi *tgbotapi.BotAPI
}

type BaitbotI interface {
	Serve() (err error)
}

func (b *baitbot) Serve() (err error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botApi.GetUpdatesChan(u)
	if err != nil {
		return
	}

	for update := range updates {

		if update.Message.Chat.IsGroup() {
			if update.Message.IsCommand() {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				if update.Message.Text == core.Bullying {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сорян, я еще не умею буллить :(")
					b.botApi.Send(msg)

					continue
				}

				if update.Message.Text == core.ChangeDesc {

					act := tgbotapi.NewChatDescription(update.Message.Chat.ID, "Test dgfg ghh sdfgh fggh")
					if _, err := b.botApi.Send(act); err != nil {
						fmt.Println(err)
					}
				}
			}
		}

		if update.Message.Chat.IsPrivate() {
			if update.Message.IsCommand() {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				continue
			}
		}

		// if update.Message.Text == core.Start {
		// 	fmt.Println(update.Message.Chat.IsGroup())
		// }

		// if update.Message != nil {
		// 	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// 	msg.ReplyToMessageID = update.Message.MessageID

		// 	b.botApi.Send(msg)
		// 	continue
		// }
	}

	return nil
}

func CreateBaitbot(botApi *tgbotapi.BotAPI) BaitbotI {
	return &baitbot{
		botApi: botApi,
	}
}
