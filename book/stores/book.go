package stores

import (
	"byfood-interview/book"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Book struct {
	db *sqlx.DB
}

func NewBook(db *sqlx.DB) *Book {
	return &Book{db: db}
}

func (b *Book) GetByID(ctx context.Context, id int64) (*book.Book, error) {
	var bookData book.Book
	query := "SELECT id, title, author, published_year FROM books WHERE id = $1 AND deleted_at IS NULL"
	err := b.db.GetContext(ctx, &bookData, query, id)
	if err != nil {
		return nil, err
	}
	return &bookData, nil
}

func (b *Book) GetAll(ctx context.Context) ([]book.Book, error) {
	var books []book.Book
	query := "SELECT id, title, author, published_year FROM books WHERE deleted_at IS NULL ORDER BY created_at DESC"
	err := b.db.SelectContext(ctx, &books, query)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (b *Book) Create(ctx context.Context, bookData *book.Book) error {
	query := "INSERT INTO books (title, author, published_year) VALUES ($1, $2, $3)"
	_, err := b.db.ExecContext(ctx, query, bookData.Title, bookData.Author, bookData.PublishedYear)
	if err != nil {
		log.Error().Err(err).Msg("failed to insert book")
		return err
	}

	return nil
}

func (b *Book) Update(ctx context.Context, bookData *book.Book) error {
	query := "UPDATE books SET title = $1, author = $2, published_year = $3, updated_at = NOW() WHERE id = $4"
	_, err := b.db.ExecContext(ctx, query, bookData.Title, bookData.Author, bookData.PublishedYear, bookData.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to update book")
		return err
	}

	return nil
}

func (b *Book) Delete(ctx context.Context, id int64) error {
	query := "UPDATE books SET deleted_at = NOW() WHERE id = $1"
	_, err := b.db.ExecContext(ctx, query, id)
	return err
}
