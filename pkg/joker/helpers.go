package joker

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
)

func (j *joker) getDeliveryText(t *sTarget, resp map[string]interface{}) string {
	if len(t.Read) == 2 {
		return resp[t.Read[1]].(string)
	}

	return ""
}

func (j *joker) isNewJoke(ctx context.Context, id string) bool {
	if _, err := j.redis.Get(ctx, id); err != nil {
		return true
	}

	return false
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
