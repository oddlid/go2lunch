package server

import (
	stdLog "log"

	"github.com/rs/zerolog"
)

type httpErrorLogger struct {
	log zerolog.Logger
}

// Write implementa the io.Writer interface
func (l httpErrorLogger) Write(data []byte) (int, error) {
	l.log.Error().Msg(string(data))
	return len(data), nil
}

func newHTTPErrorLogger(logger zerolog.Logger) *stdLog.Logger {
	return stdLog.New(httpErrorLogger{log: logger}, "", 0)
}
