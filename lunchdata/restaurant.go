package lunchdata

import (
	"sync"
	"time"
)

type Restaurant struct {
	sync.RWMutex
	Name    string    `json:"restaurant_name"`
	ID      string    `json:"restaurant_id"`
	URL     string    `json:"url,omitempty"`
	Gtag    string    `json:"-"`
	Address string    `json:"address"`
	MapURL  string    `json:"map_url"`
	Parsed  time.Time `json:"scrape_date"`
	Dishes  Dishes    `json:"dishes"`
}

type Restaurants []*Restaurant

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
		URL:    url,
		Parsed: parsed,
		Dishes: make(Dishes, 0),
	}
}

func (r *Restaurant) Len() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.Dishes)
}

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
	return r.Parsed.Format(dateFormat)
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

func (r *Restaurant) NumDishes() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.Dishes)
}

func (r *Restaurant) HasDishes() bool {
	return r.NumDishes() > 0
}
