package lunchdata

type Site struct {
	// Scraper     SiteScraper `json:"-"`
	Name        string      `json:"name"`
	ID          string      `json:"id"`
	Comment     string      `json:"comment,omitempty"`
	URL         string      `json:"url,omitempty"`
	Restaurants Restaurants `json:"restaurants"`
}

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
