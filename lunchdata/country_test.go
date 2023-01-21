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

func Test_Country_setGTag(t *testing.T) {
	assert.NotPanics(t,
		func() {
			(*Country)(nil).setGTag("")
		},
	)
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
	tag := "sometag"
	c.setGTag(tag)
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

func Test_Country_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(*Country)(nil).setIDIfEmpty()
	})
	c := Country{}
	c.setIDIfEmpty()
	assert.NotEmpty(t, c.ID)
}
