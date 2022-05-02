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
		if _, err := b.botApi.Send(c.Sender(), c.Text(), &tele.SendOptions{
			DisableNotification: true,
		}); err != nil {
			return err
		}

		return nil
	}
}
