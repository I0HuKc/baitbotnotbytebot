package baitbot

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
Command: /antre
Private: false
*/
func (b *baitbot) CommandAntreHandle(ctx context.Context, update tgbotapi.Update) error {
	for {
		pf := model.Performance{GroupId: int(update.Message.Chat.ID)}
		if err := b.store.Performance().GetByGroupId(ctx, &pf); err != nil {
			if err != sql.ErrNoRows {
				b.AdminNotify(err.Error())
				break
			}
		}

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
		b.Send(b.botApi.Send, msg)

		itr := rand.Intn(5-1+1) + 1
		nt := time.Now().UTC()

		fmt.Println(nt.Format("2006-01-02T15:04:05.00000000"))

		if !b.IsLocal() {
			pf.NextJoke = nt.Add(time.Hour * time.Duration(itr)).Format("2006-01-02T15:04:05.00000000")
			if err := b.store.Performance().Update(ctx, &pf); err != nil {
				return err
			}

			b.AdminNotify(fmt.Sprintf("Сделующая шутка через *%dч*", itr))
			time.Sleep(time.Duration(itr) * time.Hour)
			continue
		}

		pf.NextJoke = nt.Add(time.Second * time.Duration(itr)).Format("2006-01-02T15:04:05.00000000")
		if err := b.store.Performance().Update(ctx, &pf); err != nil {
			return err
		}

		b.AdminNotify(fmt.Sprintf("Сделующая шутка через *%dc*", itr))
		time.Sleep(time.Duration(itr) * time.Second)
	}

	return nil
}

/*
Command: /sa
Private: false
*/
func (b *baitbot) CommandStopAntreHandle(ctx context.Context, update tgbotapi.Update) error {
	if err := b.store.Performance().Delete(ctx,
		&model.Performance{GroupId: int(update.Message.Chat.ID)}); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я вне игры!")
	return b.Send(b.botApi.Send, msg)
}

/*
Command: /joke
Private: false
*/
func (b *baitbot) CommandJokeHandle(ctx context.Context, update tgbotapi.Update) error {
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
		return b.Send(b.botApi.Send, msg)

	}
}

/*
Command: /ping
Private: false
*/
func (b *baitbot) CommandPingHandle(ctx context.Context, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
	return b.Send(b.botApi.Send, msg)
}

/*
Command: /scd
Private: false
*/
func (b *baitbot) CommandStartChangeDescHandle(ctx context.Context, update tgbotapi.Update) error {
	for {
		rand.Seed(time.Now().UnixNano())

		// Получение ID первой и последней записи из БД
		min, max, err := b.store.Desc().FistLast(ctx)
		if err != nil {
			b.AdminNotify(err.Error())
		}

		if min < 1 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Строки закончились :)")
			b.botApi.Send(msg)
			break
		}

		// Достаю случайную запись из БД, в рамках допустимого диапазона
		d := model.Desc{Id: rand.Intn(max-min+1) + min}
		if err := b.store.Desc().Get(ctx, &d); err != nil {
			if err != sql.ErrNoRows {
				b.AdminNotify(err.Error())
			}

			continue
		}

		// Устанавливаю новое описание
		act := tgbotapi.NewChatDescription(update.Message.Chat.ID, d.Text)
		b.botApi.Send(act)

		// Отправляю уведомление о том, что новое описание установлено
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Он волен взять и поменять строку и с ней смысл темы всей!")
		b.botApi.Send(msg)

		// Удаляю установленный статус (чтобы он больше никогда не повторился)
		if err := b.store.Desc().Delete(ctx, &d); err != nil {
			b.AdminNotify(err.Error())
		}

		// Рандомный интервал через который будет установлен новый статус
		interval := rand.Intn(48-24+24) + 24

		if !b.IsLocal() {
			b.AdminNotify(fmt.Sprintf("Сделующая смена статуса через *%dч*", interval))
			time.Sleep(time.Duration(interval) * time.Hour)
			continue
		}

		b.AdminNotify(fmt.Sprintf("Сделующая смена статуса через *%dc*", interval))
		time.Sleep(time.Duration(interval) * time.Second)
	}

	return nil
}
