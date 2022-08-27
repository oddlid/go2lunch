package server

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
		Internal / built-in scrapers should be able to register themselves with the engine, preferably
		without the engine having to know anything about it, other than being able to tell all registered
		scrapers to run at given intervals, and then collect their results. Look at go-chat-bot for how
		to solve this.
		For external scrapers, all they need is to post the correct format to the correct URL while
		providing the correct API header key. How the ext scraper is implemented should not be relevant.

	- The above thoughts leaves the question: Should we have internal scrapers at all? Wouldn't it be
	  easier to just skip them and depend fully on external scraping? Advantages and disadvantages?
*/

// const (
// 	logTimeStampLayout  = `2006-01-02T15:04:05.999-07:00`
// 	defaultPubAddr      = ":20666"
// 	defaultAdmAddr      = ":20667"
// 	defaultReadTimeout  = 5
// 	defaultWriteTimeout = 10
// 	defaultIdleTimeout  = 15
// )

// exit codes
// const (
// 	exitCron int = iota
// 	exitReadJSON
// 	exitWriteJSON
// 	exitInitTemplate
// )

// var (
// 	Version    string
// 	BuildDate  string
// 	CommitID   string
// 	_gtag      string
// 	_noScrape  bool
// 	_lunchList *lunchdata.LunchList
// )

/*
Since we call registerSiteScraper() here from init, getLunchList() will get called in staticlunchlist.go
with _lunchlist being nil, resulting in initialization from the static JSON in that file.
If we at a point after that load content from JSON via a file, site.Scraper will get overwritten to nil,
and no further scraping will happen.
For now, that is acceptable, but it's not very obvious and could lead to hard to find bugs later on.
We should rather not use init() at all, but register scrapers at a later point.
*/
// func init() {
// 	lhs := lindholmen.LHScraper{
// 		Logger: zerolog.New(os.Stdout),
// 		URL:    "https://lindholmen.uit.se/omradet/dagens-lunch?embed-mode=iframe",
// 		// URL: "http://localhost:8080",
// 	}
// 	registerSiteScraper(lhs.GetCountryID(), lhs.GetCityID(), lhs.GetSiteID(), &lhs)
// }

// I'd like to find a more flexible and dynamic way of including scrapers, but for now
// we'll use this
// func registerSiteScraper(countryID, cityID, siteID string, scraper lunchdata.SiteScraper) {
// 	lsite := getLunchList().GetSiteByID(countryID, cityID, siteID)
// 	if lsite == nil {
// 		log.Error().
// 			Str("countryID", countryID).
// 			Str("cityID", cityID).
// 			Str("siteID", siteID).
// 			Msg("Site not found")
// 		return
// 	}
// 	lsite.Scraper = scraper
// }

// func lunchListFromFile(filename string) error {
// 	ll, err := lunchdata.LunchListFromFile(filename)
// 	if err != nil {
// 		return err
// 	}
// 	// If we load the LunchList from JSON, and we only have a top-level
// 	// Gtag, we need to propagate it now after load
// 	if ll.Gtag != "" {
// 		ll.PropagateGtag(ll.Gtag)
// 	}
// 	_lunchList = ll
// 	return nil
// }

// hack
// func logInventory() {
// 	//var b strings.Builder

// 	//fmt.Fprintf(
// 	//	&b,
// 	//	"Total subitems: %d, Countries: %d, Cities: %d, Sites: %d, Restaurants: %d, Dishes: %d\n",
// 	//	getLunchList().SubItems(),
// 	//	getLunchList().NumCountries(),
// 	//	getLunchList().NumCities(),
// 	//	getLunchList().NumSites(),
// 	//	getLunchList().NumRestaurants(),
// 	//	getLunchList().NumDishes(),
// 	//)

// 	//sls := getLunchList().GetSiteLinks()
// 	//for _, sl := range sls {
// 	//	fmt.Fprintf(&b, "%s: %s\n", sl.Url, sl.SiteKey)
// 	//}

// 	//log.Debug(b.String())

