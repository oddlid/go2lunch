package server

import (
	"context"
	"errors"
	htpl "html/template"
	"net/http"
	"sync"
	ttpl "text/template"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/rs/zerolog"
)

type publicServer struct {
	lunchData *lunchdata.LunchList
	hTpl      *htpl.Template
	tTpl      *ttpl.Template
	// endpoints  map[string]http.HandlerFunc
	log        zerolog.Logger
	httpServer http.Server
}

func (ls *publicServer) start() error {
	if ls.lunchData == nil {
		return errors.New("lunchData is nil")
	}
	if err := ls.loadTemplates(); err != nil {
		return err
	}
	router, err := ls.setupRouter()
	if err != nil {
		return err
	}

	// TODO: replace hardcoded values
	ls.httpServer = http.Server{
		Addr:              ":20666",
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       1 * time.Minute,
		ErrorLog:          newHttpErrorLogger(ls.log),
	}

	startWG := sync.WaitGroup{}
	startWG.Add(1)
	go func() {
		startWG.Done()
		if err := ls.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ls.log.Error().Err(err).Msg("HTTP server failed")
		} else {
			ls.log.Debug().Msg("Lunch server shut down cleanly")
		}
	}()
	startWG.Wait()

	ls.log.Debug().Str("addr", ls.httpServer.Addr).Msg("Server listening")

	return nil
}

func (ls *publicServer) stop(ctx context.Context) error {
	return ls.httpServer.Shutdown(ctx)
}

func (ls *publicServer) loadTemplates() error {
	ls.log.Debug().Msg("Loading templates...")

	box, err := rice.FindBox("tmpl")
	if err != nil {
		return err
	}
	htmlTemplateString, err := box.String("allhtml.go.tpl")
	if err != nil {
		return err
	}
	htmlTemplate, err := htpl.New("html").Parse(htmlTemplateString)
	if err != nil {
		return err
	}
	textTemplateString, err := box.String("alltext.go.tpl")
	if err != nil {
		return err
	}
	textTemplate, err := ttpl.New("text").Parse(textTemplateString)
	if err != nil {
		return err
	}

	ls.hTpl = htmlTemplate
	ls.tTpl = textTemplate

	ls.log.Debug().Msg("Templates loaded successfully")

	return nil
}

func (ls *publicServer) setupRouter() (*mux.Router, error) {
	box, err := rice.FindBox("static")
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	static := slashWrap(pathStatic) // '/static/'
	router.PathPrefix(static).Handler(http.StripPrefix(static, http.FileServer(box.HTTPBox())))
	router.HandleFunc(slash, ls.htmlRootHandler).Methods(http.MethodGet)

	// ...

	return router, nil
}

func (ls *publicServer) htmlRootHandler(w http.ResponseWriter, r *http.Request) {
	if ls.hTpl == nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		ls.log.Error().Msg("html template pointer is nil")
		return
	}
	if err := ls.hTpl.ExecuteTemplate(w, idRoot.fileName(outputHTML), ls.lunchData); err != nil {
		ls.log.Error().Err(err).Msg("Failed to execute template for /")
	}
}
