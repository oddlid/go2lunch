package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSiteMap_Len(t *testing.T) {
	var nilMap SiteMap
	assert.Equal(t, 0, nilMap.Len())

	sMap := SiteMap{"1": {}}
	assert.Equal(t, 1, sMap.Len())
}

func TestSiteMap_Empty(t *testing.T) {
	var nilMap SiteMap
	assert.True(t, nilMap.Empty())

	sMap := SiteMap{"1": {}}
	assert.False(t, sMap.Empty())
}

func TestSiteMap_NumRestaurants(t *testing.T) {
	var nilMap SiteMap
	assert.Equal(t, 0, nilMap.NumRestaurants())

	sm := SiteMap{
		"1": {Restaurants: RestaurantMap{"1": {}}},
		"2": {Restaurants: RestaurantMap{"1": {}, "2": {}}},
	}
	assert.Equal(t, 3, sm.NumRestaurants())
}

func TestSiteMap_NumDishes(t *testing.T) {
	var nilMap SiteMap
	assert.Equal(t, 0, nilMap.NumDishes())

	sm := SiteMap{
		"1": {
			Restaurants: RestaurantMap{
				"1": {
					Dishes: Dishes{{}, {}},
				},
			},
		},
		"2": {
			Restaurants: RestaurantMap{
				"1": {
					Dishes: Dishes{{}, {}},
				},
			},
		},
	}
	assert.Equal(t, 4, sm.NumDishes())
}

func TestSiteMap_Total(t *testing.T) {
	var nilMap SiteMap
	assert.Equal(t, 0, nilMap.Total())

	sm := SiteMap{
		"1": {
			Restaurants: RestaurantMap{
				"1": {
					Dishes: Dishes{{}, {}},
				},
			},
		},
		"2": {
			Restaurants: RestaurantMap{
				"1": {
					Dishes: Dishes{{}, {}},
				},
			},
		},
	}
	assert.Equal(t, 8, sm.Total())
}
