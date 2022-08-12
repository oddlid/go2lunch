package lunchdata

import (
	"encoding/json"
	"io"
	"regexp"
	"sync"
	"time"
)

type Restaurant struct {
	sync.RWMutex
	Name    string    `json:"restaurant_name"`
	ID      string    `json:"restaurant_id"`
	Url     string    `json:"url,omitempty"`
	Gtag    string    `json:"-"`
	Address string    `json:"address"`
	MapURL  string    `json:"map_url"`
	Parsed  time.Time `json:"scrape_date"`
	Dishes  Dishes    `json:"dishes"`
}

type Restaurants []*Restaurant

// func (rs *Restaurants) Add(r *Restaurant) {
// 	rs = append(rs, r)
// }

func (rs Restaurants) Len() int {
	return len(rs)
}

func (rs Restaurants) NumDishes() int {
	total := 0
	for _, r := range rs {
		total += r.NumDishes()
	}
	return total
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

func (r *Restaurant) Len() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.Dishes)
}

func (r *Restaurant) SubItems() int {
	return r.Len() // just a wrap here. We only have it for name consistency
}

// GetMapUrl returns empty string or .Address as a Google Maps URL
// func (r *Restaurant) GetMapUrl() string {
// 	r.RLock()
// 	defer r.RUnlock()
// 	if r.Address == "" {
// 		return ""
// 	}
// 	return fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%s", url.QueryEscape(r.Address))
// }

// ParsedRFC3339 returns the date in RFC3339 format
func (r *Restaurant) ParsedRFC3339() string {
	r.RLock()
	defer r.RUnlock()
	return r.Parsed.Format(time.RFC3339)
}

// ParsedHumanDate returns a more human readable date/time format, without too much detail
func (r *Restaurant) ParsedHumanDate() string {
	r.RLock()
	defer r.RUnlock()
	return r.Parsed.Format(DATE_FORMAT)
}

func (rs Restaurants) PropagateGtag(tag string) {
	for i := range rs {
		rs[i].PropagateGtag(tag)
	}
}

func (r *Restaurant) PropagateGtag(tag string) {
	r.Lock()
	defer r.Unlock()
	r.Gtag = tag
	for i := range r.Dishes {
		r.Dishes[i].Gtag = tag
	}
}

func (r *Restaurant) AddDish(d *Dish) {
	r.Lock()
	r.Dishes = append(r.Dishes, d)
	r.Unlock()
}

func (r *Restaurant) SetDishes(ds Dishes) {
	r.Lock()
	r.Dishes = ds
	r.Unlock()
}

func (r *Restaurant) ClearDishes() {
	r.Lock()
	r.Dishes = nil
	r.Unlock()
}

func (r *Restaurant) ResetDishes() {
	r.Lock()
	r.Dishes = make(Dishes, 0)
	r.Unlock()
}

func (r *Restaurant) NumDishes() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.Dishes)
}

func (r *Restaurant) HasDishes() bool {
	return r.NumDishes() > 0
}

func (r *Restaurant) GetDishByIndex(idx int) *Dish {
	if idx < 0 || idx >= len(r.Dishes) {
		return nil
	}
	r.Lock()
	dish := r.Dishes[idx]
	r.Unlock()
	return dish
}

func (r *Restaurant) FilterDishesByName(pattern string) (Dishes, error) {
	var ds Dishes
	rx, err := regexp.Compile(pattern)
	if err != nil {
		return ds, err
	}
	for i := range r.Dishes {
		if rx.MatchString(r.Dishes[i].Name) {
			ds = append(ds, r.Dishes[i])
		}
	}
	return ds, nil
}

func (r *Restaurant) FilterDishesByDesc(pattern string) (Dishes, error) {
	var ds Dishes
	rx, err := regexp.Compile(pattern)
	if err != nil {
		return ds, err
	}
	for i := range r.Dishes {
		if rx.MatchString(r.Dishes[i].Desc) {
			ds = append(ds, r.Dishes[i])
		}
	}
	return ds, nil
}

// Takes a function that receives the Dish Price as an argument
// The provided function should return true if the price is considered matching, false if not
func (r *Restaurant) FilterDishesByPrice(f func(int) bool) Dishes {
	var ds Dishes
	for i := range r.Dishes {
		if f(r.Dishes[i].Price) {
			ds = append(ds, r.Dishes[i])
		}
	}
	return ds
}

// This should be a func that combines all other filter funcs as a convenience, but I'm not
// sure how to best solve it atm
func (r *Restaurant) FilterDishes() Dishes {
	// restaurantLog.Debug("Not implemented yet")
	return nil
}

func (r *Restaurant) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}

func (r *Restaurant) Decode(rdr io.Reader) error {
	return json.NewDecoder(rdr).Decode(r)
}

func (rs *Restaurants) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(rs)
}

func (rs *Restaurants) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(rs)
}

func RestaurantFromJSON(r io.Reader) (*Restaurant, error) {
	rs := &Restaurant{}
	if err := rs.Decode(r); err != nil {
		return nil, err
	}
	return rs, nil
}

func RestaurantsFromJSON(r io.Reader) (Restaurants, error) {
	rs := make(Restaurants, 0)
	if err := rs.Decode(r); err != nil {
		return nil, err
	}
	return rs, nil
}
