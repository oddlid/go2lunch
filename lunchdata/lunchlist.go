package lunchdata

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// A giant list of everything
type LunchList struct {
	sync.RWMutex
	Countries map[string]*Country `json:"countries"`
	Gtag      string              `json:"-"`
}

func NewLunchList() *LunchList {
	return &LunchList{
		Countries: make(map[string]*Country),
	}
}

func (ll *LunchList) Len() int {
	ll.RLock()
	defer ll.RUnlock()
	return len(ll.Countries)
}

func (ll *LunchList) SubItems() int {
	total := 0
	ll.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].SubItems() + 1 // +1 to count the Country itself as well
	}
	ll.RUnlock()
	return total
}

func (ll *LunchList) PropagateGtag(tag string) *LunchList {
	//debugLunchList("PropagateGtag(): tag=%s", tag)
	ll.Lock()
	ll.Gtag = tag
	for k := range ll.Countries {
		ll.Countries[k].PropagateGtag(tag)
	}
	ll.Unlock()
	return ll
}

func (ll *LunchList) AddCountry(c Country) *LunchList {
	ll.Lock()
	ll.Countries[c.ID] = &c
	ll.Unlock()
	return ll
}

func (ll *LunchList) DeleteCountry(id string) *LunchList {
	ll.Lock()
	delete(ll.Countries, id)
	ll.Unlock()
	return ll
}

func (ll *LunchList) HasCountries() bool {
	ll.RLock()
	defer ll.RUnlock()
	return len(ll.Countries) > 0
}

func (ll *LunchList) HasCountry(countryID string) bool {
	ll.RLock()
	_, found := ll.Countries[countryID]
	ll.RUnlock()
	return found
}

func (ll *LunchList) HasCity(countryID, cityID string) bool {
	if !ll.HasCountry(countryID) {
		return false
	}
	return ll.GetCountryById(countryID).HasCity(cityID)
}

func (ll *LunchList) HasSite(countryID, cityID, siteID string) bool {
	if !ll.HasCity(countryID, cityID) {
		return false
	}
	return ll.GetCityById(countryID, cityID).HasSite(siteID)
}

func (ll *LunchList) HasRestaurant(countryID, cityID, siteID, restaurantID string) bool {
	if !ll.HasSite(countryID, cityID, siteID) {
		return false
	}
	return ll.GetSiteById(countryID, cityID, siteID).HasRestaurant(restaurantID)
}

func (ll *LunchList) ClearCountries() *LunchList {
	ll.Lock()
	ll.Countries = make(map[string]*Country)
	ll.Unlock()
	return ll
}

func (ll *LunchList) ClearCities() *LunchList {
	ll.Lock()
	for k := range ll.Countries {
		ll.Countries[k].ClearCities()
	}
	ll.Unlock()
	return ll
}

func (ll *LunchList) ClearSites() *LunchList {
	ll.Lock()
	for k := range ll.Countries {
		ll.Countries[k].ClearSites()
	}
	ll.Unlock()
	return ll
}

func (ll *LunchList) ClearRestaurants() *LunchList {
	ll.Lock()
	for k := range ll.Countries {
		ll.Countries[k].ClearRestaurants()
	}
	ll.Unlock()
	return ll
}

func (ll *LunchList) ClearDishes() *LunchList {
	ll.Lock()
	for k := range ll.Countries {
		ll.Countries[k].ClearDishes()
	}
	ll.Unlock()
	return ll
}

func (ll *LunchList) GetCountryById(id string) *Country {
	ll.RLock()
	c, found := ll.Countries[id]
	ll.RUnlock()
	if !found {
		debugLunchList("GetCountryById: %q not found", id)
	}
	return c
}

func (ll *LunchList) GetCityById(countryID, cityID string) *City {
	c := ll.GetCountryById(countryID)
	if nil == c {
		return nil
	}
	return c.GetCityById(cityID)
}

func (ll *LunchList) GetSiteById(countryID, cityID, siteID string) *Site {
	c := ll.GetCountryById(countryID)
	if nil == c {
		return nil
	}
	return c.GetSiteById(cityID, siteID)
}

