package lunchdata

import (
	"github.com/google/uuid"
)

type City struct {
	Name  string `json:"name"`
	ID    string `json:"id"` // e.g. osl, gbg or something like the airlines use
	GTag  string `json:"-"`
	Sites Sites  `json:"sites"`
}

func (c *City) NumSites() int {
	if c == nil {
		return 0
	}
	return c.Sites.Len()
}

func (c *City) NumRestaurants() int {
	if c == nil {
		return 0
	}
	return c.Sites.NumRestaurants()
}

func (c *City) NumDishes() int {
	if c == nil {
		return 0
	}
	return c.Sites.NumDishes()
}

func (c *City) Get(f SiteMatch) *Site {
	if c == nil {
		return nil
	}
	return c.Sites.Get(f)
}

func (c *City) GetByID(id string) *Site {
	if c == nil {
		return nil
	}
	return c.Sites.GetByID(id)
}

// func (c *City) RunSiteScrapers(wg *sync.WaitGroup, errChan chan<- error) {
// 	if c == nil {
// 		return
// 	}
// 	// One would think doing a lock here would be good, but since SiteMap.RunSiteScrapers()
// 	// starts one goroutine for each site and then return, the unlock here would come long before
// 	// the scraping is actually done, and so not really give any protection.
// 	// It's probably best to just lock at the top level, in LunchList.
// 	c.Sites.RunSiteScrapers(wg, errChan)
// }

func (c *City) setIDIfEmpty() {
	if c == nil {
		return
	}
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	c.Sites.setIDIfEmpty()
}
