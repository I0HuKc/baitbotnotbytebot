package baitbot

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
	"github.com/I0HuKc/baitbotnotbytebot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *baitbot) NewAction(update tgbotapi.Update, act core.Action, expiration time.Duration) {
	b.acts[update.Message.Chat.ID] = ActionSendJoke

	go func() {
		time.Sleep(expiration)
		if _, ok := b.acts[update.Message.Chat.ID]; ok {
			delete(b.acts, update.Message.Chat.ID)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Действие `%s` отменено.", act.GetName()))
			msg.ParseMode = tgbotapi.ModeMarkdown
			b.Send(b.botApi.Send, msg)
		}
	}()
}

func (b *baitbot) isNewJoke(ctx context.Context, joke string) (string, error) {
	hash, err := b.joker.GetJokeHash(joke)
	if err != nil {
		return "", err
	}

	if _, err := b.redis.Get(ctx, string(hash)); err != nil {
		return string(hash), nil
	}

	return "", nil
}

func (b *baitbot) formatJoke(joke *model.Joke) string {
	if joke.Delivery != "" {
		return fmt.Sprintf("%s\n\n||%s||", joke.Setup, joke.Delivery)
	}

	return joke.Setup
}

func (b *baitbot) messageEscapeFormat(str string) string {
	str = strings.Replace(str, ".", `\.`, -1)
	str = strings.Replace(str, ",", `\,`, -1)
	str = strings.Replace(str, "!", `\!`, -1)
	str = strings.Replace(str, "?", `\?`, -1)
	str = strings.Replace(str, "[", `\[`, -1)
	str = strings.Replace(str, "]", `\]`, -1)
	str = strings.Replace(str, "-", `\-`, -1)
	str = strings.Replace(str, "=", `\=`, -1)
	str = strings.Replace(str, "+", `\+`, -1)
	str = strings.Replace(str, ";", `\;`, -1)
	str = strings.Replace(str, ":", `\:`, -1)
	str = strings.Replace(str, "'", `\'`, -1)
	str = strings.Replace(str, ")", `\)`, -1)
	str = strings.Replace(str, "(", `\(`, -1)
	str = strings.Replace(str, "*", `\*`, -1)
	str = strings.Replace(str, "^", `\^`, -1)
	str = strings.Replace(str, ">", `\>`, -1)
	str = strings.Replace(str, "<", `\<`, -1)

	return str
}

func (b *baitbot) IsLocal() bool {
	if os.Getenv("APP_ENV") == core.LocalEnv {
		return true
	}

	return false
}

// Вспомогательная функция, чтобы после отправки
// сообщения получить только ошибку без лишней информации
func (b *baitbot) Send(send func(c tgbotapi.Chattable) (tgbotapi.Message, error), msg tgbotapi.Chattable) error {
	_, err := send(msg)
	return err
}

// Отправить сообщение админу
func (b *baitbot) AdminNotify(about string) error {
	id, err := strconv.Atoi(os.Getenv("APP_BOT_ADMID_ID"))
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(int64(id), about)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.DisableNotification = true
	return b.Send(b.botApi.Send, msg)
}

// Верезать из сообщения только значение флага
func (b *baitbot) TrimFlagCommandValue(flag, text string) string {
	return strings.TrimLeft(strings.Split(text, flag)[1], " ")
}

// Метод валидации команды которая должна содержать флаг и значение флага
func (b *baitbot) CommandFlagValidation(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	if len(strings.Split(update.Message.Text, " ")) < 2 {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Необходимо указать флаг."),
			fmt.Errorf("[%s] — flag for command /%s isn't set",
				update.Message.From.UserName,
				update.Message.Command(),
			)

	}

	if len(strings.Split(update.Message.Text, " ")) < 3 {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите значение флага"),
			fmt.Errorf("[%s] — value for command /%s isn't set",
				update.Message.From.UserName,
				update.Message.Command(),
			)
	}

	return tgbotapi.MessageConfig{}, nil
}
