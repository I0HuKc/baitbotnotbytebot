package baitbot

import (
	"context"
	"os"
	"strconv"
	"strings"

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
