package lunchdata

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockSiteScraper struct {
	err         error
	countryID   string
	cityID      string
	siteID      string
	restaurants Restaurants
}

func (s *mockSiteScraper) Scrape() (Restaurants, error) {
	return s.restaurants, s.err
}

func (s *mockSiteScraper) GetCountryID() string {
	return s.countryID
}

func (s *mockSiteScraper) GetCityID() string {
	return s.cityID
}

func (s *mockSiteScraper) GetSiteID() string {
	return s.siteID
}

func TestNewSite(t *testing.T) {
	id := "id"
	name := "name"
	comment := "comment"
	s := NewSite(name, id, comment)
	assert.NotNil(t, s)
	assert.IsType(t, (*Site)(nil), s)
	assert.Equal(t, id, s.ID)
	assert.Equal(t, name, s.Name)
	assert.Equal(t, comment, s.Comment)
	assert.NotNil(t, s.Restaurants)
}

func TestSite_NumRestaurants(t *testing.T) {
	var nilSite *Site
	assert.Equal(t, 0, nilSite.NumRestaurants())

	s := Site{
		Restaurants: RestaurantMap{
			"1": {},
			"2": {},
		},
	}
	assert.Equal(t, 2, s.NumRestaurants())
}

func TestSite_Empty(t *testing.T) {
	var nilSite *Site
	assert.True(t, nilSite.Empty())

	s := Site{Restaurants: RestaurantMap{"1": {}}}
	assert.False(t, s.Empty())
}

func TestSite_getRndRestaurant(t *testing.T) {
	var nilSite *Site
	assert.Nil(t, nilSite.getRndRestaurant())

	// map access order is random, so we test with only 1 item here to ensure we can assert correctly
	r := Restaurant{}
	s := Site{}
	assert.Nil(t, s.getRndRestaurant())
	s.Restaurants = RestaurantMap{"1": &r}
	assert.Same(t, &r, s.getRndRestaurant())
}

func TestSite_SetGTag(t *testing.T) {
	var nilSite *Site
	assert.Nil(t, nilSite.SetGTag(""))

	gtag := "sometag"
	s := Site{
		Restaurants: RestaurantMap{
			"1": {Dishes: Dishes{{}, {}}},
			"2": {Dishes: Dishes{{}, {}}},
		},
	}
	ret := s.SetGTag(gtag)
	assert.Same(t, &s, ret)
	assert.Equal(t, gtag, ret.GTag)
	for _, r := range s.Restaurants {
		assert.Equal(t, gtag, r.GTag)
		for _, d := range r.Dishes {
			assert.Equal(t, gtag, d.GTag)
		}
	}
}

func TestSite_ParsedHumanDate(t *testing.T) {
	var nilSite *Site
	assert.Equal(t, dateFormat, nilSite.ParsedHumanDate())

	now := time.Now()
	r := Restaurant{Parsed: now}
	s := Site{}
	assert.Equal(t, dateFormat, s.ParsedHumanDate())
	s.Restaurants = RestaurantMap{"1": &r}
	assert.Equal(t, r.ParsedHumanDate(), s.ParsedHumanDate())
}

func TestSite_AddRestaurants(t *testing.T) {
	var nilSite *Site
	assert.Nil(t, nilSite.AddRestaurants(nil))

	id := "id"
	r := Restaurant{ID: id}
	s := Site{}
	ret := s.AddRestaurants(nil)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 0)
	ret = s.AddRestaurants(&r)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 1)
}

func TestSite_DeleteRestaurants(t *testing.T) {
	var nilSite *Site
	assert.Nil(t, nilSite.DeleteRestaurants(""))

	s := Site{}
	// Test delete on nil map
	assert.NotPanics(t, func() { s.DeleteRestaurants("") })

	s.Restaurants = RestaurantMap{"id": {}}
	s.DeleteRestaurants("id")
	assert.Len(t, s.Restaurants, 0)
}

func TestSite_SetRestaurants(t *testing.T) {
	var nilSite *Site
	assert.Nil(t, nilSite.SetRestaurants(nil))

	id := "id"
	s := Site{
		Restaurants: RestaurantMap{id: {}},
	}
	ret := s.SetRestaurants(nil)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 0)

	s.Restaurants[id] = &Restaurant{}
	rs := Restaurants{{ID: "1"}, {ID: "2"}}
	s.SetRestaurants(rs)
	assert.Len(t, s.Restaurants, 2)
	_, found := s.Restaurants[id]
	assert.False(t, found)
}

func TestSite_GetRestaurantByID(t *testing.T) {
	id := "id"
	var nilSite *Site
	r, err := nilSite.GetRestaurantByID(id)
	assert.Nil(t, r)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilSite)

	s := Site{}
	r, err = s.GetRestaurantByID(id)
	assert.Nil(t, r)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errRestaurantNotFound)

	restaurant := Restaurant{}
	s.Restaurants = RestaurantMap{id: &restaurant}
	r, err = s.GetRestaurantByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Same(t, &restaurant, r)
}

func TestSite_NumDishes(t *testing.T) {
	var nilSite *Site
	assert.Equal(t, 0, nilSite.NumDishes())

	s := Site{
		Restaurants: RestaurantMap{
			"1": {Dishes: Dishes{{}, {}}},
			"2": {Dishes: Dishes{{}, {}}},
		},
	}
	assert.Equal(t, 4, s.NumDishes())
}

func TestSite_SetScraper(t *testing.T) {
	var nilSite *Site
	assert.Nil(t, nilSite.SetScraper(nil))

	s := Site{}
	scraper := &mockSiteScraper{}
	ret := s.SetScraper(scraper)
	assert.NotNil(t, ret)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Scraper)
	assert.Same(t, scraper, s.Scraper)

	assert.NotNil(t, s.SetScraper(nil))
	assert.Nil(t, s.Scraper)
}

func TestSite_RunScraper(t *testing.T) {
	var nilSite *Site
	err := nilSite.RunScraper(nil)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilSite)

	s := Site{}
	err = s.RunScraper(nil)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilWaitGroup)

	wg := sync.WaitGroup{}
	wg.Add(1)
	err = s.RunScraper(&wg)
	wg.Wait()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNoScraper)

	scrapeErr := errors.New("scrape error")
	s.Scraper = &mockSiteScraper{err: scrapeErr}
	wg.Add(1)
	err = s.RunScraper(&wg)
	wg.Wait()
	assert.Error(t, err)
	assert.ErrorIs(t, err, scrapeErr)
	assert.Nil(t, s.Restaurants)

	rs := Restaurants{{ID: "1"}, {ID: "2"}}
	s.Scraper = &mockSiteScraper{restaurants: rs}
	wg.Add(1)
	err = s.RunScraper(&wg)
	wg.Wait()
	assert.NoError(t, err)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 2)
}
