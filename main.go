package main

import "github.com/Alturino/url-shortener/pkg"

func main() {
	pkg.InitLogger()

	appConfig := pkg.InitConfig("application")
	pkg.NewPostgreSQLClient(appConfig.MigrationPath, appConfig.Database)
}
