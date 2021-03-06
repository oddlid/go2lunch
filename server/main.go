package main

/*
2019-08-15:
Getting a lot closer to the original goal of having an engine able to serve menus for any site
in the world. There's still a lot of legacy in the code though, referring to Lindholmen instead
of thinking about what is configured, and showing just that.
TODO:
  - Separate all scraping into separate modules. Define an interface for use by scrapers that
	  are to be built in to the binary, so that the engine doesn't need to know anything about
		the scrapers except which site a scraper is for, and an API to run a scrape.
		Internal / built-in scrapers should run asynchronously, delivering results via channels.
		Internal / built-in scrapers should be able to register themselves with the engine, preferrably
		without the engine having to know anything about it, other than being able to tell all registered
		scrapers to run at given intervals, and then collect their results. Look at go-chat-bot for how
		to solve this.
		For external scrapers, all they need is to post the correct format to the correct URL while
		providing the correct API header key. How the ext scraper is implemented should not be relevant.

	- The above thoughts leaves the question: Should we have internal scrapers at all? Wouldn't it be
	  easier to just skip them and depend fully on external scraping? Advantages and disadvantages?
*/

import (
	//"bytes"
	"context"
	"fmt"
	//"io/ioutil"
	"net/http"
	"os"
	//"strconv"
	//"strings"
	"sync"
	"time"

	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	// Internal scrapers
	"github.com/oddlid/go2lunch/scraper/se/gbg/lindholmen"
)

const (
	DEF_WEB_ADR       string = ":20666"
	DEF_ADM_ADR       string = ":20667"
	DEF_READ_TIMEOUT         = 5
	DEF_WRITE_TIMEOUT        = 10
	DEF_IDLE_TIMEOUT         = 15

	//DEF_COUNTRY_ID    string = "se"
	//DEF_CITY_ID       string = "gbg"
	//DEF_SITE_ID       string = "lindholmen"
	//GTAG_ID           string = "UA-126840341-2" // used for Google Analytics in generated pages - TODO: replace!
	//DEF_URL           string = "https://www.lindholmen.se/pa-omradet/dagens-lunch"
	//DEF_COUNTRY_NAME  string = "Sweden"
	//DEF_CITY_NAME     string = "Gothenburg"
	//DEF_SITE_NAME     string = "Lindholmen"
	//DEF_ID            string = "se/gbg/lindholmen"
	//DEF_COMMENT       string = "Gruvan"
	//JWT_TOKEN_SECRET  string = "this secret should not be part of the code"
)

// exit codes
const (
	E_OK int = iota
	//E_UPDATE
	//E_READPID
	E_WRITEPID
	//E_NOTIFYPID
	//E_WRITEJSON
	E_READJSON
	E_INITTMPL
	//E_WRITEHTML
)

var (
	VERSION    string
	BUILD_DATE string
	COMMIT_ID  string
	_gtag      string
	_noScrape  bool
	_lunchList *lunchdata.LunchList
)

/*
Since we call registerSiteScraper() here from init, getLunchList() will get called in staticlunchlist.go
with _lunchlist being nil, resulting in initialization from the static JSON in that file.
If we at a point after that load content from JSON via a file, site.Scraper will get overwritten to nil,
and no further scraping will happen.
For now, that is acceptable, but it's not very obvious and could lead to hard to find bugs later on.
We should rather not use init() at all, but register scrapers at a later point.
*/
func init() {
	lhs := lindholmen.LHScraper{}
	registerSiteScraper(lhs.GetCountryID(), lhs.GetCityID(), lhs.GetSiteID(), lhs)
}

// I'd like to find a more flexible and dynamic way of including scrapers, but for now
// we'll use this
func registerSiteScraper(countryID, cityID, siteID string, scraper lunchdata.SiteScraper) {
	lsite := getLunchList().GetSiteById(countryID, cityID, siteID)
	if nil == lsite {
		log.WithFields(log.Fields{
			"countryID": countryID,
			"cityID":    cityID,
			"siteID":    siteID,
		}).Warn("Could not find site")
		return
	}
	lsite.Scraper = scraper
}

func lunchListFromFile(filename string) error {
	ll, err := lunchdata.LunchListFromFile(filename)
	if err != nil {
		log.WithFields(log.Fields{
			"ErrMSG": err.Error(),
		}).Error("Unable to load site from JSON file")
		return err
	}
	// If we load the LunchList from JSON, and we only have a top-level
	// Gtag, we need to propagate it now after load
	if ll.Gtag != "" {
		ll.PropagateGtag(ll.Gtag)
	}
	_lunchList = ll
	return nil
}

