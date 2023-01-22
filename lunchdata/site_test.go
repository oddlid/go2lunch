package lunchdata

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSite_NumRestaurants(t *testing.T) {
	assert.Equal(t, 0, (*Site)(nil).NumRestaurants())

	s := Site{
		Restaurants: Restaurants{{}, {}},
	}
	assert.Equal(t, 2, s.NumRestaurants())
}

// func TestSite_getRndRestaurant(t *testing.T) {
// 	assert.Nil(t, (*Site)(nil).getRndRestaurant())

// 	// map access order is random, so we test with only 1 item here to ensure we can assert correctly
// 	r := Restaurant{}
// 	s := Site{}
// 	assert.Nil(t, s.getRndRestaurant())
// 	s.Restaurants = RestaurantMap{"1": &r}
// 	assert.Same(t, &r, s.getRndRestaurant())
// }

// func TestSite_ParsedHumanDate(t *testing.T) {
// 	assert.Equal(t, dateFormat, (*Site)(nil).ParsedHumanDate())

// 	now := time.Now()
// 	r := Restaurant{ParsedAt: now}
// 	s := Site{}
// 	assert.Equal(t, dateFormat, s.ParsedHumanDate())
// 	s.Restaurants = RestaurantMap{"1": &r}
// 	assert.Equal(t, r.ParsedHumanDate(), s.ParsedHumanDate())
// }

// func TestSite_Get(t *testing.T) {
// 	id := "id"
// 	r := (*Site)(nil).Get(id)
// 	assert.Nil(t, r)

// 	s := Site{}
// 	r = s.Get(id)
// 	assert.Nil(t, r)

// 	restaurant := Restaurant{}
// 	s.Restaurants = RestaurantMap{id: &restaurant}
// 	r = s.Get(id)
// 	assert.NotNil(t, r)
// 	assert.Same(t, &restaurant, r)
// }

func TestSite_NumDishes(t *testing.T) {
	assert.Equal(t, 0, (*Site)(nil).NumDishes())

	s := Site{
		Restaurants: Restaurants{
			{Dishes: Dishes{{}, {}}},
			{Dishes: Dishes{{}, {}}},
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

	rm := Restaurants{{}, {}}
	s.Scraper = &mockSiteScraper{restaurants: rm}
	err = s.RunScraper()
	assert.NoError(t, err)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 2)
}

func Test_Site_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(*Site)(nil).setIDIfEmpty()
	})
	s := Site{}
	s.setIDIfEmpty()
	assert.NotEmpty(t, s.ID)
}
