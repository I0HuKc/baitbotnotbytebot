package baitbot

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	gt "github.com/bas24/googletranslatefree"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *baitbot) GroupCmdHandler(ctx context.Context, update tgbotapi.Update) error {
	switch update.Message.Command() {

	/*
		/sa
	*/
	case core.CommandStropAntre.GetName():
		if ok, err := b.IsAdmin(update); !ok {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Как же %s ты хитёр", update.Message.Chat.FirstName),
			)
			b.botApi.Send(msg)

			return err
		}

		b.antre <- true

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я вне игры!")
		return b.Send(b.botApi.Send, msg)

	/*
		/antre
	*/
	case core.CommandAntre.GetName():
		if ok, err := b.IsAdmin(update); !ok {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Как же %s ты хитёр", update.Message.Chat.FirstName),
			)
			b.botApi.Send(msg)

			return err
		}

		b.antre = make(chan bool)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Джокер в игре!")
		if err := b.Send(b.botApi.Send, msg); err != nil {
			return err
		}

		go func(ctx context.Context) {
			for {
				select {
				case <-b.antre:
					return
				default:
					fmt.Println(1)
				}

				joke, err := b.joker.Joke(ctx)
				if err != nil {
					b.AdminNotify(err.Error())
				}

				jtext := b.formatJoke(joke)
				if joke.Lang.Translate {
					jtext, err = gt.Translate(jtext, joke.Lang.Source, joke.Lang.Target)
					if err != nil {
						b.AdminNotify(err.Error())
					}
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, b.messageEscapeFormat(jtext))
				msg.ParseMode = "MarkdownV2"
				b.Send(b.botApi.Send, msg)

				interval := rand.Intn(5-1+1) + 1

				if !b.IsLocal() {
					b.AdminNotify(fmt.Sprintf("Сделующая шутка через *%dч*", interval))
					time.Sleep(time.Duration(interval) * time.Hour)
					return
				}

				b.AdminNotify(fmt.Sprintf("Сделующая шутка через *%dc*", interval))
				time.Sleep(time.Duration(interval) * time.Second)
			}

		}(ctx)

		return nil

	/*
		/joke
	*/
	case core.CommandJoke.GetName():
		joke, err := b.joker.Joke(ctx)
		if err != nil {
			return err
		}

		jtext := b.formatJoke(joke)
		if joke.Lang.Translate {
			jtext, err = gt.Translate(jtext, joke.Lang.Source, joke.Lang.Target)
			if err != nil {
				return err
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, b.messageEscapeFormat(jtext))
		msg.ParseMode = "MarkdownV2"
		return b.Send(b.botApi.Send, msg)

	/*
		/ping
	*/
	case core.CommandPing.GetName():
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
		return b.Send(b.botApi.Send, msg)

	/*
		/bll
	*/
	case core.CommandBullying.GetName():
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сорян, я еще не умею буллить :(")
		return b.Send(b.botApi.Send, msg)

	/*
		/scd
	*/
	case core.CommandStatChangeDesc.GetName():
		if ok, err := b.IsAdmin(update); !ok {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Как же %s ты хитёр", update.Message.Chat.FirstName),
			)
			b.botApi.Send(msg)

			return err
		}

		go func() {
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

					return
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
					return
				}

				b.AdminNotify(fmt.Sprintf("Сделующая смена статуса через *%dc*", interval))
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}()

		return nil
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Данная команда не поддерживается :(")
		return b.Send(b.botApi.Send, msg)
	}
}
