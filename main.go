package main

import (
	"github.com/Alturino/url-shortener/pkg/config"
	"github.com/Alturino/url-shortener/pkg/log"
)

func main() {
	log.InitLogger()

	config.InitConfig("application")
}
