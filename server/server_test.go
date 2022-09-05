package server

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestServer_misc(t *testing.T) {
	s := Server{}
	// Test that unset/zero value logger works as no-op
	s.log.Info().Msg("This should not show up in the logs, and neither cause any problems")

	s.log = zerolog.New(zerolog.NewTestWriter(t))
	s.log.Info().Msg("But this should absolutely show up!")
}
