package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSites_Len(t *testing.T) {
	assert.Equal(t, 0, (Sites)(nil).Len())

	s := Sites{{}, {}}
	assert.Equal(t, 2, s.Len())
}

func TestSites_NumRestaurants(t *testing.T) {
	assert.Equal(t, 0, (Sites)(nil).NumRestaurants())

	ss := Sites{{Restaurants: Restaurants{{}}}}
	assert.Equal(t, 1, ss.NumRestaurants())
}

func TestSites_NumDishes(t *testing.T) {
	assert.Equal(t, 0, (Sites)(nil).NumDishes())

	ss := Sites{
		{Restaurants: Restaurants{{Dishes: Dishes{{}, {}}}}},
		{Restaurants: Restaurants{{Dishes: Dishes{{}, {}}}}},
	}
	assert.Equal(t, 4, ss.NumDishes())
}

func TestSites_Total(t *testing.T) {
	assert.Equal(t, 0, (Sites)(nil).Total())

	ss := Sites{
		{Restaurants: Restaurants{{Dishes: Dishes{{}, {}}}}},
		{Restaurants: Restaurants{{Dishes: Dishes{{}, {}}}}},
	}
	assert.Equal(t, 8, ss.Total())
}

func TestSites_setGTag(t *testing.T) {
	assert.NotPanics(t, func() { (Sites)(nil).setGTag("") })

	ss := Sites{
		{Restaurants: Restaurants{{Dishes: Dishes{{}, {}}}}},
		{Restaurants: Restaurants{{Dishes: Dishes{{}, {}}}}},
	}
	tag := "sometag"
	ss.setGTag(tag)
	for _, s := range ss {
		assert.Equal(t, tag, s.GTag)
		for _, r := range s.Restaurants {
			assert.Equal(t, tag, r.GTag)
			for _, d := range r.Dishes {
				assert.Equal(t, tag, d.GTag)
			}
		}
	}
}
