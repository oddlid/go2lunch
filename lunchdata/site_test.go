package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSites_Len(t *testing.T) {
	s := Sites{{}, {}}
	assert.Equal(t, 2, s.Len())
}

func TestNewSite(t *testing.T) {
	id := "id"
	name := "name"
	comment := "comment"
	s := NewSite(name, id, comment)
	assert.NotNil(t, s)
	assert.IsType(t, (*Site)(nil), s)
	assert.Equal(t, id, s.ID)
	assert.Equal(t, name, s.Name)
	assert.Equal(t, comment, s.Comment)
	assert.NotNil(t, s.Restaurants)
}

func TestSite_Len(t *testing.T) {
	s := Site{
		Restaurants: RestaurantMap{
			"1": {},
			"2": {},
		},
	}
	assert.Equal(t, 2, s.Len())
}

func TestSite_SubItems(t *testing.T) {
	s := Site{
		Restaurants: RestaurantMap{
			"1": { // first item
				Dishes: Dishes{{}, {}}, // second and third item
			},
			"2": { // fourth item
				Dishes: Dishes{{}, {}}, // fifth and sixth item
			},
		},
	}
	assert.Equal(t, 6, s.SubItems())
}

func TestSite_getRndRestaurant(t *testing.T) {
	// map access order is random, so we test with only 1 item here to ensure we can assert correctly
	r := Restaurant{}
	s := Site{}
	assert.Nil(t, s.getRndRestaurant())
	s.Restaurants = RestaurantMap{"1": &r}
	assert.Same(t, &r, s.getRndRestaurant())
}

func TestSite_PropagateGtag(t *testing.T) {
}

func TestSite_ParsedHumanDate(t *testing.T) {
}

func TestSite_AddRestaurant(t *testing.T) {
}

func TestSite_DeleteRestaurant(t *testing.T) {
}

func TestSite_HasRestaurants(t *testing.T) {
}

func TestSite_HasRestaurant(t *testing.T) {
}

func TestSite_SetRestaurants(t *testing.T) {
}

func TestSite_ClearRestaurants(t *testing.T) {
}

func TestSite_ClearDishes(t *testing.T) {
}

func TestSite_GetRestaurantByID(t *testing.T) {
}

func TestSite_NumRestaurants(t *testing.T) {
}

func TestSite_NumDishes(t *testing.T) {
}

func TestSite_RunScraper(t *testing.T) {
}
