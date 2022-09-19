package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCountry(t *testing.T) {
	id := "id"
	name := "name"
	c := NewCountry(name, id)
	assert.NotNil(t, c)
	assert.IsType(t, (*Country)(nil), c)
	assert.Equal(t, id, c.ID)
	assert.Equal(t, name, c.Name)
	assert.NotNil(t, c.Cities)
}

func TestCountry_NumCities(t *testing.T) {
	assert.Zero(t, (*Country)(nil).NumCities())
	c := Country{Cities: CityMap{"1": {}}}
	assert.Equal(t, 1, c.NumCities())
}

func TestCountry_NumSites(t *testing.T) {
	assert.Zero(t, (*Country)(nil).NumSites())
	c := Country{
		Cities: CityMap{
			"1": {
				Sites: SiteMap{
					"1": {},
				},
			},
		},
	}
	assert.Equal(t, 1, c.NumSites())
}

func TestCountry_NumRestaurants(t *testing.T) {
	assert.Zero(t, (*Country)(nil).NumRestaurants())
	c := Country{
		Cities: CityMap{
			"1": {
				Sites: SiteMap{
					"1": {
						Restaurants: RestaurantMap{
							"1": {},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 1, c.NumRestaurants())
}

func TestCountry_NumDishes(t *testing.T) {
	assert.Zero(t, (*Country)(nil).NumDishes())
	c := Country{
		Cities: CityMap{
			"1": {
				Sites: SiteMap{
					"1": {
						Restaurants: RestaurantMap{
							"1": {
								Dishes: Dishes{{}},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 1, c.NumDishes())
}

func TestCountry_SetGTag(t *testing.T) {
	assert.Nil(t, (*Country)(nil).SetGTag(""))
	c := Country{
		Cities: CityMap{
			"1": {
				Sites: SiteMap{
					"1": {
						Restaurants: RestaurantMap{
							"1": {
								Dishes: Dishes{{}},
							},
						},
					},
				},
			},
		},
	}
	tag := "sometag"
	got := c.SetGTag(tag)
	assert.NotNil(t, got)
	assert.Same(t, &c, got)
	assert.Equal(t, tag, c.GTag)
	for _, city := range c.Cities {
		assert.Equal(t, tag, city.GTag)
		for _, s := range city.Sites {
			assert.Equal(t, tag, c.GTag)
			for _, r := range s.Restaurants {
				assert.Equal(t, tag, r.GTag)
				for _, d := range r.Dishes {
					assert.Equal(t, tag, d.GTag)
				}
			}
		}
	}
}

func TestCountry_Add(t *testing.T) {
	assert.Nil(t, (*Country)(nil).Add(&City{}))
	c := Country{}
	assert.Nil(t, c.Cities)
	ret := c.Add(nil)
	assert.Same(t, &c, ret)
	assert.NotNil(t, c.Cities)

	city := City{ID: "1"}
	c.Add(&city)
	assert.Len(t, c.Cities, 1)
	assert.Same(t, &city, c.Cities["1"])
}

func TestCountry_Delete(t *testing.T) {
	assert.Nil(t, (*Country)(nil).Delete())
	c := Country{
		Cities: CityMap{
			"1": {},
			"2": {},
		},
	}
	ret := c.Delete("1")
	assert.Same(t, &c, ret)
	assert.Len(t, c.Cities, 1)
	c.Delete("3")
	assert.Len(t, c.Cities, 1)
}

func TestCountry_Get(t *testing.T) {
	assert.Nil(t, (*Country)(nil).Get(""))
	c := Country{
		Cities: CityMap{
			"1": {},
			"2": {},
		},
	}
	ret := c.Get("1")
	assert.NotNil(t, ret)
	assert.Same(t, c.Cities["1"], ret)
	assert.Nil(t, c.Get("3"))
}
