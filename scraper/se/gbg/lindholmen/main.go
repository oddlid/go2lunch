package main

import (
	//"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/oddlid/go2lunch/engine"
	//"github.com/oddlid/go2lunch/site"
	"github.com/oddlid/go2lunch/urlworker"
	"os"
	"bufio"
	//"time"
)

/* 
See: http://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
for how to override VERSION at compile time. 
E.g.:
go build -ldflags "-X 'main.VERSION=$(date -u '+%Y-%m-%d %H:%M:%S')'"
*/
var VERSION string = "2016-02-25"

func LoadRequestsFromJSON(c *cli.Context) (*urlworker.Requests, error) {
	f, err := os.Open(c.String("config"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	reqs := &urlworker.Requests{}
	err = reqs.NewFromJSON(r)
	if err != nil {
		return nil, err
	}
	return reqs, nil
}

func Parse(c *cli.Context) {
	log.Info("Not quite there yet...\n")
}

func main() {
	app := cli.NewApp()
	app.Version = VERSION
	app.Name = "lunch_lindholmen_scraper"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Odd E. Ebbesen",
			Email: "oddebb@gmail.com",
		},
	}
	app.Usage = "Scrape lunch for Lindholmen"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Path to JSON file containing name/url pairs for scraping",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Path to file to write JSON result to, or STDOUT if '-'",
		},
		cli.StringFlag{
			Name:  "rpc-addr",
			Usage: "Address for posting results via Go2Lunch Engine RPC. Format: host:port",
			Value: engine.DEFAULT_DSN_HOST + engine.DEFAULT_DSN_PORT,
		},
		cli.BoolFlag{
			Name:   "debug, d",
			Usage:  "Run in debug mode",
			EnvVar: "DEBUG",
		},
	}

	app.Before = func(c *cli.Context) error {
		log.SetOutput(os.Stdout)
		level, err := log.ParseLevel(c.String("log-level"))
		if err != nil {
			log.Fatal(err.Error())
		}
		log.SetLevel(level)
		if !c.IsSet("log-level") && c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	app.Action = Parse
	app.Run(os.Args)
}
