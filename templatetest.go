package main

import (
	"time"
	"html/template"
	"net/http"
	"github.com/oddlid/go2lunch/site"
)

func siteHandler(w http.ResponseWriter, r *http.Request) {
	site1 := &site.Site{Name:"Lindholmen", Comment:"Where I work"}
	rest := &site.Restaurant{Name:"LHMS", Url:"http://lhms.se", Parsed: time.Now()}
	rest.Add(&site.Dish{"Meatballs", "with mashed potatoes", "85"})
	rest.Add(&site.Dish{"Pancakes", "with jam and whipped cream", "80"})
	site1.Add(rest)

	tmpl := template.Must(template.ParseFiles("site.tmpl"))
	tmpl.Execute(w, site1)
}

func main() {
	http.HandleFunc("/", siteHandler)
	http.ListenAndServe(":10666", nil)
}
