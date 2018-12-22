package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/go2lunch/site"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
)

const (
	VERSION string = "2018-12-21"
	DEF_URL string = "https://www.lindholmen.se/pa-omradet/dagens-lunch"
	DEF_ADR string = ":20666"
)

// exit codes
const (
	E_OK int = iota
	E_UPDATE
	E_READPID
	E_WRITEPID
	E_NOTIFYPID
	E_WRITEJSON
	E_INITTMPL
)

var BUILD_DATE string
var COMMIT_ID string

type LHSite struct {
	sync.Mutex
	s   *site.Site
	url string
}

var _site *LHSite

func init() {
	_site = &LHSite{s: &site.Site{Name: "Lindholmen", ID: "se/gbg/lindholmen", Comment: "Gruvan"}, url: DEF_URL}
}

func writePid(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "%d", os.Getpid())

	return nil
}

func setUrl(ctx *cli.Context) {
	url := ctx.String("url")
	if url != "" {
		_site.url = url
	}
}

func entryPointServe(ctx *cli.Context) error {
	setUrl(ctx)

	log.Infof("LHLunch PID: %d", os.Getpid())

	pidfile := ctx.String("writepid")
	if pidfile != "" {
		log.Debugf("Got pidfile arg: %q", pidfile)
		err := writePid(pidfile)
		if err != nil {
			return cli.NewExitError(err.Error(), E_WRITEPID)
		}
		log.Infof("Wrote PID ( %d ) to %q", os.Getpid(), pidfile)
	} else {
		log.Debugf("No PIDfile given")
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
				log.Debug("Caught unhandled signal, exiting...")
				os.Exit(255)
			}
		}
	}()
	// END signal handling

	// cron-like scheduling.
	// Turns out app/docker hangs sometimes, when triggering scrape from regular cron with "docker exec" (even with timeout),
	// so trying to use an internal solution instead
	cronspec := ctx.String("cron")
	if cronspec != "" {
		log.Infof("Auto-updating via built-in cron @ %q", cronspec)
		cr := cron.New()

		err := cr.AddFunc(cronspec, func() {
			log.Info("Re-scraping on request from internal cron...")
			if err := update(); err != nil {
				log.Errorf("Update via internal cron failed: %q", err)
			}
		})

		if err != nil {
			log.Errorf("Failed to add cronjob: %q", err)
		} else {
			cr.Start()
			log.Info("Built-in cron started")
		}
	}
	// END cron

	err := initTmpl()
	if err != nil {
		return cli.NewExitError(err.Error(), E_INITTMPL)
	}
	return http.ListenAndServe(adr, setupRouter())
}

func notifyPid(pid int) error {
	return syscall.Kill(pid, syscall.SIGUSR1)
}

func readPid(filename string) (int, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return -1, err
	}
	pid, err := strconv.Atoi(string(bytes.TrimSpace(b)))
	if err != nil {
		return -1, fmt.Errorf("Error parsing PID from %q: %s", filename, err.Error())
	}
	return pid, nil
}

func entryPointScrape(ctx *cli.Context) error {
	setUrl(ctx)

	pidf := ctx.String("notify-pid")
	if pidf != "" {
		pid, err := readPid(pidf)
		if err != nil {
			return cli.NewExitError(err.Error(), E_READPID)
		}
		if pid > 0 {
			err := notifyPid(pid)
			if err != nil {
				return cli.NewExitError(err.Error(), E_NOTIFYPID)
			} else {
				log.Infof("Told PID %d to re-scrape", pid)
				return nil
			}
		}
	}

	outfile := ctx.String("outfile")
	log.Debugf("Outfile: %q", outfile)

	err := update()
	if err != nil {
		return cli.NewExitError(err.Error(), E_UPDATE)
	}

	if outfile == "-" || outfile == "" {
		err := _site.s.Encode(os.Stdout)
		if err != nil {
			return cli.NewExitError(err.Error(), E_WRITEJSON)
		}
	} else {
		err := _site.s.SaveJSON(outfile)
		if err != nil {
			return cli.NewExitError(err.Error(), E_WRITEJSON)
		}
		log.Debugf("Wrote JSON result to %q", outfile)
	}
	return nil
}

// setCustomAppHelpTmpl slightly changes the help text to include BUILD_DATE
// See https://github.com/urfave/cli/blob/master/help.go
func setCustomAppHelpTmpl() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION / BUILD_DATE:
   {{.Version}} / {{.Compiled}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
`
}

func main() {
	setCustomAppHelpTmpl()
	app := cli.NewApp()
	app.Name = "Lindholmen Lunch Scraper/Server"
	app.Version = fmt.Sprintf("%s_%s", VERSION, COMMIT_ID)
	app.Copyright = "(c) 2017 Odd Eivind Ebbesen"
	app.Compiled, _ = time.Parse(time.RFC3339, BUILD_DATE)
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
				cli.StringFlag{
					Name:  "cron",
					Usage: "Specify intervals for re-scrape in cron format with `TIMESPEC`",
					// what I had in regular cron: "0 0,30 08-12 * * 1-5"
					// That is: on second 0 of minute 0 and 30 of hour 08-12 of weekday mon-fri on any day of month any month
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
				cli.StringFlag{
					Name:  "notify-pid, p",
					Usage: "Read PID from `FILE` and tell the process with that PID to re-scrape",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url, u",
			Usage: "`URL` to scrape",
			Value: DEF_URL,
		},
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
