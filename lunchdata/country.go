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

func (c *Country) AddCity(city *City) *Country {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	c.Cities[city.ID] = city
	c.mu.Unlock()
	return c
}

func (c *Country) DeleteCity(id string) *Country {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	delete(c.Cities, id)
	c.mu.Unlock()
	return c
}

func (c *Country) GetCityByID(id string) *City {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Cities[id]
}
