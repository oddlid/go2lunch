package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSites_Len(t *testing.T) {
	var nilSites Sites
	assert.Equal(t, 0, nilSites.Len())

	s := Sites{{}, {}}
	assert.Equal(t, 2, s.Len())
}

func TestSites_Empty(t *testing.T) {
	var nilSites Sites
	assert.True(t, nilSites.Empty())

	ss := Sites{{}}
	assert.False(t, ss.Empty())
}

func TestSites_NumRestaurants(t *testing.T) {
	var nilSites Sites
	assert.Equal(t, 0, nilSites.NumRestaurants())

	ss := Sites{{Restaurants: RestaurantMap{"1": {}}}}
	assert.Equal(t, 1, ss.NumRestaurants())
}

func TestSites_NumDishes(t *testing.T) {
	var nilSites Sites
	assert.Equal(t, 0, nilSites.NumDishes())

	ss := Sites{
		{Restaurants: RestaurantMap{"1": {Dishes: Dishes{{}, {}}}}},
		{Restaurants: RestaurantMap{"2": {Dishes: Dishes{{}, {}}}}},
	}
	assert.Equal(t, 4, ss.NumDishes())
}

func TestSites_Total(t *testing.T) {
	var nilSites Sites
	assert.Equal(t, 0, nilSites.Total())

	ss := Sites{
		{Restaurants: RestaurantMap{"1": {Dishes: Dishes{{}, {}}}}},
		{Restaurants: RestaurantMap{"2": {Dishes: Dishes{{}, {}}}}},
	}
	assert.Equal(t, 8, ss.Total())
}

func TestSites_SetGTag(t *testing.T) {
	var nilSites Sites
	assert.NotPanics(t, func() { nilSites.SetGTag("") })

	ss := Sites{
		{Restaurants: RestaurantMap{"1": {Dishes: Dishes{{}, {}}}}},
		{Restaurants: RestaurantMap{"2": {Dishes: Dishes{{}, {}}}}},
	}
	tag := "sometag"
	ss.SetGTag(tag)
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

func TestSites_AsMap(t *testing.T) {
	var nilSites Sites
	emptyMap := nilSites.AsMap()
	assert.Empty(t, emptyMap)

	ids := []string{"0", "1"}
	ss := Sites{
		{ID: ids[0], Restaurants: RestaurantMap{ids[0]: {Dishes: Dishes{{}, {}}}}},
		{ID: ids[1], Restaurants: RestaurantMap{ids[1]: {Dishes: Dishes{{}, {}}}}},
	}
	sMap := ss.AsMap()
	assert.Len(t, sMap, 2)
	assert.Same(t, ss[0], sMap[ids[0]])
	assert.Same(t, ss[1], sMap[ids[1]])
}
