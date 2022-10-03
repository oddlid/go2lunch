package lunchdata

import (
	"errors"
	"sync"
)

type Site struct {
	Restaurants RestaurantMap `json:"restaurants"`
	Scraper     SiteScraper   `json:"-"`
	Name        string        `json:"site_name"`
	ID          string        `json:"site_id"` // something unique within the parent city
	Comment     string        `json:"site_comment,omitempty"`
	URL         string        `json:"url,omitempty"`
	GTag        string        `json:"-"`
	mu          sync.RWMutex
}

var (
	errNilSite      = errors.New("site is nil")
	errNilScraper   = errors.New("scraper is nil")
	errNilWaitGroup = errors.New("passed sync.WaitGroup is nil")
)

func NewSite(name, id, comment string) *Site {
	return &Site{
		Name:        name,
		ID:          id,
		Comment:     comment,
		Restaurants: make(RestaurantMap),
	}
}

func (s *Site) NumRestaurants() int {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Restaurants.Len()
}

func (s *Site) Empty() bool {
	return s.NumRestaurants() == 0
}

func (s *Site) NumDishes() int {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Restaurants.NumDishes()
}

// Just deliver the first restaurant we find.
// Convenience method for inheriting timestamp
func (s *Site) getRndRestaurant() *Restaurant {
	if s == nil {
		return nil
	}
	for _, v := range s.Restaurants {
		return v
	}
	return nil
}

func (s *Site) SetGTag(tag string) *Site {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	s.GTag = tag
	s.Restaurants.SetGTag(tag)
	s.mu.Unlock()
	return s
}

func (s *Site) ParsedHumanDate() string {
	r := s.getRndRestaurant() // safe to call on nil receiver
	if r != nil {
		return r.ParsedHumanDate()
	}
	return dateFormat
}

func (s *Site) Add(restaurants ...*Restaurant) *Site {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	if s.Restaurants == nil {
		s.Restaurants = make(RestaurantMap)
	}
	s.Restaurants.Add(restaurants...)
	s.mu.Unlock()
	return s
}

func (s *Site) Delete(ids ...string) *Site {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	s.Restaurants.Delete(ids...)
	s.mu.Unlock()
	return s
}

func (s *Site) Set(rs Restaurants) *Site {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	s.Restaurants = rs.AsMap()
	s.mu.Unlock()
	return s
}

func (s *Site) Get(id string) *Restaurant {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Restaurants.Get(id)
}

func (s *Site) SetScraper(scraper SiteScraper) *Site {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	s.Scraper = scraper
	s.mu.Unlock()
	return s
}

func (s *Site) RunScraper() error {
	if s == nil {
		return errNilSite
	}
	if s.Scraper == nil {
		return errNilScraper
	}
	rs, err := s.Scraper.Scrape()
	if err != nil {
		return err
	}
	s.Set(rs)

	return nil
}
