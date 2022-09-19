package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestaurantMap_Len_whenNil(t *testing.T) {
	var rm RestaurantMap
	assert.Equal(t, 0, rm.Len())
}

func TestRestaurantMap_Len(t *testing.T) {
	rm := make(RestaurantMap)
	rm["one"] = &Restaurant{}
	assert.Equal(t, 1, rm.Len())
}

func TestRestaurantMap_Empty(t *testing.T) {
	assert.True(t, (RestaurantMap)(nil).Empty())

	rm := RestaurantMap{"1": {}}
	assert.False(t, rm.Empty())
}

func TestRestaurantMap_NumDishes(t *testing.T) {
	assert.Zero(t, (RestaurantMap)(nil).NumDishes())

	rm := RestaurantMap{
		"1": {Dishes: Dishes{{}, {}}},
		"2": {Dishes: Dishes{{}, {}}},
		"3": {Dishes: Dishes{{}}},
	}
	assert.Equal(t, 5, rm.NumDishes())
}

func TestRestaurantMap_Total(t *testing.T) {
	assert.Zero(t, (RestaurantMap)(nil).Total())

	rm := RestaurantMap{
		"1": {Dishes: Dishes{{}, {}}},
		"2": {Dishes: Dishes{{}, {}}},
		"3": {Dishes: Dishes{{}}},
	}
	assert.Equal(t, 8, rm.Total())
}

func TestRestaurantMap_Add(t *testing.T) {
	assert.NotPanics(t, func() { (RestaurantMap)(nil).Add(&Restaurant{}) })

	rm := make(RestaurantMap)
	rm.Add(nil)
	assert.Equal(t, 0, rm.Len())

	rm.Add(&Restaurant{})
	assert.Equal(t, 1, rm.Len())
}

func TestResturantMap_Delete(t *testing.T) {
	assert.NotPanics(t, func() {
		(RestaurantMap)(nil).Delete("")
	})

	r := Restaurant{
		ID: "test",
	}
	rm := make(RestaurantMap)
	rm[r.ID] = &r
	assert.Equal(t, 1, len(rm))

	rm.Delete(r.ID)
	assert.Equal(t, 0, len(rm))
}

func TestRestaurantMap_Get(t *testing.T) {
	assert.Nil(t, (RestaurantMap)(nil).Get(""))

	id := "id"
	r := Restaurant{}
	rm := RestaurantMap{id: &r}
	got := rm.Get(id)
	assert.NotNil(t, got)
	assert.Same(t, &r, got)

	assert.Nil(t, rm.Get("otherid"))
}

func TestRestaurantMap_SetGTag(t *testing.T) {
	assert.NotPanics(t, func() { (RestaurantMap)(nil).SetGTag("") })

	rm := RestaurantMap{
		"1": {Dishes: Dishes{{}, {}}},
		"2": {Dishes: Dishes{{}, {}}},
		"3": {Dishes: Dishes{{}}},
	}
	tag := "sometag"
	rm.SetGTag(tag)
	for _, r := range rm {
		assert.Equal(t, tag, r.GTag)
		for _, d := range r.Dishes {
			assert.Equal(t, tag, d.GTag)
		}
	}
}
