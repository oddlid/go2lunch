package server

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestServer_logging(t *testing.T) {
	s := LunchServer{}
	// Test that unset/zero value logger works as no-op
	s.Log.Info().Msg("This should not show up in the logs, and neither cause any problems")

	s.Log = zerolog.New(zerolog.NewTestWriter(t))
	s.Log.Info().Msg("But this should absolutely show up!")
}
