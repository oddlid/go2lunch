package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sliceIndex(t *testing.T) {
	idx := sliceIndex([]int(nil), func(n int) bool { return n == 0 })
	assert.Equal(t, -1, idx)

	assert.Equal(
		t,
		1,
		sliceIndex(
			[]int{0, 1, 2},
			func(n int) bool {
				return n == 1
			},
		),
	)

	assert.Equal(
		t,
		-1,
		sliceIndex(
			[]int{0, 1, 2},
			func(n int) bool {
				return n > 2
			},
		),
	)
}

func Test_deleteByIndex_outOfBounds(t *testing.T) {
	slice := []int{0, 1, 2}
	ret := deleteByIndex(slice, 3)
	assert.Equal(t, slice, ret)
}

func Test_deleteByIndex(t *testing.T) {
	slice := []int{0, 1, 2}
	ret := deleteByIndex(slice, 0)
	assert.Len(t, ret, 2)
}
