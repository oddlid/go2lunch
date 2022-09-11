package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestaurants_Len_whenNil(t *testing.T) {
	var rs Restaurants
	assert.Equal(t, 0, rs.Len())
}

func TestRestaurants_Len(t *testing.T) {
	rs := Restaurants{{}, {}}
	assert.Equal(t, 2, rs.Len())
}

func TestRestaurants_NumDishes(t *testing.T) {
	rs := Restaurants{
		{Dishes: Dishes{{}, {}}},
		{Dishes: Dishes{{}}},
		{},
	}
	assert.Equal(t, 3, rs.NumDishes())
}

func TestRestaurants_AsMap(t *testing.T) {
	rs := Restaurants{
		{
			ID: "test",
		},
	}
	rm := rs.AsMap()
	assert.Equal(t, 1, len(rm))
	assert.Equal(t, rs[0], rm[rs[0].ID])
}

func TestRestaurants_SetGTag(t *testing.T) {
	gtag := "sometag"
	rs := Restaurants{
		{
			Name: "Bistrot",
			Dishes: Dishes{
				{Name: "Köttbullar"},
				{Name: "Pasta"},
			},
		},
		{
			Name: "Kooperativet",
			Dishes: Dishes{
				{Name: "Kyckling"},
				{Name: "Fisk"},
			},
		},
	}
	rs.SetGTag(gtag)
	for _, r := range rs {
		assert.Equal(t, gtag, r.GTag)
		for _, d := range r.Dishes {
			assert.Equal(t, gtag, d.GTag)
		}
	}
}
