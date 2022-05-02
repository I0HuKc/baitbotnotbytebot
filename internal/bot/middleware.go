package bot

import (
	"context"
	"os"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	tele "gopkg.in/telebot.v3"
)

// Если в кеше есть начатое действие, обработка запроса
// перекидывается в обработчик найденого кеш-кейса
func (b *baitbot) CacheHandle(ctx context.Context) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			val, err := b.cache.Get(strconv.Itoa(int(c.Message().Chat.ID)))
			if err != nil {
				if err != memcache.ErrCacheMiss {
					return err
				}

				return next(c)
			}
			defer b.cache.Delete(string(val.Value))

			return b.heap.Get(string(val.Value))(c)
		}
	}
}

func (b *baitbot) AdminOnly(ctx context.Context) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			id, err := strconv.Atoi(os.Getenv("APP_BOT_ADMID_ID"))
			if err != nil {
				return err
			}

			if c.Sender().ID != int64(id) {
				c.Send(ErrNotAvailableFoYou.Error())
				return nil
			}

			return next(c)
		}
	}
}
