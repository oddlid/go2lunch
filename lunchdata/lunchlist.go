package lunchdata

import (
	"fmt"

	"github.com/google/uuid"
)

// A giant list of everything
type LunchList struct {
	Countries Countries `json:"countries"`
	ID        string    `json:"id"`
	GTag      string    `json:"-"`
}

func (l *LunchList) NumCountries() int {
	if l == nil {
		return 0
	}
	return l.Countries.Len()
}

func (l *LunchList) NumCities() int {
	if l == nil {
		return 0
	}
	return l.Countries.NumCities()
}

func (l *LunchList) NumSites() int {
	if l == nil {
		return 0
	}
	return l.Countries.NumSites()
}

func (l *LunchList) NumRestaurants() int {
	if l == nil {
		return 0
	}
	return l.Countries.NumRestaurants()
}

func (l *LunchList) NumDishes() int {
	if l == nil {
		return 0
	}
	return l.Countries.NumDishes()
}

func (l *LunchList) Get(f CountryMatch) *Country {
	if l == nil {
		return nil
	}
	return l.Countries.Get(f)
}

func (l *LunchList) GetByID(id string) *Country {
	if l == nil {
		return nil
	}
	return l.Countries.GetByID(id)
}

func (l *LunchList) RegisterSiteScraper(s SiteScraper) error {
	if l == nil {
		return nil
	}
	if s == nil {
		return errNilScraper
	}
	if site := l.GetByID(s.CountryID()).GetByID(s.CityID()).GetByID(s.SiteID()).SetScraper(s); site == nil {
		return fmt.Errorf(
			"%w: Not found: Country: %q City: %q Site: %q",
			errNilSite,
			s.CountryID(),
			s.CityID(),
			s.SiteID(),
		)
	}
	return nil
}

func (l *LunchList) RunSiteScrapers() {
	if l == nil {
		return
	}
	// TODO: Think about how to best solve this. Do we want this func to be blocking or not?
	// If we want to lock, then we need to create our own WaitGroup to pass in here, so we don't unlock
	// until all is done.
	// What might be a good way, is to create both the wg and the error channel here, pass them in,
	// then wait on the wg, and after that close the error channel and return it. That way, the caller can range
	// over any returned errors. Downside to that, is that this func is then blocking.
}

// func (ll *LunchList) RunSiteScrapers(wg *sync.WaitGroup) {
// 	// I _think_ we might not need to lock the whole LunchList... Should be
// 	// enough that each site locks itself before updating contents, as we're
// 	// not adding or removing any countries/cities/sites from the list, only
// 	// changing the content in each site, if it has a registered scraper
// 	//ll.Lock()
// 	for _, country := range ll.Countries {
// 		for _, city := range country.Cities {
// 			for _, site := range city.Sites {
// 				wg.Add(1)
// 				go site.RunScraper(wg)
// 			}
// 		}
// 	}
// 	//ll.Unlock()
// }

func (l *LunchList) SetIDIfEmpty() {
	if l == nil {
		return
	}
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	l.Countries.setIDIfEmpty()
}
