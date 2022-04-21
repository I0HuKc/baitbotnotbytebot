package baitbot

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *baitbot) GroupCmdHandler(ctx context.Context, update tgbotapi.Update) error {
	switch update.Message.Command() {

	case core.CommandGetEvilinsultJoke.GetName():
		joke, err := b.joker.Joke(ctx)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, b.joker.FormatJoke(joke))
		msg.ParseMode = tgbotapi.ModeMarkdown
		return b.Send(b.botApi.Send, msg)

	// Обработка команды /ping
	case core.CommandPing.GetName():
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
		return b.Send(b.botApi.Send, msg)

	// Обработка команды /bll
	case core.CommandBullying.GetName():
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сорян, я еще не умею буллить :(")
		return b.Send(b.botApi.Send, msg)

	// Обработка команды /scd
	case core.CommandStatChangeDesc.GetName():
		if ok, err := b.IsAdmin(update); !ok {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ты не Егор")
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

func (b *baitbot) PrivateCmdHandler(ctx context.Context, update tgbotapi.Update) error {
	switch update.Message.Command() {

	// Обработка команды /start
	case core.CommandStart.GetName():
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Зачем вообще нужно приветствие? Сразу к делу!\n%s", helpInfo))
		msg.ParseMode = tgbotapi.ModeMarkdown
		return b.Send(b.botApi.Send, msg)

	// Обработка команды /ad
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

	// Обработка команды /gd
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

	case core.CommandHelp.GetName():
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpInfo)
		msg.ParseMode = tgbotapi.ModeMarkdown
		return b.Send(b.botApi.Send, msg)

	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Данная команда не поддерживается :(")
		return b.Send(b.botApi.Send, msg)
	}
}
