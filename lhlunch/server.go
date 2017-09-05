package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/go2lunch/site"
	htmpl "html/template"
	"net/http"
	ttmpl "text/template"
)

const (
	htmpl_ID     string = "LH_HTML"
	ttmpl_ID     string = "LH_TEXT"
	urlpath_base string = "/lindholmen"
)

type lhHandler struct{}

var tmpl_html = htmpl.Must(htmpl.New(htmpl_ID).Parse(lhlunch_html_tmpl_str))
var tmpl_text = ttmpl.Must(ttmpl.New(ttmpl_ID).Parse(lhlunch_text_tmpl_str))
var mux map[string]func(http.ResponseWriter, *http.Request)

func setupMux() {
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = lhHandlerHTMLIndex
	mux[urlpath_base+".html"] = lhHandlerHTML
	mux[urlpath_base+".txt"] = lhHandlerTXT
	mux[urlpath_base+".json"] = lhHandlerJSON
}

func (*lhHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r) // pass on to registered handler
		return
	}
	//fmt.Fprintf(w, "LHLunch server: %q", r.URL.String())
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

func initSite(w http.ResponseWriter) error {
	if _site == nil {
		http.Error(w, "Site is uninitialised", http.StatusInternalServerError)
		return fmt.Errorf("Site not initialised")
	}
	if _site.s.Restaurants == nil {
		log.Debug("No content yet, scraping...")
		err := update()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}
	return nil
}

func renderHTMLTemplate(w http.ResponseWriter, tpl string, s *site.Site) {
	err := tmpl_html.ExecuteTemplate(w, tpl, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderTextTemplate(w http.ResponseWriter, tpl string, s *site.Site) {
	err := tmpl_text.ExecuteTemplate(w, tpl, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func lhHandlerHTMLIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(lhlunch_html_tmpl_str_def))
}

func lhHandlerHTML(w http.ResponseWriter, r *http.Request) {
	err := initSite(w)
	if err != nil {
		return
	}
	renderHTMLTemplate(w, htmpl_ID, _site.s)
}

func lhHandlerTXT(w http.ResponseWriter, r *http.Request) {
	err := initSite(w)
	if err != nil {
		return
	}
	renderTextTemplate(w, ttmpl_ID, _site.s)
}

func lhHandlerJSON(w http.ResponseWriter, r *http.Request) {
	err := initSite(w)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = _site.s.Encode(w)
	if err != nil {
		log.Errorf("Error serving JSON: %q", err.Error())
	}
}
