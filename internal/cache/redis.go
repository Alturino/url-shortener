package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"

	"github.com/Alturino/url-shortener/internal/config"
	"github.com/Alturino/url-shortener/internal/log"
)

func NewCacheClient(
	c context.Context,
	config config.Cache,
) *redis.Client {
	logger := zerolog.Ctx(c).With().Str(log.KeyProcess, "main NewCacheClient").Logger()

	logger.Info().Msg("intializing redis client")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       0,
	})
	logger.Info().Msg("initialized redis client")

	logger.Info().Msg("pinging redis client")
	err := redisClient.Ping(c).Err()
	if err != nil {
		err = fmt.Errorf("failed pinging redis client with error=%w", err)
		logger.Fatal().Err(err).Msg(err.Error())
	}
	logger.Info().Msg("pinged redis client")

	logger.Info().Msg("attach instrumentation to redis client")
	err = redisotel.InstrumentTracing(redisClient, redisotel.WithAttributes(semconv.DBSystemRedis))
	if err != nil {
		err = fmt.Errorf("failed attaching instrumentation to redis client with error=%w", err)
		logger.Fatal().Err(err).Msg(err.Error())
	}
	logger.Info().Msg("attach instrumentation to redis client")

	logger.Info().Msg("successed connecting to redis")

	return redisClient
}
