package main

import (
	"encoding/json"
	"fmt"
	htmpl "html/template"
	"net/http"
	"strings"
	ttmpl "text/template"

	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/oddlid/go2lunch/lunchdata"
)

var (
	htmlTemplates *htmpl.Template
	textTemplates *ttmpl.Template
	htmlFiles     = []string{"lunchlist.html", "country.html", "city.html", "site.html", "default.html"} // virtual files
	textFiles     = []string{"lunchlist.txt", "country.txt", "city.txt", "site.txt"}                     // virtual files
	urlIds        = []string{"country_id", "city_id", "site_id"}
)

func initTmpl() error {
	log.Debug("Looking for template folder...")
	tBox, err := rice.FindBox("tmpl")
	if err != nil {
		return err
	}

	htmplStr, err := tBox.String("allhtml.tmpl")
	if err != nil {
		return err
	}
	htmlTemplates, err = htmpl.New("html").Parse(htmplStr)
	if err != nil {
		return err
	}

	ttmplStr, err := tBox.String("alltext.tmpl")
	if err != nil {
		return err
	}
	textTemplates, err = ttmpl.New("text").Parse(ttmplStr)
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
	const (
		MGET   = "GET"
		ppStat = "/static/"
		ppJson = "/json/"
		ppHtml = "/html/"
		ppText = "/text/"
		slash  = "/"
	)

	ppath := func(upto int) string {
		var b strings.Builder
		for i := 0; i <= upto; i++ {
			fmt.Fprintf(&b, "/{%s}", urlIds[i])
		}
		b.WriteString(slash) // add final slash
		return b.String()
	}

	box := rice.MustFindBox("static")
	r := mux.NewRouter()
	r.PathPrefix(ppStat).Handler(http.StripPrefix(ppStat, http.FileServer(box.HTTPBox())))
	r.HandleFunc(slash, htmlIndexHandler).Methods(MGET)

	jsubr := r.PathPrefix(ppJson).Subrouter().StrictSlash(true)
	jsubr.HandleFunc(slash, jsonApiHandler).Methods(MGET)
	jsubr.HandleFunc(ppath(0), jsonApiHandler).Methods(MGET)
	jsubr.HandleFunc(ppath(1), jsonApiHandler).Methods(MGET)
	jsubr.HandleFunc(ppath(2), jsonApiHandler).Methods(MGET)
	//jsubr.HandleFunc("/{country_id}/{city_id}/{site_id}/{restaurant_id}", jsonApiHandler).Methods(MGET)

	hsubr := r.PathPrefix(ppHtml).Subrouter().StrictSlash(true)
	hsubr.HandleFunc(slash, htmlTmplHandler).Methods(MGET)
	hsubr.HandleFunc(ppath(0), htmlTmplHandler).Methods(MGET)
	hsubr.HandleFunc(ppath(1), htmlTmplHandler).Methods(MGET)
	hsubr.HandleFunc(ppath(2), htmlTmplHandler).Methods(MGET)

	tsubr := r.PathPrefix(ppText).Subrouter().StrictSlash(true)
	tsubr.HandleFunc(slash, textTmplHandler).Methods(MGET)
	tsubr.HandleFunc(ppath(0), textTmplHandler).Methods(MGET)
	tsubr.HandleFunc(ppath(1), textTmplHandler).Methods(MGET)
	tsubr.HandleFunc(ppath(2), textTmplHandler).Methods(MGET)

	// Redirects just to not break old urls
	r.HandleFunc("/lindholmen.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "html/se/gbg/lindholmen", http.StatusMovedPermanently)
	}).Methods(MGET)
	r.HandleFunc("/lindholmen", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "html/se/gbg/lindholmen", http.StatusMovedPermanently)
	}).Methods(MGET)
	r.HandleFunc("/lindholmen.json", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "json/se/gbg/lindholmen", http.StatusMovedPermanently)
	}).Methods(MGET)
	r.HandleFunc("/lindholmen.txt", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "text/se/gbg/lindholmen", http.StatusMovedPermanently)
	}).Methods(MGET)

	return r
}

func setCTJsonHdr(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func htmlIndexHandler(w http.ResponseWriter, r *http.Request) {
	htmlTemplates.ExecuteTemplate(w, htmlFiles[4], _site.ll)
}

func genericTmplHandler(
	w http.ResponseWriter,
	r *http.Request,
	f func(tmplIndex int, w http.ResponseWriter, obj interface{})) {

	vars := mux.Vars(r)

	var country *lunchdata.Country
	var city *lunchdata.City
	var site *lunchdata.Site

	countryID, found := vars[urlIds[0]]
	if !found {
		// show list of countries and return
		f(0, w, _site.ll)
		return
	}
	country = _site.ll.GetCountryById(countryID)
	if nil == country {
		http.NotFound(w, r) // 404
		return
	}
	cityID, found := vars[urlIds[1]]
	if !found {
		// show list of cities below above country and return
		f(1, w, country)
		return
	}
	city = country.GetCityById(cityID)
	if nil == city {
		http.NotFound(w, r) // 404
		return
	}
	siteID, found := vars[urlIds[2]]
	if !found {
		// show list of sites below above city and return
		f(2, w, city)
		return
	}
	site = city.GetSiteById(siteID)
	if nil == site {
		http.NotFound(w, r) // 404
		return
	}
	// at this point, we might need to do a scrape
	if !site.HasRestaurants() {
		err := update()
		if err != nil {
			log.Error(err)
		}
	}

	f(3, w, site)
}

func textTmplHandler(w http.ResponseWriter, r *http.Request) {
	genericTmplHandler(w, r, func(tmplIdx int, w http.ResponseWriter, obj interface{}) {
		textTemplates.ExecuteTemplate(w, textFiles[tmplIdx], obj)
	})
}

func htmlTmplHandler(w http.ResponseWriter, r *http.Request) {
	genericTmplHandler(w, r, func(tmplIdx int, w http.ResponseWriter, obj interface{}) {
		htmlTemplates.ExecuteTemplate(w, htmlFiles[tmplIdx], obj)
	})
}

func jsonApiHandler(w http.ResponseWriter, r *http.Request) {
	setCTJsonHdr(w)
	genericTmplHandler(w, r, func(tmplIdx int, w http.ResponseWriter, obj interface{}) {
		json.NewEncoder(w).Encode(obj)
	})
}
