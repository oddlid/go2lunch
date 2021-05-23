package lunchdata

import (
	"encoding/json"
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Country struct {
	sync.RWMutex
	Name   string           `json:"country_name"`
	ID     string           `json:"country_id"` // preferrably international country code, like "se", "no", and so on
	Gtag   string           `json:"-"`
	Cities map[string]*City `json:"cities"`
}

type Countries []*Country

func (cs *Countries) Add(c *Country) {
	*cs = append(*cs, c)
}

func (cs *Countries) Len() int {
	return len(*cs)
}

func NewCountry(name, id string) *Country {
	return &Country{
		Name:   name,
		ID:     id,
		Cities: make(map[string]*City),
	}
}

func (c *Country) Len() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.Cities)
}

func (c *Country) SubItems() int {
	total := 0
	c.RLock()
	for k := range c.Cities {
		total += c.Cities[k].SubItems() + 1 // +1 to count the City itself as well
	}
	c.RUnlock()
	return total
}

func (c *Country) PropagateGtag(tag string) *Country {
	c.Lock()
	c.Gtag = tag
	for k := range c.Cities {
		c.Cities[k].PropagateGtag(tag)
	}
	c.Unlock()
	return c
}

func (c *Country) AddCity(city *City) *Country {
	c.Lock()
	c.Cities[city.ID] = city
	c.Unlock()
	return c
}

func (c *Country) DeleteCity(id string) *Country {
	c.Lock()
	delete(c.Cities, id)
	c.Unlock()
	return c
}

func (c *Country) HasCities() bool {
	c.RLock()
	defer c.RUnlock()
	return len(c.Cities) > 0
}

func (c *Country) HasCity(cityID string) bool {
	c.RLock()
	_, found := c.Cities[cityID]
	c.RUnlock()
	return found
}

func (c *Country) HasSite(cityID, siteID string) bool {
	if !c.HasCity(cityID) {
		return false
	}
	return c.GetCityById(cityID).HasSite(siteID)
}

func (c *Country) HasRestaurant(cityID, siteID, restaurantID string) bool {
	if !c.HasSite(cityID, siteID) {
		return false
	}
	return c.GetSiteById(cityID, siteID).HasRestaurant(restaurantID)
}

func (c *Country) ClearCities() *Country {
	c.Lock()
	c.Cities = make(map[string]*City)
	c.Unlock()
	return c
}

func (c *Country) ClearSites() *Country {
	c.Lock()
	for k := range c.Cities {
		c.Cities[k].ClearSites()
	}
	c.Unlock()
	return c
}

func (c *Country) ClearRestaurants() *Country {
	c.Lock()
	for k := range c.Cities {
		c.Cities[k].ClearRestaurants()
	}
	c.Unlock()
	return c
}

func (c *Country) ClearDishes() *Country {
	c.Lock()
	for k := range c.Cities {
		c.Cities[k].ClearDishes()
	}
	c.Unlock()
	return c
}

func (c *Country) GetCityById(id string) *City {
	c.RLock()
	city, found := c.Cities[id]
	c.RUnlock()
	if !found {
		countryLog.WithFields(log.Fields{
			"func": "GetCityById",
			"id":   id,
		}).Debug("Not found")
	}
	return city
}

func (c *Country) GetSiteById(cityID, siteID string) *Site {
	city := c.GetCityById(cityID)
	if nil == city {
		return nil
	}
	return city.GetSiteById(siteID)
}

func (c *Country) GetRestaurantById(cityID, siteID, restaurantID string) *Restaurant {
	city := c.GetCityById(cityID)
	if nil == city {
		return nil
	}
	return city.GetRestaurantById(siteID, restaurantID)
}

func (c *Country) NumCities() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.Cities)
}

func (c *Country) NumSites() int {
	total := 0
	c.RLock()
	for k := range c.Cities {
		total += c.Cities[k].NumSites()
	}
	c.RUnlock()
	return total
}

func (c *Country) NumRestaurants() int {
	total := 0
	c.RLock()
	for k := range c.Cities {
		total += c.Cities[k].NumRestaurants()
	}
	c.RUnlock()
	return total
}

func (c *Country) NumDishes() int {
	total := 0
	c.RLock()
	for k := range c.Cities {
		total += c.Cities[k].NumDishes()
	}
	c.RUnlock()
	return total
}

func (c *Country) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(c)
}

func (c *Country) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(c)
}

func CountryFromJSON(r io.Reader) (*Country, error) {
	c := &Country{}
	if err := c.Decode(r); err != nil {
		return nil, err
	}
	return c, nil
}
