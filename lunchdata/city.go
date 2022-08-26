package lunchdata

import (
	"sync"
)

type City struct {
	sync.RWMutex
	Name  string           `json:"city_name"`
	ID    string           `json:"city_id"` // e.g. osl, gbg or something like the airlines use
	Gtag  string           `json:"-"`
	Sites map[string]*Site `json:"sites"`
}

type Cities []*City

func (cs *Cities) Len() int {
	return len(*cs)
}

func NewCity(name, id string) *City {
	return &City{
		Name:  name,
		ID:    id,
		Sites: make(map[string]*Site),
	}
}

func (c *City) Len() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.Sites)
}

func (c *City) SubItems() int {
	total := 0
	c.RLock()
	for k := range c.Sites {
		total += c.Sites[k].SubItems() + 1 // +1 to count the Site itself as well
	}
	c.RUnlock()
	return total
}

func (c *City) PropagateGtag(tag string) *City {
	c.Lock()
	c.Gtag = tag
	for k := range c.Sites {
		c.Sites[k].PropagateGtag(tag)
	}
	c.Unlock()
	return c
}

func (c *City) AddSite(s *Site) *City {
	c.Lock()
	c.Sites[s.ID] = s
	c.Unlock()
	return c
}

func (c *City) DeleteSite(id string) *City {
	c.Lock()
	delete(c.Sites, id)
	c.Unlock()
	return c
}

func (c *City) HasSites() bool {
	c.RLock()
	defer c.RUnlock()
	return len(c.Sites) > 0
}

func (c *City) HasSite(siteID string) bool {
	c.RLock()
	_, found := c.Sites[siteID]
	c.RUnlock()
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
	c.Lock()
	c.Sites = make(map[string]*Site)
	c.Unlock()
	return c
}

func (c *City) ClearRestaurants() *City {
	c.Lock()
	for k := range c.Sites {
		c.Sites[k].ClearRestaurants()
	}
	c.Unlock()
	return c
}

func (c *City) ClearDishes() *City {
	c.Lock()
	for k := range c.Sites {
		c.Sites[k].ClearDishes()
	}
	c.Unlock()
	return c
}

func (c *City) GetSiteByID(id string) *Site {
	c.RLock()
	defer c.RUnlock()
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
	c.RLock()
	defer c.RUnlock()
	return len(c.Sites)
}

func (c *City) NumRestaurants() int {
	total := 0
	c.RLock()
	for k := range c.Sites {
		total += c.Sites[k].NumRestaurants()
	}
	c.RUnlock()
	return total
}

func (c *City) NumDishes() int {
	total := 0
	c.RLock()
	for k := range c.Sites {
		total += c.Sites[k].NumDishes()
	}
	c.RUnlock()
	return total
}
