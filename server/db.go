package server

import (
	internalDb "byfood-interview/internal/db"
	"byfood-interview/migration"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func db() *sqlx.DB {
	config := &internalDb.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	db, err := internalDb.NewDB(config)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	migration := migration.NewMigration(db)
	if err := migration.Run("./migration/file"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return db
}
