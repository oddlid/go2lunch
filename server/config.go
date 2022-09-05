package server

import (
	"fmt"
	"time"
)

type Config struct {
	Host              string // use to limit listening, leave unset as default
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	Port              uint16
}

func (sc Config) addr() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}
