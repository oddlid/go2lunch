package lunchdata

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSite_NumRestaurants(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, (*Site)(nil).NumRestaurants())

	s := Site{
		Restaurants: Restaurants{{}, {}},
	}
	assert.Equal(t, 2, s.NumRestaurants())
}

func Test_Site_ParsedHumanDate(t *testing.T) {
	t.Parallel()
	assert.Equal(t, dateFormat, (*Site)(nil).ParsedHumanDate())

	s := Site{}
	assert.Equal(t, dateFormat, s.ParsedHumanDate())

	s.Restaurants = Restaurants{{}}
	assert.Equal(t, time.Time{}.Format(dateFormat), s.ParsedHumanDate())
}

func TestSite_NumDishes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, (*Site)(nil).NumDishes())

	s := Site{
		Restaurants: Restaurants{
			{Dishes: Dishes{{}, {}}},
			{Dishes: Dishes{{}, {}}},
		},
	}
	assert.Equal(t, 4, s.NumDishes())
}

func Test_Site_Get(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (*Site)(nil).Get(nil))

	const id = `blah`
	s := Site{Restaurants: Restaurants{{ID: id}}}
	assert.Same(t, &s.Restaurants[0], s.Get(func(r Restaurant) bool { return r.ID == id }))
}

func Test_Site_GetByID(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (*Site)(nil).GetByID(""))

	const id = `blah`
	s := Site{Restaurants: Restaurants{{ID: id}}}
	assert.Same(t, &s.Restaurants[0], s.GetByID(id))
}

func TestSite_SetScraper(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	err := (*Site)(nil).RunScraper()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilSite)

	mockScraper := mockSiteScraper{
		err: errors.New("scrape error"),
	}
	s := Site{}
	err = s.RunScraper()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilScraper)

	s.Scraper = &mockScraper
	err = s.RunScraper()
	assert.Error(t, err)
	assert.ErrorIs(t, err, mockScraper.err)
	assert.Nil(t, s.Restaurants)

	mockScraper.restaurants = Restaurants{{}, {}}
	mockScraper.err = nil
	err = s.RunScraper()
	assert.NoError(t, err)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 2)
}

func Test_Site_setIDIfEmpty(t *testing.T) {
	t.Parallel()
	assert.NotPanics(t, func() {
		(*Site)(nil).setIDIfEmpty()
	})
	s := Site{}
	s.setIDIfEmpty()
	assert.NotEmpty(t, s.ID)
}
