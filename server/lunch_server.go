package server

import (
	htpl "html/template"
	"net/http"
	ttpl "text/template"

	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/rs/zerolog"
)

type lunchServer struct {
	lunchData  *lunchdata.LunchList
	hTpl       *htpl.Template
	tTpl       *ttpl.Template
	endpoints  map[string]http.HandlerFunc
	log        zerolog.Logger
	httpServer http.Server
}
