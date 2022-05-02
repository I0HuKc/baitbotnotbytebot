package db

import (
	"fmt"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

func SetMemcacheConn() (*memcache.Client, error) {
	mc := memcache.New(fmt.Sprintf("%s:%s",
		os.Getenv("APP_MEMCACHE_HOST"), os.Getenv("APP_MEMCACHE_PORT")))

	if err := mc.Ping(); err != nil {
		return nil, err
	}

	return mc, nil
}
