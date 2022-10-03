package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	// logger := zerolog.New(os.Stdout)
	// s := server.LunchServer{
	// 	Log:       logger,
	// 	LunchList: getEmptyLunchList(logger),
	// 	Config:    server.DefaultConfig(),
	// }

	// if err := s.Start(ctx); err != nil {
	// 	logger.Error().Err(err).Msg("Failed to start Lunch Server")
	// 	return
	// }

	// <-ctx.Done()

	// shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer shutdownCancel()
	// if err := s.Stop(shutdownCtx); err != nil {
	// 	logger.Error().Err(err).Msg("Server failed to shut down cleanly")
	// }

	app := newApp()
	if err := app.RunContext(ctx, os.Args); err != nil {
		if !errors.Is(err, context.Canceled) {
			cancel()
			log.Fatal().Err(err).Msg("Execution failed")
		}
	}
}
