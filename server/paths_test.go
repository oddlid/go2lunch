package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_urlID_fileName(t *testing.T) {
	assert.Equal(t, "lunchlist.html", idLunchList.fileName(outputHTML))
	assert.Equal(t, "lunchlist.txt", idLunchList.fileName(outputTXT))
	assert.Equal(t, "lunchlist.json", idLunchList.fileName(outputJSON))

	assert.Equal(t, "country.html", idCountry.fileName(outputHTML))
	assert.Equal(t, "country.txt", idCountry.fileName(outputTXT))
	assert.Equal(t, "country.json", idCountry.fileName(outputJSON))

	assert.Equal(t, "city.html", idCity.fileName(outputHTML))
	assert.Equal(t, "city.txt", idCity.fileName(outputTXT))
	assert.Equal(t, "city.json", idCity.fileName(outputJSON))

	assert.Equal(t, "site.html", idSite.fileName(outputHTML))
	assert.Equal(t, "site.txt", idSite.fileName(outputTXT))
	assert.Equal(t, "site.json", idSite.fileName(outputJSON))

	assert.Equal(t, "restaurant.html", idRestaurant.fileName(outputHTML))
	assert.Equal(t, "restaurant.txt", idRestaurant.fileName(outputTXT))
	assert.Equal(t, "restaurant.json", idRestaurant.fileName(outputJSON))

	assert.Equal(t, "dish.html", idDish.fileName(outputHTML))
	assert.Equal(t, "dish.txt", idDish.fileName(outputTXT))
	assert.Equal(t, "dish.json", idDish.fileName(outputJSON))
}

func Test_buildRouterPathArgs(t *testing.T) {
	format, args := buildRouterPathArgs("one")
	assert.Equal(t, "/{%s}/", format)
	assert.Len(t, args, 1)
	assert.Equal(t, "one", args[0])

	format, args = buildRouterPathArgs("one", "two")
	assert.Equal(t, "/{%s}/{%s}/", format)
	assert.Len(t, args, 2)
	assert.Equal(t, "one", args[0])
	assert.Equal(t, "two", args[1])
}

func Test_urlID_routerPath(t *testing.T) {
	assert.Equal(t, slash, idLunchList.routerPath())
	assert.Equal(t, "/{country}/", idCountry.routerPath())
	assert.Equal(t, "/{country}/{city}/", idCity.routerPath())
	assert.Equal(t, "/{country}/{city}/{site}/", idSite.routerPath())
	assert.Equal(t, "/{country}/{city}/{site}/{restaurant}/", idRestaurant.routerPath())
	assert.Equal(t, "/{country}/{city}/{site}/{restaurant}/{dish}/", idDish.routerPath())
}
