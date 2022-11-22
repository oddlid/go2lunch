package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CountryMap_Len(t *testing.T) {
	assert.Zero(t, (CountryMap)(nil).Len())
	cm := CountryMap{"1": {}}
	assert.Equal(t, 1, cm.Len())
}

func Test_CountryMap_Empty(t *testing.T) {
	assert.True(t, (CountryMap)(nil).Empty())
	cm := CountryMap{"1": {}}
	assert.False(t, cm.Empty())
}

func Test_CountryMap_NumCities(t *testing.T) {
	assert.Zero(t, (CountryMap)(nil).NumCities())
	cm := CountryMap{
		"1": {
			Cities: CityMap{
				"1": {},
			},
		},
	}
	assert.Equal(t, 1, cm.NumCities())
}

func Test_CountryMap_NumSites(t *testing.T) {
	assert.Zero(t, (CountryMap)(nil).NumSites())
	cm := CountryMap{
		"1": {
			Cities: CityMap{
				"1": {
					Sites: SiteMap{"1": {}},
				},
			},
		},
	}
	assert.Equal(t, 1, cm.NumSites())
}

func Test_CountryMap_NumRestaurants(t *testing.T) {
	assert.Zero(t, (CountryMap)(nil).NumRestaurants())
	cm := CountryMap{
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
	}
	assert.Equal(t, 1, cm.NumRestaurants())
}

func Test_CountryMap_NumDishes(t *testing.T) {
	assert.Zero(t, (CountryMap)(nil).NumDishes())
	cm := CountryMap{
		"1": {
			Cities: CityMap{
				"1": {
					Sites: SiteMap{
						"1": {
							Restaurants: RestaurantMap{
								"1": {
									Dishes: Dishes{{}},
								},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 1, cm.NumDishes())
}

func Test_CountryMap_Total(t *testing.T) {
	assert.Zero(t, (CountryMap)(nil).Total())
	cm := CountryMap{
		"1": {
			Cities: CityMap{
				"1": {
					Sites: SiteMap{
						"1": {
							Restaurants: RestaurantMap{
								"1": {
									Dishes: Dishes{{}},
								},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 5, cm.Total())
}

func Test_CountryMap_Add(t *testing.T) {
	assert.NotPanics(t, func() {
		(CountryMap)(nil).Add(&Country{})
	})
	cs := Countries{{ID: "1"}, {ID: "2"}}
	cm := CountryMap{}
	cm.Add(cs...)
	assert.Len(t, cm, 2)
	assert.Same(t, cs[0], cm["1"])
	assert.Same(t, cs[1], cm["2"])
}

func Test_CountryMap_Delete(t *testing.T) {
	assert.NotPanics(t, func() {
		(CountryMap)(nil).Delete("")
	})
	cm := CountryMap{
		"1": {},
		"2": {},
	}
	cm.Delete("1")
	assert.Len(t, cm, 1)
	cm.Delete("3")
	assert.Len(t, cm, 1)
}

func Test_CountryMap_Get(t *testing.T) {
	assert.Nil(t, (CountryMap)(nil).Get(""))
	cm := CountryMap{
		"1": {},
		"2": {},
	}
	ret := cm.Get("1")
	assert.NotNil(t, ret)
	assert.Same(t, cm["1"], ret)
	assert.Nil(t, cm.Get("3"))
}

func Test_CountryMap_setGTag(t *testing.T) {
	assert.NotPanics(t, func() { (CountryMap)(nil).setGTag("") })
	cm := CountryMap{
		"1": {
			Cities: CityMap{
				"1": {
					Sites: SiteMap{
						"1": {
							Restaurants: RestaurantMap{
								"1": {
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
	cm.setGTag(tag)
	for _, country := range cm {
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

func Test_CountryMap_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(CountryMap)(nil).setIDIfEmpty()
	})
	cm := CountryMap{
		"1": {},
	}
	cm.setIDIfEmpty()
	assert.NotEmpty(t, cm["1"].ID)
}
