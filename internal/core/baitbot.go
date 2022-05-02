package core

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

type Baitbot interface {
	Configure(ctx context.Context) Baitbot
	Serve()
}

type Heap interface {
	Handle(endpoint string, hf tele.HandlerFunc)
	Get(endpoint string) tele.HandlerFunc
}