// 	// log.WithFields(log.Fields{
// 	// 	"Total subitems": getLunchList().SubItems(),
// 	// 	"Countries":      getLunchList().NumCountries(),
// 	// 	"Cities":         getLunchList().NumCities(),
// 	// 	"Sites":          getLunchList().NumSites(),
// 	// 	"Restaurants":    getLunchList().NumRestaurants(),
// 	// 	"Dishes":         getLunchList().NumDishes(),
// 	// }).Debug("Inventory")
// }

// func setGtag(cCtx *cli.Context) {
// 	gtag := cCtx.String(optGTag)
// 	if gtag != "" {
// 		_gtag = gtag
// 		getLunchList().PropagateGtag(_gtag)
// 	}
// }

// func entryPointServe(cCtx *cli.Context) error {
// 	// cron-like scheduling.
// 	// When post updates is fully ready, this should be removed..?
// 	cronspec := cCtx.String(optCron)
// 	if cronspec != "" {
// 		log.Info().
// 			Str("cronspec", cronspec).
// 			Msg("Auto-updating via built-in cron")

// 		cr := cron.New()
// 		if _, err := cr.AddFunc(cronspec, func() {
// 			log.Debug().Msg("Re-scraping on request from internal cron...")
// 			wg := sync.WaitGroup{}
// 			getLunchList().RunSiteScrapers(&wg)
// 			wg.Wait()
// 			// we don't need to propagate _gtag here, as SiteScrapers can only set new Restaurants
// 			// for an already configured Site, so a given Gtag will not be removed by scraping
// 			logInventory()
// 		}); err != nil {
// 			return cli.Exit(err.Error(), exitCron)
// 		}
// 		cr.Start()
// 		log.Info().Msg("Built-in cron started")
// 	}

// 	if err := initTmpl(); err != nil {
// 		return cli.Exit(err.Error(), exitInitTemplate)
// 	}

// 	// Loading the lunchlist from a JSON file is mostly useful for testing.
// 	// As the entire structure is replaced, every Site's Scraper instance will be set to nil,
// 	// and so no auto-scraping will happen, unless we re-register scrapers after JSON load.
// 	// Unless testing scraping, one can load a copy of PROD data in a local instance like this:
// 	//
// 	// $ go run . -d serve --load <(curl -sS https://lunch.oddware.net/json/)

// 	// handle "load" argument here before serving
// 	jfile := cCtx.String(optLoad)
// 	if jfile != "" {
// 		if err := lunchListFromFile(jfile); err != nil {
// 			return cli.Exit(err.Error(), exitReadJSON)
// 		}
// 		log.Info().
// 			Str("file", jfile).
// 			Msg("Lunch list loaded from file. Scraping is now disabled")
// 		_noScrape = true
// 	}

// 	// Important that this call comes after anything that sets content, like lunchListFromFile above
// 	// We should probably make a hook that calls this after any update of content as well
// 	// If we did load from JSON, this gives us the possibility to override the Gtag from CLI
// 	setGtag(cCtx)

// 	listenAdr := cCtx.String(optListenAdr)
// 	listenAdm := cCtx.String(optListenAdm)

// 	pubRouter, admRouter := setupRouter()
// 	pubSrv := createServer(listenAdr, pubRouter)
// 	admSrv := createServer(listenAdm, admRouter)

// 	serverWG := sync.WaitGroup{}
// 	serverWG.Add(2) // 2 = number of servers to wait for
// 	go gracefulShutdown(cCtx.Context, "PubSRV", pubSrv, &serverWG)
// 	go gracefulShutdown(cCtx.Context, "AdmSRV", admSrv, &serverWG)
// 	go func() {
// 		log.Info().
// 			Str("addr", listenAdr).
// 			Msg("Starting public http server")
// 		if err := pubSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatal().Err(err).Msg("Failed to start public http server")
// 		} else {
// 			log.Info().Msg("Public http server shut down cleanly")
// 		}
// 	}()
// 	go func() {
// 		log.Info().
// 			Str("addr", listenAdm).
// 			Msg("Starting admin http server")
// 		if err := admSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatal().Err(err).Msg("Failed to start admin http server")
// 		} else {
// 			log.Info().Msg("Admin http server shut down cleanly")
// 		}
// 	}()

