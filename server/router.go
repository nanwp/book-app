package server

import (
	_ "byfood-interview/docs"
	"byfood-interview/process-url/handler"
	"byfood-interview/server/middleware"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) routes() {
	s.Router.Use(middleware.Logger)
	s.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := s.Router.PathPrefix("/api/v1/").Subrouter()

	api.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	// book routes
	api.HandleFunc("/books", s.BookHandler.CreateBook()).Methods(http.MethodPost)
	api.HandleFunc("/books/{id}", s.BookHandler.GetBookByID()).Methods(http.MethodGet)
	api.HandleFunc("/books", s.BookHandler.GetAllBooks()).Methods(http.MethodGet)
	api.HandleFunc("/books/{id}", s.BookHandler.UpdateBook()).Methods(http.MethodPut)
	api.HandleFunc("/books/{id}", s.BookHandler.DeleteBook()).Methods(http.MethodDelete)

	// URL cleanup routes
	api.HandleFunc("/process-url", handler.ProcessURLHandler()).Methods(http.MethodPost)
}