// hack
func logInventory() {
	//var b strings.Builder

	//fmt.Fprintf(
	//	&b,
	//	"Total subitems: %d, Countries: %d, Cities: %d, Sites: %d, Restaurants: %d, Dishes: %d\n",
	//	getLunchList().SubItems(),
	//	getLunchList().NumCountries(),
	//	getLunchList().NumCities(),
	//	getLunchList().NumSites(),
	//	getLunchList().NumRestaurants(),
	//	getLunchList().NumDishes(),
	//)

	//sls := getLunchList().GetSiteLinks()
	//for _, sl := range sls {
	//	fmt.Fprintf(&b, "%s: %s\n", sl.Url, sl.SiteKey)
	//}

	//log.Debug(b.String())

	log.WithFields(log.Fields{
		"Total subitems": getLunchList().SubItems(),
		"Countries":      getLunchList().NumCountries(),
		"Cities":         getLunchList().NumCities(),
		"Sites":          getLunchList().NumSites(),
		"Restaurants":    getLunchList().NumRestaurants(),
		"Dishes":         getLunchList().NumDishes(),
	}).Debug("Inventory")
}

// This might not be needed anymore when on next level
func writePid(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "%d", os.Getpid())

	return nil
}

func setGtag(ctx *cli.Context) {
	gtag := ctx.String("gtag")
	if gtag != "" {
		_gtag = gtag
		getLunchList().PropagateGtag(_gtag)
	}
}

func entryPointServe(ctx *cli.Context) error {
	pid := os.Getpid()

	log.WithFields(log.Fields{
		"PID": pid,
	}).Debug("Startup info")

	pidfile := ctx.String("writepid")
	if pidfile != "" {
		log.WithFields(log.Fields{
			"pidfile": pidfile,
		}).Debug("Got argument")
		err := writePid(pidfile)
		if err != nil {
			return cli.Exit(err.Error(), E_WRITEPID)
		}
		log.WithFields(log.Fields{
			"PID":     pid,
			"PIDFile": pidfile,
		}).Info("Wrote pidfile")
	} else {
		log.Debug("No PIDfile given")
	}

	// cron-like scheduling.
	// Turns out app/docker hangs sometimes, when triggering scrape from regular cron with "docker exec" (even with timeout),
	// so trying to use an internal solution instead
	// When post updates is fully ready, this should be removed..?
	cronspec := ctx.String("cron")
	if cronspec != "" {
		log.WithFields(log.Fields{
			"cronspec": cronspec,
		}).Info("Auto-updating via built-in cron")
		cr := cron.New()

		_, err := cr.AddFunc(cronspec, func() {
			log.Info("Re-scraping on request from internal cron...")
			var wg sync.WaitGroup
			getLunchList().RunSiteScrapers(&wg)
			wg.Wait()
			// we don't need to propagate _gtag here, as SiteScrapers can only set new Restaurants
			// for an already configured Site, so a given Gtag will not be removed by scraping
			logInventory()
		})

		if err != nil {
			log.WithFields(log.Fields{
				"ErrMSG": err.Error(),
			}).Error("Failed to add cronjob")
		} else {
			cr.Start()
			log.Info("Built-in cron started")
		}
	}
	// END cron

	err := initTmpl()
	if err != nil {
		return cli.Exit(err.Error(), E_INITTMPL)
	}

	// Loading the lunchlist from a JSON file is mostly useful for testing.
	// As the entire structure is replaced, every Site's Scraper instance will be set to nil,
	// and so no auto-scraping will happen, unless we re-register scrapers after JSON load.
	// Unless testing scraping, one can load a copy of PROD data in a local instance like this:
	//
	// $ go run . -d serve --load <(curl -sS https://lunch.oddware.net/json/)

	// handle "load" argument here before serving
	jfile := ctx.String("load")
	if jfile != "" {
		err := lunchListFromFile(jfile)
		if err != nil {
			return cli.Exit(err.Error(), E_READJSON)
		}
		log.WithFields(log.Fields{
			"JSONFile": jfile,
		}).Info("Load lunch list from file successful")
		log.Info("Auto-scraping is now disabled")
		//_noScrape = true
	}

	// Important that this call comes after anything that sets content, like lunchListFromFile above
	// We should probably make a hook that calls this after any update of content as well
	// If we did load from JSON, this gives us the possibility to override the Gtag from CLI
	setGtag(ctx)

	numServers := 2
	quit := make(chan bool, numServers)
	// signal handling
	// On windows this only logs a message about skipping signal handling
	// See: signalhandling.go and signalhandling_windows.go
	setupSignalHandling(quit, numServers)
	// END signal handling

	listenAdr := ctx.String("listen-adr")
	listenAdm := ctx.String("listen-adm")

	pubR, admR := setupRouter()
	pubSrv := createServer(listenAdr, pubR)
	admSrv := createServer(listenAdm, admR)

	var wg sync.WaitGroup
	wg.Add(numServers)
	go gracefulShutdown("PubSRV", pubSrv, quit, &wg)
	go gracefulShutdown("AdmSRV", admSrv, quit, &wg)
	go func() {
		ctxlog := log.WithFields(log.Fields{
			"Port": listenAdr,
		})
		ctxlog.Info("Public server listening")
		if err := pubSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctxlog.Error("Error starting public server")
		}
	}()
	go func() {
		ctxlog := log.WithFields(log.Fields{
			"Port": listenAdm,
		})
		ctxlog.Info("Admin server listening")
		if err := admSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ctxlog.Error("Error starting admin server")
		}
	}()

	// now, run registered scrapers
	// By running 2 levels of goroutines and using waitgroups, we can do all scraping in the background,
	// and still wait until all are done before displaying updated stats
	_noScrape = ctx.Bool("noscrape")
	if !_noScrape {
		go func() {
			var wg2 sync.WaitGroup
			getLunchList().RunSiteScrapers(&wg2) // each scraper runs in its own goroutine, incrementing wg2
			wg2.Wait()
			logInventory()
		}()
	}

	// Here we block until we get a signal to quit
	// Might be months until we reach the next line after this
	// Given that thought, maybe this would be a good place to display some uptime stats or something...
	wg.Wait()

	// If we did load contents from a file, let's save it back
	saveOnExit := ctx.Bool("save-on-exit")
	stripOnSave := ctx.Bool("strip-menus-on-save")
	if saveOnExit && jfile != "" {
		if stripOnSave {
			getLunchList().ClearRestaurants()
		}
		err := getLunchList().SaveJSON(jfile)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		log.WithFields(log.Fields{
			"JSONFile": jfile,
		}).Info("Wrote config")
	}

	return nil
}

