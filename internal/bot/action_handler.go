package baitbot

import (
	"context"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *baitbot) ActionSendJokeHandle(ctx context.Context, update tgbotapi.Update) error {
	defer delete(b.acts, update.Message.Chat.ID)

	groupId, err := strconv.Atoi(os.Getenv("APP_BOT_GROUP"))
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(int64(groupId), update.Message.Text)
	msg.DisableNotification = true
	return b.Send(b.botApi.Send, msg)
}
