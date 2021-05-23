package lunchdata

import (
	log "github.com/sirupsen/logrus"
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
	PKG            = "pkg"
	MOD            = "module"
)

var (
	pkgLog        = log.WithField(PKG, PKG_NAME)
	llLog         = pkgLog.WithField(MOD, TAG_LUNCHLIST)
	countryLog    = pkgLog.WithField(MOD, TAG_COUNTRY)
	cityLog       = pkgLog.WithField(MOD, TAG_CITY)
	siteLog       = pkgLog.WithField(MOD, TAG_SITE)
	restaurantLog = pkgLog.WithField(MOD, TAG_RESTAURANT)
	//dishLog       = pkgLog.WithField(MOD, TAG_DISH)
)
