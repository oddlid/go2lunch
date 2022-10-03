package lunchdata

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, 0, (*Site)(nil).NumRestaurants())

	s := Site{
		Restaurants: RestaurantMap{
			"1": {},
			"2": {},
		},
	}
	assert.Equal(t, 2, s.NumRestaurants())
}

func TestSite_Empty(t *testing.T) {
	assert.True(t, (*Site)(nil).Empty())

	s := Site{Restaurants: RestaurantMap{"1": {}}}
	assert.False(t, s.Empty())
}

func TestSite_getRndRestaurant(t *testing.T) {
	assert.Nil(t, (*Site)(nil).getRndRestaurant())

	// map access order is random, so we test with only 1 item here to ensure we can assert correctly
	r := Restaurant{}
	s := Site{}
	assert.Nil(t, s.getRndRestaurant())
	s.Restaurants = RestaurantMap{"1": &r}
	assert.Same(t, &r, s.getRndRestaurant())
}

func TestSite_SetGTag(t *testing.T) {
	assert.Nil(t, (*Site)(nil).SetGTag(""))

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
	assert.Equal(t, dateFormat, (*Site)(nil).ParsedHumanDate())

	now := time.Now()
	r := Restaurant{Parsed: now}
	s := Site{}
	assert.Equal(t, dateFormat, s.ParsedHumanDate())
	s.Restaurants = RestaurantMap{"1": &r}
	assert.Equal(t, r.ParsedHumanDate(), s.ParsedHumanDate())
}

func TestSite_Add(t *testing.T) {
	assert.Nil(t, (*Site)(nil).Add(&Restaurant{}))

	id := "id"
	r := Restaurant{ID: id}
	s := Site{}
	ret := s.Add(nil)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 0)
	ret = s.Add(&r)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 1)
}

func TestSite_Delete(t *testing.T) {
	assert.Nil(t, (*Site)(nil).Delete(""))

	s := Site{}
	// Test delete on nil map
	assert.NotPanics(t, func() { s.Delete("") })

	s.Restaurants = RestaurantMap{"id": {}}
	s.Delete("id")
	assert.Len(t, s.Restaurants, 0)
}

func TestSite_Set(t *testing.T) {
	assert.Nil(t, (*Site)(nil).Set(nil))

	id := "id"
	s := Site{
		Restaurants: RestaurantMap{id: {}},
	}
	ret := s.Set(nil)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 0)

	s.Restaurants[id] = &Restaurant{}
	rs := Restaurants{{ID: "1"}, {ID: "2"}}
	s.Set(rs)
	assert.Len(t, s.Restaurants, 2)
	_, found := s.Restaurants[id]
	assert.False(t, found)
}

func TestSite_Get(t *testing.T) {
	id := "id"
	r := (*Site)(nil).Get(id)
	assert.Nil(t, r)

	s := Site{}
	r = s.Get(id)
	assert.Nil(t, r)

	restaurant := Restaurant{}
	s.Restaurants = RestaurantMap{id: &restaurant}
	r = s.Get(id)
	assert.NotNil(t, r)
	assert.Same(t, &restaurant, r)
}

func TestSite_NumDishes(t *testing.T) {
	assert.Equal(t, 0, (*Site)(nil).NumDishes())

	s := Site{
		Restaurants: RestaurantMap{
			"1": {Dishes: Dishes{{}, {}}},
			"2": {Dishes: Dishes{{}, {}}},
		},
	}
	assert.Equal(t, 4, s.NumDishes())
}

func TestSite_SetScraper(t *testing.T) {
	assert.Nil(t, (*Site)(nil).SetScraper(nil))

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
	err := (*Site)(nil).RunScraper()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilSite)

	s := Site{}
	err = s.RunScraper()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilScraper)

	scrapeErr := errors.New("scrape error")
	s.Scraper = &mockSiteScraper{err: scrapeErr}
	err = s.RunScraper()
	assert.Error(t, err)
	assert.ErrorIs(t, err, scrapeErr)
	assert.Nil(t, s.Restaurants)

	rs := Restaurants{{ID: "1"}, {ID: "2"}}
	s.Scraper = &mockSiteScraper{restaurants: rs}
	err = s.RunScraper()
	assert.NoError(t, err)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 2)
}
