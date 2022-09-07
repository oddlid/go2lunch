package server

import (
	"context"

	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/rs/zerolog"
)

type LunchServer struct {
	LunchList *lunchdata.LunchList
	Log       zerolog.Logger
	server    publicServer
}

func (s *LunchServer) Start() error {
	s.server = publicServer{
		lunchData: s.LunchList,
		log:       s.Log,
	}
	return s.server.start()
}

func (s *LunchServer) Stop(ctx context.Context) error {
	return s.server.stop(ctx)
}
