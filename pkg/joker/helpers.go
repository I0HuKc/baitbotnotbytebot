package joker

import (
	"math/rand"
)

func (j *joker) getDeliveryText(t *sTarget, resp map[string]interface{}) string {
	if len(t.Read) == 2 {
		return resp[t.Read[1]].(string)
	}

	return ""
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
