package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LunchList_NumCountries(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (*LunchList)(nil).NumCountries())
	l := LunchList{
		Countries: Countries{{}},
	}
	assert.Equal(t, 1, l.NumCountries())
}

func Test_LunchList_NumCities(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (*LunchList)(nil).NumCities())
	l := LunchList{
		Countries: Countries{
			{
				Cities: Cities{{}},
			},
		},
	}
	assert.Equal(t, 1, l.NumCities())
}

func Test_LunchList_NumSites(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (*LunchList)(nil).NumSites())
	l := LunchList{
		Countries: Countries{
			{
				Cities: Cities{
					{
						Sites: Sites{{}},
					},
				},
			},
		},
	}
	assert.Equal(t, 1, l.NumSites())
}

func Test_LunchList_NumRestaurants(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (*LunchList)(nil).NumRestaurants())
	l := LunchList{
		Countries: Countries{
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
		},
	}
	assert.Equal(t, 1, l.NumRestaurants())
}

func Test_LunchList_NumDishes(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (*LunchList)(nil).NumDishes())
	l := LunchList{
		Countries: Countries{
			{
				Cities: Cities{
					{
						Sites: Sites{
							{
								Restaurants: Restaurants{
									{
										Dishes: Dishes{{}},
									}},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 1, l.NumDishes())
}

func Test_LunchList_Get(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (*LunchList)(nil).Get(nil))

	const id = `blah`
	ll := LunchList{Countries: Countries{{ID: id}}}
	assert.Same(t, &ll.Countries[0], ll.Get(func(c Country) bool { return c.ID == id }))
}

func Test_LunchList_GetByID(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (*LunchList)(nil).GetByID(""))
	const id = `blah`
	ll := LunchList{Countries: Countries{{ID: id}}}
	assert.Same(t, &ll.Countries[0], ll.GetByID(id))
}

func Benchmark_LunchList_GetByID(b *testing.B) {
	const id = `id`
	ll := LunchList{
		Countries: Countries{
			{
				ID: id,
				Cities: Cities{
					{
						ID: id,
						Sites: Sites{
							{
								ID: id,
								Restaurants: Restaurants{
									{
										ID: id,
										Dishes: Dishes{
											{
												ID: id,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i := 0; i < b.N; i++ {
		assert.NotNil(b, ll.GetByID(id).GetByID(id).GetByID(id).GetByID(id).GetByID(id))
	}
}
