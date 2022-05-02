package bot

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	gt "github.com/bas24/googletranslatefree"
	"github.com/bradfitz/gomemcache/memcache"
	tele "gopkg.in/telebot.v3"
)

func (b *baitbot) CommandPingHandle(ctx context.Context) tele.HandlerFunc {
	return func(c tele.Context) error {
		return c.Send("pong")
	}
}

func (b *baitbot) CommandForwardHandle(ctx context.Context) tele.HandlerFunc {
	return func(c tele.Context) error {
		if err := b.cache.Set(&memcache.Item{
			Key:        strconv.Itoa(int(c.Sender().ID)),
			Value:      []byte(CacheForward),
			Expiration: 60,
		}); err != nil {
			return c.Send(err.Error())
		}

		return c.Send("Хорошо, пришли шутку которую необходимо отправить в беседу.")
	}
}

func (b *baitbot) CommandStartHandle(ctx context.Context) tele.HandlerFunc {
	return func(c tele.Context) error {
		return c.Send("Зачем вообще нужно приветствие? Сразу к делу!")
	}
}

func (b *baitbot) CommandAntreHandle(ctx context.Context) tele.HandlerFunc {
	return func(c tele.Context) error {
		for {
			joke, err := b.joker.Joke(ctx)
			if err != nil {
				c.Send(err.Error())
				continue
			}

			jtext := b.formatJoke(joke)
			hash, err := b.isNewJoke(ctx, jtext)
			if err != nil {
				continue
			}

			if err := b.redis.Set(ctx, hash, joke.Target, (24*time.Hour)*30).Err(); err != nil {
				c.Send(err.Error())
				continue
			}

			if joke.Lang.Translate {
				jtext, err = gt.Translate(jtext, joke.Lang.Source, joke.Lang.Target)
				if err != nil {
					c.Send(err.Error())
					continue
				}
			}

			if err := c.Send(b.escapeFormat(jtext), &tele.SendOptions{
				ParseMode:           tele.ModeMarkdownV2,
				DisableNotification: true,
			}); err != nil {
				c.Send(err.Error())
				continue
			}

			itr := int(float32(10 * rand.Float32()))
			if err := c.Send(
				fmt.Sprintf("Сделующая шутка через *%s*", time.Duration(itr)*core.GetInterval()),
				&tele.SendOptions{
					ParseMode:           tele.ModeMarkdownV2,
					DisableNotification: true,
				},
			); err != nil {
				c.Send(err.Error())
				continue
			}

			time.Sleep(time.Duration(itr) * core.GetInterval())
		}
	}
}

func (b *baitbot) CommandJokeHandle(ctx context.Context) tele.HandlerFunc {
	return func(c tele.Context) (err error) {
		for {
			joke, err := b.joker.Joke(ctx)
			if err != nil {
				return c.Send(err.Error())
			}

			jtext := b.formatJoke(joke)
			hash, err := b.isNewJoke(ctx, jtext)
			if err != nil {
				return c.Send(err.Error())
			}

			if joke.Lang.Translate {
				jtext, err = gt.Translate(jtext, joke.Lang.Source, joke.Lang.Target)
				if err != nil {
					return c.Send(err.Error())
				}
			}

			if err := b.redis.Set(ctx, hash, joke.Target, (24*time.Hour)*30).Err(); err != nil {
				return c.Send(err.Error())
			}

			return c.Send(b.escapeFormat(jtext), &tele.SendOptions{
				ParseMode: tele.ModeMarkdownV2,
			})
		}
	}
}
