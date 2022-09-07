package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/oddlid/go2lunch/server"
	"github.com/rs/zerolog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger := zerolog.New(os.Stdout)
	s := server.LunchServer{
		Log:       logger,
		LunchList: &lunchdata.LunchList{},
	}

	if err := s.Start(); err != nil {
		logger.Error().Err(err).Msg("Failed to start Lunch Server")
		return
	}

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := s.Stop(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Server failed to shut down cleanly")
	}
}
