package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCities_Len(t *testing.T) {
	var nilCities Cities
	assert.Zero(t, nilCities.Len())

	cs := Cities{{}, {}}
	assert.Equal(t, 2, cs.Len())
}

func TestCities_Empty(t *testing.T) {
	var nilCities Cities
	assert.True(t, nilCities.Empty())

	cs := Cities{{}}
	assert.False(t, cs.Empty())
}

func TestCities_NumSites(t *testing.T) {
	var nilCities Cities
	assert.Zero(t, nilCities.NumSites())

	cs := Cities{
		{Sites: SiteMap{"1": {}}},
	}
	assert.Equal(t, 1, cs.NumSites())
}

func TestCities_NumRestaurants(t *testing.T) {
	var nilCities Cities
	assert.Zero(t, nilCities.NumRestaurants())

	cs := Cities{
		{
			Sites: SiteMap{
				"1": {
					Restaurants: RestaurantMap{
						"1": {},
						"2": {},
					},
				},
			},
		},
	}
	assert.Equal(t, 2, cs.NumRestaurants())
}

func TestCities_NumDishes(t *testing.T) {
	var nilCities Cities
	assert.Zero(t, nilCities.NumDishes())

	cs := Cities{
		{
			Sites: SiteMap{
				"1": {
					Restaurants: RestaurantMap{
						"1": {
							Dishes: Dishes{{}, {}},
						},
						"2": {
							Dishes: Dishes{{}},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 3, cs.NumDishes())
}

func TestCities_Total(t *testing.T) {
	var nilCities Cities
	assert.Zero(t, nilCities.Total())

	cs := Cities{
		{
			Sites: SiteMap{
				"1": {
					Restaurants: RestaurantMap{
						"1": {
							Dishes: Dishes{{}, {}},
						},
						"2": {
							Dishes: Dishes{{}},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 7, cs.Total())
}

func TestCities_SetGTag(t *testing.T) {
	var nilCities Cities
	assert.NotPanics(t, func() { nilCities.SetGTag("") })

	cs := Cities{
		{
			Sites: SiteMap{
				"1": {
					Restaurants: RestaurantMap{
						"1": {
							Dishes: Dishes{{}, {}},
						},
						"2": {
							Dishes: Dishes{{}},
						},
					},
				},
			},
		},
	}
	tag := "sometag"
	cs.SetGTag(tag)
	for _, c := range cs {
		assert.Equal(t, tag, c.GTag)
		for _, s := range c.Sites {
			assert.Equal(t, tag, s.GTag)
			for _, r := range s.Restaurants {
				assert.Equal(t, tag, r.GTag)
				for _, d := range r.Dishes {
					assert.Equal(t, tag, d.GTag)
				}
			}
		}
	}
}

func TestCities_AsMap(t *testing.T) {
	var nilCities Cities
	assert.Empty(t, nilCities.AsMap())

	cs := Cities{
		{
			ID: "1",
			Sites: SiteMap{
				"1": {
					Restaurants: RestaurantMap{
						"1": {
							Dishes: Dishes{{}, {}},
						},
						"2": {
							Dishes: Dishes{{}},
						},
					},
				},
			},
		},
	}
	cMap := cs.AsMap()
	assert.NotEmpty(t, cMap)
	assert.Same(t, cs[0], cMap["1"])
}
