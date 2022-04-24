package baitbot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
Command: /start
Private: true
*/
func (b *baitbot) CommandStartHandle(ctx context.Context, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("Зачем вообще нужно приветствие? Сразу к делу!\n%s", helpInfo))
	msg.ParseMode = tgbotapi.ModeMarkdown
	return b.Send(b.botApi.Send, msg)
}

/*
Command: /gd
Private: true
*/
func (b *baitbot) CommandGetDescHandle(ctx context.Context, update tgbotapi.Update) error {
	if msg, err := b.CommandFlagValidation(update); err != nil {
		return b.Send(b.botApi.Send, msg)
	}

	id, err := strconv.Atoi(b.TrimFlagCommandValue("-id", update.Message.Text))
	if err != nil {
		return err
	}

	desc := model.Desc{Id: id}
	if err := b.store.Desc().Get(ctx, &desc); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, desc.Text)
	return b.Send(b.botApi.Send, msg)
}

/*
Command: /ad
Private: true;
*/
func (b *baitbot) CommandAddDescHandle(ctx context.Context, update tgbotapi.Update) error {
	if msg, err := b.CommandFlagValidation(update); err != nil {
		return b.Send(b.botApi.Send, msg)
	}

	if err := b.store.Desc().Create(ctx, &model.Desc{
		Text: b.TrimFlagCommandValue("-v", update.Message.Text),
	}); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Статус успешно добавлен.")
	return b.Send(b.botApi.Send, msg)
}

/*
Command: /gad
Private: true
*/
func (b *baitbot) CommandGetDescListHandle(ctx context.Context, update tgbotapi.Update) error {
	// Получение всех записей
	arr, err := b.store.Desc().List(ctx)
	if err != nil {
		return err
	}

	// Подготовка сообщения
	var mtext string
	for _, d := range arr {
		mtext += fmt.Sprintf("*%d*. %s\n\n", d.Id, d.Text)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, mtext)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return b.Send(b.botApi.Send, msg)
}

/*
Command: /help
Private: true
*/
func (b *baitbot) CommandHelpDescHandle(ctx context.Context, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpInfo)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return b.Send(b.botApi.Send, msg)
}
