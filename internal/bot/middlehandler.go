package baitbot

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Пропускает к след обработчикам, только если команда
// используется в личной переписке с ботом.
func (b *baitbot) OnlyPrivateChat(ctx context.Context, update tgbotapi.Update) error {
	if update.Message.Chat.IsPrivate() {
		return nil
	}

	return ErrOnlyForPrivateChat
}

// Пропускает к след обработчикам, только если команда
// используется в групповой переписке.
func (b *baitbot) OnlyGroupChat(ctx context.Context, update tgbotapi.Update) error {
	if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() {
		return nil
	}

	return ErrNotSupportedInGroup
}

// Пропускает к след обработчикам, только если chat_id
// пользователя добавлен в .env файле в поле админа.
func (b *baitbot) OnlyForAdmin(ctx context.Context, update tgbotapi.Update) error {
	id, err := strconv.Atoi(os.Getenv("APP_BOT_ADMID_ID"))
	if err != nil {
		return err
	}

	if update.Message.From.ID != int64(id) {
		return ErrNotAvailableFoYou
	}

	return nil
}

// Пропускает к след обработчикам, только если chat_id
// пользователя добавлен в .env файле в поле авторов.
func (b *baitbot) OnlyForAuthor(ctx context.Context, update tgbotapi.Update) error {
	for _, v := range strings.Split(os.Getenv("APP_BOT_AUTHOR"), ",") {
		id, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		if update.Message.From.ID == int64(id) {
			return nil
		}
	}

	return ErrNotAvailableFoYou
}

// Пропускает к след обработчикам, только если
// антре еще НЕ создана для этой группы.
func (b *baitbot) NoRunPerformance(ctx context.Context, update tgbotapi.Update) error {
	if err := b.store.Performance().GetByGroupId(ctx,
		&model.Performance{GroupId: int(update.Message.Chat.ID)}); err != nil {
		if err != sql.ErrNoRows {
			return err
		}

		return nil
	}

	return ErrAlreadyCreated
}

// Пропускает к след обработчикам, только если
// антре уже создана для этой группы.
func (b *baitbot) RunPerformance(ctx context.Context, update tgbotapi.Update) error {
	if err := b.store.Performance().GetByGroupId(ctx,
		&model.Performance{GroupId: int(update.Message.Chat.ID)}); err != nil {
		if err == sql.ErrNoRows {
			return ErrNoAntreCreated
		}

		return err
	}

	return nil
}

// Предварительное сохранение данных
// перед запуском антре
func (b *baitbot) SaveAntre(ctx context.Context, update tgbotapi.Update) error {
	if err := b.store.Performance().Create(ctx, &model.Performance{
		GroupId:   int(update.Message.Chat.ID),
		GroupName: update.Message.Chat.Title,
		NextJoke:  time.Now().UTC().Format("2006-01-02T15:04:05.00000000"),
	}); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Джокер в игре!")
	return b.Send(b.botApi.Send, msg)
}
