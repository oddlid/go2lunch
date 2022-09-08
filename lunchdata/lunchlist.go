package lunchdata

import (
	"sync"
)

// A giant list of everything
type LunchList struct {
	Countries CountryMap `json:"countries"`
	Gtag      string     `json:"gtag,omitempty"`
	mu        sync.RWMutex
}

func NewLunchList() *LunchList {
	return &LunchList{
		Countries: make(CountryMap),
	}
}

func (ll *LunchList) Len() int {
	ll.mu.RLock()
	defer ll.mu.RUnlock()
	return len(ll.Countries)
}

func (ll *LunchList) SubItems() int {
	total := 0
	ll.mu.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].SubItems() + 1 // +1 to count the Country itself as well
	}
	ll.mu.RUnlock()
	return total
}

func (ll *LunchList) PropagateGtag(tag string) *LunchList {
	ll.mu.Lock()
	ll.Gtag = tag
	for k := range ll.Countries {
		ll.Countries[k].PropagateGtag(tag)
	}
	ll.mu.Unlock()
	return ll
}

func (ll *LunchList) AddCountry(c *Country) *LunchList {
	ll.mu.Lock()
	ll.Countries[c.ID] = c
	ll.mu.Unlock()
	return ll
}

func (ll *LunchList) DeleteCountry(id string) *LunchList {
	ll.mu.Lock()
	delete(ll.Countries, id)
	ll.mu.Unlock()
	return ll
}

func (ll *LunchList) HasCountries() bool {
	ll.mu.RLock()
	defer ll.mu.RUnlock()
	return len(ll.Countries) > 0
}

func (ll *LunchList) HasCountry(countryID string) bool {
	ll.mu.RLock()
	_, found := ll.Countries[countryID]
	ll.mu.RUnlock()
	return found
}

func (ll *LunchList) HasCity(countryID, cityID string) bool {
	if !ll.HasCountry(countryID) {
		return false
	}
	return ll.GetCountryByID(countryID).HasCity(cityID)
}

func (ll *LunchList) HasSite(countryID, cityID, siteID string) bool {
	if !ll.HasCity(countryID, cityID) {
		return false
	}
	return ll.GetCityByID(countryID, cityID).HasSite(siteID)
}

func (ll *LunchList) HasRestaurant(countryID, cityID, siteID, restaurantID string) bool {
	if !ll.HasSite(countryID, cityID, siteID) {
		return false
	}
	return ll.GetSiteByID(countryID, cityID, siteID).HasRestaurant(restaurantID)
}

func (ll *LunchList) ClearCountries() *LunchList {
	ll.mu.Lock()
	ll.Countries = make(map[string]*Country)
	ll.mu.Unlock()
	return ll
}

func (ll *LunchList) ClearCities() *LunchList {
	ll.mu.Lock()
	for k := range ll.Countries {
		ll.Countries[k].ClearCities()
	}
	ll.mu.Unlock()
	return ll
}

func (ll *LunchList) ClearSites() *LunchList {
	ll.mu.Lock()
	for k := range ll.Countries {
		ll.Countries[k].ClearSites()
	}
	ll.mu.Unlock()
	return ll
}

func (ll *LunchList) ClearRestaurants() *LunchList {
	ll.mu.Lock()
	for k := range ll.Countries {
		ll.Countries[k].ClearRestaurants()
	}
	ll.mu.Unlock()
	return ll
}

func (ll *LunchList) ClearDishes() *LunchList {
	ll.mu.Lock()
	for k := range ll.Countries {
		ll.Countries[k].ClearDishes()
	}
	ll.mu.Unlock()
	return ll
}

func (ll *LunchList) GetCountryByID(id string) *Country {
	ll.mu.RLock()
	defer ll.mu.RUnlock()
	return ll.Countries[id]
}

func (ll *LunchList) GetCityByID(countryID, cityID string) *City {
	c := ll.GetCountryByID(countryID)
	if c == nil {
		return nil
	}
	return c.GetCityByID(cityID)
}

func (ll *LunchList) GetSiteByID(countryID, cityID, siteID string) *Site {
	c := ll.GetCountryByID(countryID)
	if c == nil {
		return nil
	}
	return c.GetSiteByID(cityID, siteID)
}

func (ll *LunchList) GetSiteByLink(sl SiteLink) *Site {
	return ll.GetSiteByID(sl.CountryID, sl.CityID, sl.SiteID)
}

func (ll *LunchList) GetRestaurantByID(countryID, cityID, siteID, restaurantID string) *Restaurant {
	c := ll.GetCountryByID(countryID)
	if c == nil {
		return nil
	}
	return c.GetRestaurantByID(cityID, siteID, restaurantID)
}

func (ll *LunchList) NumCountries() int {
	ll.mu.RLock()
	defer ll.mu.RUnlock()
	return len(ll.Countries)
}

func (ll *LunchList) NumCities() int {
	total := 0
	ll.mu.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].NumCities()
	}
	ll.mu.RUnlock()
	return total
}

func (ll *LunchList) NumSites() int {
	total := 0
	ll.mu.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].NumSites()
	}
	ll.mu.RUnlock()
	return total
}

func (ll *LunchList) NumRestaurants() int {
	total := 0
	ll.mu.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].NumRestaurants()
	}
	ll.mu.RUnlock()
	return total
}

func (ll *LunchList) NumDishes() int {
	total := 0
	ll.mu.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].NumDishes()
	}
	ll.mu.RUnlock()
	return total
}

func (ll *LunchList) RunSiteScrapers(wg *sync.WaitGroup) {
	// I _think_ we might not need to lock the whole LunchList... Should be
	// enough that each site locks itself before updating contents, as we're
	// not adding or removing any countries/cities/sites from the list, only
	// changing the content in each site, if it has a registered scraper
	//ll.Lock()
	for _, country := range ll.Countries {
		for _, city := range country.Cities {
			for _, site := range city.Sites {
				wg.Add(1)
				go site.RunScraper(wg)
			}
		}
	}
	//ll.Unlock()
}
