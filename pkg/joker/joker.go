package joker

import (
	"context"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/api"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db/rdstore"
	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
)

type joker struct {
	schema JSchema
	redis  rdstore.RedisStore
}

type Joker interface {
	Joke(ctx context.Context) (*model.Joke, error)
}

func CallJoker(s JSchema, r rdstore.RedisStore) Joker {
	return &joker{
		schema: s,
		redis:  r,
	}
}

func (j *joker) Joke(ctx context.Context) (*model.Joke, error) {
	// Получаю случайный ресурс из доступных
	rt := j.getRandomTarget()

	req := api.CreateApiReq(rt.Target, nil)
	for {
		resp, err := req.Repeater(req.MakeGetReq, 3, time.Second)(ctx)
		if err != nil {
			return nil, err
		}

		id, err := j.jokeIdToStr(resp[rt.Id.Field], rt.Id.Type)
		if err != nil {
			return nil, err
		}

		if ok := j.isNewJoke(ctx, id); ok {
			if rt.Id.Save {
				if err := j.redis.Save(ctx, id, rt.Name, (24*time.Hour)*30); err != nil {
					return nil, err
				}
			}

			return &model.Joke{
				Id:       id,
				Target:   rt.Name,
				Setup:    resp[rt.Read[0]].(string),
				Delivery: j.getDeliveryText(rt, resp),
				Lang: model.JLang{
					Source:    rt.Lang.Source,
					Target:    rt.Lang.Target,
					Translate: rt.Lang.Translate,
				},
			}, nil
		}
	}
}
