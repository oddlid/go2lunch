package lunchdata

import (
	"sync"
	// log "github.com/sirupsen/logrus"
)

type Country struct {
	sync.RWMutex
	Name   string           `json:"country_name"`
	ID     string           `json:"country_id"` // preferably international country code, like "se", "no", and so on
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
	return c.GetCityByID(cityID).HasSite(siteID)
}

func (c *Country) HasRestaurant(cityID, siteID, restaurantID string) bool {
	if !c.HasSite(cityID, siteID) {
		return false
	}
	return c.GetSiteByID(cityID, siteID).HasRestaurant(restaurantID)
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

func (c *Country) GetCityByID(id string) *City {
	c.RLock()
	defer c.RUnlock()
	return c.Cities[id]
}

func (c *Country) GetSiteByID(cityID, siteID string) *Site {
	city := c.GetCityByID(cityID)
	if city == nil {
		return nil
	}
	return city.GetSiteByID(siteID)
}

func (c *Country) GetRestaurantByID(cityID, siteID, restaurantID string) *Restaurant {
	city := c.GetCityByID(cityID)
	if city == nil {
		return nil
	}
	return city.GetRestaurantByID(siteID, restaurantID)
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
