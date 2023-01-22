package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dishes_Len_whenNil(t *testing.T) {
	var dishes Dishes
	assert.Equal(t, 0, dishes.Len())
}

func Test_Dishes_Len(t *testing.T) {
	dishes := Dishes{{}, {}}
	assert.Equal(t, 2, dishes.Len())
}

func Test_Dishes_Get_whenReceiverIsNil(t *testing.T) {
	assert.Nil(t, (Dishes)(nil).Get(func(_ Dish) bool { return false }))
}

func Test_Dishes_GetByID(t *testing.T) {
	const id = `blaH`
	ds := Dishes{Dish{ID: id}}
	ref := ds.GetByID(id)
	assert.NotNil(t, ref)
	assert.Same(t, &ds[0], ref)
}

func Test_Dishes_setIDIFEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		var nilDishes Dishes
		nilDishes.setIDIfEmpty()
	})
	ds := Dishes{{}, {}}
	ds.setIDIfEmpty()
	assert.NotEmpty(t, ds[0].ID)
	assert.NotEmpty(t, ds[1].ID)
}
