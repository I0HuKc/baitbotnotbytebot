package core

import (
	"strings"
)

var _ Command = (*command)(nil)

type command struct {
	name    string
	private bool
}

type Command interface {
	GetName() string
	IsThisCommand(msg string) bool
	IsPrivate() bool
}

func (c *command) GetName() string {
	return c.name
}

func (c *command) IsThisCommand(msg string) bool {
	if strings.TrimSpace(msg) == c.name {
		return true
	}

	return false
}

func (c *command) IsPrivate() bool {
	return c.private
}

var (
	// Старт бота
	CommandStart = command{
		name:    "start",
		private: true,
	}

	// Старт процесса смены описания
	CommandStatChangeDesc = command{
		name:    "scd",
		private: false,
	}

	CommandAntre = command{
		name:    "antre",
		private: false,
	}

	CommandStropAntre = command{
		name:    "sa",
		private: false,
	}

	// Начать буллить человека
	CommandBullying = command{
		name:    "bll",
		private: false,
	}

	// Принудительное изменение статуса в группе
	CommandChangeDesc = command{
		name:    "cd",
		private: false,
	}

	// Добавить новый статус в библиотеку
	CommandAddDesc = command{
		name:    "ad",
		private: true,
	}

	// Получить статус
	CommandGetDesc = command{
		name:    "gd",
		private: true,
	}

	// Получить все описания из БД
	CommandGetAllDesc = command{
		name:    "gad",
		private: true,
	}

	// Проверить, работает ли бот
	CommandPing = command{
		name:    "ping",
		private: false,
	}

	CommandHelp = command{
		name:    "help",
		private: true,
	}

	CommandJoke = command{
		name:    "joke",
		private: false,
	}
)
