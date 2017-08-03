package main

import (
	"github.com/oddlid/go2lunch/site"
	"html/template"
	"net/http"
)

const (
	lhtmpl string = "lhlunch.html"
)

var tmpl = template.Must(template.ParseFiles(lhtmpl))

func renderTemplate(w http.ResponseWriter, tpl string, s *site.Site) {
	err := tmpl.ExecuteTemplate(w, tpl, s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
