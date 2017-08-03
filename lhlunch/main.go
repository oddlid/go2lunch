package main

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/go2lunch/site"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

const (
	VERSION string = "2017-08-03"
	DEF_URL string = "https://www.lindholmen.se/pa-omradet/dagens-lunch"
	DEF_ADR string = ":20666"
)

const (
	E_OK int = iota
	E_UPDATE
	E_READPID
	E_WRITEPID
	E_NOTIFYPID
)

var BUILD_TIME string

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

	setupMux()
	server := http.Server{
		Addr:    adr,
		Handler: &lhHandler{},
	}
	return server.ListenAndServe()
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

	outfile := ctx.String("outfile")

	err = update()
	if err != nil {
		return cli.NewExitError(err.Error(), E_UPDATE)
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
	app.Version = VERSION
	app.Copyright = "(c) 2017 Odd Eivind Ebbesen"
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
