package baitbot

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
Command: /sjoke
Private: true
*/
func (b *baitbot) CommandSendJokeHandle(ctx context.Context, update tgbotapi.Update) error {
	// Сохраняю действие
	b.NewAction(update, ActionSendJoke, time.Minute)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хорошо, пришли шутку которую необходимо отправить в беседу.")
	return b.Send(b.botApi.Send, msg)
}

/*
Command: /antre
Private: true
*/
func (b *baitbot) CommandAntreHandle(ctx context.Context, update tgbotapi.Update) error {
	for {
		joke, err := b.joker.Joke(ctx)
		if err != nil {
			b.AdminNotify(err.Error())
			continue
		}

		jtext := b.formatJoke(joke)
		hash, err := b.isNewJoke(ctx, jtext)
		if err != nil {
			continue
		}

		if err := b.redis.Save(ctx, hash, joke.Target, (24*time.Hour)*30); err != nil {
			b.AdminNotify(err.Error())
			continue
		}

		if joke.Lang.Translate {
			jtext, err = gt.Translate(jtext, joke.Lang.Source, joke.Lang.Target)
			if err != nil {
				b.AdminNotify(err.Error())
				continue
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, b.messageEscapeFormat(jtext))
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		msg.DisableNotification = true
		b.Send(b.botApi.Send, msg)

		itr := int(float32(10 * rand.Float32()))
		b.AdminNotify(fmt.Sprintf("Сделующая шутка через *%s*", time.Duration(itr)*core.GetInterval()))
		time.Sleep(time.Duration(itr) * core.GetInterval())
	}
}

/*
Command: /start
Private: true
*/
func (b *baitbot) CommandStartHandle(ctx context.Context, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Зачем вообще нужно приветствие? Сразу к делу!")
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
	// if msg, err := b.CommandFlagValidation(update); err != nil {
	// 	return b.Send(b.botApi.Send, msg)
	// }

	if err := b.store.Desc().Create(ctx, &model.Desc{
		Text: b.hub.GetCommand(update).ReadFlag(update.Message.Text)["-v"],
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
