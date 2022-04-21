package joker

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/api"
	"github.com/I0HuKc/baitbotnotbytebot/internal/db/rdstore"
	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"
)

type joker struct {
	schema JSchema
	redis  rdstore.RedisStore
}

type Joker interface {
	Joke(ctx context.Context) (*model.Joke, error)
	FormatJoke(joke *model.Joke) string
}

func CallJoker(s JSchema, r rdstore.RedisStore) Joker {
	return &joker{
		schema: s,
		redis:  r,
	}
}

func (j *joker) FormatJoke(joke *model.Joke) string {
	if joke.Delivery != "" {
		return fmt.Sprintf("%s\n\n||%s||", joke.Setup, joke.Delivery)
	}

	return joke.Setup
}

func (j *joker) Joke(ctx context.Context) (*model.Joke, error) {
	t := j.getRandomTarget()
	if t == nil {
		return nil, fmt.Errorf("failed to get joker target")
	}

	errs, ctx := errgroup.WithContext(ctx)
	req := api.CreateApiReq(t.Target, nil)
	i := 0
	cDone := make(chan bool, 1)
	cJoke := make(chan model.Joke, 1)

	defer close(cDone)
	defer close(cJoke)

	for {
		select {
		case <-cDone:
			fmt.Println("done")
			// if errs.Wait() != nil {
			// 	return nil, errs.Wait()
			// }

			joke := <-cJoke
			fmt.Println(11111)
			if err := j.redis.Save(ctx, joke.Id, joke, (24*time.Hour)*30); err != nil {
				return nil, err
			}

			fmt.Println(joke)

			return &joke, nil

		default:
			// Do other stuff
		}

		if i < 20 {
			time.Sleep(100 * time.Millisecond)

			errs.Go(func() error {
				i++
				defer func() {
					i--
				}()

				fmt.Printf("new loop %d\n", i)
				resp, err := req.Repeater(req.MakeGetReq, 3, time.Second)(ctx)
				if err != nil {
					return err
				}

				id, err := j.jokeIdToStr(resp[t.Id.Field], t.Id.Type)
				if err != nil {
					return nil
				}

				fmt.Println(id)
				_, rerr := j.redis.Get(ctx, id)
				switch rerr {
				// Значит эта шутка уже была
				case nil:
					fmt.Println("Эта шутка была")
					return nil

				// Новая шутка
				case redis.Nil:
					fmt.Println("Новая шутка")

					var d string
					if len(t.Read) == 2 {
						d = resp[t.Read[1]].(string)
					}

					cDone <- true
					cJoke <- model.Joke{
						Id:       id,
						Target:   t.Name,
						Setup:    resp[t.Read[0]].(string),
						Delivery: d,
					}

					return nil

				default:
					return rerr
				}
			})
		}

	}
}

func (j *joker) checkDelivery(v any, ok bool) string {
	if ok {
		return v.(string)
	}
	return ""
}

func (j *joker) jokeIdToStr(id interface{}, idType string) (string, error) {
	switch idType {
	case "str":
		return id.(string), nil

	case "int":
		if v, ok := id.(float64); ok {
			return strconv.Itoa(int(v)), nil
		}

		return "", errors.New("id_type isn't int")

	default:
		return "", errors.New("unknown id_type")
	}
}

func (j *joker) getRandomTarget() *sTarget {
	i := int(float32(len(j.schema.Joker.Targets)) * rand.Float32())
	for _, v := range j.schema.Joker.Targets {
		if i == 0 {
			return &v
		}
		i--
	}

	return nil
}
