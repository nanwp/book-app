package services

import (
	"byfood-interview/book"
	"byfood-interview/helper"
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
)

type BookRepository interface {
	Create(ctx context.Context, bookData *book.Book) error
	GetByID(ctx context.Context, id int64) (*book.Book, error)
	GetAll(ctx context.Context) ([]book.Book, error)
	Update(ctx context.Context, bookData *book.Book) error
	Delete(ctx context.Context, id int64) error
}

type Book struct {
	BookRepository BookRepository
}

func (s *Book) Create(ctx context.Context, bookData *book.Book) error {
	log := log.Ctx(ctx).With().Str("service", "book").Logger()

	if err := bookData.Validate(); err != nil {
		log.Error().Err(err).Msg("invalid book data")
		return helper.NewErrBadRequest(err.Error())
	}

	if err := s.BookRepository.Create(ctx, bookData); err != nil {
		log.Error().Err(err).Msg("failed to create book")
		return err
	}

	return nil
}

func (s *Book) GetByID(ctx context.Context, id int64) (*book.Book, error) {
	log := log.Ctx(ctx).With().Str("service", "book").Logger()

	data, err := s.BookRepository.GetByID(ctx, id)

	if err != nil {
		log.Error().Err(err).Msg("failed to get book by ID")
		if err == sql.ErrNoRows {
			return nil, helper.NewErrNotFound("book not found")
		}
		return nil, err
	}

	return data, nil
}

func (s *Book) GetAll(ctx context.Context) ([]book.Book, error) {
	log := log.Ctx(ctx).With().Str("service", "book").Logger()

	books, err := s.BookRepository.GetAll(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all books")
		return nil, err
	}

	return books, nil
}

func (s *Book) Update(ctx context.Context, bookData *book.Book) error {
	log := log.Ctx(ctx).With().Str("service", "book").Logger()

	bookExisting, err := s.BookRepository.GetByID(ctx, bookData.ID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get existing book for update")
		if err == sql.ErrNoRows {
			return helper.NewErrNotFound("book not found")
		}

		return err
	}

	if bookData.Title != "" {
		bookExisting.Title = bookData.Title
	}
	if bookData.Author != "" {
		bookExisting.Author = bookData.Author
	}
	if bookData.PublishedYear != 0 {
		bookExisting.PublishedYear = bookData.PublishedYear
	}

	if err := s.BookRepository.Update(ctx, bookExisting); err != nil {
		log.Error().Err(err).Msg("failed to update book")
		return err
	}

	return nil
}

func (s *Book) Delete(ctx context.Context, id int64) error {
	log := log.Ctx(ctx).With().Str("service", "book").Logger()

	if err := s.BookRepository.Delete(ctx, id); err != nil {
		log.Error().Err(err).Msg("failed to delete book")
		return err
	}

	return nil
}
