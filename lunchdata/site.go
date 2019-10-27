package lunchdata

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"sync"
	//"time"
	//log "github.com/Sirupsen/logrus"
)

type Site struct {
	sync.RWMutex
	Name        string                 `json:"site_name"`
	ID          string                 `json:"site_id"` // something unique within the parent city
	Comment     string                 `json:"site_comment,omitempty"`
	Url         string                 `json:"url,omitempty"`
	Gtag        string                 `json:"-"`
	Key         string                 `json:"-"` // validation against submitting scrapers
	Restaurants map[string]*Restaurant `json:"restaurants"`
	Scraper     SiteScraper            `json:"-"`
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
	s.RLock()
	defer s.RUnlock()
	return len(s.Restaurants)
}

func (s *Site) SubItems() int {
	total := 0
	s.RLock()
	for k := range s.Restaurants {
		total += s.Restaurants[k].SubItems() + 1 // +1 to count the restaurant itself as well
	}
	s.RUnlock()
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
	return DATE_FORMAT
}

func (s *Site) AddRestaurant(r Restaurant) *Site {
	s.Lock()
	s.Restaurants[r.ID] = &r
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
	//s.ClearRestaurants() // otherwise, old restaurants in the new set might linger on
	//for _, r := range rs {
	//	// maybe we should re-implment AddRestaurant here, to avoid lock/unlock for each call...?
	//	s.AddRestaurant(r)
	//}
	// So, here we do exactly the same as ClearRestaurants() followed by a series of AddRestaurant(),
	// but with only one sequence of Lock()/Unlock(), as this should be faster and safer
	// 2019-10-23 22:09: The implementation below did not work. Got duplicates of first or last Restaurant
	// all over... Weird....
	//s.Lock()
	//s.Restaurants = make(map[string]*Restaurant)
	//for _, r := range rs {
	//	debugSite("SetRestaurants(): Adding %q", r.Name)
	//	s.Restaurants[r.ID] = &r
	//}
	//s.Unlock()

	// New attempt, just with an inner func
	// 2019-10-23 22:23
	// And this made all the difference!
	// Seems there is some black magic going on behind the scenes in regard to pointer dereferencing or
	// something.
	add := func(r Restaurant) { // just the same as AddRestaurant, but without locks
		s.Restaurants[r.ID] = &r
	}
	s.Lock()
	s.Restaurants = make(map[string]*Restaurant) // just the same as ClearRestaurants, but without locks
	for _, r2 := range rs {
		add(r2)
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
		s.Restaurants[k].ClearDishes()
	}
	s.Unlock()
	return s
}

func (s *Site) GetRestaurantById(id string) *Restaurant {
	s.RLock()
	r, found := s.Restaurants[id]
	s.RUnlock()
	if !found {
		debugSite("GetRestaurantById: %q not found", id)
	}
	return r
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

func (s *Site) RunScraper(wg *sync.WaitGroup) {
	defer wg.Done()
	if nil == s.Scraper {
		debugSite("RunScraper(): %q has no scraper instance configured. Returning.", s.ID)
		return
	}
	rs, err := s.Scraper.Scrape()
	if nil != err {
		errorSite("%q: Error running scraper: %s", s.ID, err)
		return
	}
	s.SetRestaurants(rs)
}
