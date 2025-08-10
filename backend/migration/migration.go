package migration

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Migration struct {
	db *sqlx.DB
}

func NewMigration(db *sqlx.DB) *Migration {
	return &Migration{db: db}
}

func (m *Migration) Run(path string) error {
	driver, err := postgres.WithInstance(m.db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	mig, err := migrate.NewWithDatabaseInstance(
		"file://"+path,
		"postgres", driver)

	if err != nil {
		return err
	}

	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Info().Msg("migrations applied successfully")

	return nil
}
