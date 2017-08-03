package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/go2lunch/site"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	VERSION string = "2017-08-03"
	SRC_URL string = "https://www.lindholmen.se/pa-omradet/dagens-lunch"
	DEF_ADR string = ":20666"
)

type LHSite struct {
	sync.Mutex
	s *site.Site
}

var _site *LHSite

func init() {
	_site = &LHSite{s: &site.Site{Name: "Lindholmen", ID: "se/gbg/lindholmen", Comment: "Gruvan"}}
}

func entryPointServe(ctx *cli.Context) error {
	adr := ctx.String("listen-adr")
	if adr == "" {
		adr = DEF_ADR
	}

	// signal handling
	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for sig := range sig_chan {
			switch sig {
			case syscall.SIGUSR1:
				err := update()
				if err != nil {
					log.Error(err.Error())
				}
			case syscall.SIGUSR2:
				log.Info("Dumping parsed contents as JSON to STDOUT:")
				err := _site.s.Encode(os.Stdout)
				if err != nil {
					log.Error(err.Error())
				}
			default:
				log.Info("Caught unhandled signal, exiting...")
				os.Exit(0)
			}
		}
	}()
	// END signal handling

	log.Infof("LHLunch PID: %d", os.Getpid())

	http.HandleFunc("/", lhHandler)
	err := http.ListenAndServe(adr, nil)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func entryPointScrape(ctx *cli.Context) error {
	outfile := ctx.String("outfile")

	err := update()
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	if outfile == "-" || outfile == "" {
		err := _site.s.Encode(os.Stdout)
		if err != nil {
			return err
		}
	} else {
		err := _site.s.SaveJSON(outfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "Lindholmen Lunch Scraper/Server"
	app.Version = VERSION
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Odd E. Ebbesen",
			Email: "oddebb@gmail.com",
		},
	}
	app.Usage = "Scrape and/or serve results of todays lunch from lindholmen.se"

	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"srv"},
			Usage:   "Start server",
			Action:  entryPointServe,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen-adr, l",
					Value: DEF_ADR,
					Usage: "(hostname|IP):port to listen on",
				},
				cli.StringFlag{
					Name:  "writepid, p",
					Usage: "Write PID to `FILE`",
				},
			},
		},
		{
			Name:    "scrape",
			Aliases: []string{"scr"},
			Usage:   "Scrape source and output JSON, then exit",
			Action:  entryPointScrape,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "outfile, o",
					Usage: "Write JSON result to `FILE` ('-' for STDOUT)",
					Value: "-",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level, l",
			Value: "info",
			Usage: "Log `level` (options: debug, info, warn, error, fatal, panic)",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Run in debug mode",
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stderr)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatal(err.Error())
		}
		log.SetLevel(level)
		if !c.IsSet("log-level") && !c.IsSet("l") && c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: false,
			FullTimestamp:    true,
		})
		return nil
	}

	app.Run(os.Args)
}
