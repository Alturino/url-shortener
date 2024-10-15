package pkg

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string `mapstructure:"env"`
	Database `mapstructure:"db"`
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

func InitConfig(filename string) Config {
	startTime := time.Now()

	config := Config{}
	Logger.Info().
		Str(KeyProcess, "InitConfig").
		Time(KeyStartTime, startTime).
		Msg("starting InitConfig")
	defer func() {
		Logger.Info().
			Str(KeyProcess, "InitConfig").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Time(KeyEndTime, time.Now()).
			Time(KeyStartTime, startTime).
			Interface("config", config).
			Msg("finished InitConfig")
	}()

	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		Logger.Err(err).
			Str("filename", filename).
			Str(KeyProcess, "InitConfig").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Msgf("error when reading config with error=%s", err.Error())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		Logger.Err(err).
			Str(KeyProcess, "InitConfig").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Msgf("error unmarshaling config with error=%s", err.Error())
	}

	return config
}
