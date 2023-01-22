package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LunchList_NumCountries(t *testing.T) {
	assert.Zero(t, (*LunchList)(nil).NumCountries())
	l := LunchList{
		Countries: Countries{{}},
	}
	assert.Equal(t, 1, l.NumCountries())
}

func Test_LunchList_NumCities(t *testing.T) {
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
	assert.Nil(t, (*LunchList)(nil).Get(nil))

	const id = `blah`
	ll := LunchList{Countries: Countries{{ID: id}}}
	assert.Same(t, &ll.Countries[0], ll.Get(func(c Country) bool { return c.ID == id }))
}

func Test_LunchList_GetByID(t *testing.T) {
	assert.Nil(t, (*LunchList)(nil).GetByID(""))
	const id = `blah`
	ll := LunchList{Countries: Countries{{ID: id}}}
	assert.Same(t, &ll.Countries[0], ll.GetByID(id))
}

func Test_LunchList_RegisterSiteScraper(t *testing.T) {
	assert.Nil(t, (*LunchList)(nil).RegisterSiteScraper(nil))

	ll := LunchList{}
	err := ll.RegisterSiteScraper(nil)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilScraper)

	scraper := mockSiteScraper{
		countryID: "se",
		cityID:    "gbg",
		siteID:    "lh",
	}
	err = ll.RegisterSiteScraper(&scraper)
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilSite)

	ll = LunchList{
		Countries: Countries{
			{
				ID: "se",
				Cities: Cities{
					{
						ID: "gbg",
						Sites: Sites{
							{
								ID: "lh",
							},
						},
					},
				},
			},
		},
	}
	err = ll.RegisterSiteScraper(&scraper)
	assert.NoError(t, err)
	assert.NotNil(t, ll.Countries[0].Cities[0].Sites[0].Scraper)
}

func Test_LunchList_RunSiteScrapers(t *testing.T) {
	assert.NotPanics(t, func() {
		(*LunchList)(nil).RunSiteScrapers()
	})
}

func Test_LunchList_SetIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(*LunchList)(nil).SetIDIfEmpty()
	})
	ll := LunchList{Countries: Countries{{Cities: Cities{{Sites: Sites{{Restaurants: Restaurants{{Dishes: Dishes{}}}}}}}}}}
	ll.SetIDIfEmpty()
	assert.NotEmpty(t, ll.ID)
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
						Sites: Sites{{
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
						}},
					},
				},
			},
		},
	}

	for i := 0; i < b.N; i++ {
		assert.NotNil(b, ll.GetByID(id).GetByID(id).GetByID(id).GetByID(id).GetByID(id))
	}
}
