package lunchdata

import (
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
