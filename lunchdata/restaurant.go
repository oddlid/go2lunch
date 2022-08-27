package lunchdata

import (
	"sync"
	"time"
)

type Restaurant struct {
	mu      sync.RWMutex
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
type RestaurantMap map[string]*Restaurant

/*** funcs for Restaurant ***/

func NewRestaurant(name, id, url string, parsed time.Time) *Restaurant {
	return &Restaurant{
		Name:   name,
		ID:     id,
		URL:    url,
		Parsed: parsed,
		Dishes: make(Dishes, 0),
	}
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

func (r *Restaurant) PropagateGtag(tag string) {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Gtag = tag
	for i := range r.Dishes {
		r.Dishes[i].Gtag = tag
	}
}

func (r *Restaurant) AddDishes(ds ...*Dish) {
	if r == nil || len(ds) == 0 {
		return
	}
	r.mu.Lock()
	r.Dishes = append(r.Dishes, ds...)
	r.mu.Unlock()
}

func (r *Restaurant) SetDishes(ds Dishes) {
	if r == nil {
		return
	}
	r.mu.Lock()
	r.Dishes = ds
	r.mu.Unlock()
}

func (r *Restaurant) NumDishes() int {
	if r == nil {
		return 0
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Dishes.Len()
}

/*** funcs for Restaurants ***/

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

func (rs Restaurants) PropagateGtag(tag string) {
	for i := range rs {
		rs[i].PropagateGtag(tag)
	}
}

func (rs Restaurants) AsMap() RestaurantMap {
	rMap := make(RestaurantMap)
	for i := range rs {
		rMap.Add(rs[i])
	}
	return rMap
}

/*** funcs for RestaurantMap ***/

func (rm RestaurantMap) Len() int {
	return len(rm)
}

func (rm RestaurantMap) Add(r *Restaurant) {
	if rm == nil || r == nil {
		return
	}
	rm[r.ID] = r
}

func (rm RestaurantMap) Delete(id string) {
	delete(rm, id)
}
