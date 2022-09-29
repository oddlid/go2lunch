package server

import (
	"context"
	"encoding/json"
	"errors"
	htpl "html/template"
	"net"
	"net/http"
	"sync"
	ttpl "text/template"

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
	config     Config
}

type templateHandler func(w http.ResponseWriter, id urlID, data any)

var (
	errNilHTMLTemplate = errors.New("html template is nil")
	errNilTextTemplate = errors.New("text template is nil")
)

func respondTemplateError(w http.ResponseWriter) {
	http.Error(w, "Template error", http.StatusInternalServerError)
}

func (s *publicServer) start(ctx context.Context) error {
	if s.lunchData == nil {
		return errors.New("lunchData is nil")
	}
	if err := s.loadTemplates(); err != nil {
		return err
	}
	router, err := s.setupRouter()
	if err != nil {
		return err
	}

	s.httpServer = http.Server{
		Handler:           router,
		Addr:              s.config.addr(),
		ReadTimeout:       s.config.ReadTimeout,
		ReadHeaderTimeout: s.config.ReadHeaderTimeout,
		WriteTimeout:      s.config.WriteTimeout,
		IdleTimeout:       s.config.IdleTimeout,
		ErrorLog:          newHTTPErrorLogger(s.log),
		BaseContext:       func(l net.Listener) context.Context { return ctx },
	}

	startWG := sync.WaitGroup{}
	startWG.Add(1)
	go func() {
		startWG.Done()
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error().Err(err).Msg("HTTP server failed")
		} else {
			s.log.Debug().Msg("Lunch server shut down cleanly")
		}
	}()
	startWG.Wait()

	s.log.Debug().Str("addr", s.httpServer.Addr).Msg("Server listening")

	return nil
}

func (s *publicServer) stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *publicServer) loadTemplates() error {
	s.log.Debug().Msg("Loading templates...")

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

	s.hTpl = htmlTemplate
	s.tTpl = textTemplate

	s.log.Debug().Msg("Templates successfully loaded")

	return nil
}

func (s *publicServer) checkHTMLTemplate(w http.ResponseWriter) error {
	if s.hTpl == nil {
		respondTemplateError(w)
		s.log.Error().Err(errNilHTMLTemplate).Send()
		return errNilHTMLTemplate
	}
	return nil
}

func (s *publicServer) checkTextTemplate(w http.ResponseWriter) error {
	if s.tTpl == nil {
		respondTemplateError(w)
		s.log.Error().Err(errNilTextTemplate).Send()
		return errNilTextTemplate
	}
	return nil
}

func (s *publicServer) setupRouter() (*mux.Router, error) {
	box, err := rice.FindBox("static")
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	static := slashWrap(pathStatic) // '/static/'
	router.PathPrefix(static).Handler(http.StripPrefix(static, http.FileServer(box.HTTPBox())))
	router.HandleFunc(idLunchList.routerPath(), s.htmlRootHandler).Methods(http.MethodGet)

	htmlSubRouter := router.PathPrefix(prefixHTML).Subrouter().StrictSlash(true)
	htmlSubRouter.HandleFunc(idLunchList.routerPath(), s.htmlTemplateHandler).Methods(http.MethodGet)
	htmlSubRouter.HandleFunc(idCountry.routerPath(), s.htmlTemplateHandler).Methods(http.MethodGet)
	htmlSubRouter.HandleFunc(idCity.routerPath(), s.htmlTemplateHandler).Methods(http.MethodGet)
	htmlSubRouter.HandleFunc(idSite.routerPath(), s.htmlTemplateHandler).Methods(http.MethodGet)

	textSubRouter := router.PathPrefix(prefixTXT).Subrouter().StrictSlash(true)
	textSubRouter.HandleFunc(idLunchList.routerPath(), s.textTemplateHandler).Methods(http.MethodGet)
	textSubRouter.HandleFunc(idCountry.routerPath(), s.textTemplateHandler).Methods(http.MethodGet)
	textSubRouter.HandleFunc(idCity.routerPath(), s.textTemplateHandler).Methods(http.MethodGet)
	textSubRouter.HandleFunc(idSite.routerPath(), s.textTemplateHandler).Methods(http.MethodGet)

	jsonSubRouter := router.PathPrefix(prefixJSON).Subrouter().StrictSlash(true)
	jsonSubRouter.HandleFunc(idLunchList.routerPath(), s.jsonHandler).Methods(http.MethodGet)
	jsonSubRouter.HandleFunc(idCountry.routerPath(), s.jsonHandler).Methods(http.MethodGet)
	jsonSubRouter.HandleFunc(idCity.routerPath(), s.jsonHandler).Methods(http.MethodGet)
	jsonSubRouter.HandleFunc(idSite.routerPath(), s.jsonHandler).Methods(http.MethodGet)

	return router, nil
}

func (s *publicServer) htmlRootHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.checkHTMLTemplate(w); err != nil {
		return
	}
	if err := s.hTpl.ExecuteTemplate(w, idRoot.fileName(outputHTML), s.lunchData); err != nil {
		s.log.Error().Err(err).Msg("Failed to execute template for /")
	}
}

func (s *publicServer) pathHandler(w http.ResponseWriter, r *http.Request, f templateHandler) {
	v := mux.Vars(r)

	countryID, found := v[country]
	if !found {
		f(w, idLunchList, s.lunchData)
		return
	}
	countryPtr := s.lunchData.Get(countryID)
	if countryPtr == nil {
		http.NotFound(w, r)
		return
	}
	cityID, found := v[city]
	if !found {
		f(w, idCountry, countryPtr)
		return
	}
	cityPtr := countryPtr.Get(cityID)
	if cityPtr == nil {
		http.NotFound(w, r)
		return
	}
	siteID, found := v[site]
	if !found {
		f(w, idCity, cityPtr)
		return
	}
	sitePtr := cityPtr.Get(siteID)
	if sitePtr == nil {
		http.NotFound(w, r)
		return
	}
	restaurantID, found := v[restaurant]
	if !found {
		f(w, idSite, sitePtr)
		return
	}
	restaurantPtr := sitePtr.Get(restaurantID)
	if restaurantPtr == nil {
		http.NotFound(w, r)
		return
	}
	f(w, idRestaurant, restaurantPtr)
}

func (s *publicServer) htmlTemplateHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.checkHTMLTemplate(w); err != nil {
		return
	}
	s.pathHandler(w, r, func(w http.ResponseWriter, id urlID, data any) {
		if err := s.hTpl.ExecuteTemplate(w, id.fileName(outputHTML), data); err != nil {
			respondTemplateError(w)
			s.log.Error().Str("id", id.String()).Err(err).Msg("Failed to execute HTML template")
		}
	})
}

func (s *publicServer) textTemplateHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.checkTextTemplate(w); err != nil {
		return
	}
	s.pathHandler(w, r, func(w http.ResponseWriter, id urlID, data any) {
		if err := s.tTpl.ExecuteTemplate(w, id.fileName(outputTXT), data); err != nil {
			respondTemplateError(w)
			s.log.Error().Str("id", id.String()).Err(err).Msg("Failed to execute TEXT template")
		}
	})
}

func (s *publicServer) jsonHandler(w http.ResponseWriter, r *http.Request) {
	s.pathHandler(w, r, func(w http.ResponseWriter, id urlID, data any) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			respondTemplateError(w)
			s.log.Error().Str("id", id.String()).Err(err).Msg("Failed to encode JSON")
		}
	})
}
