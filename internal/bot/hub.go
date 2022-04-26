package baitbot

import (
	"context"
	"errors"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type hub map[core.Command][]core.Handler

func (h hub) SetHandleFunc(cmd core.Command, handler ...core.Handler) {
	h[cmd] = handler
}

func (h hub) GetCommand(update tgbotapi.Update) core.Command {
	for cmd := range h {
		if cmd.IsThisCommand(update.Message.Command()) {
			return cmd
		}
	}

	return nil
}

func (h hub) IsExistingCommand(update tgbotapi.Update) ([]core.Handler, bool) {
	for cmd, handlers := range h {
		if cmd.IsThisCommand(update.Message.Command()) {
			return handlers, true
		}
	}

	return nil, false
}

func (h hub) HandleFunc(ctx context.Context, update tgbotapi.Update) error {
	if arr, ok := h.IsExistingCommand(update); ok {
		for _, handler := range arr {
			if err := handler(ctx, update); err != nil {
				return err
			}
		}

		return nil
	}

	return errors.New("unsupported command")
}
