package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCityMap_Len(t *testing.T) {
	var nilMap CityMap
	assert.Zero(t, nilMap.Len())

	cm := CityMap{"1": {}}
	assert.Equal(t, 1, cm.Len())
}

func TestCityMap_Empty(t *testing.T) {
	var nilMap CityMap
	assert.True(t, nilMap.Empty())

	cm := CityMap{"1": {}}
	assert.False(t, cm.Empty())
}

func TestCityMap_NumSites(t *testing.T) {
	var nilMap CityMap
	assert.Zero(t, nilMap.NumSites())

	cm := CityMap{"1": {Sites: SiteMap{"1": {}}}}
	assert.Equal(t, 1, cm.NumSites())
}

func TestCityMap_NumRestaurants(t *testing.T) {
	var nilMap CityMap
	assert.Zero(t, nilMap.NumRestaurants())

	cm := CityMap{"1": {Sites: SiteMap{"1": {Restaurants: RestaurantMap{"1": {}}}}}}
	assert.Equal(t, 1, cm.NumRestaurants())
}

func TestCityMap_NumDishes(t *testing.T) {
	var nilMap CityMap
	assert.Zero(t, nilMap.NumDishes())

	cm := CityMap{
		"1": {
			Sites: SiteMap{
				"1": {
					Restaurants: RestaurantMap{
						"1": {
							Dishes: Dishes{
								{}, {},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 2, cm.NumDishes())
}

func TestCityMap_Total(t *testing.T) {
	var nilMap CityMap
	assert.Zero(t, nilMap.Total())

	cm := CityMap{
		"1": {
			Sites: SiteMap{
				"1": {
					Restaurants: RestaurantMap{
						"1": {
							Dishes: Dishes{
								{}, {},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, 5, cm.Total())
}

func TestCityMap_Add(t *testing.T) {
	var nilMap CityMap
	assert.NotPanics(t, func() { nilMap.Add(nil) })

	ids := []string{"1", "2"}
	c1 := City{ID: ids[0]}
	c2 := City{ID: ids[1]}
	cm := CityMap{}
	cm.Add(&c1, &c2)
	assert.Len(t, cm, 2)
	assert.Same(t, &c1, cm[ids[0]])
	assert.Same(t, &c2, cm[ids[1]])
}

func TestCityMap_Delete(t *testing.T) {
	var nilMap CityMap
	assert.NotPanics(t, func() { nilMap.Delete("") })

	ids := []string{"1", "2"}
	c1 := City{ID: ids[0]}
	c2 := City{ID: ids[1]}
	cm := CityMap{ids[0]: &c1, ids[1]: &c2}
	assert.Len(t, cm, 2)
	cm.Delete(ids[0])
	assert.Len(t, cm, 1)
	assert.Nil(t, cm[ids[0]])
	assert.Same(t, &c2, cm[ids[1]])
}

func TestCityMap_SetGTag(t *testing.T) {
	var nilMap CityMap
	assert.NotPanics(t, func() { nilMap.SetGTag("") })

	cm := CityMap{
		"1": {
			Sites: SiteMap{
				"1": {
					Restaurants: RestaurantMap{
						"1": {
							Dishes: Dishes{
								{}, {},
							},
						},
					},
				},
			},
		},
	}
	tag := "sometag"
	cm.SetGTag(tag)
	for _, c := range cm {
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
