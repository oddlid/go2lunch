package server

/*
The URL paths to serve should be these:

/
/{country}
/{country}/{city}
/{country}/{city}/{site}
/{country}/{city}/{site}/{restaurant}
/{country}/{city}/{site}/{restaurant}/{dish}

The virtual files to serve should be these:

default.{html,txt,json}
lunchlist.{html,txt,json}
country.{html,txt,json}
city.{html,txt,json}
site.{html,txt,json}
restaurant.{html,txt,json}
dish.{html,txt,json}

NOTE: For paths and files, we currently don't support restaurant and dish, but we might later on

It would have been more clear to type all of the above out verbatim where they're used,
but in order to minimize the risk of typos and subtle bugs, I've gone for this not so
obvious way of doing it, where I use constants and types to build all of the above
dynamically on the fly, so that the output would get updated if any constant is changed.

*/

import (
	"fmt"
	"strings"
)

type urlID uint8
type outputFormat uint8

const (
	slash            = `/`
	htmlTemplateFile = `allhtml.go.tpl`
	textTemplateFile = `alltext.go.tpl`
	home             = `default`
	lunchList        = `lunchlist`
	country          = `country`
	city             = `city`
	site             = `site`
	restaurant       = `restaurant`
	dish             = `dish`
	extHTML          = `html`
	extTXT           = `txt`
	extJSON          = `json`
	pathStatic       = `static`
	pathTemplates    = `tmpl`
	prefixHTML       = `/html/`
	prefixTXT        = `/txt/`
	prefixJSON       = `/json/`
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

func (id urlID) String() string {
	switch id {
	case idRoot:
		return slash
	case idLunchList:
		return lunchList
	case idCountry:
		return country
	case idCity:
		return city
	case idSite:
		return site
	case idRestaurant:
		return restaurant
	case idDish:
		return dish
	default:
		return ""
	}
}

func (id urlID) fileName(format outputFormat) string {
	var base string
	switch id {
	case idRoot:
		base = home
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

func buildRouterPathArgs(elements ...any) (string, []any) {
	var buf strings.Builder
	for range elements {
		buf.WriteString(slash)
		buf.WriteString(`{%s}`)
	}
	buf.WriteString(slash)
	return buf.String(), append([]any{}, elements...)
}

func (id urlID) routerPath() string {
	switch id {
	case idCountry:
		format, args := buildRouterPathArgs(country)
		return fmt.Sprintf(format, args...)
	case idCity:
		format, args := buildRouterPathArgs(country, city)
		return fmt.Sprintf(format, args...)
	case idSite:
		format, args := buildRouterPathArgs(country, city, site)
		return fmt.Sprintf(format, args...)
	case idRestaurant:
		format, args := buildRouterPathArgs(country, city, site, restaurant)
		return fmt.Sprintf(format, args...)
	case idDish:
		format, args := buildRouterPathArgs(country, city, site, restaurant, dish)
		return fmt.Sprintf(format, args...)
	default:
		return slash
	}
}

func slashWrap(str string) string {
	return fmt.Sprintf("%s%s%s", slash, str, slash)
}
