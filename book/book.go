package book

import (
	"errors"
	"time"
)

var (
	ErrBookNotFound          = errors.New("book not found")
	ErrBookExists            = errors.New("book already exists")
	ErrTitleRequired         = errors.New("title is required")
	ErrAuthorRequired        = errors.New("author is required")
	ErrPublishedYearRequired = errors.New("published year is required")
)

type Book struct {
	ID            int64      `json:"id" db:"id"`
	Title         string     `json:"title" db:"title"`
	Author        string     `json:"author" db:"author"`
	PublishedYear int        `json:"published_year" db:"published_year"`
	CreatedAt     time.Time  `json:"-" db:"created_at"`
	UpdatedAt     time.Time  `json:"-" db:"updated_at"`
	DeletedAt     *time.Time `json:"-" db:"deleted_at"`
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return ErrTitleRequired
	}
	if b.Author == "" {
		return ErrAuthorRequired
	}
	if b.PublishedYear <= 0 {
		return ErrPublishedYearRequired
	}
	return nil
}
