package baitbot

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *baitbot) IsLocal() bool {
	if os.Getenv("APP_ENV") == "local" {
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
	return b.Send(b.botApi.Send, msg)
}

// Проверка на админа
func (b *baitbot) IsAdmin(update tgbotapi.Update) (bool, error) {
	id, err := strconv.Atoi(os.Getenv("APP_BOT_ADMID_ID"))
	if err != nil {
		return false, b.AdminNotify(
			fmt.Sprintf(
				"[%s] — попытка доступа к защищенным командам!\n[error]: %s",
				update.Message.Chat.UserName,
				err.Error(),
			),
		)
	}

	if update.Message.From.ID == int64(id) {
		return true, nil
	}

	return false, b.AdminNotify(
		fmt.Sprintf(
			"[%s] — попытка доступа к защищенным командам!",
			update.Message.Chat.UserName,
		),
	)
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
