package main

import (
	"github.com/Alturino/url-shortener/pkg/config"
	"github.com/Alturino/url-shortener/pkg/log"
	"github.com/Alturino/url-shortener/pkg/postgres"
)

func main() {
	log.InitLogger()

	appConfig := config.InitConfig("application")
	postgreClient := postgres.NewPostgreSQLClient(appConfig.MigrationPath, appConfig.Database)
}
