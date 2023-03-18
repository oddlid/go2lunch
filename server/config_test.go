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

func Test_DefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.Equal(t, defaultHost, cfg.Host)
	assert.Equal(t, uint16(DefaultPort), cfg.Port)
	assert.Equal(t, defaultReadTimeout, cfg.ReadTimeout)
	assert.Equal(t, defaultReadHeaderTimeout, cfg.ReadHeaderTimeout)
	assert.Equal(t, defaultWriteTimeout, cfg.WriteTimeout)
	assert.Equal(t, defaultIdleTimeout, cfg.IdleTimeout)
}
