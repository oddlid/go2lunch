package server

import (
	"fmt"
	"strings"
)

type urlID uint8
type outputFormat uint8

const (
	slash            = `/`
	htmlTemplateFile = `allhtml.go.tpl`
	lunchList        = `lunchlist`
	country          = `country`
	city             = `city`
	site             = `site`
	restaurant       = `restaurant`
	dish             = `dish`
	extHTML          = `html`
	extTXT           = `txt`
	extJSON          = `json`
)

const (
	outputHTML outputFormat = iota
	outputTXT
	outputJSON
)

const (
	idRoot urlID = iota
	idLunchList
	idCountry
	idCity
	idSite
	idRestaurant
	idDish
)

func (id urlID) fileName(format outputFormat) string {
	var base string
	switch id {
	case idLunchList:
		base = lunchList
	case idCountry:
		base = country
	case idCity:
		base = city
	case idSite:
		base = site
	case idRestaurant:
		base = restaurant
	case idDish:
		base = dish
	}
	fileNameFormat := "%s.%s"
	switch format {
	case outputHTML:
		return fmt.Sprintf(fileNameFormat, base, extHTML)
	case outputJSON:
		return fmt.Sprintf(fileNameFormat, base, extJSON)
	default:
		return fmt.Sprintf(fileNameFormat, base, extTXT)
	}
}

func buildRouterPathArgs(elements ...string) []string {
	var buf strings.Builder
	for range elements {
		buf.WriteString(slash)
		buf.WriteString(`{%s}`)
	}
	return append([]string{buf.String()}, elements...)
}

func (id urlID) routerPath() string {
	switch id {
	case idCountry:
		a := buildRouterPathArgs(country)
		return fmt.Sprintf(a[0], a[1:])
	case idCity:
		a := buildRouterPathArgs(country, city)
		return fmt.Sprintf(a[0], a[1:])
	case idSite:
		a := buildRouterPathArgs(country, city, site)
		return fmt.Sprintf(a[0], a[1:])
	default:
		return slash
	}
}

func surroundWithSlash(word string) string {
	return slash + word + slash
}
