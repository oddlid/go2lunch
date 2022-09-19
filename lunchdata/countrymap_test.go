package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountryMap_Len(t *testing.T) {
	assert.Zero(t, (CountryMap)(nil).Len())
	cm := CountryMap{"1": {}}
	assert.Equal(t, 1, cm.Len())
}

func TestCountryMap_Empty(t *testing.T) {
	assert.True(t, (CountryMap)(nil).Empty())
	cm := CountryMap{"1": {}}
	assert.False(t, cm.Empty())
}

func TestCountryMap_NumCities(t *testing.T) {
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

func TestCountryMap_NumSites(t *testing.T) {
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

func TestCountryMap_NumRestaurants(t *testing.T) {
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

func TestCountryMap_NumDishes(t *testing.T) {
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

func TestCountryMap_Total(t *testing.T) {
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

func TestCountryMap_Add(t *testing.T) {
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

func TestCountryMap_Delete(t *testing.T) {
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

func TestCountryMap_Get(t *testing.T) {
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

func TestCountryMap_SetGTag(t *testing.T) {
	assert.NotPanics(t, func() { (CountryMap)(nil).SetGTag("") })
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
	cm.SetGTag(tag)
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
