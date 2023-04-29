package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Countries_Len(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (Countries)(nil).Len())
	cs := Countries{{}}
	assert.Equal(t, 1, cs.Len())
}

func Test_Countries_NumCities(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (Countries)(nil).NumCities())
	cs := Countries{{Cities: Cities{{}}}}
	assert.Equal(t, 1, cs.NumCities())
}

func Test_Countries_NumSites(t *testing.T) {
	t.Parallel()
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

func Test_Countries_NumRestaurants(t *testing.T) {
	t.Parallel()
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

func Test_Countries_NumDishes(t *testing.T) {
	t.Parallel()
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

func Test_Countries_Total(t *testing.T) {
	t.Parallel()
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

func Test_Countries_GetByID(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (Countries)(nil).GetByID(""))

	const id = `id`
	cs := Countries{{ID: id}}
	assert.Same(t, &cs[0], cs.GetByID(id))
}
