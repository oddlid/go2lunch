package lunchdata

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

const (
	PKG_NAME       = "lunchdata"
	TAG_DISH       = "Dish"
	TAG_RESTAURANT = "Restaurant"
	TAG_SITE       = "Site"
	TAG_CITY       = "City"
	TAG_COUNTRY    = "Country"
	TAG_LUNCHLIST  = "LunchList"
	DATE_FORMAT    = "2006-01-02 15:04"
)

func debugf(moduleName, format string, args ...interface{}) {
	if !log.IsLevelEnabled(log.DebugLevel) {
		return
	}
	prefix := fmt.Sprintf("%s.%s: ", PKG_NAME, moduleName)
	log.Debugf(prefix + format, args...)
}

func errorf(moduleName, format string, args ...interface{}) {
	if !log.IsLevelEnabled(log.ErrorLevel) {
		return
	}
	prefix := fmt.Sprintf("%s.%s: ", PKG_NAME, moduleName)
	log.Errorf(prefix + format, args...)
}


func debugDish(format string, args ...interface{}) {
	debugf(TAG_DISH, format, args...)
}

func debugRestaurant(format string, args ...interface{}) {
	debugf(TAG_RESTAURANT, format, args...)
}

func debugSite(format string, args ...interface{}) {
	debugf(TAG_SITE, format, args...)
}

func errorSite(format string, args ...interface{}) {
	errorf(TAG_SITE, format, args...)
}

func debugCity(format string, args ...interface{}) {
	debugf(TAG_CITY, format, args...)
}

func debugCountry(format string, args ...interface{}) {
	debugf(TAG_COUNTRY, format, args...)
}

func debugLunchList(format string, args ...interface{}) {
	debugf(TAG_LUNCHLIST, format, args...)
}

