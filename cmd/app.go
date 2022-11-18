package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/oddlid/go2lunch/server"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

var (
	Version  string // to be set by the linker
	Compiled string // to be set by the linker
)

func showLogLevels() string {
	levels := []string{
		zerolog.TraceLevel.String(),
		zerolog.DebugLevel.String(),
		zerolog.InfoLevel.String(),
		zerolog.WarnLevel.String(),
		zerolog.ErrorLevel.String(),
		zerolog.PanicLevel.String(),
		zerolog.FatalLevel.String(),
		zerolog.Disabled.String(),
	}
	return strings.Join(levels, ", ")
}

func getCompileTime() time.Time {
	ts, err := time.Parse(time.RFC3339, Compiled)
	if err != nil {
		return time.Time{}
	}
	return ts
}

func logSetup(cCtx *cli.Context) error {
	zerolog.TimeFieldFormat = logTimeStampLayout
	level, err := zerolog.ParseLevel(cCtx.String(optLogLevel))
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)
	return nil
}

func newApp() *cli.App {
	return &cli.App{
		Version:   Version,
		Compiled:  getCompileTime(),
		Copyright: "(C) 2017 Odd Eivind Ebbesen",
		Authors: []*cli.Author{
			{
				Name:  "Odd E. Ebbesen",
				Email: "oddebb@gmail.com",
			},
		},
		Usage: "Scrape and serve lunch menus from configured sites",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  optLogLevel,
				Usage: fmt.Sprintf("Log `level` (options: %s)", showLogLevels()),
				Value: zerolog.InfoLevel.String(),
			},
		},
		Before: logSetup,
		Commands: []*cli.Command{
			{
				Action: actionScrape,
				Name:   cmdScrape,
				Usage:  "Scrape and output data only",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: optOutput,
					},
				},
			},
			{
				Action: actionServe,
				Name:   cmdServe,
				Usage:  "Start lunch server with optional background scraping",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  optBindHost,
						Usage: "Bind to `host`",
					},
					&cli.UintFlag{
						Name:  optPort,
						Usage: "Listen on `port`",
						Value: server.DefaultPort,
					},
					&cli.StringFlag{
						Name:  optCron,
						Usage: "Cron `spec` for background scraping",
					},
					&cli.StringFlag{
						Name:  optGTag,
						Usage: "`Tag` for Google analytics",
					},
					&cli.PathFlag{
						Name:  optLoad,
						Usage: "Load (initial) content from JSON `file`",
					},
				},
			},
		},
	}
}
