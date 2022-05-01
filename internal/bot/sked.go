package baitbot

import (
	"context"
	"errors"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type sked map[core.Action][]core.Handler

func (s sked) SetHandleFunc(act core.Action, handler ...core.Handler) {
	s[act] = handler
}

func (s sked) Handle(ctx context.Context, update tgbotapi.Update, act core.Action) error {
	if arr, ok := s.IsExistingAction(act); ok {
		for _, handler := range arr {
			if err := handler(ctx, update); err != nil {
				return err
			}
		}

		return nil
	}

	return errors.New("unsupported action")
}

func (s sked) IsExistingAction(act core.Action) ([]core.Handler, bool) {
	for a, handlers := range s {
		if act.GetName() == a.GetName() {
			return handlers, true
		}
	}

	return nil, false
}
