package baitbot

import (
	"github.com/I0HuKc/baitbotnotbytebot/internal/core"
)

type action struct {
	name string
}

func (a action) GetName() string {
	return a.name
}

func CreateAction(n string) core.Action {
	return &action{
		name: n,
	}
}
