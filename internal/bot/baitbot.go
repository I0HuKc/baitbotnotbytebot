package bot

import (
	"context"
	"os"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db"
	"github.com/I0HuKc/baitbotnotbytebot/pkg/joker"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type baitbot struct {
	botApi *tele.Bot
	redis  *redis.Client
	cache  *memcache.Client
	heap   core.Heap
	joker  joker.Joker
}

func (b *baitbot) Serve() {
	b.botApi.Start()
}

func (b *baitbot) Configure(ctx context.Context) core.Baitbot {
	b.botApi.Use(middleware.AutoRespond())

	c := b.botApi.Group()
	c.Use(b.CacheHandle(ctx))

	c.Handle("/start", b.CommandStartHandle(ctx))
	c.Handle("/joke", b.CommandJokeHandle(ctx))
	c.Handle("/ping", b.CommandPingHandle(ctx), b.AdminOnly(ctx))
	c.Handle("/antre", b.CommandAntreHandle(ctx), b.AdminOnly(ctx))
	c.Handle("/forward", b.CommandForwardHandle(ctx), b.AdminOnly(ctx))

	c.Handle(tele.OnText, nil)
	c.Handle(tele.OnPhoto, nil)

	b.heap.Handle(CacheForward, b.CacheForwardHandle(ctx))

	return b
}

func NewBaitbot(rc *redis.Client) (core.Baitbot, error) {
	pref := tele.Settings{
		Token:  os.Getenv("APP_BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	botApi, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	var js joker.JSchema
	if err := js.ParseSchema(os.Getenv("JOKER_SCHEMA_PATH")); err != nil {
		return nil, err
	}

	mc, err := db.SetMemcacheConn()
	if err != nil {
		return nil, err
	}

	return &baitbot{
		botApi: botApi,
		redis:  rc,
		cache:  mc,
		heap:   make(heap),
		joker:  joker.CallJoker(js),
	}, nil
}
