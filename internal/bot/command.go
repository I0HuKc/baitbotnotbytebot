package baitbot

import (
	"strings"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
)

type command struct {
	name string
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

func CreateCommand(n string) core.Command {
	return &command{
		name: n,
	}
}
