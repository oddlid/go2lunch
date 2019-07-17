package lunchdata

import (
	"encoding/json"
	"io"
	"time"
)

const (
	DATE_FORMAT string = "2006-01-02 15:04"
)

type Restaurant struct {
	Name   string    `json:"restaurant_name"`
	ID     string    `json:"restaurant_id"`
	Url    string    `json:"url"`
	Parsed time.Time `json:"scrape_date"`
	Dishes Dishes    `json:"dishes"`
}

type Restaurants []Restaurant

func (rs *Restaurants) Add(r Restaurant) {
	*rs = append(*rs, r)
}

func NewRestaurant(name, id, url string, parsed time.Time) *Restaurant {
	return &Restaurant{
		Name:   name,
		ID:     id,
		Url:    url,
		Parsed: parsed,
		Dishes: make(Dishes, 0),
	}
}

// ParsedRFC3339 returns the date in RFC3339 format
func (r Restaurant) ParsedRFC3339() string {
	return r.Parsed.Format(time.RFC3339)
}

// ParsedHumanDate returns a more human readable date/time format, without too much detail
func (r Restaurant) ParsedHumanDate() string {
	return r.Parsed.Format(DATE_FORMAT)
}

func (r *Restaurant) AddDish(d Dish) *Restaurant {
	r.Dishes = append(r.Dishes, d)
	return r
}

func (r *Restaurant) SetDishes(ds Dishes) *Restaurant {
	r.Dishes = ds
	return r
}

func (r *Restaurant) NumDishes() int {
	return len(r.Dishes)
}

func (rs *Restaurant) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(rs)
}

func (rs *Restaurant) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(rs)
}

func RestaurantFromJSON(r io.Reader) (*Restaurant, error) {
	rs := &Restaurant{}
	if err := rs.Decode(r); err != nil {
		return nil, err
	}
	return rs, nil
}