// 	// now, run registered scrapers
// 	// By running 2 levels of goroutines and using waitgroups, we can do all scraping in the background,
// 	// and still wait until all are done before displaying updated stats
// 	_noScrape = cCtx.Bool(optNoScrape)
// 	if !_noScrape {
// 		go func() {
// 			wg := sync.WaitGroup{}
// 			getLunchList().RunSiteScrapers(&wg)
// 			wg.Wait()
// 			logInventory()
// 		}()
// 	}

// 	// Here we block until we get a signal to quit
// 	// Might be months until we reach the next line after this
// 	// Given that thought, maybe this would be a good place to display some uptime stats or something...
// 	serverWG.Wait()

// 	// If we did load contents from a file, let's save it back
// 	saveOnExit := cCtx.Bool(optSaveOnExit)
// 	stripOnSave := cCtx.Bool(optStripMenusOnSave)
// 	if saveOnExit && jfile != "" {
// 		if stripOnSave {
// 			getLunchList().ClearRestaurants()
// 		}
// 		if err := getLunchList().SaveJSON(jfile); err != nil {
// 			return cli.Exit(err.Error(), exitWriteJSON)
// 		}
// 		log.Info().
// 			Str("file", jfile).
// 			Msg("Config saved")
// 	}

// 	return nil
// }

// func createServer(addr string, handler http.Handler) *http.Server {
// 	return &http.Server{
// 		Addr:         addr,
// 		Handler:      handler,
// 		ReadTimeout:  defaultReadTimeout * time.Second,
// 		WriteTimeout: defaultWriteTimeout * time.Second,
// 		IdleTimeout:  defaultIdleTimeout * time.Second,
// 	}
// }

// func gracefulShutdown(ctx context.Context, tag string, srv *http.Server, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	<-ctx.Done()
// 	log.Debug().Str("server", tag).Msg("Shutting down server")

// 	// need a new context here, since if we'd inherit from the passed in context, the shutdown would be immediate and not graceful
// 	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := srv.Shutdown(shutdownCtx); err != nil {
// 		log.Error().Str("server", tag).Err(err).Msg("Failed to shutdown server cleanly")
// 	}
// }

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
// func setCustomAppHelpTmpl() {
// 	cli.AppHelpTemplate = `NAME:
//    {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

// USAGE:
//    {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

// VERSION / BUILD_DATE:
//    {{.Version}} / {{.Compiled}}{{end}}{{end}}{{if .Description}}

// DESCRIPTION:
//    {{.Description}}{{end}}{{if len .Authors}}

// AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
//    {{range $index, $author := .Authors}}{{if $index}}
//    {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

// COMMANDS:{{range .VisibleCategories}}{{if .Name}}
//    {{.Name}}:{{end}}{{range .VisibleCommands}}
//      {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

// GLOBAL OPTIONS:
//    {{range $index, $option := .VisibleFlags}}{{if $index}}
//    {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

// COPYRIGHT:
//    {{.Copyright}}{{end}}
// `
// }

// func setCompileTime(app *cli.App) {
// 	ts, err := time.Parse(time.RFC3339, BuildDate)
// 	if err != nil {
// 		return
// 	}
// 	app.Compiled = ts
// }

// func main() {
// 	setCustomAppHelpTmpl()
// 	app := cli.NewApp()
// 	setCompileTime(app)
// 	app.Name = "go2lunch server"
// 	app.Version = fmt.Sprintf("%s_%s", Version, CommitID)
// 	app.Copyright = "(c) 2017 Odd Eivind Ebbesen"
// 	app.Authors = []*cli.Author{
// 		{
// 			Name:  "Odd E. Ebbesen",
// 			Email: "oddebb@gmail.com",
// 		},
// 	}
// 	app.Usage = "Serve lunch menus from configured sites"
// 	app.EnableBashCompletion = true

