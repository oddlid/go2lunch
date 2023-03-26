package server

import (
	"context"

	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/rs/zerolog"
)

type LunchServer struct {
	LunchList lunchdata.LunchList
	Log       zerolog.Logger
	server    publicServer
	Config    Config
}

func (s *LunchServer) Start(ctx context.Context) error {
	s.server = publicServer{
		lunchData: &s.LunchList,
		log:       s.Log,
		config:    s.Config,
	}
	return s.server.start(ctx)
}

func (s *LunchServer) Stop(ctx context.Context) error {
	return s.server.httpServer.Shutdown(ctx)
}
