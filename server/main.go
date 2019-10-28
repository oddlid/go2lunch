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
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/robfig/cron"
	"github.com/urfave/cli"

	// Internal scrapers
	"github.com/oddlid/go2lunch/scraper/se/gbg/lindholmen"
)

const (
	VERSION           string = "2019-10-28"
	DEF_WEB_ADR       string = ":20666"
	DEF_ADM_ADR       string = ":20667"
	DEF_COUNTRY_ID    string = "se"
	DEF_CITY_ID       string = "gbg"
	DEF_SITE_ID       string = "lindholmen"
	DEF_READ_TIMEOUT         = 5
	DEF_WRITE_TIMEOUT        = 10
	DEF_IDLE_TIMEOUT         = 15

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
	BUILD_DATE string
	COMMIT_ID  string
	_gtag      string
	_noScrape  bool
	_lunchList *lunchdata.LunchList
)


func init() {
	lhs := lindholmen.LHScraper{}
	registerSiteScraper(lhs.GetCountryID(), lhs.GetCityID(), lhs.GetSiteID(), lhs)
}

// I'd like to find a more flexible and dynamic way of including scrapers, but for now
// we'll use this
func registerSiteScraper(countryID, cityID, siteID string, scraper lunchdata.SiteScraper) {
	lsite := getLunchList().GetSiteById(countryID, cityID, siteID)
	if nil == lsite {
		log.Warnf("registerSiteScraper(): Could not find site @ %s/%s/%s", countryID, cityID, siteID)
		return
	}
	lsite.Scraper = scraper
}

func lunchListFromJSON(filename string) error {
	ll, err := lunchdata.LunchListFromFile(filename)
	if err != nil {
		log.Errorf("Unable to load site from JSON: %q", err.Error())
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
	var b strings.Builder

	fmt.Fprintf(
		&b,
		"Total subitems: %d, Countries: %d, Cities: %d, Sites: %d, Restaurants: %d, Dishes: %d\n",
		getLunchList().SubItems(),
		getLunchList().NumCountries(),
		getLunchList().NumCities(),
		getLunchList().NumSites(),
		getLunchList().NumRestaurants(),
		getLunchList().NumDishes(),
	)

	sls := getLunchList().GetSiteLinks()
	for _, sl := range sls {
		fmt.Fprintf(&b, "%s: %s\n", sl.Url, sl.SiteKey)
	}

	log.Debug(b.String())
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
	log.Debugf("PID: %d", os.Getpid())

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

	// cron-like scheduling.
	// Turns out app/docker hangs sometimes, when triggering scrape from regular cron with "docker exec" (even with timeout),
	// so trying to use an internal solution instead
	// When post updates is fully ready, this should be removed..?
	cronspec := ctx.String("cron")
	if cronspec != "" {
		log.Infof("Auto-updating via built-in cron @ %q", cronspec)
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
		err := lunchListFromJSON(jfile)
		if err != nil {
			return cli.NewExitError(err.Error(), E_READJSON)
		}
		log.Debugf("Load site from %q successful!", jfile)
		//_noScrape = true
	}

	// Important that this call comes after anything that sets content, like lunchListFromJSON above
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
	go gracefulShutdown(pubSrv, quit, &wg)
	go gracefulShutdown(admSrv, quit, &wg)
	go func() {
		log.Infof("Public server listening on port %s", listenAdr)
		if err := pubSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Error starting pub server on port %s", listenAdr)
		}
	}()
	go func() {
		log.Infof("Admin server listening on port %s", listenAdm)
		if err := admSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Error starting admin server on port %s", listenAdm)
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
		log.Infof("Wrote config back to %q", jfile)
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

func gracefulShutdown(srv *http.Server, quit <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	<-quit
	log.Debug("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Error shutting down server: %v", err)
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
//			return cli.NewExitError(err.Error(), E_READPID)
//		}
//		if pid > 0 {
//			err := notifyPid(pid)
//			if err != nil {
//				return cli.NewExitError(err.Error(), E_NOTIFYPID)
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
//		return cli.NewExitError(err.Error(), E_UPDATE)
//	}
//
//	//	// write html output if requested, otherwise json
//	//	dumpHtml := ctx.Bool("html")
//	//	if dumpHtml {
//	//		err := initTmpl() // does more than we need here, so maybe rewrite some time...
//	//		if err != nil {
//	//			return cli.NewExitError(err.Error(), E_INITTMPL)
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
//			return cli.NewExitError(err.Error(), E_WRITEJSON)
//		}
//	} else {
//		err := _site.ll.SaveJSON(outfile)
//		if err != nil {
//			return cli.NewExitError(err.Error(), E_WRITEJSON)
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
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Odd E. Ebbesen",
			Email: "oddebb@gmail.com",
		},
	}
	app.Usage = "Serve lunch menus from configured sites"

	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"srv"},
			Usage:   "Start server",
			Action:  entryPointServe,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen-adr, l",
					Value: DEF_WEB_ADR,
					Usage: "(hostname|IP):port to listen on for the regular/public site",
				},
				cli.StringFlag{
					Name:  "listen-adm",
					Value: DEF_ADM_ADR,
					Usage: "(hostname|IP):port to listen on for the admin api site",
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
					// Update @ 2019-10-28: New version of cron-lib does not use second field
				},
				cli.StringFlag{
					Name:  "load",
					Usage: "Load initial data from `JSONFILE`",
				},
				cli.BoolFlag{
					Name:  "save-on-exit",
					Usage: "If config was loaded from file, save it back to the same file on exit",
				},
				cli.BoolFlag{
					Name:  "strip-menus-on-save",
					Usage: "Do not save restaurants and their dishes when saving on exit. Only save structure.",
				},
				cli.BoolFlag{
					Name:  "noscrape",
					Usage: "Disable scraping",
				},
				cli.StringFlag{
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
		//				cli.StringFlag{
		//					Name:  "outfile, o",
		//					Usage: "Write JSON result to `FILE` ('-' for STDOUT)",
		//					Value: "-",
		//				},
		//				//				cli.BoolFlag{
		//				//					Name:  "html",
		//				//					Usage: "Write HTML result to STDOUT",
		//				//				},
		//				cli.StringFlag{
		//					Name:  "notify-pid, p",
		//					Usage: "Read PID from `FILE` and tell the process with that PID to re-scrape",
		//				},
		//			},
		//		},
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
