package db

import (
	"context"
	"fmt"

	"github.com/geoo115/property-manager/config"
	"github.com/geoo115/property-manager/logger"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.LogInfo("Connected to Redis", logrus.Fields{
		"addr": cfg.Redis.Addr,
		"db":   cfg.Redis.DB,
	})

	return nil
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
