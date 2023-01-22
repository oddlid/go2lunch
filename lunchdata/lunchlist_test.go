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

// func Test_LunchList_Get(t *testing.T) {
// 	assert.Nil(t, (*LunchList)(nil).Get(""))
// 	l := LunchList{
// 		Countries: Countries{{}, {}},
// 	}
// 	ret := l.Get("1")
// 	assert.Same(t, l.Countries["1"], ret)
// 	assert.Nil(t, l.Get("3"))
// }

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
				Cities: Cities{
					{
						Sites: Sites{
							{},
						},
					},
				},
			},
		},
	}
	err = ll.RegisterSiteScraper(&scraper)
	assert.NoError(t, err)
}

func Test_LunchList_SetIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(*LunchList)(nil).SetIDIfEmpty()
	})
	ll := LunchList{}
	ll.SetIDIfEmpty()
	assert.NotEmpty(t, ll.ID)
}
