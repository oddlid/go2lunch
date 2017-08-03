package main

import (
	"fmt"
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

func writePid(filename string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "%d", os.Getpid())

	return nil
}

func entryPointServe(ctx *cli.Context) error {
	pidfile := ctx.String("writepid")
	if pidfile != "" {
		err := writePid(pidfile)
		if err != nil {
			return cli.NewExitError(err.Error(), 5)
		}
		log.Infof("Wrote PID ( %d ) to %q", os.Getpid(), pidfile)
	}

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
			case syscall.SIGUSR1: // re-scrape and update internal DB
				err := update()
				if err != nil {
					log.Error(err.Error())
				}
			case syscall.SIGUSR2: // dump internal DB to stdout
				log.Info("Dumping parsed contents as JSON to STDOUT:")
				err := _site.s.Encode(os.Stdout)
				if err != nil {
					log.Error(err.Error())
				}
			default:
				log.Info("Caught unhandled signal, exiting...")
				os.Exit(255)
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

func notifyPid(pid int) error {
	return syscall.Kill(pid, syscall.SIGUSR1)
}

func entryPointScrape(ctx *cli.Context) error {
	pid := ctx.Int("notify-pid")
	if pid > 0 {
		err := notifyPid(pid)
		if err != nil {
			return cli.NewExitError(err.Error(), 3)
		} else {
			log.Infof("Told PID %d to re-scrape", pid)
			return nil
		}
	}

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
				cli.IntFlag{
					Name:  "notify-pid, p",
					Usage: "Make `PID` re-scrape instead of doing it in this process",
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
