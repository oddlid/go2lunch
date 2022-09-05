package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_addr(t *testing.T) {
	c := Config{
		Host: "localhost",
		Port: 1234,
	}
	exp := "localhost:1234"
	assert.Equal(t, exp, c.addr())
}
