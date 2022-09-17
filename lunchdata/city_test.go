package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCity(t *testing.T) {
	id := "id"
	name := "name"
	c := NewCity(name, id)
	assert.NotNil(t, c)
	assert.IsType(t, (*City)(nil), c)
	assert.Equal(t, id, c.ID)
	assert.Equal(t, name, c.Name)
	assert.NotNil(t, c.Sites)
}

func TestCity_NumSites(t *testing.T) {
	var nilCity *City
	assert.Zero(t, nilCity.NumSites())

	c := City{Sites: SiteMap{"1": {}}}
	assert.Equal(t, 1, c.NumSites())
}

func TestCity_NumRestaurants(t *testing.T) {
	var nilCity *City
	assert.Zero(t, nilCity.NumRestaurants())

	c := City{
		Sites: SiteMap{
			"1": {
				Restaurants: RestaurantMap{
					"1": {},
				},
			},
		},
	}
	assert.Equal(t, 1, c.NumRestaurants())
}

func TestCity_NumDishes(t *testing.T) {
	var nilCity *City
	assert.Zero(t, nilCity.NumDishes())

	c := City{
		Sites: SiteMap{
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
		},
	}
	assert.Equal(t, 4, c.NumDishes())
}

func TestCity_SetGTag(t *testing.T) {
	var nilCity *City
	assert.NotPanics(t, func() { nilCity.SetGTag("") })

	c := City{
		Sites: SiteMap{
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
		},
	}
	tag := "sometag"
	ret := c.SetGTag(tag)
	assert.Same(t, &c, ret)
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

func TestCity_AddSites(t *testing.T) {
	var nilCity *City
	assert.Nil(t, nilCity.AddSites(nil))

	c := City{}
	assert.Nil(t, c.Sites)
	ret := c.AddSites(nil)
	assert.Same(t, &c, ret)
	assert.NotNil(t, c.Sites)

	s := Site{ID: "1"}
	c.AddSites(&s, nil)
	assert.Len(t, c.Sites, 1)
	assert.Same(t, &s, c.Sites["1"])
}

func TestCity_DeleteSites(t *testing.T) {
	var nilCity *City
	assert.Nil(t, nilCity.DeleteSites(""))

	c := City{
		Sites: SiteMap{
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
		},
	}
	c.DeleteSites("3")
	assert.Len(t, c.Sites, 2)
	c.DeleteSites("2")
	assert.Len(t, c.Sites, 1)
}

func TestCity_GetSiteByID(t *testing.T) {
	var nilCity *City
	assert.Nil(t, nilCity.GetSiteByID(""))

	c := City{
		Sites: SiteMap{
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
		},
	}
	s := c.GetSiteByID("1")
	assert.NotNil(t, s)
	assert.Same(t, c.Sites["1"], s)

	s = c.GetSiteByID("blah")
	assert.Nil(t, s)
}
