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
	var nilRM RestaurantMap
	assert.True(t, nilRM.Empty())

	rm := RestaurantMap{"1": {}}
	assert.False(t, rm.Empty())
}

func TestRestaurantMap_NumDishes(t *testing.T) {
	var nilRM RestaurantMap
	assert.Zero(t, nilRM.NumDishes())

	rm := RestaurantMap{
		"1": {Dishes: Dishes{{}, {}}},
		"2": {Dishes: Dishes{{}, {}}},
		"3": {Dishes: Dishes{{}}},
	}
	assert.Equal(t, 5, rm.NumDishes())
}

func TestRestaurantMap_Total(t *testing.T) {
	var nilRM RestaurantMap
	assert.Zero(t, nilRM.Total())

	rm := RestaurantMap{
		"1": {Dishes: Dishes{{}, {}}},
		"2": {Dishes: Dishes{{}, {}}},
		"3": {Dishes: Dishes{{}}},
	}
	assert.Equal(t, 8, rm.Total())
}

func TestRestaurantMap_Add(t *testing.T) {
	var nilRM RestaurantMap
	nilRM.Add(&Restaurant{})
	assert.Equal(t, 0, nilRM.Len())

	rm := make(RestaurantMap)
	rm.Add(nil)
	assert.Equal(t, 0, rm.Len())

	rm.Add(&Restaurant{})
	assert.Equal(t, 1, rm.Len())
}

func TestResturantMap_Delete(t *testing.T) {
	var nilRM RestaurantMap
	assert.NotPanics(t, func() {
		nilRM.Delete("test")
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

func TestRestaurantMap_SetGTag(t *testing.T) {
	var nilRM RestaurantMap
	assert.NotPanics(t, func() { nilRM.SetGTag("") })

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
