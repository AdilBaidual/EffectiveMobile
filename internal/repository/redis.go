package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisCache interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}

func NewRedisClient(options *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(options)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
