package server

import (
	internalDb "byfood-interview/internal/db"
	"byfood-interview/migration"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func db(migrationPath string) *sqlx.DB {
	config := &internalDb.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	db, err := internalDb.NewDB(config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	migration := migration.NewMigration(db)
	if err := migration.Run(migrationPath); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	return db
}