// 	app.Commands = []*cli.Command{
// 		{
// 			Name:    "serve",
// 			Aliases: []string{"srv"},
// 			Usage:   "Start server",
// 			Action:  entryPointServe,
// 			Flags: []cli.Flag{
// 				&cli.StringFlag{
// 					Name:    optListenAdr,
// 					Aliases: []string{"l"},
// 					Value:   defaultPubAddr,
// 					Usage:   "[hostname|IP]:port to listen on for the regular/public site",
// 				},
// 				&cli.StringFlag{
// 					Name:  optListenAdm,
// 					Value: defaultAdmAddr,
// 					Usage: "[hostname|IP]:port to listen on for the admin api site",
// 				},
// 				&cli.StringFlag{
// 					Name:    optWritePid,
// 					Aliases: []string{"p"},
// 					Usage:   "Write PID to `FILE`",
// 				},
// 				&cli.StringFlag{
// 					Name:  optCron,
// 					Usage: "Specify intervals for re-scrape in cron format with `TIMESPEC`",
// 					// what I had in regular cron: "0 0,30 08-12 * * 1-5"
// 					// That is: on second 0 of minute 0 and 30 of hour 08-12 of weekday mon-fri on any day of month any month
// 					// Update @ 2019-10-28: New version of cron-lib does not use second field
// 				},
// 				&cli.StringFlag{
// 					Name:  optLoad,
// 					Usage: "Load initial data from `JSONFILE`",
// 				},
// 				&cli.BoolFlag{
// 					Name:  optSaveOnExit,
// 					Usage: "If config was loaded from file, save it back to the same file on exit",
// 				},
// 				&cli.BoolFlag{
// 					Name:  optStripMenusOnSave,
// 					Usage: "Do not save restaurants and their dishes when saving on exit. Only save structure.",
// 				},
// 				&cli.BoolFlag{
// 					Name:  optNoScrape,
// 					Usage: "Disable scraping",
// 				},
// 				&cli.StringFlag{
// 					Name:  optGTag,
// 					Usage: "GTAG for Google Analytics in generated pages",
// 					//Value: GTAG_ID,
// 				},
// 			},
// 		},
// 		//		{
// 		//			Name:    "scrape",
// 		//			Aliases: []string{"scr"},
// 		//			Usage:   "Scrape source and output JSON or HTML, then exit",
// 		//			Action:  entryPointScrape,
// 		//			Flags: []cli.Flag{
// 		//				&cli.StringFlag{
// 		//					Name:  "outfile, o",
// 		//					Usage: "Write JSON result to `FILE` ('-' for STDOUT)",
// 		//					Value: "-",
// 		//				},
// 		//				//				&cli.BoolFlag{
// 		//				//					Name:  "html",
// 		//				//					Usage: "Write HTML result to STDOUT",
// 		//				//				},
// 		//				&cli.StringFlag{
// 		//					Name:  "notify-pid, p",
// 		//					Usage: "Read PID from `FILE` and tell the process with that PID to re-scrape",
// 		//				},
// 		//			},
// 		//		},
// 	}

// 	app.Flags = []cli.Flag{
// 		&cli.StringFlag{
// 			Name:    optLogLevel,
// 			Aliases: []string{"l"},
// 			Value:   zerolog.InfoLevel.String(),
// 			Usage:   "Log `level` (options: trace, debug, info, warn, error, panic, fatal)",
// 		},
// 		&cli.BoolFlag{
// 			Name:    optDebug,
// 			Aliases: []string{"d"},
// 			Usage:   "Set log level to debug",
// 		},
// 	}

// 	app.Before = func(c *cli.Context) error {
// 		zerolog.TimeFieldFormat = logTimeStampLayout
// 		if c.Bool(optDebug) {
// 			zerolog.SetGlobalLevel(zerolog.DebugLevel)
// 		} else if c.IsSet(c.String(optLogLevel)) {
// 			level, err := zerolog.ParseLevel(c.String(optLogLevel))
// 			if err != nil {
// 				return err
// 			}
// 			zerolog.SetGlobalLevel(level)
// 		} else {
// 			zerolog.SetGlobalLevel(zerolog.InfoLevel)
// 		}
// 		return nil
// 	}

// 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
// 	defer cancel()

// 	if err := app.RunContext(ctx, os.Args); err != nil {
// 		cancel()
// 		log.Fatal().Err(err).Msg("Execution failed")
// 	}
// }
