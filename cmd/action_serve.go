package main

import (
	"context"
	"os"
	"time"

	"github.com/oddlid/go2lunch/server"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

func actionServe(cCtx *cli.Context) error {
	// If load param is given, we try to load a LunchList object from the given file.
	// If load fails, we use the default empty LunchList with the predefined structure.
	// If load param is not given, we use the default empty LunchList right away.
	// If cron param is given, we set up background scraping at the specified schedule.
	// If cron param is not given, content will be static with whatever we have for the
	// lunchlist.

	s := server.LunchServer{
		Log:       zerolog.New(os.Stdout).With().Timestamp().Logger(),
		LunchList: getEmptyLunchList(),
		Config:    server.DefaultConfig(),
	}

	if err := s.Start(cCtx.Context); err != nil {
		return err
	}

	<-cCtx.Done()

	subCtx, cancel := context.WithTimeout(cCtx.Context, 5*time.Second)
	defer cancel()

	return s.Stop(subCtx)
}
