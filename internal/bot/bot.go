package baitbot

import (
	"context"
	"unicode"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db/rdstore"
	"github.com/I0HuKc/baitbotnotbytebot/pkg/joker"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type baitbot struct {
	botApi *tgbotapi.BotAPI
	hub    core.Hub
	store  db.SqlStore
	redis  rdstore.RedisStore
	joker  joker.Joker

	scd chan bool
}

func (b *baitbot) Serve(ctx context.Context) (err error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botApi.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				go func() {
					b.ResponseHandler(update, b.hub.HandleFunc(ctx, update))
				}()
			}
		}
	}

	return nil
}

func (b *baitbot) SetHub() core.Baitbot {
	// Private commands
	b.hub.SetHandleFunc(CommandStart, b.OnlyPrivateChat, b.CommandStartHandle)
	b.hub.SetHandleFunc(CommandGetDecription, b.OnlyPrivateChat, b.OnlyForAdmin, b.CommandGetDescHandle)
	b.hub.SetHandleFunc(CommandAddDecription, b.OnlyPrivateChat, b.OnlyForAuthor, b.CommandAddDescHandle)
	b.hub.SetHandleFunc(CommandGetDescriptionList, b.OnlyPrivateChat, b.OnlyForAdmin, b.CommandGetDescListHandle)
	b.hub.SetHandleFunc(CommandHelp, b.OnlyPrivateChat, b.CommandHelpDescHandle)

	// Group commands
	b.hub.SetHandleFunc(CommandAntre, b.OnlyGroupChat, b.OnlyForAdmin, b.NoRunPerformance, b.SaveAntre, b.CommandAntreHandle)
	b.hub.SetHandleFunc(CommandStopAntre, b.OnlyGroupChat, b.OnlyForAdmin, b.RunPerformance, b.CommandStopAntreHandle)
	b.hub.SetHandleFunc(CommandJoke, b.OnlyGroupChat, b.OnlyForAdmin, b.CommandJokeHandle)
	b.hub.SetHandleFunc(CommandPing, b.OnlyGroupChat, b.OnlyForAdmin, b.CommandPingHandle)
	b.hub.SetHandleFunc(CommandStartChangeDecription, b.OnlyGroupChat, b.OnlyForAdmin, b.CommandStartChangeDescHandle)

	return b
}

func (b *baitbot) ResponseHandler(update tgbotapi.Update, err error) {
	if err != nil {
		r := []rune(err.Error())
		r[0] = unicode.ToUpper(r[0])

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(r))
		b.botApi.Send(msg)
		return
	}
}

func CreateBaitbot(b *tgbotapi.BotAPI, s db.SqlStore, r rdstore.RedisStore, j joker.Joker) core.Baitbot {
	return &baitbot{
		botApi: b,
		store:  s,
		redis:  r,
		joker:  j,

		hub: make(hub),
	}
}
