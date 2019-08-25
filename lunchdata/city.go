package lunchdata

import (
	"encoding/json"
	"io"
)

type City struct {
	Name  string           `json:"city_name"`
	ID    string           `json:"city_id"` // e.g. osl, gbg or something like the airlines use
	Gtag  string           `json:"-"`
	Sites map[string]*Site `json:"sites"`
}

type Cities []City

func (cs *Cities) Add(c City) {
	*cs = append(*cs, c)
}

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
	return len(c.Sites)
}

func (c *City) SubItems() int {
	total := 0
	for k := range c.Sites {
		total += c.Sites[k].SubItems() + 1 // +1 to count the Site itself as well
	}
	return total
}

func (c *City) PropagateGtag(tag string) *City {
	c.Gtag = tag
	for k := range c.Sites {
		c.Sites[k].PropagateGtag(tag)
	}
	return c
}

func (c *City) AddSite(s Site) *City {
	c.Sites[s.ID] = &s
	return c
}

func (c *City) DeleteSite(id string) *City {
	delete(c.Sites, id)
	return c
}

func (c *City) HasSites() bool {
	return len(c.Sites) > 0
}

func (c *City) HasSite(siteID string) bool {
	_, found := c.Sites[siteID]
	return found
}

func (c *City) HasRestaurant(siteID, restaurantID string) bool {
	if !c.HasSite(siteID) {
		return false
	}
	// We should only get here if there is a Site with siteID, so this should not crash
	return c.GetSiteById(siteID).HasRestaurant(restaurantID)
}

func (c *City) ClearSites() *City {
	c.Sites = make(map[string]*Site)
	return c
}

func (c *City) ClearRestaurants() *City {
	for k := range c.Sites {
		c.Sites[k].ClearRestaurants()
	}
	return c
}

func (c *City) ClearDishes() *City {
	for k := range c.Sites {
		c.Sites[k].ClearDishes()
	}
	return c
}

func (c *City) GetSiteById(id string) *Site {
	s, found := c.Sites[id]
	if !found {
		debugCity("GetSiteById: %q not found", id)
	}
	return s
}

func (c *City) GetRestaurantById(siteID, restaurantID string) *Restaurant {
	s := c.GetSiteById(siteID)
	if nil == s {
		return nil
	}
	return s.GetRestaurantById(restaurantID)
}

func (c *City) NumSites() int {
	return len(c.Sites)
}

func (c *City) NumRestaurants() int {
	total := 0
	for k := range c.Sites {
		total += c.Sites[k].NumRestaurants()
	}
	return total
}

func (c *City) NumDishes() int {
	total := 0
	for k := range c.Sites {
		total += c.Sites[k].NumDishes()
	}
	return total
}

func (c *City) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(c)
}

func (c *City) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(c)
}

func CityFromJSON(r io.Reader) (*City, error) {
	c := &City{}
	if err := c.Decode(r); err != nil {
		return nil, err
	}
	return c, nil
}
