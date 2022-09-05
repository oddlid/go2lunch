package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildRouterPathArgs(t *testing.T) {
	assert.Equal(
		t,
		[]string{"/{%s}", "one"},
		buildRouterPathArgs("one"),
	)
	assert.Equal(
		t,
		[]string{"/{%s}/{%s}", "one", "two"},
		buildRouterPathArgs("one", "two"),
	)
}
