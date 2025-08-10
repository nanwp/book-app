package stores

import (
	"byfood-interview/book"
	"byfood-interview/migration"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/lib/pq"
)

func postgresC(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "book",
			"POSTGRES_USER":     "book",
			"POSTGRES_PASSWORD": "book",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithPollInterval(1 * time.Second),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func dbFromContainer(ctx context.Context, container testcontainers.Container) (*sqlx.DB, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("%v %v", host, port.Port())

	dsn := fmt.Sprintf("postgres://book:book@%s:%s/book?sslmode=disable",
		host,
		port.Port())

	var db *sqlx.DB

	for i := 0; i < 5; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	migration := migration.NewMigration(db)
	if err := migration.Run("../../migration/file"); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Debug().Msg("Database successful connect")

	return db, nil
}

var (
	testContainer testcontainers.Container
	testDB        *sqlx.DB
)

func TestMain(m *testing.M) {
	ctx := context.TODO()
	var err error

	testContainer, err = postgresC(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create postgres container")
		os.Exit(1)
	}

	testDB, err = dbFromContainer(ctx, testContainer)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		os.Exit(1)
	}

	exitCode := m.Run()

	testDB.Close()
	testContainer.Terminate(ctx)
	log.Debug().Msg("TestMain completed successfully")
	os.Exit(exitCode)
}

func TestNewBook(t *testing.T) {
	ctx := context.TODO()

	bookStore := NewBook(testDB)

	book := book.Book{
		Title:         "Test Book",
		Author:        "Test Author",
		PublishedYear: 2023,
	}

	err := bookStore.Create(ctx, &book)
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}
}

func TestGetByID(t *testing.T) {
	ctx := context.TODO()

	bookStore := NewBook(testDB)

	book, err := bookStore.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("failed to get book by ID: %v", err)
	}

	if book == nil {
		t.Fatal("expected book to be not nil")
	}

	if book.ID != 1 {
		t.Errorf("expected book ID to be 1, got %d", book.ID)
	}
}

func TestGetAll(t *testing.T) {
	ctx := context.TODO()

	bookStore := NewBook(testDB)

	books, err := bookStore.GetAll(ctx)
	if err != nil {
		t.Fatalf("failed to get all books: %v", err)
	}

	if len(books) == 0 {
		t.Fatal("expected books to not be empty")
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.TODO()

	bookStore := NewBook(testDB)

	book := book.Book{
		ID:            1,
		Title:         "Updated Book",
		Author:        "Updated Author",
		PublishedYear: 2024,
	}

	err := bookStore.Update(ctx, &book)
	if err != nil {
		t.Fatalf("failed to update book: %v", err)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.TODO()

	bookStore := NewBook(testDB)

	err := bookStore.Delete(ctx, 1)
	if err != nil {
		t.Fatalf("failed to delete book: %v", err)
	}

	book, err := bookStore.GetByID(ctx, 1)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			t.Fatalf("unexpected error when getting book by ID after deletion: %v", err)
		}
	}

	if book != nil {
		t.Fatal("expected book to be nil after deletion")
	}
}
