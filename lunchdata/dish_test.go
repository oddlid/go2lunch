package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dish_String(t *testing.T) {
	t.Parallel()
	assert.Empty(t, (*Dish)(nil).String())
	d := Dish{
		Name: "name",
		Desc: "desc",
	}
	assert.Equal(t, "name desc", d.String())
}
