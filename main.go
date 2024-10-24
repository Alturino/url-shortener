package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/Alturino/url-shortener/internal/config"
	"github.com/Alturino/url-shortener/internal/controller"
	"github.com/Alturino/url-shortener/internal/database"
	"github.com/Alturino/url-shortener/internal/log"
	"github.com/Alturino/url-shortener/internal/middleware"
	"github.com/Alturino/url-shortener/internal/repository"
	"github.com/Alturino/url-shortener/internal/service"
)

func main() {
	startTime := time.Now()
	logger := log.InitLogger()

	logger.Info().
		Str(log.KeyProcess, "main").
		Time(log.KeyStartTime, startTime).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Msg("initializing config")
	appConfig := config.InitConfig("application", logger)
	logger.Info().
		Str(log.KeyProcess, "main").
		Time(log.KeyStartTime, startTime).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Any("config", appConfig).
		Msg("initialized config")

	logger.Info().
		Str(log.KeyProcess, "main").
		Time(log.KeyStartTime, startTime).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Any("config", appConfig).
		Msg("initializing postgresql client")
	db := database.NewPostgreSQLClient(appConfig.MigrationPath, appConfig.Database, logger)
	logger.Info().
		Str(log.KeyProcess, "main").
		Time(log.KeyStartTime, startTime).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Any("config", appConfig).
		Msg("initialized postgresql client")

	logger.Info().
		Str(log.KeyProcess, "main").
		Time(log.KeyStartTime, startTime).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Any("config", appConfig).
		Msg("initializing urlService")
	queries := repository.New(db)
	encoder := base64.StdEncoding
	urlService := service.NewUrlService(db, queries, encoder)
	logger.Info().
		Str(log.KeyProcess, "main").
		Time(log.KeyStartTime, startTime).
		Any("config", appConfig).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Msg("initialized urlService")

	mux := http.NewServeMux()
	middlewares := middleware.CreateStack(middleware.Logging)
	controller.AttachUrlController(mux, urlService)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", appConfig.Application.Host, appConfig.Application.Port),
		Handler: middlewares(mux),
	}

	logger.Info().
		Str(log.KeyProcess, "main").
		Time(log.KeyStartTime, startTime).
		Any("config", appConfig).
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Msgf("listening to address=%s", server.Addr)
	server.ListenAndServe()
}
