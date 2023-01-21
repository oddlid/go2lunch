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
		Countries: lunchdata.Countries{
			{
				Name: nameSE,
				ID:   idSE,
				Cities: lunchdata.Cities{
					{
						Name: nameGBG,
						ID:   idGBG,
						Sites: lunchdata.Sites{
							{
								Name:        nameLH,
								ID:          idLH,
								Comment:     "Gruvan",
								Restaurants: lunchdata.Restaurants{},
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