func (ll *LunchList) GetSiteByLink(sl SiteLink) *Site {
	return ll.GetSiteById(sl.CountryID, sl.CityID, sl.SiteID)
}

func (ll *LunchList) GetRestaurantById(countryID, cityID, siteID, restaurantID string) *Restaurant {
	c := ll.GetCountryById(countryID)
	if nil == c {
		return nil
	}
	return c.GetRestaurantById(cityID, siteID, restaurantID)
}

func (ll *LunchList) NumCountries() int {
	ll.RLock()
	defer ll.RUnlock()
	return len(ll.Countries)
}

func (ll *LunchList) NumCities() int {
	total := 0
	ll.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].NumCities()
	}
	ll.RUnlock()
	return total
}

func (ll *LunchList) NumSites() int {
	total := 0
	ll.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].NumSites()
	}
	ll.RUnlock()
	return total
}

func (ll *LunchList) NumRestaurants() int {
	total := 0
	ll.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].NumRestaurants()
	}
	ll.RUnlock()
	return total
}

func (ll *LunchList) NumDishes() int {
	total := 0
	ll.RLock()
	for k := range ll.Countries {
		total += ll.Countries[k].NumDishes()
	}
	ll.RUnlock()
	return total
}

func (ll *LunchList) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(ll)
}

func (ll *LunchList) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(ll)
}

func (ll *LunchList) SaveJSON(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = ll.Encode(w)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func LunchListFromJSON(r io.Reader) (*LunchList, error) {
	ll := &LunchList{}
	if err := ll.Decode(r); err != nil {
		return nil, err
	}
	return ll, nil
}

func LunchListFromFile(fileName string) (*LunchList, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	return LunchListFromJSON(r)
}

// GetSiteLinks returns a SiteLinks slice for all configured sites
func (ll *LunchList) GetSiteLinks() SiteLinks {
	sl := make(SiteLinks, 0)

	for _, country := range ll.Countries {
		for _, city := range country.Cities {
			for _, site := range city.Sites {
				sl = append(sl, SiteLink{
					CountryName: country.Name,
					CountryID:   country.ID,
					CityName:    city.Name,
					CityID:      city.ID,
					SiteName:    site.Name,
					SiteID:      site.ID,
					//SiteKey:     site.Key,
					Url:     fmt.Sprintf("%s/%s/%s/", country.ID, city.ID, site.ID),
					Comment: fmt.Sprintf("%s / %s / %s", country.Name, city.Name, site.Name),
				})
			}
		}
	}

	return sl
}

func (ll *LunchList) GetSiteKeyLinks() SiteKeyLinks {
	skls := make(SiteKeyLinks, 0)

	for _, country := range ll.Countries {
		for _, city := range country.Cities {
			for _, site := range city.Sites {
				//if site.Key != "" {
				skls = append(skls, SiteKeyLink{
					CountryID: country.ID,
					CityID:    city.ID,
					SiteID:    site.ID,
					SiteKey:   site.Key,
				})
				//}
			}
		}
	}

	return skls
}

func (ll *LunchList) SetSiteKeys(skls SiteKeyLinks) {
	for _, skl := range skls {
		site := ll.GetSiteById(skl.CountryID, skl.CityID, skl.SiteID)
		if nil != site {
			site.Lock()
			site.Key = skl.SiteKey
			site.Unlock()
		}
	}
}

func (ll *LunchList) GetSiteLinkById(countryID, cityID, siteID string) *SiteLink {
	country := ll.GetCountryById(countryID)
	if nil == country {
		return nil
	}
	city := country.GetCityById(cityID)
	if nil == city {
		return nil
	}
	site := city.GetSiteById(siteID)
	if nil == site {
		return nil
	}

	return &SiteLink{
		CountryName: country.Name,
		CountryID:   country.ID,
		CityName:    city.Name,
		CityID:      city.ID,
		SiteName:    site.Name,
		SiteID:      site.ID,
		SiteKey:     site.Key,
		Url:         fmt.Sprintf("%s/%s/%s/", country.ID, city.ID, site.ID),
		Comment:     fmt.Sprintf("%s / %s / %s", country.Name, city.Name, site.Name),
	}
}
