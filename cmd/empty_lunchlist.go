package main

import (
	"github.com/oddlid/go2lunch/lunchdata"
)

const (
	idLunchList = `lunchList`
	idSE        = `se`
	nameSE      = `Sweden`
	idGBG       = `gbg`
	nameGBG     = `Gothenburg`
	idLH        = `lindholmen`
	nameLH      = `Lindholmen`
)

func getEmptyLunchList() *lunchdata.LunchList {
	return &lunchdata.LunchList{
		ID: idLunchList,
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
								// Scraper: &lindholmen.Scraper{
								// 	Logger: logger,
								// 	URL:    lindholmen.DefaultScrapeURL,
								// },
							},
						},
					},
				},
			},
		},
	}
}
