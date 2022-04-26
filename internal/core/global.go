package core

import (
	"os"
	"time"
)

const (
	LocalEnv = "local"
	TestEnv  = "test"
	ProdEnv  = "prod"
)

func GetInterval() time.Duration {
	if os.Getenv("APP_ENV") != LocalEnv {
		return time.Hour
	}

	return time.Second
}
