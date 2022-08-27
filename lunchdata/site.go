package lunchdata

import (
	"sync"
)

type Site struct {
	sync.RWMutex
	Name        string                 `json:"site_name"`
	ID          string                 `json:"site_id"` // something unique within the parent city
	Comment     string                 `json:"site_comment,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Gtag        string                 `json:"-"`
	Key         string                 `json:"-"` // validation against submitting scrapers
	Restaurants map[string]*Restaurant `json:"restaurants"`
	Scraper     SiteScraper            `json:"-"`
}

type Sites []*Site

func (ss Sites) Len() int {
	return len(ss)
}

func NewSite(name, id, comment string) *Site {
	return &Site{
		Name:        name,
		ID:          id,
		Comment:     comment,
		Restaurants: make(map[string]*Restaurant),
	}
}

func (s *Site) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.Restaurants)
}

func (s *Site) SubItems() int {
	total := 0
	s.RLock()
	for k := range s.Restaurants {
		total += s.Restaurants[k].NumDishes() + 1 // +1 to count the restaurant itself as well
	}
	s.RUnlock()
	return total
}

// Just deliver the first restaurant we find.
// Convenience method for inheriting timestamp
func (s *Site) getRndRestaurant() *Restaurant {
	for _, v := range s.Restaurants {
		return v
	}
	return nil
}

func (s *Site) PropagateGtag(tag string) *Site {
	s.Lock()
	s.Gtag = tag
	for k := range s.Restaurants {
		s.Restaurants[k].PropagateGtag(tag)
	}
	s.Unlock()
	return s
}

func (s *Site) ParsedHumanDate() string {
	r := s.getRndRestaurant()
	if r != nil {
		return r.ParsedHumanDate()
	}
	return dateFormat
}

func (s *Site) AddRestaurant(r *Restaurant) *Site {
	s.Lock()
	s.Restaurants[r.ID] = r
	s.Unlock()
	return s
}

func (s *Site) DeleteRestaurant(id string) *Site {
	s.Lock()
	delete(s.Restaurants, id)
	s.Unlock()
	return s
}

func (s *Site) HasRestaurants() bool {
	s.RLock()
	defer s.RUnlock()
	return len(s.Restaurants) > 0
}

func (s *Site) HasRestaurant(restaurantID string) bool {
	s.RLock()
	_, found := s.Restaurants[restaurantID]
	s.RUnlock()
	return found
}

// Replace existing restaurants with these new given ones
func (s *Site) SetRestaurants(rs Restaurants) *Site {
	s.Lock()
	s.Restaurants = make(map[string]*Restaurant)
	for i := range rs {
		s.Restaurants[rs[i].ID] = rs[i]
	}
	s.Unlock()
	return s
}

func (s *Site) ClearRestaurants() *Site {
	s.Lock()
	s.Restaurants = make(map[string]*Restaurant)
	s.Unlock()
	return s
}

func (s *Site) ClearDishes() *Site {
	s.Lock()
	for k := range s.Restaurants {
		s.Restaurants[k].SetDishes(nil)
	}
	s.Unlock()
	return s
}

func (s *Site) GetRestaurantByID(id string) *Restaurant {
	s.RLock()
	defer s.RUnlock()
	return s.Restaurants[id]
}

func (s *Site) NumRestaurants() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.Restaurants)
}

func (s *Site) NumDishes() int {
	total := 0
	s.RLock()
	for k := range s.Restaurants {
		total += s.Restaurants[k].NumDishes()
	}
	s.RUnlock()
	return total
}

func (s *Site) RunScraper(wg *sync.WaitGroup) {
	defer wg.Done()
	if s.Scraper == nil {
		return
	}
	rs, err := s.Scraper.Scrape()
	if err != nil {
		return
	}
	s.SetRestaurants(rs)
}
