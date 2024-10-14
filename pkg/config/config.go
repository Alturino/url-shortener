package config

import (
	"time"

	"github.com/spf13/viper"

	"github.com/Alturino/url-shortener/pkg/log"
)

type Config struct {
	Env      string `mapstructure:"env"`
	Database `mapstructure:"db"`
}

type Database struct {
	Host           string `mapstructure:"host"`
	DbName         string `mapstructure:"name"`
	Password       string `mapstructure:"password"`
	Port           uint16 `mapstructure:"port"`
	Username       string `mapstructure:"username"`
	MaxConnections byte   `mapstructure:"max_connections"`
	MinConnections byte   `mapstructure:"min_connections"`
	TimeZone       string `mapstructure:"timezone"`
}

func InitConfig(filename string) Config {
	startTime := time.Now()

	config := Config{}
	log.Logger.Info().
		Str(log.KeyProcess, "InitConfig").
		Time(log.KeyStartTime, startTime).
		Msg("starting InitConfig")
	defer func() {
		log.Logger.Info().
			Str(log.KeyProcess, "InitConfig").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Time(log.KeyEndTime, time.Now()).
			Time(log.KeyStartTime, startTime).
			Interface("config", config).
			Msg("finished InitConfig")
	}()

	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Logger.Err(err).
			Str("filename", filename).
			Str(log.KeyProcess, "InitConfig").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Msgf("error when reading config with error=%s", err.Error())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Logger.Err(err).
			Str(log.KeyProcess, "InitConfig").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Msgf("error unmarshaling config with error=%s", err.Error())
	}

	return config
}
