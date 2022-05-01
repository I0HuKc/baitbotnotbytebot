package core

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler func(ctx context.Context, update tgbotapi.Update) error
type HandlerFunc func(ctx context.Context, update tgbotapi.Update, h ...Handler) error

type Baitbot interface {
	Serve(ctx context.Context) (err error)
	Fuse() Baitbot
}

type Hub interface {
	Handle(ctx context.Context, update tgbotapi.Update) error
	SetHandleFunc(cmd Command, h ...Handler)
	IsExistingCommand(update tgbotapi.Update) ([]Handler, bool)
	GetCommand(update tgbotapi.Update) Command
}

type Command interface {
	GetName() string
	IsThisCommand(msg string) bool
	ReadFlag(text string) map[string]string
}

type Sked interface {
	Handle(ctx context.Context, update tgbotapi.Update, act Action) error
	SetHandleFunc(act Action, handler ...Handler)
	IsExistingAction(act Action) ([]Handler, bool)
}

type Action interface {
	GetName() string
}
