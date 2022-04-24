package core

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler func(ctx context.Context, update tgbotapi.Update) error
type HandlerFunc func(ctx context.Context, update tgbotapi.Update, h ...Handler) error

type Baitbot interface {
	Serve(ctx context.Context) (err error)
	SetHub() Baitbot
}

type Hub interface {
	// Go(h Handler) error
	HandleFunc(ctx context.Context, update tgbotapi.Update) error
	SetHandleFunc(cmd Command, h ...Handler)
	IsExistingCommand(update tgbotapi.Update) ([]Handler, bool)
}

type Command interface {
	GetName() string
	IsThisCommand(msg string) bool
}
