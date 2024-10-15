package main

import (
	"context"

	"github.com/Alturino/url-shortener/internal/repository"
	"github.com/Alturino/url-shortener/pkg"
)

func main() {
	pkg.InitLogger()

	appConfig := pkg.InitConfig("application")
	db := pkg.NewPostgreSQLClient(appConfig.MigrationPath, appConfig.Database)

	queries := repository.New(db)
	queries.GetAllUrls(context.Background())
}
