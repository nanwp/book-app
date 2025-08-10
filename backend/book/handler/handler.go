package handler

import (
	"byfood-interview/book"
	"byfood-interview/helper"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BookService interface {
	Create(ctx context.Context, bookData *book.Book) (*book.Book, error)
	GetByID(ctx context.Context, id int64) (*book.Book, error)
	GetAll(ctx context.Context) ([]book.Book, error)
	Update(ctx context.Context, bookData *book.Book) (*book.Book, error)
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	Service BookService
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with title, author, and published year
// @Tags books
// @Accept json
// @Produce json
// @Param book body book.Book true "Book data (without id)"
// @Success 200 {object} helper.Response{}
// @Failure 400 {object} helper.Response{errors=string}
// @Failure 500 {object} helper.Response{errors=string}
// @Router /api/v1/books [post]
// CreateBook handles the creation of a new book
func (h *Handler) CreateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request book.Book
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			helper.WriteResponse(w, helper.NewErrBadRequest("invalid JSON body"), nil)
			return
		}

		data, err := h.Service.Create(r.Context(), &request)
		if err != nil {
			helper.WriteResponse(w, err, nil)
			return
		}

		helper.WriteResponse(w, nil, data)
	}
}

// GetBookByID godoc
// @Summary Get a book by ID
// @Description Get a book by its ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} helper.Response{data=book.Book}
// @Failure 400 {object} helper.Response{errors=string}
// @Failure 404 {object} helper.Response{errors=string}
// @Failure 500 {object} helper.Response{errors=string}
// @Router /api/v1/books/{id} [get]
// GetBookByID handles fetching a book by its ID
func (h *Handler) GetBookByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs := mux.Vars(r)

		idStr := rs["id"]
		idInt, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			helper.WriteResponse(w, helper.NewErrBadRequest("invalid book ID"), nil)
			return
		}

		bookData, err := h.Service.GetByID(r.Context(), idInt)
		if err != nil {
			helper.WriteResponse(w, err, nil)
			return
		}

		if bookData == nil {
			helper.WriteResponse(w, nil, "Book not found")
			return
		}

		helper.WriteResponse(w, nil, bookData)
	}
}

// GetAllBooks godoc
// @Summary Get all books
// @Description Get a list of all books
// @Tags books
// @Produce json
// @Success 200 {object} helper.Response{data=[]book.Book}
// @Failure 500 {object} helper.Response{errors=string}
// @Router /api/v1/books [get]
// GetAllBooks handles fetching all books
func (h *Handler) GetAllBooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := h.Service.GetAll(r.Context())
		if err != nil {
			helper.WriteResponse(w, err, nil)
			return
		}

		helper.WriteResponse(w, nil, books)
	}
}

// UpdateBook godoc
// @Summary Update a book by ID
// @Description Update a book's details by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body book.Book true "Updated book data (with ID)"
// @Success 200 {object} helper.Response{}
// @Failure 400 {object} helper.Response{errors=string}
// @Failure 404 {object} helper.Response{errors=string}
// @Failure 500 {object} helper.Response{errors=string}
// @Router /api/v1/books/{id} [put]
// UpdateBook handles updating a book's details
func (h *Handler) UpdateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs := mux.Vars(r)
		idStr := rs["id"]
		idInt, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			helper.WriteResponse(w, err, nil)
			return
		}
		var request book.Book
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			helper.WriteResponse(w, err, nil)
			return
		}

		request.ID = idInt

		data, err := h.Service.Update(r.Context(), &request)
		if err != nil {
			helper.WriteResponse(w, err, nil)
			return
		}

		helper.WriteResponse(w, nil, data)
	}
}

// DeleteBook godoc
// @Summary Delete a book by ID
// @Description Delete a book by its ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} helper.Response{}
// @Failure 400 {object} helper.Response{errors=string}
// @Failure 404 {object} helper.Response{errors=string}
// @Failure 500 {object} helper.Response{errors=string}
// @Router /api/v1/books/{id} [delete]
// DeleteBook handles deleting a book by its ID
func (h *Handler) DeleteBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs := mux.Vars(r)
		id := rs["id"]
		if id == "" {
			helper.WriteResponse(w, helper.NewErrBadRequest("id is required"), nil)
			return
		}

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			helper.WriteResponse(w, err, nil)
			return
		}

		if err := h.Service.Delete(r.Context(), idInt); err != nil {
			helper.WriteResponse(w, err, nil)
			return
		}

		helper.WriteResponse(w, nil, "Book deleted successfully")
	}
}
