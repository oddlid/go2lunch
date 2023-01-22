package lunchdata

import (
	"errors"

	"github.com/google/uuid"
)

type Site struct {
	Scraper     SiteScraper `json:"-"`
	Name        string      `json:"name"`
	ID          string      `json:"id"`
	Comment     string      `json:"comment,omitempty"`
	URL         string      `json:"url,omitempty"`
	GTag        string      `json:"-"`
	Restaurants Restaurants `json:"restaurants"`
}

var (
	errNilSite    = errors.New("site is nil")
	errNilScraper = errors.New("scraper is nil")
)

func (s *Site) NumRestaurants() int {
	if s == nil {
		return 0
	}
	return s.Restaurants.Len()
}

func (s *Site) NumDishes() int {
	if s == nil {
		return 0
	}
	return s.Restaurants.NumDishes()
}

func (s *Site) Get(f RestaurantMatch) *Restaurant {
	if s == nil {
		return nil
	}
	return s.Restaurants.Get(f)
}

func (s *Site) GetByID(id string) *Restaurant {
	if s == nil {
		return nil
	}
	return s.Restaurants.GetByID(id)
}

func (s *Site) ParsedHumanDate() string {
	if s == nil {
		return dateFormat
	}
	if r := s.Restaurants.first(); r != nil {
		return r.ParsedHumanDate()
	}
	return dateFormat
}

func (s *Site) SetScraper(scraper SiteScraper) *Site {
	if s == nil {
		return nil
	}
	s.Scraper = scraper
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
	s.Restaurants = rs

	return nil
}

func (s *Site) setIDIfEmpty() {
	if s == nil {
		return
	}
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	s.Restaurants.setIDIfEmpty()
}
