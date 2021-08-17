package config

import (
	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func init() {
	opt, err := redis.ParseURL("redis://redis-service:6379/0")
	if err != nil {
		panic(err)
	}

	redis := redis.NewClient(opt)
	Redis = redis
}