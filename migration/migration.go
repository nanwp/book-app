package migration

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
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

	return nil
}

// func runMigrations(db *sqlx.DB, migrationsPath string) error {
//     driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
//     if err != nil {
//         return fmt.Errorf("failed to create migrate db driver: %w", err)
//     }

//     m, err := migrate.NewWithDatabaseInstance(
//         "file://"+migrationsPath,
//         "postgres", driver)
//     if err != nil {
//         return fmt.Errorf("failed to init migrate: %w", err)
//     }

//     if err := m.Up(); err != nil && err != migrate.ErrNoChange {
//         return fmt.Errorf("failed to run migration: %w", err)
//     }
//     return nil
// }
