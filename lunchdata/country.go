package lunchdata

import (
	"encoding/json"
	"io"
)

type Country struct {
	Name   string           `json:"country_name"`
	ID     string           `json:"country_id"` // preferrably international country code, like "se", "no", and so on
	Gtag   string           `json:"-"`
	Cities map[string]*City `json:"cities"`
}

type Countries []Country

func (cs *Countries) Add(c Country) {
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
	return len(c.Cities)
}

func (c *Country) SubItems() int {
	total := 0
	for k := range c.Cities {
		total += c.Cities[k].SubItems() + 1 // +1 to count the City itself as well
	}
	return total
}

func (c *Country) PropagateGtag(tag string) *Country {
	c.Gtag = tag
	for k := range c.Cities {
		c.Cities[k].PropagateGtag(tag)
	}
	return c
}

func (c *Country) AddCity(city City) *Country {
	c.Cities[city.ID] = &city
	return c
}

func (c *Country) DeleteCity(id string) *Country {
	delete(c.Cities, id)
	return c
}

func (c *Country) HasCities() bool {
	return len(c.Cities) > 0
}

func (c *Country) HasCity(cityID string) bool {
	_, found := c.Cities[cityID]
	return found
}

func (c *Country) HasSite(cityID, siteID string) bool {
	if !c.HasCity(cityID) {
		return false
	}
	return c.GetCityById(cityID).HasSite(siteID)
}

func (c *Country) HasRestaurant(cityID, siteID, restaurantID string) bool {
	if !c.HasSite(cityID, siteID) {
		return false
	}
	return c.GetSiteById(cityID, siteID).HasRestaurant(restaurantID)
}

func (c *Country) ClearCities() *Country {
	c.Cities = make(map[string]*City)
	return c
}

func (c *Country) ClearSites() *Country {
	for k := range c.Cities {
		c.Cities[k].ClearSites()
	}
	return c
}

func (c *Country) ClearRestaurants() *Country {
	for k := range c.Cities {
		c.Cities[k].ClearRestaurants()
	}
	return c
}

func (c *Country) ClearDishes() *Country {
	for k := range c.Cities {
		c.Cities[k].ClearDishes()
	}
	return c
}

func (c *Country) GetCityById(id string) *City {
	city, found := c.Cities[id]
	if !found {
		debugCountry("GetCityById: %q not found", id)
	}
	return city
}

func (c *Country) GetSiteById(cityID, siteID string) *Site {
	city := c.GetCityById(cityID)
	if nil == city {
		return nil
	}
	return city.GetSiteById(siteID)
}

func (c *Country) GetRestaurantById(cityID, siteID, restaurantID string) *Restaurant {
	city := c.GetCityById(cityID)
	if nil == city {
		return nil
	}
	return city.GetRestaurantById(siteID, restaurantID)
}

func (c *Country) NumCities() int {
	return len(c.Cities)
}

func (c *Country) NumSites() int {
	total := 0
	for k := range c.Cities {
		total += c.Cities[k].NumSites()
	}
	return total
}

func (c *Country) NumRestaurants() int {
	total := 0
	for k := range c.Cities {
		total += c.Cities[k].NumRestaurants()
	}
	return total
}

func (c *Country) NumDishes() int {
	total := 0
	for k := range c.Cities {
		total += c.Cities[k].NumDishes()
	}
	return total
}

func (c *Country) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(c)
}

func (c *Country) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(c)
}

func CountryFromJSON(r io.Reader) (*Country, error) {
	c := &Country{}
	if err := c.Decode(r); err != nil {
		return nil, err
	}
	return c, nil
}
