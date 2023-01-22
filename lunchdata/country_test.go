package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Country_NumCities(t *testing.T) {
	assert.Zero(t, (*Country)(nil).NumCities())
	c := Country{Cities: Cities{{}}}
	assert.Equal(t, 1, c.NumCities())
}

func Test_Country_NumSites(t *testing.T) {
	assert.Zero(t, (*Country)(nil).NumSites())
	c := Country{
		Cities: Cities{
			{
				Sites: Sites{
					{},
				},
			},
		},
	}
	assert.Equal(t, 1, c.NumSites())
}

func Test_Country_NumRestaurants(t *testing.T) {
	assert.Zero(t, (*Country)(nil).NumRestaurants())
	c := Country{
		Cities: Cities{
			{
				Sites: Sites{
					{
						Restaurants: Restaurants{
							{},
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
		Cities: Cities{
			{
				Sites: Sites{
					{
						Restaurants: Restaurants{
							{
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

func Test_Country_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(*Country)(nil).setIDIfEmpty()
	})
	c := Country{}
	c.setIDIfEmpty()
	assert.NotEmpty(t, c.ID)
}

func Test_Country_Get(t *testing.T) {
	assert.Nil(t, (*Country)(nil).Get(nil))
	const id = `blah`
	c := Country{Cities: Cities{{ID: id}}}
	assert.Same(t, &c.Cities[0], c.Get(func(c City) bool { return c.ID == id }))
}
