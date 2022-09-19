package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSiteMap_Len(t *testing.T) {
	assert.Equal(t, 0, (SiteMap)(nil).Len())

	sMap := SiteMap{"1": {}}
	assert.Equal(t, 1, sMap.Len())
}

func TestSiteMap_Empty(t *testing.T) {
	assert.True(t, (SiteMap)(nil).Empty())

	sMap := SiteMap{"1": {}}
	assert.False(t, sMap.Empty())
}

func TestSiteMap_NumRestaurants(t *testing.T) {
	assert.Equal(t, 0, (SiteMap)(nil).NumRestaurants())

	sm := SiteMap{
		"1": {Restaurants: RestaurantMap{"1": {}}},
		"2": {Restaurants: RestaurantMap{"1": {}, "2": {}}},
	}
	assert.Equal(t, 3, sm.NumRestaurants())
}

func TestSiteMap_NumDishes(t *testing.T) {
	assert.Equal(t, 0, (SiteMap)(nil).NumDishes())

	sm := SiteMap{
		"1": {
			Restaurants: RestaurantMap{
				"1": {
					Dishes: Dishes{{}, {}},
				},
			},
		},
		"2": {
			Restaurants: RestaurantMap{
				"1": {
					Dishes: Dishes{{}, {}},
				},
			},
		},
	}
	assert.Equal(t, 4, sm.NumDishes())
}

func TestSiteMap_Total(t *testing.T) {
	assert.Equal(t, 0, (SiteMap)(nil).Total())

	sm := SiteMap{
		"1": {
			Restaurants: RestaurantMap{
				"1": {
					Dishes: Dishes{{}, {}},
				},
			},
		},
		"2": {
			Restaurants: RestaurantMap{
				"1": {
					Dishes: Dishes{{}, {}},
				},
			},
		},
	}
	assert.Equal(t, 8, sm.Total())
}

func TestSiteMap_Add(t *testing.T) {
	assert.NotPanics(t, func() { (SiteMap)(nil).Add(&Site{}) })

	sm := SiteMap{}
	ids := []string{"1", "2"}
	s1 := Site{ID: ids[0]}
	s2 := Site{ID: ids[1]}
	sm.Add(&s1, &s2)
	assert.Len(t, sm, 2)
	assert.Same(t, &s1, sm[ids[0]])
	assert.Same(t, &s2, sm[ids[1]])
}

func TestSiteMap_Delete(t *testing.T) {
	assert.NotPanics(t, func() { (SiteMap)(nil).Delete("") })

	ids := []string{"1", "2"}
	s1 := Site{ID: ids[0]}
	s2 := Site{ID: ids[1]}
	sm := SiteMap{ids[0]: &s1, ids[1]: &s2}
	assert.Len(t, sm, 2)
	sm.Delete(ids[0])
	assert.Len(t, sm, 1)
	assert.Nil(t, sm[ids[0]])
	assert.Same(t, &s2, sm[ids[1]])
}

func TestSiteMap_Get(t *testing.T) {
	assert.Nil(t, (SiteMap)(nil).Get(""))

	id := "id"
	s := Site{}
	sm := SiteMap{id: &s}
	got := sm.Get(id)
	assert.NotNil(t, got)
	assert.Same(t, &s, got)

	assert.Nil(t, sm.Get("otherid"))
}

func TestSiteMap_SetGTag(t *testing.T) {
	assert.NotPanics(t, func() { (SiteMap)(nil).SetGTag("") })

	sm := SiteMap{
		"1": {
			Restaurants: RestaurantMap{
				"1": {
					Dishes: Dishes{{}, {}},
				},
			},
		},
	}
	tag := "sometag"
	sm.SetGTag(tag)
	for _, s := range sm {
		assert.Equal(t, tag, s.GTag)
		for _, r := range s.Restaurants {
			assert.Equal(t, tag, r.GTag)
			for _, d := range r.Dishes {
				assert.Equal(t, tag, d.GTag)
			}
		}
	}
}
