package lunchdata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Restaurant_NumDishes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0, (*Restaurant)(nil).NumDishes())

	r := Restaurant{
		Dishes: Dishes{{}, {}},
	}
	assert.Equal(t, 2, r.NumDishes())
}

func Test_Restaurant_ParsedRFC3339(t *testing.T) {
	t.Parallel()
	assert.NotEmpty(t, (*Restaurant)(nil).ParsedRFC3339())

	now := time.Now()
	r := Restaurant{ParsedAt: now}
	assert.Equal(t, now.Format(time.RFC3339), r.ParsedRFC3339())
}

func Test_Restaurant_ParsedHumanDate(t *testing.T) {
	t.Parallel()
	assert.NotEmpty(t, (*Restaurant)(nil).ParsedHumanDate())

	now := time.Now()
	r := Restaurant{ParsedAt: now}
	assert.Equal(t, now.Format(dateFormat), r.ParsedHumanDate())
}

func Test_Restaurant_Get(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (*Restaurant)(nil).Get(nil))
	const id = `blah`
	r := Restaurant{Dishes: Dishes{{ID: id}}}
	assert.Same(t, &r.Dishes[0], r.Get(func(d Dish) bool { return d.ID == id }))
}

func Test_Restaurant_GetByID(t *testing.T) {
	t.Parallel()
	assert.Nil(t, (*Restaurant)(nil).GetByID(""))

	const id = `blah`
	r := Restaurant{Dishes: Dishes{{ID: id}}}
	assert.Same(t, &r.Dishes[0], r.GetByID(id))
}
