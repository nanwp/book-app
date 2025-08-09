package main

import (
	"byfood-interview/server"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
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

	api.Run(ctx, 8080)
}
