package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCities_Len(t *testing.T) {
	assert.Zero(t, (Cities)(nil).Len())

	cs := Cities{{}, {}}
	assert.Equal(t, 2, cs.Len())
}

func TestCities_NumSites(t *testing.T) {
	assert.Zero(t, (Cities)(nil).NumSites())

	cs := Cities{{Sites: Sites{{}}}}
	assert.Equal(t, 1, cs.NumSites())
}

func TestCities_NumRestaurants(t *testing.T) {
	assert.Zero(t, (Cities)(nil).NumRestaurants())

	cs := Cities{
		{
			Sites: Sites{
				{
					Restaurants: Restaurants{{}, {}},
				},
			},
		},
	}
	assert.Equal(t, 2, cs.NumRestaurants())
}

func TestCities_NumDishes(t *testing.T) {
	assert.Zero(t, (Cities)(nil).NumDishes())

	cs := Cities{
		{
			Sites: Sites{
				{
					Restaurants: Restaurants{
						{
							Dishes: Dishes{{}, {}},
						},
						{
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
	assert.Zero(t, (Cities)(nil).Total())

	cs := Cities{
		{
			Sites: Sites{
				{
					Restaurants: Restaurants{
						{
							Dishes: Dishes{{}, {}},
						},
						{
							Dishes: Dishes{{}},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 7, cs.Total())
}

func TestCities_setGTag(t *testing.T) {
	assert.NotPanics(t, func() { (Cities)(nil).setGTag("") })

	cs := Cities{
		{
			Sites: Sites{
				{
					Restaurants: Restaurants{
						{
							Dishes: Dishes{{}, {}},
						},
						{
							Dishes: Dishes{{}},
						},
					},
				},
			},
		},
	}
	tag := "sometag"
	cs.setGTag(tag)
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
