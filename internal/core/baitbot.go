package core

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler func(ctx context.Context, update tgbotapi.Update) error

type Baitbot interface {
	Serve(ctx context.Context) (err error)
}
