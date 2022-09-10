package lunchdata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRestaurant(t *testing.T) {
	name := "Bistrot"
	id := "id"
	url := "url"
	parsed := time.Now()

	r := NewRestaurant(name, id, url, parsed)
	assert.NotNil(t, r)
	assert.IsType(t, (*Restaurant)(nil), r)
	assert.Equal(t, name, r.Name)
	assert.Equal(t, id, r.ID)
	assert.Equal(t, url, r.URL)
	assert.NotNil(t, r.Dishes)
	assert.Len(t, r.Dishes, 0)
}
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

func TestRestaurantMap_Len_whenNil(t *testing.T) {
	var rm RestaurantMap
	assert.Equal(t, 0, rm.Len())
}

func TestRestaurantMap_Len(t *testing.T) {
	rm := make(RestaurantMap)
	rm["one"] = &Restaurant{}
	assert.Equal(t, 1, rm.Len())
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

func TestRestaurant_NumDishes(t *testing.T) {
	var nilRestaurant *Restaurant
	assert.Equal(t, 0, nilRestaurant.NumDishes())

	r := Restaurant{
		Dishes: Dishes{{}, {}},
	}
	assert.Equal(t, 2, r.NumDishes())
}

func TestRestaurant_SetDishes(t *testing.T) {
	var nilRestaurant *Restaurant
	assert.NotPanics(t, func() {
		nilRestaurant.SetDishes(nil)
	})

	r := Restaurant{}
	assert.Nil(t, r.Dishes)

	ds := Dishes{}
	ret := r.SetDishes(ds)
	assert.Same(t, &r, ret)
	assert.NotNil(t, r.Dishes)
	assert.Equal(t, ds, r.Dishes)

	r.SetDishes(nil)
	assert.Nil(t, r.Dishes)
}

func TestRestaurant_AddDishes(t *testing.T) {
	var nilRestaurant *Restaurant
	assert.NotPanics(t, func() {
		nilRestaurant.AddDishes()
	})

	r := Restaurant{}
	ret := r.AddDishes()
	assert.Same(t, &r, ret)
	assert.Nil(t, r.Dishes)
	ds := Dishes{{}, {}}
	ret = r.AddDishes(ds...)
	assert.Same(t, &r, ret)
	assert.NotNil(t, r.Dishes)
	assert.Equal(t, len(ds), len(r.Dishes))
	ret = r.AddDishes(nil, nil, nil)
	assert.Same(t, &r, ret)
	assert.Equal(t, len(ds), len(r.Dishes))
}

func TestRestaurant_ParsedRFC3339(t *testing.T) {
	var nilRestaurant *Restaurant
	assert.NotEmpty(t, nilRestaurant.ParsedRFC3339())

	now := time.Now()
	r := Restaurant{Parsed: now}
	assert.Equal(t, now.Format(time.RFC3339), r.ParsedRFC3339())
}

func TestRestaurant_ParsedHumanDate(t *testing.T) {
	var nilRestaurant *Restaurant
	assert.NotEmpty(t, nilRestaurant.ParsedHumanDate())

	now := time.Now()
	r := Restaurant{Parsed: now}
	assert.Equal(t, now.Format(dateFormat), r.ParsedHumanDate())
}

func TestRestaurant_SetGTag(t *testing.T) {
	var nilRestaurant *Restaurant
	assert.Nil(t, nilRestaurant.SetGTag(""))

	gtag := "sometag"
	r := Restaurant{
		Dishes: Dishes{
			{Name: "Middag"},
			{Name: "Lunch"},
		},
	}
	ret := r.SetGTag(gtag)
	assert.Same(t, &r, ret)
	for _, dish := range r.Dishes {
		assert.Equal(t, gtag, dish.GTag)
	}
	assert.Equal(t, gtag, r.GTag)
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
