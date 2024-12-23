package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/Alturino/url-shortener/internal/cache"
	"github.com/Alturino/url-shortener/internal/config"
	"github.com/Alturino/url-shortener/internal/controller"
	"github.com/Alturino/url-shortener/internal/database"
	"github.com/Alturino/url-shortener/internal/log"
	"github.com/Alturino/url-shortener/internal/middleware"
	"github.com/Alturino/url-shortener/internal/repository"
	"github.com/Alturino/url-shortener/internal/service"
)

func main() {
	c := context.Background()

	logger := log.InitLogger()
	c = logger.WithContext(c)

	c, stop := signal.NotifyContext(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer func() {
		logger.Info().
			Str(log.KeyProcess, "main").
			Msg("Received SIGINT or SIGKILL shutting down")
		stop()
		logger.Info().
			Str(log.KeyProcess, "main").
			Msg("shutdown")
	}()

	logger.Info().
		Str(log.KeyProcess, "main").
		Msg("initializing otelsdk")
	otelShutdown, err := log.InitOtelSdk(c)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str(log.KeyProcess, "main").
			Msgf("failed initialized otelsdk with error=%s", err.Error())
	}
	logger.Info().
		Str(log.KeyProcess, "main").
		Msg("initalized otelsdk")
	defer func() {
		logger.Info().Str(log.KeyProcess, "main").Msgf("shutting down otelsdk")
		err := otelShutdown(c)
		if err != nil {
			logger.Fatal().
				Err(err).
				Str(log.KeyProcess, "main").
				Msgf("failed shutdown otelsdk with error=%s", err.Error())
		}
		logger.Info().Str(log.KeyProcess, "main").Msgf("shutdown otelsdk")
	}()

	logger.Info().
		Str(log.KeyProcess, "main").
		Msg("initializing config")
	appConfig := config.InitConfig("application", logger)
	logger.Info().
		Str(log.KeyProcess, "main").
		Any(log.KeyConfig, appConfig).
		Msg("initialized config")

	logger.Info().
		Str(log.KeyProcess, "main").
		Any(log.KeyConfig, appConfig).
		Msg("initializing postgresql client")
	db := database.NewPostgreSQLClient(appConfig.MigrationPath, appConfig.Database, logger)
	logger.Info().
		Str(log.KeyProcess, "main").
		Any(log.KeyConfig, appConfig).
		Msg("initialized postgresql client")

	logger.Info().
		Str(log.KeyProcess, "main").
		Any(log.KeyConfig, appConfig).
		Msg("initializing redis client")
	redis := cache.NewCacheClient(c, appConfig.Cache)
	logger.Info().
		Str(log.KeyProcess, "main").
		Any(log.KeyConfig, appConfig).
		Msg("initialized redis client")

	logger.Info().
		Str(log.KeyProcess, "main").
		Any(log.KeyConfig, appConfig).
		Msg("initializing urlService")
	queries := repository.New(db)
	encoder := base64.StdEncoding
	urlService := service.NewUrlService(redis, db, encoder, queries)
	logger.Info().
		Str(log.KeyProcess, "main").
		Any(log.KeyConfig, appConfig).
		Msg("initialized urlService")

	mux := http.NewServeMux()
	middlewares := middleware.CreateStack(middleware.Logging, middleware.Otlp)
	otelhttpHandler := otelhttp.NewHandler(
		middlewares(mux),
		"url-shortener",
	)
	controller.AttachUrlController(mux, urlService)

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%d", appConfig.Application.Host, appConfig.Application.Port),
		Handler:      otelhttpHandler,
		BaseContext:  func(net.Listener) context.Context { return c },
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	srvErr := make(chan error, 1)
	go func() {
		logger.Info().
			Str(log.KeyProcess, "main").
			Any(log.KeyConfig, appConfig).
			Msgf("listening to address=%s", server.Addr)
		srvErr <- server.ListenAndServe()
	}()

	select {
	case err := <-srvErr:
		logger.Fatal().
			Err(err).
			Str(log.KeyProcess, "main").
			Any(log.KeyConfig, appConfig).
			Msgf("ListenAndServe with error=%s", err.Error())
	case <-c.Done():
		logger.Info().
			Str(log.KeyProcess, "main").
			Any(log.KeyConfig, appConfig).
			Msg("shutting down server")
		err := server.Shutdown(c)
		if err != nil {
			logger.Fatal().
				Err(err).
				Str(log.KeyProcess, "main").
				Any(log.KeyConfig, appConfig).
				Msgf("failed shutting down server with error=%s", err.Error())
		}
		stop()
		logger.Info().
			Str(log.KeyProcess, "main").
			Any(log.KeyConfig, appConfig).
			Msg("shutdown server")
	}
}
