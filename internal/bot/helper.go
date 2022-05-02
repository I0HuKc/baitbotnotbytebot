package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
)

func (b *baitbot) isNewJoke(ctx context.Context, joke string) (string, error) {
	hash, err := b.joker.GetJokeHash(joke)
	if err != nil {
		return "", err
	}

	if err := b.redis.Get(ctx, string(hash)).Err(); err != nil {
		return string(hash), nil
	}

	return "", nil
}

func (b *baitbot) formatJoke(joke *model.Joke) string {
	if joke.Delivery != "" {
		return fmt.Sprintf("%s\n\n||%s||", joke.Setup, joke.Delivery)
	}

	return joke.Setup
}

func (b *baitbot) escapeFormat(str string) string {
	str = strings.Replace(str, ".", `\.`, -1)
	str = strings.Replace(str, ",", `\,`, -1)
	str = strings.Replace(str, "!", `\!`, -1)
	str = strings.Replace(str, "?", `\?`, -1)
	str = strings.Replace(str, "[", `\[`, -1)
	str = strings.Replace(str, "]", `\]`, -1)
	str = strings.Replace(str, "-", `\-`, -1)
	str = strings.Replace(str, "=", `\=`, -1)
	str = strings.Replace(str, "+", `\+`, -1)
	str = strings.Replace(str, ";", `\;`, -1)
	str = strings.Replace(str, ":", `\:`, -1)
	str = strings.Replace(str, "'", `\'`, -1)
	str = strings.Replace(str, ")", `\)`, -1)
	str = strings.Replace(str, "(", `\(`, -1)
	str = strings.Replace(str, "*", `\*`, -1)
	str = strings.Replace(str, "^", `\^`, -1)
	str = strings.Replace(str, ">", `\>`, -1)
	str = strings.Replace(str, "<", `\<`, -1)

	return str
}
