package redis

import (
	"order/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.REDIS_HOST + ":" + cfg.REDIS_PORT,
	})

}
