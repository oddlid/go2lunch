package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_City_NumSites(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (*City)(nil).NumSites())

	c := City{Sites: Sites{{}}}
	assert.Equal(t, 1, c.NumSites())
}

func Test_City_NumRestaurants(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (*City)(nil).NumRestaurants())

	c := City{
		Sites: Sites{
			{
				Restaurants: Restaurants{{}},
			},
		},
	}
	assert.Equal(t, 1, c.NumRestaurants())
}

func Test_City_NumDishes(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (*City)(nil).NumDishes())

	c := City{
		Sites: Sites{
			{
				Restaurants: Restaurants{
					{
						Dishes: Dishes{{}, {}},
					},
				},
			},
			{
				Restaurants: Restaurants{
					{
						Dishes: Dishes{{}, {}},
					},
				},
			},
		},
	}
	assert.Equal(t, 4, c.NumDishes())
}

func Test_City_Get(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (*City)(nil).Get(nil))
	const id = `blah`
	c := City{Sites: Sites{{ID: id}}}
	assert.Same(t, &c.Sites[0], c.Get(func(s Site) bool { return s.ID == id }))
}

func Test_City_setIDIfEmpty(t *testing.T) {
	t.Parallel()
	assert.NotPanics(t, func() {
		(*City)(nil).setIDIfEmpty()
	})
	c := City{}
	c.setIDIfEmpty()
	assert.NotEmpty(t, c.ID)
}
