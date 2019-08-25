package lunchdata

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	//"time"
	//log "github.com/Sirupsen/logrus"
)

type Site struct {
	Name        string                 `json:"site_name"`
	ID          string                 `json:"site_id"` // something unique within the parent city
	Comment     string                 `json:"site_comment,omitempty"`
	Url         string                 `json:"url,omitempty"`
	Gtag        string                 `json:"-"`
	Key         string                 `json:"-"` // validation against submitting scrapers
	Restaurants map[string]*Restaurant `json:"restaurants"`
}

type Sites []Site

func (ss *Sites) Add(s Site) {
	*ss = append(*ss, s)
}

func (ss *Sites) Len() int {
	return len(*ss)
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
	return len(s.Restaurants)
}

func (s *Site) SubItems() int {
	total := 0
	for k := range s.Restaurants {
		total += s.Restaurants[k].SubItems() + 1 // +1 to count the restaurant itself as well
	}
	return total
}

// Reminder of sort of what I think I want...
// A Site instance should be able to calculate its key, given the signing key.
func (s *Site) CalcKey(signKey string) string {
	return ""
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
	s.Gtag = tag
	for k := range s.Restaurants {
		s.Restaurants[k].PropagateGtag(tag)
	}
	return s
}

func (s *Site) ParsedHumanDate() string {
	r := s.getRndRestaurant()
	if r != nil {
		return r.ParsedHumanDate()
	}
	return DATE_FORMAT
}

func (s *Site) AddRestaurant(r Restaurant) *Site {
	s.Restaurants[r.ID] = &r
	return s
}

func (s *Site) DeleteRestaurant(id string) *Site {
	delete(s.Restaurants, id)
	return s
}

func (s *Site) HasRestaurants() bool {
	return len(s.Restaurants) > 0
}

func (s *Site) HasRestaurant(restaurantID string) bool {
	_, found := s.Restaurants[restaurantID]
	return found
}

// Replace existing restaurants with these new given ones
func (s *Site) SetRestaurants(rs Restaurants) *Site {
	s.ClearRestaurants() // otherwise, old restaurants in the new set might linger on
	for _, r := range rs {
		s.AddRestaurant(r)
	}
	return s
}

func (s *Site) ClearRestaurants() *Site {
	s.Restaurants = make(map[string]*Restaurant)
	return s
}

func (s *Site) ClearDishes() *Site {
	for k := range s.Restaurants {
		s.Restaurants[k].ClearDishes()
	}
	return s
}

func (s *Site) GetRestaurantById(id string) *Restaurant {
	r, found := s.Restaurants[id]
	if !found {
		debugSite("GetRestaurantById: %q not found", id)
	}
	return r
}

func (s *Site) NumRestaurants() int {
	return len(s.Restaurants)
}

func (s *Site) NumDishes() int {
	total := 0
	for k := range s.Restaurants {
		total += s.Restaurants[k].NumDishes()
	}
	return total
}

func (s *Site) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(s)
}

func (s *Site) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(s)
}

func (s *Site) SaveJSON(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = s.Encode(w)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func SiteFromJSON(r io.Reader) (*Site, error) {
	s := &Site{}
	if err := s.Decode(r); err != nil {
		return nil, err
	}
	return s, nil
}