func createServer(addr string, hnd http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      hnd,
		ReadTimeout:  DEF_READ_TIMEOUT * time.Second,
		WriteTimeout: DEF_WRITE_TIMEOUT * time.Second,
		IdleTimeout:  DEF_IDLE_TIMEOUT * time.Second,
	}
}

func gracefulShutdown(tag string, srv *http.Server, quit <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	<-quit
	log.WithFields(log.Fields{
		"Server": tag,
	}).Debug("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"Server": tag,
			"ErrMSG": err.Error(),
		}).Error("Error shutting down server")
	}
}

//func readPid(filename string) (int, error) {
//	b, err := ioutil.ReadFile(filename)
//	if err != nil {
//		return -1, err
//	}
//	pid, err := strconv.Atoi(string(bytes.TrimSpace(b)))
//	if err != nil {
//		return -1, fmt.Errorf("Error parsing PID from %q: %s", filename, err.Error())
//	}
//	return pid, nil
//}

//func entryPointScrape(ctx *cli.Context) error {
//	log.Debugf("Scrape sub-command deprecated, exiting")
//	return nil
//
//	// ***************** //
//
//	setUrl(ctx)
//
//	pidf := ctx.String("notify-pid")
//	if pidf != "" {
//		pid, err := readPid(pidf)
//		if err != nil {
//			return cli.Exit(err.Error(), E_READPID)
//		}
//		if pid > 0 {
//			err := notifyPid(pid)
//			if err != nil {
//				return cli.Exit(err.Error(), E_NOTIFYPID)
//			} else {
//				log.Infof("Told PID %d to re-scrape", pid)
//				return nil
//			}
//		}
//	}
//
//	// get content
//	err := update()
//	if err != nil {
//		return cli.Exit(err.Error(), E_UPDATE)
//	}
//
//	//	// write html output if requested, otherwise json
//	//	dumpHtml := ctx.Bool("html")
//	//	if dumpHtml {
//	//		err := initTmpl() // does more than we need here, so maybe rewrite some time...
//	//		if err != nil {
//	//			return cli.Exit(err.Error(), E_INITTMPL)
//	//		}
//	//		tmpl_lhlunch_html.Execute(os.Stdout, _site.getLHSite()) // TODO: Rethink and replace
//	//		return nil // be sure to not proceed when done here
//	//	}
//
//	outfile := ctx.String("outfile")
//	log.Debugf("Outfile: %q", outfile)
//
//	if outfile == "-" || outfile == "" {
//		err := _site.ll.Encode(os.Stdout)
//		if err != nil {
//			return cli.Exit(err.Error(), E_WRITEJSON)
//		}
//	} else {
//		err := _site.ll.SaveJSON(outfile)
//		if err != nil {
//			return cli.Exit(err.Error(), E_WRITEJSON)
//		}
//		log.Debugf("Wrote JSON result to %q", outfile)
//	}
//	return nil
//}

