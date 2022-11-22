package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCountry(t *testing.T) {
	id := "id"
	name := "name"
	c := NewCountry(name, id)
	assert.NotNil(t, c)
	assert.IsType(t, (*Country)(nil), c)
	assert.Equal(t, id, c.ID)
	assert.Equal(t, name, c.Name)
	assert.NotNil(t, c.Cities)
}

func Test_Country_NumCities(t *testing.T) {
	assert.Zero(t, (*Country)(nil).NumCities())
	c := Country{Cities: CityMap{"1": {}}}
	assert.Equal(t, 1, c.NumCities())
}

func Test_Country_NumSites(t *testing.T) {
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

func Test_Country_NumRestaurants(t *testing.T) {
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

func Test_Country_NumDishes(t *testing.T) {
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

func Test_Country_setGTag(t *testing.T) {
	assert.Nil(t, (*Country)(nil).setGTag(""))
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
	got := c.setGTag(tag)
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

func Test_Country_Add(t *testing.T) {
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

func Test_Country_Delete(t *testing.T) {
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

func Test_Country_Get(t *testing.T) {
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

func Test_Country_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(*Country)(nil).setIDIfEmpty()
	})
	c := Country{}
	c.setIDIfEmpty()
	assert.NotEmpty(t, c.ID)
}
