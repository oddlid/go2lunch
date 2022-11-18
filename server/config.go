package server

import (
	"fmt"
	"time"
)

const (
	DefaultPort              = 20666
	defaultHost              = ``
	defaultReadTimeout       = 5 * time.Second
	defaultReadHeaderTimeout = 5 * time.Second
	defaultWriteTimeout      = 5 * time.Second
	defaultIdleTimeout       = 2 * time.Minute
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
		Port:              DefaultPort,
		ReadTimeout:       defaultReadTimeout,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		WriteTimeout:      defaultWriteTimeout,
		IdleTimeout:       defaultIdleTimeout,
	}
}
