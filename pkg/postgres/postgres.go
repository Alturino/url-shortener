package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/Alturino/url-shortener/pkg/config"
	"github.com/Alturino/url-shortener/pkg/log"
)

func NewPostgreSQLClient(migrationPath string, databaseConfig config.Database) *sql.DB {
	startTime := time.Now()
	log.Logger.Info().
		Str(log.KeyProcess, "NewPostgreSQLClient").
		Dur(log.KeyProcessingTime, time.Since(startTime)).
		Time(log.KeyEndTime, time.Now()).
		Time(log.KeyStartTime, startTime).
		Msgf("initiate connection to database")
	defer func() {
		log.Logger.Info().
			Str(log.KeyProcess, "NewPostgreSQLClient").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Time(log.KeyEndTime, time.Now()).
			Time(log.KeyStartTime, startTime).
			Msgf("successed connecting to database")
	}()
	postgresUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		databaseConfig.Username,
		databaseConfig.Password,
		databaseConfig.Host,
		int(databaseConfig.Port),
		databaseConfig.DbName,
	)
	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		log.Logger.Fatal().
			Err(err).
			Str(log.KeyProcess, "NewPostgreSQLClient").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Time(log.KeyEndTime, time.Now()).
			Time(log.KeyStartTime, startTime).
			Msgf("failed opening connection to postgres with error=%s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Logger.Fatal().
			Err(err).
			Str(log.KeyProcess, "NewPostgreSQLClient").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Time(log.KeyEndTime, time.Now()).
			Time(log.KeyStartTime, startTime).
			Msgf("failed pinging connection to postgres with error=%s", err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Logger.Fatal().
			Err(err).
			Str(log.KeyProcess, "NewPostgreSQLClient").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Time(log.KeyEndTime, time.Now()).
			Time(log.KeyStartTime, startTime).
			Msgf("failed creating postgres driver to do migration with error=%s", err.Error())
	}

	migration, err := migrate.NewWithDatabaseInstance(migrationPath, postgresUrl, driver)
	if err != nil {
		log.Logger.Fatal().
			Err(err).
			Str(log.KeyProcess, "NewPostgreSQLClient").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Time(log.KeyEndTime, time.Now()).
			Time(log.KeyStartTime, startTime).
			Msgf("failed migration postgres with error=%s", err.Error())
	}

	err = migration.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Logger.Fatal().
			Err(err).
			Str(log.KeyProcess, "NewPostgreSQLClient").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Time(log.KeyEndTime, time.Now()).
			Time(log.KeyStartTime, startTime).
			Msgf("failed migration down postgres with error=%s", err.Error())
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Logger.Fatal().
			Err(err).
			Str(log.KeyProcess, "NewPostgreSQLClient").
			Dur(log.KeyProcessingTime, time.Since(startTime)).
			Time(log.KeyEndTime, time.Now()).
			Time(log.KeyStartTime, startTime).
			Msgf("failed migration up postgres with error=%s", err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 15)
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetMaxOpenConns(int(databaseConfig.MaxConnections))
	db.SetMaxIdleConns(int(databaseConfig.MinConnections))

	return db
}
