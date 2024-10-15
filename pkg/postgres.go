package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func NewPostgreSQLClient(migrationPath string, dbConfig Database) *sql.DB {
	startTime := time.Now()
	Logger.Info().
		Str(KeyProcess, "NewPostgreSQLClient").
		Dur(KeyProcessingTime, time.Since(startTime)).
		Time(KeyEndTime, time.Now()).
		Time(KeyStartTime, startTime).
		Msgf("initiate connection to database")
	defer func() {
		Logger.Info().
			Str(KeyProcess, "NewPostgreSQLClient").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Time(KeyEndTime, time.Now()).
			Time(KeyStartTime, startTime).
			Msgf("successed connecting to database")
	}()
	postgresUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		int(dbConfig.Port),
		dbConfig.DbName,
	)

	db, err := sql.Open("postgres", postgresUrl)
	if err != nil {
		Logger.Fatal().
			Err(err).
			Str(KeyProcess, "NewPostgreSQLClient").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Time(KeyEndTime, time.Now()).
			Time(KeyStartTime, startTime).
			Msgf("failed opening connection to postgres with error=%s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		Logger.Fatal().
			Err(err).
			Str(KeyProcess, "NewPostgreSQLClient").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Time(KeyEndTime, time.Now()).
			Time(KeyStartTime, startTime).
			Msgf("failed pinging connection to postgres with error=%s", err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		Logger.Fatal().
			Err(err).
			Str(KeyProcess, "NewPostgreSQLClient").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Time(KeyEndTime, time.Now()).
			Time(KeyStartTime, startTime).
			Msgf("failed creating postgres driver to do migration with error=%s", err.Error())
	}

	migration, err := migrate.NewWithDatabaseInstance(migrationPath, postgresUrl, driver)
	if err != nil {
		Logger.Fatal().
			Err(err).
			Str(KeyProcess, "NewPostgreSQLClient").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Time(KeyEndTime, time.Now()).
			Time(KeyStartTime, startTime).
			Msgf("failed migration postgres with error=%s", err.Error())
	}

	err = migration.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		Logger.Fatal().
			Err(err).
			Str(KeyProcess, "NewPostgreSQLClient").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Time(KeyEndTime, time.Now()).
			Time(KeyStartTime, startTime).
			Msgf("failed migration down postgres with error=%s", err.Error())
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		Logger.Fatal().
			Err(err).
			Str(KeyProcess, "NewPostgreSQLClient").
			Dur(KeyProcessingTime, time.Since(startTime)).
			Time(KeyEndTime, time.Now()).
			Time(KeyStartTime, startTime).
			Msgf("failed migration up postgres with error=%s", err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 15)
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetMaxOpenConns(int(dbConfig.MaxConnections))
	db.SetMaxIdleConns(int(dbConfig.MinConnections))

	return db
}
