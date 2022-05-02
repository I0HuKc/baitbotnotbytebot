package bot

import (
	"context"
	"os"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

func (b *baitbot) CacheForwardHandle(ctx context.Context) tele.HandlerFunc {
	return func(c tele.Context) error {
		id, err := strconv.Atoi(os.Getenv("APP_BOT_GROUP"))
		if err != nil {
			return err
		}

		c.Sender().ID = int64(id)
		so := tele.SendOptions{
			DisableNotification: true,
		}

		// Проверяю, было ли отправлено фото
		if c.Message().Photo != nil {
			if _, err := c.Message().Photo.Send(b.botApi, c.Sender(), &so); err != nil {
				return err
			}
		}

		// Проверяю, был ли указан текст
		if len(c.Text()) > 0 {
			if _, err := b.botApi.Send(c.Sender(), c.Text(), &so); err != nil {
				return err
			}
		}

		return nil
	}
}
