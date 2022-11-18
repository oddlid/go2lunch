package lunchdata

import (
	"sync"
)

type Country struct {
	Cities CityMap `json:"cities"`
	Name   string  `json:"name"`
	ID     string  `json:"id"` // preferably international country code, like "se", "no", and so on
	GTag   string  `json:"-"`
	mu     sync.RWMutex
}

func NewCountry(name, id string) *Country {
	return &Country{
		Name:   name,
		ID:     id,
		Cities: make(CityMap),
	}
}

func (c *Country) NumCities() int {
	if c == nil {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Cities.Len()
}

func (c *Country) NumSites() int {
	if c == nil {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Cities.NumSites()
}

func (c *Country) NumRestaurants() int {
	if c == nil {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Cities.NumRestaurants()
}

func (c *Country) NumDishes() int {
	if c == nil {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Cities.NumDishes()
}

func (c *Country) SetGTag(tag string) *Country {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	c.GTag = tag
	c.Cities.SetGTag(tag)
	c.mu.Unlock()
	return c
}

func (c *Country) Add(cities ...*City) *Country {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	if c.Cities == nil {
		c.Cities = make(CityMap)
	}
	c.Cities.Add(cities...)
	c.mu.Unlock()
	return c
}

func (c *Country) Delete(ids ...string) *Country {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	c.Cities.Delete(ids...)
	c.mu.Unlock()
	return c
}

func (c *Country) Get(id string) *City {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Cities.Get(id)
}

func (c *Country) RunSiteScrapers(wg *sync.WaitGroup, errChan chan<- error) {
	if c == nil {
		return
	}
	if c.Cities != nil {
		c.Cities.RunSiteScrapers(wg, errChan)
	}
}
