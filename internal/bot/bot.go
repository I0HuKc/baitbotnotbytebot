package baitbot

import (
	"context"
	"fmt"
	"log"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type baitbot struct {
	botApi *tgbotapi.BotAPI
	store  db.SqlStoreI
}

func (b *baitbot) Serve(ctx context.Context) (err error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botApi.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
			if update.Message.IsCommand() {
				if b.IsLocal() {
					log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				}

				b.ErrorHandler(ctx, b.GroupCmdHandler, update)
			}
		}

		if update.Message.Chat.IsPrivate() {
			if b.IsLocal() {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			}

			if update.Message.IsCommand() {
				b.ErrorHandler(ctx, b.PrivateCmdHandler, update)
			}
		}
	}

	return nil
}

func (b *baitbot) ErrorHandler(ctx context.Context, handler core.Handler, update tgbotapi.Update) {
	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	if err := handler(ctx, update); err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "⚠️ Похоже что-то пошло не так ⚠️ ")
		b.botApi.Send(msg)

		fmt.Println(ctx.Err() != context.Canceled)
		if ctx.Err() != context.Canceled {
			ab := fmt.Sprintf("🔄[%s] — %s🔄", update.Message.Chat.UserName, err.Error())
			if err := b.AdminNotify(ab); err != nil {
				log.Println(err)
			}
		}
	}
}

func CreateBaitbot(b *tgbotapi.BotAPI, s db.SqlStoreI) core.Baitbot {
	return &baitbot{
		botApi: b,
		store:  s,
	}
}
