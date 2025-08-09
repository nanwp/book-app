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
	Create(ctx context.Context, bookData *book.Book) error
	GetByID(ctx context.Context, id int64) (*book.Book, error)
	GetAll(ctx context.Context) ([]book.Book, error)
	Update(ctx context.Context, bookData *book.Book) error
	Delete(ctx context.Context, id int64) error
}

type Handler struct {
	Service BookService
}

func (h *Handler) CreateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request book.Book
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}

		if err := h.Service.Create(r.Context(), &request); err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}

		helper.WriteResponse(w, r, nil, nil)
	}
}

func (h *Handler) GetBookByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs := mux.Vars(r)

		idStr := rs["id"]
		idInt, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}

		bookData, err := h.Service.GetByID(r.Context(), idInt)
		if err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}

		if bookData == nil {
			helper.WriteResponse(w, r, nil, "Book not found")
			return
		}

		helper.WriteResponse(w, r, nil, bookData)
	}
}

func (h *Handler) GetAllBooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := h.Service.GetAll(r.Context())
		if err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}

		helper.WriteResponse(w, r, nil, books)
	}
}

func (h *Handler) UpdateBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs := mux.Vars(r)
		idStr := rs["id"]
		idInt, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}
		var request book.Book
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}

		request.ID = idInt

		if err := h.Service.Update(r.Context(), &request); err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}
	}
}

func (h *Handler) DeleteBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs := mux.Vars(r)
		id := rs["id"]
		if id == "" {
			helper.WriteResponse(w, r, helper.NewErrBadRequest("id is required"), nil)
			return
		}

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}

		if err := h.Service.Delete(r.Context(), idInt); err != nil {
			helper.WriteResponse(w, r, err, nil)
			return
		}

		helper.WriteResponse(w, r, nil, "Book deleted successfully")
	}
}
