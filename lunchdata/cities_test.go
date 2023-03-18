package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCities_Len(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (Cities)(nil).Len())

	cs := Cities{{}, {}}
	assert.Equal(t, 2, cs.Len())
}

func TestCities_NumSites(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (Cities)(nil).NumSites())

	cs := Cities{{Sites: Sites{{}}}}
	assert.Equal(t, 1, cs.NumSites())
}

func TestCities_NumRestaurants(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func Test_Cities_Get(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (Cities)(nil).Get(nil))
}

func Test_Cities_GetByID(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (Cities)(nil).GetByID(""))
}
