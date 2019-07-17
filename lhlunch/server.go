package main

import (
	"fmt"
	htmpl "html/template"
	"net/http"
	ttmpl "text/template"

	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

const (
	urlpath_base string = "/lindholmen"
	lhlunch_html string = "lhlunch.html"
	lhlunch_text string = "lhlunch.txt"
	default_html string = "default.html"
	//tmpl_folder  string = "tmpl"
)

var (
	//tmpl_default_html *htmpl.Template
	tmpl_lhlunch_html *htmpl.Template
	tmpl_lhlunch_text *ttmpl.Template
	str_default_html  string
)

func initTmpl() error {
	log.Debug("Looking for template folder...")
	tBox, err := rice.FindBox("tmpl")
	if err != nil {
		return err
	}
	log.Debug("Loading default html template...")
	str_default_html, err = tBox.String(default_html)
	if err != nil {
		return err
	}
	log.Debug("Loading lunch html template...")
	tLunchHtml, err := tBox.String(lhlunch_html)
	if err != nil {
		return err
	}
	log.Debug("Loading lunch text template...")
	tLunchStr, err := tBox.String(lhlunch_text)
	if err != nil {
		return err
	}
	//log.Debug("Parsing default html template...")
	//tmpl_default_html, err = htmpl.New(default_html).Parse(str_default_html)
	//if err != nil {
	//	return err
	//}
	log.Debug("Parsing lunch html template...")
	tmpl_lhlunch_html, err = htmpl.New(lhlunch_html).Parse(tLunchHtml)
	if err != nil {
		return err
	}
	log.Debug("Parsing lunch text template...")
	tmpl_lhlunch_text, err = ttmpl.New(lhlunch_text).Parse(tLunchStr)
	if err != nil {
		return err
	}

	log.Debug("All templates loaded and parsed successfully!")
	return nil
}

func initSite(w http.ResponseWriter) error {
	if _site == nil {
		http.Error(w, "Site is uninitialised", http.StatusInternalServerError)
		return fmt.Errorf("Site not initialised")
	}
	lhs := _site.getLHSite()
	if lhs == nil || !lhs.HasRestaurants() {
		log.Debug("No content yet, scraping...")
		err := update()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}
	return nil
}

func setupRouter() *mux.Router {
	const s string = "/static/"
	box := rice.MustFindBox("static")
	r := mux.NewRouter()
	r.PathPrefix(s).Handler(http.StripPrefix(s, http.FileServer(box.HTTPBox())))
	r.HandleFunc("/", htmlIndexHandler)
	r.HandleFunc("/api/", jsonApiIndexHandler)
	r.HandleFunc(urlpath_base+"{ext:.?[a-z]+}", preHandler)
	//r.HandleFunc(urlpath_base, htmlLunchHandler) // this is needed if I want /lindholmen (no ext) to work
	return r
}

func preHandler(w http.ResponseWriter, r *http.Request) {
	err := initSite(w)
	if err != nil {
		log.Error(err)
		return
	}
	vars := mux.Vars(r)
	log.Debugf("mux vars: %+v", vars)
	if vars["ext"] == ".html" {
		htmlLunchHandler(w, r)
		return
	} else if vars["ext"] == ".json" {
		jsonLunchHandler(w, r)
		return
	} else if vars["ext"] == ".txt" {
		textLunchHandler(w, r)
		return
	}
}

func setCTJsonHdr(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func htmlIndexHandler(w http.ResponseWriter, r *http.Request) {
	//tmpl_default_html.Execute(w, nil)
	w.Write([]byte(str_default_html))
}

func htmlLunchHandler(w http.ResponseWriter, r *http.Request) {
	tmpl_lhlunch_html.Execute(w, _site.getLHSite())
}

func textLunchHandler(w http.ResponseWriter, r *http.Request) {
	tmpl_lhlunch_text.Execute(w, _site.getLHSite())
}

func jsonApiIndexHandler(w http.ResponseWriter, r *http.Request) {
	setCTJsonHdr(w)
	err := _site.ll.GetSiteLinks().Encode(w)
	if err != nil {
		log.Errorf("Error serving JSON: %q", err.Error())
	}
}

func jsonLunchHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	setCTJsonHdr(w)
	err := _site.ll.Encode(w)
	if err != nil {
		log.Errorf("Error serving JSON: %q", err.Error())
	}
}
