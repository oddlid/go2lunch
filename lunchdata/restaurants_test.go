package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestaurants_Len_whenNil(t *testing.T) {
	var rs Restaurants
	assert.Zero(t, rs.Len())
}

func TestRestaurants_Len(t *testing.T) {
	rs := Restaurants{{}, {}}
	assert.Equal(t, 2, rs.Len())
}

func TestRestaurants_Empty(t *testing.T) {
	var nilRestaurants Restaurants
	assert.True(t, nilRestaurants.Empty())

	rs := Restaurants{{}}
	assert.False(t, rs.Empty())
}

func TestRestaurants_NumDishes(t *testing.T) {
	rs := Restaurants{
		{Dishes: Dishes{{}, {}}},
		{Dishes: Dishes{{}}},
		{},
	}
	assert.Equal(t, 3, rs.NumDishes())
}

func TestRestaurants_Total(t *testing.T) {
	var nilRestaurants Restaurants
	assert.Zero(t, nilRestaurants.Total())

	rs := Restaurants{
		{Dishes: Dishes{{}, {}}},
		{Dishes: Dishes{{}}},
		{},
	}
	assert.Equal(t, 6, rs.Total())
}

func TestRestaurants_AsMap(t *testing.T) {
	rs := Restaurants{{ID: "test"}}
	rm := rs.AsMap()
	assert.Equal(t, 1, len(rm))
	assert.Equal(t, rs[0], rm[rs[0].ID])
}

func TestRestaurants_SetGTag(t *testing.T) {
	gtag := "sometag"
	rs := Restaurants{
		{Dishes: Dishes{{}, {}}},
		{Dishes: Dishes{{}, {}}},
	}
	rs.SetGTag(gtag)
	for _, r := range rs {
		assert.Equal(t, gtag, r.GTag)
		for _, d := range r.Dishes {
			assert.Equal(t, gtag, d.GTag)
		}
	}
}
