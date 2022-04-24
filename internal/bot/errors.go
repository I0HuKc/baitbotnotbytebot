package baitbot

import "errors"

var (
	ErrNotAvailableFoYou   = errors.New("this command is not available for you")
	ErrOnlyForPrivateChat  = errors.New("this command can only be used in private chat")
	ErrNotSupportedInGroup = errors.New("this command is not supported in group chats")
	ErrAlreadyCreated      = errors.New("already created")
	ErrNoAntreCreated      = errors.New("antre is not created in this group")
)
