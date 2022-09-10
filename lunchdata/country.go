package lunchdata

import (
	"sync"
)

type Country struct {
	Cities CityMap `json:"cities"`
	Name   string  `json:"country_name"`
	ID     string  `json:"country_id"` // preferably international country code, like "se", "no", and so on
	GTag   string  `json:"-"`
	mu     sync.RWMutex
}

type Countries []*Country
type CountryMap map[string]*Country

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
		Cities: make(CityMap),
	}
}

func (c *Country) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Cities)
}

// func (c *Country) SubItems() int {
// 	total := 0
// 	c.mu.RLock()
// 	for k := range c.Cities {
// 		total += c.Cities[k].SubItems() + 1 // +1 to count the City itself as well
// 	}
// 	c.mu.RUnlock()
// 	return total
// }

func (c *Country) SetGTag(tag string) *Country {
	c.mu.Lock()
	c.GTag = tag
	for k := range c.Cities {
		c.Cities[k].SetGTag(tag)
	}
	c.mu.Unlock()
	return c
}

func (c *Country) AddCity(city *City) *Country {
	c.mu.Lock()
	c.Cities[city.ID] = city
	c.mu.Unlock()
	return c
}

func (c *Country) DeleteCity(id string) *Country {
	c.mu.Lock()
	delete(c.Cities, id)
	c.mu.Unlock()
	return c
}

func (c *Country) HasCities() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Cities) > 0
}

func (c *Country) HasCity(cityID string) bool {
	c.mu.RLock()
	_, found := c.Cities[cityID]
	c.mu.RUnlock()
	return found
}

func (c *Country) HasSite(cityID, siteID string) bool {
	if !c.HasCity(cityID) {
		return false
	}
	return c.GetCityByID(cityID).HasSite(siteID)
}

// func (c *Country) HasRestaurant(cityID, siteID, restaurantID string) bool {
// 	if !c.HasSite(cityID, siteID) {
// 		return false
// 	}
// 	return c.GetSiteByID(cityID, siteID).HasRestaurant(restaurantID)
// }

func (c *Country) ClearCities() *Country {
	c.mu.Lock()
	c.Cities = make(map[string]*City)
	c.mu.Unlock()
	return c
}

// func (c *Country) ClearSites() *Country {
// 	c.mu.Lock()
// 	for k := range c.Cities {
// 		c.Cities[k].ClearSites()
// 	}
// 	c.mu.Unlock()
// 	return c
// }

// func (c *Country) ClearRestaurants() *Country {
// 	c.mu.Lock()
// 	for k := range c.Cities {
// 		c.Cities[k].ClearRestaurants()
// 	}
// 	c.mu.Unlock()
// 	return c
// }

// func (c *Country) ClearDishes() *Country {
// 	c.mu.Lock()
// 	for k := range c.Cities {
// 		c.Cities[k].ClearDishes()
// 	}
// 	c.mu.Unlock()
// 	return c
// }

func (c *Country) GetCityByID(id string) *City {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Cities[id]
}

func (c *Country) GetSiteByID(cityID, siteID string) *Site {
	city := c.GetCityByID(cityID)
	if city == nil {
		return nil
	}
	return city.GetSiteByID(siteID)
}

// func (c *Country) GetRestaurantByID(cityID, siteID, restaurantID string) *Restaurant {
// 	city := c.GetCityByID(cityID)
// 	if city == nil {
// 		return nil
// 	}
// 	return city.GetRestaurantByID(siteID, restaurantID)
// }

func (c *Country) NumCities() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Cities)
}

func (c *Country) NumSites() int {
	total := 0
	c.mu.RLock()
	for k := range c.Cities {
		total += c.Cities[k].NumSites()
	}
	c.mu.RUnlock()
	return total
}

// func (c *Country) NumRestaurants() int {
// 	total := 0
// 	c.mu.RLock()
// 	for k := range c.Cities {
// 		total += c.Cities[k].NumRestaurants()
// 	}
// 	c.mu.RUnlock()
// 	return total
// }

func (c *Country) NumDishes() int {
	total := 0
	c.mu.RLock()
	for k := range c.Cities {
		total += c.Cities[k].NumDishes()
	}
	c.mu.RUnlock()
	return total
}
