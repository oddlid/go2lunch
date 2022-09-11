package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDish_String(t *testing.T) {
	d := Dish{
		Name: "name",
		Desc: "desc",
	}
	assert.Equal(t, "name desc", d.String())
}
