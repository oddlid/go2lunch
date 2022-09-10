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

func (r *Restaurant) AddDishes(dishes ...*Dish) *Restaurant {
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

func (r *Restaurant) SetDishes(ds Dishes) *Restaurant {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	r.Dishes = ds
	r.mu.Unlock()
	return r
}

/*** funcs for Restaurants ***/

func (rs Restaurants) Len() int {
	return len(rs)
}

func (rs Restaurants) Empty() bool {
	return rs.Len() == 0
}

func (rs Restaurants) NumDishes() int {
	total := 0
	for _, r := range rs {
		total += r.NumDishes()
	}
	return total
}

func (rs Restaurants) Total() int {
	return rs.Len() + rs.NumDishes()
}

func (rs Restaurants) SetGTag(tag string) {
	for _, r := range rs {
		r.SetGTag(tag)
	}
}

func (rs Restaurants) AsMap() RestaurantMap {
	rMap := make(RestaurantMap)
	rMap.Add(rs...)
	return rMap
}

/*** funcs for RestaurantMap ***/

func (rm RestaurantMap) Len() int {
	return len(rm)
}

func (rm RestaurantMap) Empty() bool {
	return rm.Len() == 0
}

func (rm RestaurantMap) NumDishes() int {
	total := 0
	for _, r := range rm {
		total += r.NumDishes()
	}
	return total
}

func (rm RestaurantMap) Total() int {
	return rm.Len() + rm.NumDishes()
}

func (rm RestaurantMap) Add(restaurants ...*Restaurant) {
	if rm == nil {
		return
	}
	for _, r := range restaurants {
		if r != nil {
			rm[r.ID] = r
		}
	}
}

func (rm RestaurantMap) Delete(ids ...string) {
	for _, id := range ids {
		delete(rm, id)
	}
}

func (rm RestaurantMap) SetGTag(tag string) {
	for _, r := range rm {
		r.SetGTag(tag)
	}
}
