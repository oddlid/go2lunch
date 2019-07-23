package lunchdata

import (
	"encoding/json"
	"io"
)

type Country struct {
	Name   string           `json:"country_name"`
	ID     string           `json:"country_id"` // preferrably international country code, like "se", "no", and so on
	Cities map[string]*City `json:"cities"`
	Gtag   string           `json:"-"`
}

type Countries []Country

func (cs *Countries) Add(c Country) {
	*cs = append(*cs, c)
}

func NewCountry(name, id, tag string) *Country {
	return &Country{
		Name:   name,
		ID:     id,
		Cities: make(map[string]*City),
		Gtag:   tag,
	}
}

func (c *Country) AddCity(city City) *Country {
	c.Cities[city.ID] = &city
	return c
}

func (c *Country) GetCityById(id string) *City {
	return c.Cities[id]
}

func (c *Country) NumCities() int {
	return len(c.Cities)
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
