package config

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/Alturino/url-shortener/internal/log"
)

type Config struct {
	Env         string `mapstructure:"env"`
	Database    `mapstructure:"db"`
	Application `mapstructure:"application"`
}

type Application struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Database struct {
	Host           string `mapstructure:"host"`
	DbName         string `mapstructure:"name"`
	Password       string `mapstructure:"password"`
	Username       string `mapstructure:"username"`
	MigrationPath  string `mapstructure:"migration_path"`
	TimeZone       string `mapstructure:"timezone"`
	Port           uint16 `mapstructure:"port"`
	MaxConnections byte   `mapstructure:"max_connections"`
	MinConnections byte   `mapstructure:"min_connections"`
}

func InitConfig(filename string, logger *zerolog.Logger) Config {
	startTime := time.Now()

	config := Config{}
	logger.Info().
		Str(log.KeyProcess, "InitConfig").
		Time(log.KeyStartTime, startTime).
		Msg("starting InitConfig")
	defer func() {
		logger.Info().
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
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Err(err).
			Str("filename", filename).
			Str(log.KeyProcess, "InitConfig").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Msgf("error when reading config with error=%s", err.Error())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Err(err).
			Str(log.KeyProcess, "InitConfig").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Msgf("error unmarshaling config with error=%s", err.Error())
	}

	return config
}