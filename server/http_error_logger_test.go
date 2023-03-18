package server

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_httpErrorLogger_Write(t *testing.T) {
	const msg = `test`
	l := httpErrorLogger{
		log: zerolog.New(zerolog.NewTestWriter(t)),
	}
	n, err := l.Write([]byte(msg))
	assert.NoError(t, err)
	assert.Equal(t, len(msg), n)
}

func Test_newHTTPErrorLogger(t *testing.T) {
	l := newHTTPErrorLogger(zerolog.Nop())
	assert.NotNil(t, l)
}
