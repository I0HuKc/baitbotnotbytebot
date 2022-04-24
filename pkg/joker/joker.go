package joker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/api"
	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type joker struct {
	schema JSchema
}

type Joker interface {
	GetJokeHash(joke string) ([]byte, error)
	Joke(ctx context.Context) (*model.Joke, error)
}

func CallJoker(s JSchema) Joker {
	return &joker{
		schema: s,
	}
}

func (j *joker) Joke(ctx context.Context) (*model.Joke, error) {
	rt := j.getRandomTarget()

	req := api.CreateApiReq[any, map[string]any](
		(rt.Source.Target + j.schema.PrepareUrlParams(rt.Source.Params)),
		nil,
		5*time.Second,
	)

	resp, err := req.Repeater(req.NewRequest, 3, time.Second)(ctx, rt.Source.Method)
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
