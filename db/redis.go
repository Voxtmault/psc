package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/rotisserie/eris"
	"github.com/voxtmault/psc/config"
)

var redisClient *redis.Client

func InitRedis(config *config.RedisConfig) error {

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       int(config.RedisDBNum),
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return eris.Wrap(err, "Init Redis")
	}

	slog.Info("Successfully opened redis connection")
	return nil
}

func CloseRedis() error {
	if err := redisClient.Close(); err != nil {
		return eris.Wrap(err, "Closing redis connection")
	}

	return nil
}

func GetRedisCon() *redis.Client {
	return redisClient
}
