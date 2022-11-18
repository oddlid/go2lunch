package lunchdata

import (
	"sync"
)

type City struct {
	Sites SiteMap `json:"sites"`
	Name  string  `json:"name"`
	ID    string  `json:"id"` // e.g. osl, gbg or something like the airlines use
	GTag  string  `json:"-"`
	mu    sync.RWMutex
}

func NewCity(name, id string) *City {
	return &City{
		Name:  name,
		ID:    id,
		Sites: make(SiteMap),
	}
}

func (c *City) NumSites() int {
	if c == nil {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Sites.Len()
}

func (c *City) NumRestaurants() int {
	if c == nil {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Sites.NumRestaurants()
}

func (c *City) NumDishes() int {
	if c == nil {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Sites.NumDishes()
}

func (c *City) SetGTag(tag string) *City {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	c.GTag = tag
	c.Sites.SetGTag(tag)
	c.mu.Unlock()
	return c
}

func (c *City) Add(sites ...*Site) *City {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	if c.Sites == nil {
		c.Sites = make(SiteMap)
	}
	c.Sites.Add(sites...)
	c.mu.Unlock()
	return c
}

func (c *City) Delete(ids ...string) *City {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	c.Sites.Delete(ids...)
	c.mu.Unlock()
	return c
}

func (c *City) Get(id string) *Site {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Sites[id]
}

func (c *City) RunSiteScrapers(wg *sync.WaitGroup, errChan chan<- error) {
	if c == nil {
		return
	}
	// One would think doing a lock here would be good, but since SiteMap.RunSiteScrapers()
	// starts one goroutine for each site and then return, the unlock here would come long before
	// the scraping is actually done, and so not really give any protection.
	// It's probably best to just lock at the top level, in LunchList.
	if c.Sites != nil {
		c.Sites.RunSiteScrapers(wg, errChan)
	}
}
