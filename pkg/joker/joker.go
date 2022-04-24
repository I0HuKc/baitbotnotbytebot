package joker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/api"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db/rdstore"
	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type joker struct {
	schema JSchema
	redis  rdstore.RedisStore
}

type Joker interface {
	GetJokeHash(joke string) ([]byte, error)
	Joke(ctx context.Context) (*model.Joke, error)
}

func CallJoker(s JSchema, r rdstore.RedisStore) Joker {
	return &joker{
		schema: s,
		redis:  r,
	}
}

func (j *joker) Joke(ctx context.Context) (*model.Joke, error) {
	rt := j.getRandomTarget()
	fmt.Println(rt.Name)
	req := api.CreateApiReq(rt.Source.Target+j.schema.PrepareUrlParams(rt.Source.Params), nil)

	resp, err := req.Repeater(req.MakeGetReq, 3, time.Second)(ctx)
	if err != nil {
		return nil, err
	}

	var id string
	switch resp[rt.Id].(type) {
	case float64:
		id = fmt.Sprintf("%v", resp[rt.Id])
	case string:
		id = resp[rt.Id].(string)
	default:
		return nil, errors.New("unknown id type")
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

func (j *joker) GetJokeHash(joke string) ([]byte, error) {
	joke = strings.Replace(joke, "\n", "", -1)
	joke = strings.Replace(joke, "\t", "", -1)
	joke = strings.Replace(joke, " ", "", -1)

	return bcrypt.GenerateFromPassword(
		[]byte(strings.ToLower(joke)),
		bcrypt.MinCost,
	)
}
