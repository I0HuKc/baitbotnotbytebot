package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func SetRedisConn(ctx context.Context) (*redis.Client, error) {
	db, err := strconv.Atoi(os.Getenv("APP_REDIS_EVIL_JOKE_DB"))
	if err != nil {
		return nil, err
	}

	ps, err := strconv.Atoi(os.Getenv("APP_REDIS_POOL_SIZE"))
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			os.Getenv("APP_REDIS_HOST"),
			os.Getenv("APP_REDIS_PORT"),
		),
		DB:              db,
		PoolSize:        ps,
		MinRetryBackoff: 1 * time.Second,
		MaxRetryBackoff: 2 * time.Second,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
