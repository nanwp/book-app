package services

import (
	"byfood-interview/book"
	"byfood-interview/helper"
	"context"
	"database/sql"
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
	if err := bookData.Validate(); err != nil {
		return helper.NewErrBadRequest(err.Error())
	}

	if err := s.BookRepository.Create(ctx, bookData); err != nil {
		return err
	}

	return nil
}

func (s *Book) GetByID(ctx context.Context, id int64) (*book.Book, error) {
	data, err := s.BookRepository.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, helper.NewErrNotFound("book not found")
		}
		return nil, err
	}

	return data, nil
}

func (s *Book) GetAll(ctx context.Context) ([]book.Book, error) {
	books, err := s.BookRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *Book) Update(ctx context.Context, bookData *book.Book) error {
	bookExisting, err := s.BookRepository.GetByID(ctx, bookData.ID)
	if err != nil {
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

	return s.BookRepository.Update(ctx, bookExisting)
}

func (s *Book) Delete(ctx context.Context, id int64) error {
	if err := s.BookRepository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
