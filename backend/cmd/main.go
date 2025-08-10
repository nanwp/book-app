package main

import (
	"byfood-interview/server"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	configEnv "github.com/joho/godotenv"
)

// @title           Books API
// @version         1.0
// @description     Demo OpenAPI (Swagger) untuk net/http + gorilla/mux.
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.email  dev@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @BasePath  /

// @schemes http

func main() {
	err := configEnv.Load(".env")
	if err != nil {
		log.Warn().Err(err).Msg("no .env file, using real environment")
	}

	api := server.NewServer()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	go func() {
		oscall := <-ch
		log.Debug().Msgf("Received signal: %v", oscall)
		cancel()
	}()

	api.Run(ctx, os.Getenv("HTTP_PORT"))
}
