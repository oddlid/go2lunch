package main

import (
	"github.com/urfave/cli/v2"
)

func actionServe(cCtx *cli.Context) error {
	// If load param is given, we try to load a LunchList object from the given file.
	// If load fails, we use the default empty LunchList with the predefined structure.
	// If load param is not given, we use the default empty LunchList right away.
	// If cron param is given, we set up background scraping at the specified schedule.
	// If cron param is not given, content will be static with whatever we have for the
	// lunchlist.
	return nil
}
