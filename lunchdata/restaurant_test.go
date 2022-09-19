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

func TestRestaurant_NumDishes(t *testing.T) {
	assert.Equal(t, 0, (*Restaurant)(nil).NumDishes())

	r := Restaurant{
		Dishes: Dishes{{}, {}},
	}
	assert.Equal(t, 2, r.NumDishes())
}

func TestRestaurant_Empty(t *testing.T) {
	assert.True(t, (*Restaurant)(nil).Empty())

	r := Restaurant{
		Dishes: Dishes{{}, {}},
	}
	assert.False(t, r.Empty())
}

func TestRestaurant_Set(t *testing.T) {
	assert.NotPanics(t, func() {
		(*Restaurant)(nil).Set(nil)
	})

	r := Restaurant{}
	assert.Nil(t, r.Dishes)

	ds := Dishes{}
	ret := r.Set(ds)
	assert.Same(t, &r, ret)
	assert.NotNil(t, r.Dishes)
	assert.Equal(t, ds, r.Dishes)

	r.Set(nil)
	assert.Nil(t, r.Dishes)
}

func TestRestaurant_Add(t *testing.T) {
	assert.NotPanics(t, func() {
		(*Restaurant)(nil).Add()
	})

	r := Restaurant{}
	ret := r.Add()
	assert.Same(t, &r, ret)
	assert.Nil(t, r.Dishes)
	ds := Dishes{{}, {}}
	ret = r.Add(ds...)
	assert.Same(t, &r, ret)
	assert.NotNil(t, r.Dishes)
	assert.Equal(t, len(ds), len(r.Dishes))
	ret = r.Add(nil, nil, nil)
	assert.Same(t, &r, ret)
	assert.Equal(t, len(ds), len(r.Dishes))
}

func TestRestaurant_ParsedRFC3339(t *testing.T) {
	assert.NotEmpty(t, (*Restaurant)(nil).ParsedRFC3339())

	now := time.Now()
	r := Restaurant{Parsed: now}
	assert.Equal(t, now.Format(time.RFC3339), r.ParsedRFC3339())
}

func TestRestaurant_ParsedHumanDate(t *testing.T) {
	assert.NotEmpty(t, (*Restaurant)(nil).ParsedHumanDate())

	now := time.Now()
	r := Restaurant{Parsed: now}
	assert.Equal(t, now.Format(dateFormat), r.ParsedHumanDate())
}

func TestRestaurant_SetGTag(t *testing.T) {
	assert.Nil(t, (*Restaurant)(nil).SetGTag(""))

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
