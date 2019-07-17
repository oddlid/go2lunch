package lunchdata

import (
	"encoding/json"
	"io"
)

type City struct {
	Name  string           `json:"city_name"`
	ID    string           `json:"city_id"` // e.g. osl, gbg or something like the airlines use
	Sites map[string]*Site `json:"sites"`
}

type Cities []City

func (cs *Cities) Add(c City) {
	*cs = append(*cs, c)
}

func NewCity(name, id string) *City {
	return &City{
		Name:  name,
		ID:    id,
		Sites: make(map[string]*Site),
	}
}

func (c *City) AddSite(s Site) *City {
	c.Sites[s.ID] = &s
	return c
}

func (c *City) GetSiteById(id string) *Site {
	return c.Sites[id]
}

//func (c *City) HasSite() {
//}

func (c *City) NumSites() int {
	return len(c.Sites)
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
