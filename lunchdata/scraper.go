package lunchdata

/*
This file defines the interface for a scraper.
This is intended for scrapers that are to be registered internally and run on a schedule.
Scrapers that deliver their results via HTTP POST can be implemented in any way in
any language, as long as they provide the proper JSON to the proper URL.
*/

type SiteScraper interface {
	Scrape() (RestaurantMap, error)
	CountryID() string
	CityID() string
	SiteID() string
}
