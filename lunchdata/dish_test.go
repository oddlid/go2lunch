package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dish_String(t *testing.T) {
	t.Parallel()
	d := Dish{
		Name:  "name",
		Desc:  "desc",
		Price: 1.2345,
	}
	assert.Equal(t, "name desc :: 1.23", d.String())
}
