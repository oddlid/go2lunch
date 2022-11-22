package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CityMap_Len(t *testing.T) {
	assert.Zero(t, (CityMap)(nil).Len())

	cm := CityMap{"1": {}}
	assert.Equal(t, 1, cm.Len())
}

func Test_CityMap_Empty(t *testing.T) {
	assert.True(t, (CityMap)(nil).Empty())

	cm := CityMap{"1": {}}
	assert.False(t, cm.Empty())
}

func Test_CityMap_NumSites(t *testing.T) {
	assert.Zero(t, (CityMap)(nil).NumSites())

	cm := CityMap{"1": {Sites: SiteMap{"1": {}}}}
	assert.Equal(t, 1, cm.NumSites())
}

func Test_CityMap_NumRestaurants(t *testing.T) {
	assert.Zero(t, (CityMap)(nil).NumRestaurants())

	cm := CityMap{"1": {Sites: SiteMap{"1": {Restaurants: RestaurantMap{"1": {}}}}}}
	assert.Equal(t, 1, cm.NumRestaurants())
}

func Test_CityMap_NumDishes(t *testing.T) {
	assert.Zero(t, (CityMap)(nil).NumDishes())

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

func Test_CityMap_Total(t *testing.T) {
	assert.Zero(t, (CityMap)(nil).Total())

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

func Test_CityMap_Add(t *testing.T) {
	assert.NotPanics(t, func() { (CityMap)(nil).Add(&City{}) })

	ids := []string{"1", "2"}
	c1 := City{ID: ids[0]}
	c2 := City{ID: ids[1]}
	cm := CityMap{}
	cm.Add(&c1, &c2)
	assert.Len(t, cm, 2)
	assert.Same(t, &c1, cm[ids[0]])
	assert.Same(t, &c2, cm[ids[1]])
}

func Test_CityMap_Delete(t *testing.T) {
	assert.NotPanics(t, func() { (CityMap)(nil).Delete("") })

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

func Test_CityMap_Get(t *testing.T) {
	assert.Nil(t, (CityMap)(nil).Get(""))

	id := "id"
	c := City{}
	cm := CityMap{id: &c}
	got := cm.Get(id)
	assert.NotNil(t, got)
	assert.Same(t, &c, got)

	assert.Nil(t, cm.Get("otherid"))
}

func Test_CityMap_setGTag(t *testing.T) {
	assert.NotPanics(t, func() { (CityMap)(nil).setGTag("") })

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
	cm.setGTag(tag)
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

func Test_CityMap_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(CityMap)(nil).setIDIfEmpty()
	})
	cm := CityMap{
		"1": {},
	}
	cm.setIDIfEmpty()
	assert.NotEmpty(t, cm["1"].ID)
}
