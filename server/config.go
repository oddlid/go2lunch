package server

import (
	"fmt"
	"time"
)

const (
	defaultHost              = ``
	defaultReadTimeout       = 5 * time.Minute
	defaultReadHeaderTimeout = 5 * time.Minute
	defaultWriteTimeout      = 5 * time.Minute
	defaultIdleTimeout       = 5 * time.Minute
	defaultPort              = 20666
)

type Config struct {
	Host              string // use to limit listening, leave unset as default
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	Port              uint16
}

func (c Config) addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func DefaultConfig() Config {
	return Config{
		Host:              defaultHost,
		Port:              defaultPort,
		ReadTimeout:       defaultReadTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
	}
}