// new, potential variant of this:
//func entryPointScrape(ctx *cli.Context) error {
//	var wg sync.WaitGroup
//	getLunchList().RunSiteScrapers(&wg)
//	wg.Wait()
//	getLunchList().Encode(os.Stdout)
//	return nil
//}

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
	app.Name = "go2lunch server"
	app.Version = fmt.Sprintf("%s_%s", VERSION, COMMIT_ID)
	app.Copyright = "(c) 2017 Odd Eivind Ebbesen"
	app.Compiled, _ = time.Parse(time.RFC3339, BUILD_DATE)
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "Odd E. Ebbesen",
			Email: "oddebb@gmail.com",
		},
	}
	app.Usage = "Serve lunch menus from configured sites"
	app.EnableBashCompletion = true

	app.Commands = []*cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"srv"},
			Usage:   "Start server",
			Action:  entryPointServe,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen-adr",
					Aliases: []string{"l"},
					Value:   DEF_WEB_ADR,
					Usage:   "[hostname|IP]:port to listen on for the regular/public site",
				},
				&cli.StringFlag{
					Name:  "listen-adm",
					Value: DEF_ADM_ADR,
					Usage: "[hostname|IP]:port to listen on for the admin api site",
				},
				&cli.StringFlag{
					Name:    "writepid",
					Aliases: []string{"p"},
					Usage:   "Write PID to `FILE`",
				},
				&cli.StringFlag{
					Name:  "cron",
					Usage: "Specify intervals for re-scrape in cron format with `TIMESPEC`",
					// what I had in regular cron: "0 0,30 08-12 * * 1-5"
					// That is: on second 0 of minute 0 and 30 of hour 08-12 of weekday mon-fri on any day of month any month
					// Update @ 2019-10-28: New version of cron-lib does not use second field
				},
				&cli.StringFlag{
					Name:  "load",
					Usage: "Load initial data from `JSONFILE`",
				},
				&cli.BoolFlag{
					Name:  "save-on-exit",
					Usage: "If config was loaded from file, save it back to the same file on exit",
				},
				&cli.BoolFlag{
					Name:  "strip-menus-on-save",
					Usage: "Do not save restaurants and their dishes when saving on exit. Only save structure.",
				},
				&cli.BoolFlag{
					Name:  "noscrape",
					Usage: "Disable scraping",
				},
				&cli.StringFlag{
					Name:  "gtag",
					Usage: "GTAG for Google Analytics in generated pages",
					//Value: GTAG_ID,
				},
			},
		},
		//		{
		//			Name:    "scrape",
		//			Aliases: []string{"scr"},
		//			Usage:   "Scrape source and output JSON or HTML, then exit",
		//			Action:  entryPointScrape,
		//			Flags: []cli.Flag{
		//				&cli.StringFlag{
		//					Name:  "outfile, o",
		//					Usage: "Write JSON result to `FILE` ('-' for STDOUT)",
		//					Value: "-",
		//				},
		//				//				&cli.BoolFlag{
		//				//					Name:  "html",
		//				//					Usage: "Write HTML result to STDOUT",
		//				//				},
		//				&cli.StringFlag{
		//					Name:  "notify-pid, p",
		//					Usage: "Read PID from `FILE` and tell the process with that PID to re-scrape",
		//				},
		//			},
		//		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Aliases: []string{"l"},
			Value:   "info",
			Usage:   "Log `level` (options: debug, info, warn, error, fatal, panic)",
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "Run in debug mode",
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

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}

	os.Exit(E_OK)
}
