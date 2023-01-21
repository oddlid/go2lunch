package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountries_Len(t *testing.T) {
	assert.Zero(t, (Countries)(nil).Len())
	cs := Countries{{}}
	assert.Equal(t, 1, cs.Len())
}

func TestCountries_NumCities(t *testing.T) {
	assert.Zero(t, (Countries)(nil).NumCities())
	cs := Countries{{Cities: Cities{{}}}}
	assert.Equal(t, 1, cs.NumCities())
}

func TestCountries_NumSites(t *testing.T) {
	assert.Zero(t, (Countries)(nil).NumSites())
	cs := Countries{
		{
			Cities: Cities{
				{
					Sites: Sites{{}},
				},
			},
		},
	}
	assert.Equal(t, 1, cs.NumSites())
}

func TestCountries_NumRestaurants(t *testing.T) {
	assert.Zero(t, (Countries)(nil).NumRestaurants())
	cs := Countries{
		{
			Cities: Cities{
				{
					Sites: Sites{
						{
							Restaurants: Restaurants{{}},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 1, cs.NumRestaurants())
}

func TestCountries_NumDishes(t *testing.T) {
	assert.Zero(t, (Countries)(nil).NumDishes())
	cs := Countries{
		{
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
		},
	}
	assert.Equal(t, 1, cs.NumDishes())
}

func TestCountries_Total(t *testing.T) {
	assert.Zero(t, (Countries)(nil).Total())
	cs := Countries{
		{
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
		},
	}
	assert.Equal(t, 5, cs.Total())
}

func TestCountries_setGTag(t *testing.T) {
	cs := Countries{
		{
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
		},
	}
	tag := "sometag"
	cs.setGTag(tag)
	for _, country := range cs {
		assert.Equal(t, tag, country.GTag)
		for _, city := range country.Cities {
			assert.Equal(t, tag, city.GTag)
			for _, site := range city.Sites {
				assert.Equal(t, tag, site.GTag)
				for _, r := range site.Restaurants {
					assert.Equal(t, tag, r.GTag)
					for _, d := range r.Dishes {
						assert.Equal(t, tag, d.GTag)
					}
				}
			}
		}
	}
}
