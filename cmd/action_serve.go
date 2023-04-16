package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/oddlid/go2lunch/lunchdata"
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
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	var lunchList lunchdata.LunchList

	if loadPath := cCtx.Path(optLoad); loadPath != "" {
		data, err := os.ReadFile(loadPath)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, &lunchList); err != nil {
			return err
		}
		logger.Debug().Str("file", loadPath).Msg("Loaded menus from file")
	} else {
		lunchList = getEmptyLunchList()
	}

	s := server.LunchServer{
		Log:       logger,
		LunchList: lunchList,
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
