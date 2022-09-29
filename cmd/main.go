package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/oddlid/go2lunch/scraper/se/gbg/lindholmen"
	"github.com/oddlid/go2lunch/server"
	"github.com/rs/zerolog"
)

const (
	idSE    = `se`
	nameSE  = `Sweden`
	idGBG   = `gbg`
	nameGBG = `Gothenburg`
	idLH    = `lindholmen`
	nameLH  = `Lindholmen`
)

func getEmptyLunchList(logger zerolog.Logger) *lunchdata.LunchList {
	return &lunchdata.LunchList{
		Countries: lunchdata.CountryMap{
			idSE: {
				Name: nameSE,
				ID:   idSE,
				Cities: lunchdata.CityMap{
					idGBG: {
						Name: nameGBG,
						ID:   idGBG,
						Sites: lunchdata.SiteMap{
							idLH: {
								Name:        nameLH,
								ID:          idLH,
								Comment:     "Gruvan",
								Restaurants: lunchdata.RestaurantMap{},
								Scraper: &lindholmen.Scraper{
									Logger: logger,
									URL:    lindholmen.DefaultScrapeURL,
								},
							},
						},
					},
				},
			},
		},
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	logger := zerolog.New(os.Stdout)
	s := server.LunchServer{
		Log:       logger,
		LunchList: getEmptyLunchList(logger),
		Config:    server.DefaultConfig(),
	}

	if err := s.Start(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed to start Lunch Server")
		return
	}

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := s.Stop(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Server failed to shut down cleanly")
	}
}
