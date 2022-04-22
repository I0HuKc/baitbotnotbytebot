package baitbot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *baitbot) PrivateCmdHandler(ctx context.Context, update tgbotapi.Update) error {
	switch update.Message.Command() {

	/*
		/start
	*/
	case core.CommandStart.GetName():
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Зачем вообще нужно приветствие? Сразу к делу!\n%s", helpInfo))
		msg.ParseMode = tgbotapi.ModeMarkdown
		return b.Send(b.botApi.Send, msg)

	/*
		/ad
	*/
	case core.CommandAddDesc.GetName():
		if ok, err := b.IsAuthor(update); !ok {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не велено тебя сюда пускать")
			b.Send(b.botApi.Send, msg)

			return err
		}

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

	/*
		/gd
	*/
	case core.CommandGetDesc.GetName():
		if ok, err := b.IsAdmin(update); !ok {
			return err
		}

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

	/*
		/gad
	*/
	case core.CommandGetAllDesc.GetName():
		if ok, err := b.IsAdmin(update); !ok {
			return err
		}

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

	/*
		/help
	*/
	case core.CommandHelp.GetName():
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpInfo)
		msg.ParseMode = tgbotapi.ModeMarkdown
		return b.Send(b.botApi.Send, msg)

	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Данная команда не поддерживается :(")
		return b.Send(b.botApi.Send, msg)
	}
}
