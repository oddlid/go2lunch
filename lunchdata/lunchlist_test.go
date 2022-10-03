package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewLunchList(t *testing.T) {
	l := NewLunchList()
	assert.NotNil(t, l)
	assert.NotNil(t, l.Countries)
}

func Test_LunchList_NumCountries(t *testing.T) {
	assert.Zero(t, (*LunchList)(nil).NumCountries())
	l := LunchList{
		Countries: CountryMap{"1": {}},
	}
	assert.Equal(t, 1, l.NumCountries())
}

func Test_LunchList_NumCities(t *testing.T) {
	assert.Zero(t, (*LunchList)(nil).NumCities())
	l := LunchList{
		Countries: CountryMap{
			"1": {
				Cities: CityMap{"1": {}},
			},
		},
	}
	assert.Equal(t, 1, l.NumCities())
}

func Test_LunchList_NumSites(t *testing.T) {
	assert.Zero(t, (*LunchList)(nil).NumSites())
	l := LunchList{
		Countries: CountryMap{
			"1": {
				Cities: CityMap{
					"1": {
						Sites: SiteMap{"1": {}},
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
		Countries: CountryMap{
			"1": {
				Cities: CityMap{
					"1": {
						Sites: SiteMap{
							"1": {
								Restaurants: RestaurantMap{"1": {}},
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
		Countries: CountryMap{
			"1": {
				Cities: CityMap{
					"1": {
						Sites: SiteMap{
							"1": {
								Restaurants: RestaurantMap{
									"1": {
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

func Test_LunchList_SetGTag(t *testing.T) {
	assert.Nil(t, (*LunchList)(nil).SetGTag(""))
	l := LunchList{
		Countries: CountryMap{
			"1": {
				Cities: CityMap{
					"1": {
						Sites: SiteMap{
							"1": {
								Restaurants: RestaurantMap{
									"1": {
										Dishes: Dishes{{}},
									}},
							},
						},
					},
				},
			},
		},
	}
	tag := "sometag"
	l.SetGTag(tag)
	assert.Equal(t, tag, l.GTag)
	for _, country := range l.Countries {
		assert.Equal(t, tag, country.GTag)
		for _, city := range country.Cities {
			assert.Equal(t, tag, city.GTag)
			for _, s := range city.Sites {
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
}

func Test_LunchList_Add(t *testing.T) {
	assert.Nil(t, (*LunchList)(nil).Add(&Country{}))
	l := LunchList{}
	assert.Nil(t, l.Countries)
	ret := l.Add(nil)
	assert.Same(t, &l, ret)
	assert.NotNil(t, l.Countries)

	c := Country{ID: "1"}
	l.Add(&c)
	assert.Len(t, l.Countries, 1)
	assert.Same(t, &c, l.Countries["1"])
}

func Test_LunchList_Delete(t *testing.T) {
	assert.Nil(t, (*LunchList)(nil).Delete(""))
	l := LunchList{
		Countries: CountryMap{
			"1": {},
			"2": {},
		},
	}
	ret := l.Delete("1")
	assert.Same(t, &l, ret)
	assert.Len(t, l.Countries, 1)
	l.Delete("3")
	assert.Len(t, l.Countries, 1)
}

func Test_LunchList_Get(t *testing.T) {
	assert.Nil(t, (*LunchList)(nil).Get(""))
	l := LunchList{
		Countries: CountryMap{
			"1": {},
			"2": {},
		},
	}
	ret := l.Get("1")
	assert.Same(t, l.Countries["1"], ret)
	assert.Nil(t, l.Get("3"))
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
		Countries: CountryMap{
			"se": {
				Cities: CityMap{
					"gbg": {
						Sites: SiteMap{
							"lh": {},
						},
					},
				},
			},
		},
	}
	err = ll.RegisterSiteScraper(&scraper)
	assert.NoError(t, err)
}
