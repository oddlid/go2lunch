package lunchdata

/*
This file defines the interface for a scraper.
This is intended for scrapers that are to be registered internally and run on a schedule.
Scrapers that deliver their results via HTTP POST can be implemented in any way in
any language, as long as they provide the proper JSON to the proper URL.
*/

type DishScraper interface {
	Scrape() (Dish, error)
}

type RestaurantScraper interface {
}

// A site instance could have a field for a SiteScraper
// If not nil, it would be able to run and update the site contents.
// So, should the scraper return a full site instance, or just a slice of Restaurants?
// 2019-08-21 21:45: A slice of Restaurant is the answer
type SiteScraper interface {
}

type CityScraper interface {
}

type CountryScraper interface {
}

type LunchListScraper interface {
}
