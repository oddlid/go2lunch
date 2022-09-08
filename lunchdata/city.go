package lunchdata

import (
	"sync"
)

type City struct {
	Sites SiteMap `json:"sites"`
	Name  string  `json:"city_name"`
	ID    string  `json:"city_id"` // e.g. osl, gbg or something like the airlines use
	Gtag  string  `json:"-"`
	mu    sync.RWMutex
}

type Cities []*City
type CityMap map[string]*City

func (cs *Cities) Len() int {
	return len(*cs)
}

func NewCity(name, id string) *City {
	return &City{
		Name:  name,
		ID:    id,
		Sites: make(SiteMap),
	}
}

func (c *City) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Sites)
}

func (c *City) SubItems() int {
	total := 0
	c.mu.RLock()
	for k := range c.Sites {
		total += c.Sites[k].SubItems() + 1 // +1 to count the Site itself as well
	}
	c.mu.RUnlock()
	return total
}

func (c *City) PropagateGtag(tag string) *City {
	c.mu.Lock()
	c.Gtag = tag
	for k := range c.Sites {
		c.Sites[k].PropagateGtag(tag)
	}
	c.mu.Unlock()
	return c
}

func (c *City) AddSite(s *Site) *City {
	c.mu.Lock()
	c.Sites[s.ID] = s
	c.mu.Unlock()
	return c
}

func (c *City) DeleteSite(id string) *City {
	c.mu.Lock()
	delete(c.Sites, id)
	c.mu.Unlock()
	return c
}

func (c *City) HasSites() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Sites) > 0
}

func (c *City) HasSite(siteID string) bool {
	c.mu.RLock()
	_, found := c.Sites[siteID]
	c.mu.RUnlock()
	return found
}

func (c *City) HasRestaurant(siteID, restaurantID string) bool {
	if !c.HasSite(siteID) {
		return false
	}
	// We should only get here if there is a Site with siteID, so this should not crash
	return c.GetSiteByID(siteID).HasRestaurant(restaurantID)
}

func (c *City) ClearSites() *City {
	c.mu.Lock()
	c.Sites = make(map[string]*Site)
	c.mu.Unlock()
	return c
}

func (c *City) ClearRestaurants() *City {
	c.mu.Lock()
	for k := range c.Sites {
		c.Sites[k].ClearRestaurants()
	}
	c.mu.Unlock()
	return c
}

func (c *City) ClearDishes() *City {
	c.mu.Lock()
	for k := range c.Sites {
		c.Sites[k].ClearDishes()
	}
	c.mu.Unlock()
	return c
}

func (c *City) GetSiteByID(id string) *Site {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Sites[id]
}

func (c *City) GetRestaurantByID(siteID, restaurantID string) *Restaurant {
	s := c.GetSiteByID(siteID)
	if s == nil {
		return nil
	}
	return s.GetRestaurantByID(restaurantID)
}

func (c *City) NumSites() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Sites)
}

func (c *City) NumRestaurants() int {
	total := 0
	c.mu.RLock()
	for k := range c.Sites {
		total += c.Sites[k].NumRestaurants()
	}
	c.mu.RUnlock()
	return total
}

func (c *City) NumDishes() int {
	total := 0
	c.mu.RLock()
	for k := range c.Sites {
		total += c.Sites[k].NumDishes()
	}
	c.mu.RUnlock()
	return total
}
