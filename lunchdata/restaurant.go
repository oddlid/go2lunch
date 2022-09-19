package lunchdata

import (
	"sync"
	"time"
)

type Restaurant struct {
	Name    string    `json:"restaurant_name"`
	ID      string    `json:"restaurant_id"`
	URL     string    `json:"url,omitempty"`
	GTag    string    `json:"-"`
	Address string    `json:"address"`
	MapURL  string    `json:"map_url"`
	Parsed  time.Time `json:"scrape_date"`
	Dishes  Dishes    `json:"dishes"`
	mu      sync.RWMutex
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

func (r *Restaurant) NumDishes() int {
	if r == nil {
		return 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Dishes.Len()
}

func (r *Restaurant) Empty() bool {
	return r.NumDishes() == 0
}

// ParsedRFC3339 returns the date in RFC3339 format
func (r *Restaurant) ParsedRFC3339() string {
	if r == nil {
		return time.Now().Format(time.RFC3339)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Parsed.Format(time.RFC3339)
}

// ParsedHumanDate returns a more human readable date/time format, without too much detail
func (r *Restaurant) ParsedHumanDate() string {
	if r == nil {
		return time.Now().Format(dateFormat)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Parsed.Format(dateFormat)
}

func (r *Restaurant) SetGTag(tag string) *Restaurant {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	r.GTag = tag
	r.Dishes.SetGTag(tag)
	r.mu.Unlock()
	return r
}

func (r *Restaurant) Add(dishes ...*Dish) *Restaurant {
	if r == nil {
		return nil
	}
	if len(dishes) == 0 {
		return r
	}
	r.mu.Lock()
	for _, dish := range dishes {
		if dish != nil {
			r.Dishes = append(r.Dishes, dish)
		}
	}
	r.mu.Unlock()
	return r
}

func (r *Restaurant) Set(ds Dishes) *Restaurant {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	r.Dishes = ds
	r.mu.Unlock()
	return r
}
