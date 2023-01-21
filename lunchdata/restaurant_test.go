package lunchdata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Restaurant_NumDishes(t *testing.T) {
	assert.Equal(t, 0, (*Restaurant)(nil).NumDishes())

	r := Restaurant{
		Dishes: Dishes{{}, {}},
	}
	assert.Equal(t, 2, r.NumDishes())
}

func Test_Restaurant_ParsedRFC3339(t *testing.T) {
	assert.NotEmpty(t, (*Restaurant)(nil).ParsedRFC3339())

	now := time.Now()
	r := Restaurant{ParsedAt: now}
	assert.Equal(t, now.Format(time.RFC3339), r.ParsedRFC3339())
}

func Test_Restaurant_ParsedHumanDate(t *testing.T) {
	assert.NotEmpty(t, (*Restaurant)(nil).ParsedHumanDate())

	now := time.Now()
	r := Restaurant{ParsedAt: now}
	assert.Equal(t, now.Format(dateFormat), r.ParsedHumanDate())
}

func Test_Restaurant_setGTag(t *testing.T) {
	assert.NotPanics(t,
		func() {
			(*Restaurant)(nil).setGTag("")
		},
	)

	gtag := "sometag"
	r := Restaurant{
		Dishes: Dishes{
			{Name: "Middag"},
			{Name: "Lunch"},
		},
	}
	r.setGTag(gtag)
	for _, dish := range r.Dishes {
		assert.Equal(t, gtag, dish.GTag)
	}
	assert.Equal(t, gtag, r.GTag)
}

func Test_Restaurant_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		var r *Restaurant
		r.setIDIfEmpty()
	})
	r := Restaurant{}
	r.setIDIfEmpty()
	assert.NotEmpty(t, r.ID)
}
