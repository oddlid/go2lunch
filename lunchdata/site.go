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
	ID          string                 `json:"site_id"` // something unique within this site
	Comment     string                 `json:"site_comment,omitempty"`
	Restaurants map[string]*Restaurant `json:"restaurants"`
}

type Sites []Site

func (ss *Sites) Add(s Site) {
	*ss = append(*ss, s)
}

func NewSite(name, id, comment string) *Site {
	return &Site{
		Name:        name,
		ID:          id,
		Comment:     comment,
		Restaurants: make(map[string]*Restaurant),
	}
}

// Just deliver the first restaurant we find.
// Convenience method for inheriting timestamp
func (s *Site) getRndRestaurant() *Restaurant {
	for _, v := range s.Restaurants {
		return v
	}
	return nil
}

func (s *Site) ParsedHumanDate() string {
	r := s.getRndRestaurant()
	if r != nil {
		return r.ParsedHumanDate()
	}
	//return "0000-00-00 00:00" // "undefined"
	return DATE_FORMAT
}

func (s *Site) AddRestaurant(r Restaurant) *Site {
	//log.Debugf("AddRestaurant(): %+v", r)
	s.Restaurants[r.ID] = &r
	return s
}

//func (s *Site) DeleteRestaurant

func (s *Site) HasRestaurants() bool {
	return len(s.Restaurants) > 0
}

// Replace existing restaurants with these new given ones
func (s *Site) SetRestaurants(rs Restaurants) *Site {
	//s.Restaurants = make(map[string]*Restaurant)
	for _, r := range rs {
		s.AddRestaurant(r)
	}
	return s
}

//func (s *Site) HasRestaurant(id string) bool {
//	_, found := s.Restaurants[id]
//	return found
//}

func (s *Site) GetRestaurantById(id string) *Restaurant {
	return s.Restaurants[id]
}

//func (s *Site) GetRestaurantOrNew(name, id string) *Restaurant {
//	r := s.GetRestaurant(id)
//	if r == nil {
//		r = NewRestaurant(name, id, "", time.Now())
//		s.AddRestaurant(r)
//	}
//	return r
//}

func (s *Site) NumRestaurants() int {
	return len(s.Restaurants)
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
