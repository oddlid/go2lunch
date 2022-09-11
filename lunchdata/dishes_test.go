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

func TestDishes_Empty(t *testing.T) {
	var nilDishes Dishes
	assert.True(t, nilDishes.Empty())

	ds := Dishes{{}}
	assert.False(t, ds.Empty())
}

func TestDishes_SetGTag(t *testing.T) {
	ds := Dishes{{}, {}}
	tag := "sometag"
	ds.SetGTag(tag)
	for _, d := range ds {
		assert.Equal(t, tag, d.GTag)
	}
}
