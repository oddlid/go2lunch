package main

import (
	"github.com/oddlid/go2lunch/site"
	"html/template"
	"net/http"
	//"regexp"
	//"time"
)

const (
	lhtmpl string = "lhlunch.html"
)

var tmpl = template.Must(template.ParseFiles(lhtmpl))

//var validPath = regexp.MustCompile("^/(gris|hest|grevling)/([a-zA-Z0-9]+)$")
//var validPath = regexp.MustCompile("^/lindholmen$")

func renderTemplate(w http.ResponseWriter, tpl string, s *site.Site) {
	err := tmpl.ExecuteTemplate(w, tpl, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		m := validPath.FindStringSubmatch(r.URL.Path)
//		if m == nil {
//			http.NotFound(w, r)
//			return
//		}
//		fn(w, r, m[2])
//	}
//}

func lhHandler(w http.ResponseWriter, r *http.Request) {
	if _site == nil {
		http.Error(w, "Site is uninitialised", http.StatusInternalServerError)
		return
	}
	if _site.s.Restaurants == nil {
		err := update()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	renderTemplate(w, lhtmpl, _site.s)
}

//func defHandler(w http.ResponseWriter, r *http.Request) {
//}

