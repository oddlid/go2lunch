package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDishes_Len_whenNil(t *testing.T) {
	var dishes Dishes
	assert.Equal(t, 0, dishes.Len())
}

func TestDishes_Len(t *testing.T) {
	dishes := Dishes{{}, {}}
	assert.Equal(t, 2, dishes.Len())
}

func TestDish_String(t *testing.T) {
	d := Dish{
		Name: "name",
		Desc: "desc",
	}
	assert.Equal(t, "name desc", d.String())
}
