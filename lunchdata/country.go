package lunchdata

import (
	"github.com/google/uuid"
)

type Country struct {
	Cities Cities `json:"cities"`
	Name   string `json:"name"`
	ID     string `json:"id"` // preferably international country code, like "se", "no", and so on
	GTag   string `json:"-"`
}

func (c *Country) NumCities() int {
	if c == nil {
		return 0
	}
	return c.Cities.Len()
}

func (c *Country) NumSites() int {
	if c == nil {
		return 0
	}
	return c.Cities.NumSites()
}

func (c *Country) NumRestaurants() int {
	if c == nil {
		return 0
	}
	return c.Cities.NumRestaurants()
}

func (c *Country) NumDishes() int {
	if c == nil {
		return 0
	}
	return c.Cities.NumDishes()
}

func (c *Country) Get(f CityMatch) *City {
	if c == nil {
		return nil
	}
	return c.Cities.Get(f)
}

func (c *Country) GetByID(id string) *City {
	if c == nil {
		return nil
	}
	return c.Cities.GetByID(id)
}

func (c *Country) setGTag(tag string) {
	if c == nil {
		return
	}
	c.GTag = tag
	c.Cities.setGTag(tag)
}

// func (c *Country) RunSiteScrapers(wg *sync.WaitGroup, errChan chan<- error) {
// 	if c == nil {
// 		return
// 	}
// 	c.Cities.RunSiteScrapers(wg, errChan)
// }

func (c *Country) setIDIfEmpty() {
	if c == nil {
		return
	}
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	c.Cities.setIDIfEmpty()
}
