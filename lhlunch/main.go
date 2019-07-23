package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"os/signal"
	"strconv"
	"sync"
	//"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
)

const (
	VERSION          string = "2019-07-23"
	DEF_URL          string = "https://www.lindholmen.se/pa-omradet/dagens-lunch"
	DEF_ADR          string = ":20666"
	DEF_COUNTRY_NAME string = "Sweden"
	DEF_COUNTRY_ID   string = "se"
	DEF_CITY_NAME    string = "Gothenburg"
	DEF_CITY_ID      string = "gbg"
	DEF_SITE_NAME    string = "Lindholmen"
	DEF_SITE_ID      string = "lindholmen"
	DEF_ID           string = "/se/gbg/lindholmen"
	DEF_COMMENT      string = "Gruvan"
	GTAG_ID          string = "UA-126840341-2"
)

// exit codes
const (
	E_OK int = iota
	E_UPDATE
	E_READPID
	E_WRITEPID
	E_NOTIFYPID
	E_WRITEJSON
	E_READJSON
	E_INITTMPL
	E_WRITEHTML
)

var BUILD_DATE string
var COMMIT_ID string

type LHSite struct {
	sync.Mutex
	//s   *lunchdata.Site
	ll  *lunchdata.LunchList
	url string
}

var _site *LHSite

func init() {
	defaultSite()
}

func defaultSite() {
	lh := lunchdata.NewSite(DEF_SITE_NAME, DEF_SITE_ID, DEF_COMMENT, GTAG_ID)
	gbg := lunchdata.NewCity(DEF_CITY_NAME, DEF_CITY_ID, GTAG_ID)
	sthlm := lunchdata.NewCity("Stockholm", "sthlm", GTAG_ID)
	se := lunchdata.NewCountry(DEF_COUNTRY_NAME, DEF_COUNTRY_ID, GTAG_ID)
	no := lunchdata.NewCountry("Norway", "no", GTAG_ID)


	llist := lunchdata.NewLunchList(GTAG_ID)

	gbg.AddSite(*lh)
	se.AddCity(*gbg)
	se.AddCity(*sthlm)

	llist.AddCountry(*se)
	llist.AddCountry(*no)

	_site = &LHSite{
		url: DEF_URL,
		//s:   lhsite,
		ll:  llist,
	}
}

func siteFromJSON(filename string) error {
	ll, err := lunchdata.LunchListFromFile(filename)
	if err != nil {
		log.Errorf("Unable to load site from JSON: %q", err.Error())
		return err
	}
	//_site.s = s // replace default
	//_site.ll.Countries[0].Cities[0].Sites[0] = *s
	_site.ll = ll
	return nil
}

func (lhs *LHSite) getLHSite() *lunchdata.Site {
	se := lhs.ll.GetCountryById(DEF_COUNTRY_ID)
	if nil == se {
		return nil
	}

	gbg := se.GetCityById(DEF_CITY_ID)
	if nil == gbg {
		return nil
	}

	lh := gbg.GetSiteById(DEF_SITE_ID)

	return lh // might be nil
}

func (lhs *LHSite) setLHRestaurants(rs lunchdata.Restaurants) {
	lh := lhs.getLHSite()
	if nil != lh {
		lh.SetRestaurants(rs)
	}
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

	log.Debugf("LHLunch PID: %d", os.Getpid())

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
	// On windows this only logs a message about skipping signal handling
	// See: signalhandling.go and signalhandling_windows.go
	setupSignalHandling()
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

	// handle "load" argument here before serving
	jfile := ctx.String("load")
	if jfile != "" {
		err := siteFromJSON(jfile)
		if err != nil {
			return cli.NewExitError(err.Error(), E_READJSON)
		}
		log.Debugf("Load site from %q successful!", jfile)
	}

	// temporary fix for struct migration
	doDump := ctx.Bool("dump")
	if doDump {
		//_site.ll.Encode(os.Stdout)
		_site.ll.GetSiteLinks().Encode(os.Stdout)
		return nil
	}

	log.Infof("Listening on: %s", adr)
	return http.ListenAndServe(adr, setupRouter())
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

	// get content
	err := update()
	if err != nil {
		return cli.NewExitError(err.Error(), E_UPDATE)
	}

//	// write html output if requested, otherwise json
//	dumpHtml := ctx.Bool("html")
//	if dumpHtml {
//		err := initTmpl() // does more than we need here, so maybe rewrite some time...
//		if err != nil {
//			return cli.NewExitError(err.Error(), E_INITTMPL)
//		}
//		tmpl_lhlunch_html.Execute(os.Stdout, _site.getLHSite()) // TODO: Rethink and replace
//		return nil // be sure to not proceed when done here
//	}

	outfile := ctx.String("outfile")
	log.Debugf("Outfile: %q", outfile)

	if outfile == "-" || outfile == "" {
		err := _site.ll.Encode(os.Stdout)
		if err != nil {
			return cli.NewExitError(err.Error(), E_WRITEJSON)
		}
	} else {
		err := _site.ll.SaveJSON(outfile)
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
				cli.StringFlag{
					Name:  "load",
					Usage: "Load data from `JSONFILE` instead of scraping",
				},
				cli.BoolFlag{
					Name:  "dump",
					Usage: "Dump new struct to stdout",
				},
			},
		},
		{
			Name:    "scrape",
			Aliases: []string{"scr"},
			Usage:   "Scrape source and output JSON or HTML, then exit",
			Action:  entryPointScrape,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "outfile, o",
					Usage: "Write JSON result to `FILE` ('-' for STDOUT)",
					Value: "-",
				},
				cli.BoolFlag{
					Name:  "html",
					Usage: "Write HTML result to STDOUT",
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
