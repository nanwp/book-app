package server

import (
	"byfood-interview/book/handler"
	"byfood-interview/book/services"
	"byfood-interview/book/stores"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type BookHandler interface {
	CreateBook() http.HandlerFunc
	GetBookByID() http.HandlerFunc
	GetAllBooks() http.HandlerFunc
	UpdateBook() http.HandlerFunc
	DeleteBook() http.HandlerFunc
}

type Server struct {
	Router *mux.Router
	DB     *sqlx.DB

	BookHandler BookHandler
}

func NewServer(migrationPath string) *Server {
	fmt.Println("Initializing server...")

	db := db(migrationPath)

	bookService := services.Book{
		BookRepository: stores.NewBook(db),
	}

	srv := &Server{
		Router:      mux.NewRouter(),
		DB:          db,
		BookHandler: &handler.Handler{Service: &bookService},
	}

	srv.routes()

	return srv
}

func (s *Server) Run(ctx context.Context, port string) error {
	if port == "" {
		port = "8080" // default port
	}

	httpS := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: s.cors().Handler(s.Router),
	}

	fmt.Println("Server started")

	log.Info().Msgf("server serving on port %s ", port)

	go func() {
		if err := httpS.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen:%+s\n", err)
		}
	}()

	<-ctx.Done()

	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		cancel()
	}()

	err := httpS.Shutdown(ctxShutDown)
	if err != nil {
		log.Fatal().Msgf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	if err := s.DB.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close database connection")
	}

	return err
}

func (s *Server) cors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		MaxAge:             60, // 1 minutes
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              false,
	})
}
