package baitbot

import (
	"fmt"
	"strings"

	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
)

type command struct {
	name  string
	flags []string
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

func (c *command) ReadFlag(text string) map[string]string {
	m := make(map[string]string)
	for i, f := range c.flags {
		w := strings.Split(text, f)[1]
		if len(c.flags) > i+1 {
			w = strings.Split(w, c.flags[i+1])[0]
		}
		m[f] = w
	}

	fmt.Println(m)
	return m
}

func CreateCommand(n string, f ...string) core.Command {
	return &command{
		name:  n,
		flags: f,
	}
}
