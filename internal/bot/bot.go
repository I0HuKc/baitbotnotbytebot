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
	sked   core.Sked
	acts   map[int64]core.Action

	store db.SqlStore
	redis rdstore.RedisStore
	joker joker.Joker
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
			if b.acts[update.Message.Chat.ID] != nil {
				go func() {
					b.ResponseHandler(update, b.sked.Handle(ctx, update, b.acts[update.Message.Chat.ID]))
				}()

				continue
			}

			if update.Message.IsCommand() {
				go func() {
					b.ResponseHandler(update, b.hub.Handle(ctx, update))
				}()

				continue
			}
		}
	}

	return nil
}

func (b *baitbot) Fuse() core.Baitbot {
	// Actions
	b.sked.SetHandleFunc(ActionSendJoke, b.ActionSendJokeHandle)

	// Private commands
	b.hub.SetHandleFunc(CommandStart, b.OnlyPrivateChat, b.CommandStartHandle)
	b.hub.SetHandleFunc(CommandGetDecription, b.OnlyPrivateChat, b.OnlyForAdmin, b.CommandGetDescHandle)
	b.hub.SetHandleFunc(CommandAddDecription, b.OnlyPrivateChat, b.OnlyForAuthor, b.CommandAddDescHandle)
	b.hub.SetHandleFunc(CommandGetDescriptionList, b.OnlyPrivateChat, b.OnlyForAdmin, b.CommandGetDescListHandle)
	b.hub.SetHandleFunc(CommandSendJoke, b.OnlyPrivateChat, b.OnlyForAdmin, b.CommandSendJokeHandle)
	b.hub.SetHandleFunc(CommandAntre, b.OnlyPrivateChat, b.OnlyForAdmin, b.CommandAntreHandle)

	// Group commands
	b.hub.SetHandleFunc(CommandJoke, b.OnlyGroupChat, b.CommandJokeHandle)
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
	}
}

func CreateBaitbot(b *tgbotapi.BotAPI, s db.SqlStore, r rdstore.RedisStore, j joker.Joker) core.Baitbot {
	return &baitbot{
		botApi: b,
		store:  s,
		redis:  r,
		joker:  j,

		hub:  make(hub),
		sked: make(sked),
		acts: make(map[int64]core.Action),
	}
}
