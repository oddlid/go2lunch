package lunchdata

import (
	"sync"
)

// A giant list of everything
type LunchList struct {
	Countries CountryMap `json:"countries"`
	GTag      string     `json:"gtag,omitempty"`
	mu        sync.RWMutex
}

func NewLunchList() *LunchList {
	return &LunchList{
		Countries: make(CountryMap),
	}
}

func (l *LunchList) NumCountries() int {
	if l == nil {
		return 0
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Countries.Len()
}

func (l *LunchList) NumCities() int {
	if l == nil {
		return 0
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Countries.NumCities()
}

func (l *LunchList) NumSites() int {
	if l == nil {
		return 0
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Countries.NumSites()
}

func (l *LunchList) NumRestaurants() int {
	if l == nil {
		return 0
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Countries.NumRestaurants()
}

func (l *LunchList) NumDishes() int {
	if l == nil {
		return 0
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Countries.NumDishes()
}

func (l *LunchList) SetGTag(tag string) *LunchList {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	l.GTag = tag
	l.Countries.SetGTag(tag)
	l.mu.Unlock()
	return l
}

func (l *LunchList) Add(countries ...*Country) *LunchList {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	if l.Countries == nil {
		l.Countries = make(CountryMap)
	}
	l.Countries.Add(countries...)
	l.mu.Unlock()
	return l
}

func (l *LunchList) Delete(ids ...string) *LunchList {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	l.Countries.Delete(ids...)
	l.mu.Unlock()
	return l
}

func (l *LunchList) Get(id string) *Country {
	if l == nil {
		return nil
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.Countries.Get(id)
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
