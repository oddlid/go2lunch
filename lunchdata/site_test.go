package lunchdata

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func Test_Site_Clone(t *testing.T) {
	assert.Nil(t, (*Site)(nil).Clone())

	s := Site{
		Name:    "sName",
		ID:      "sID",
		Comment: "sComment",
		URL:     "sURL",
		GTag:    "sTAG",
		Restaurants: RestaurantMap{
			"rID": {
				Name:     "rName",
				ID:       "rID",
				URL:      "rURL",
				GTag:     "rTAG",
				Address:  "rAddr",
				MapURL:   "rMapUrl",
				ParsedAt: time.Now(),
				Dishes: Dishes{
					{
						Name:  "dName",
						ID:    "dID",
						Desc:  "dDesc",
						Price: 1,
						GTag:  "dTAG",
					},
				},
			},
		},
	}
	clone := s.Clone()
	assert.Equal(t, &s, clone)
}

func TestSite_NumRestaurants(t *testing.T) {
	assert.Equal(t, 0, (*Site)(nil).NumRestaurants())

	s := Site{
		Restaurants: RestaurantMap{
			"1": {},
			"2": {},
		},
	}
	assert.Equal(t, 2, s.NumRestaurants())
}

func TestSite_Empty(t *testing.T) {
	assert.True(t, (*Site)(nil).Empty())

	s := Site{Restaurants: RestaurantMap{"1": {}}}
	assert.False(t, s.Empty())
}

func TestSite_getRndRestaurant(t *testing.T) {
	assert.Nil(t, (*Site)(nil).getRndRestaurant())

	// map access order is random, so we test with only 1 item here to ensure we can assert correctly
	r := Restaurant{}
	s := Site{}
	assert.Nil(t, s.getRndRestaurant())
	s.Restaurants = RestaurantMap{"1": &r}
	assert.Same(t, &r, s.getRndRestaurant())
}

func TestSite_setGTag(t *testing.T) {
	assert.Nil(t, (*Site)(nil).setGTag(""))

	gtag := "sometag"
	s := Site{
		Restaurants: RestaurantMap{
			"1": {Dishes: Dishes{{}, {}}},
			"2": {Dishes: Dishes{{}, {}}},
		},
	}
	ret := s.setGTag(gtag)
	assert.Same(t, &s, ret)
	assert.Equal(t, gtag, ret.GTag)
	for _, r := range s.Restaurants {
		assert.Equal(t, gtag, r.GTag)
		for _, d := range r.Dishes {
			assert.Equal(t, gtag, d.GTag)
		}
	}
}

func TestSite_ParsedHumanDate(t *testing.T) {
	assert.Equal(t, dateFormat, (*Site)(nil).ParsedHumanDate())

	now := time.Now()
	r := Restaurant{ParsedAt: now}
	s := Site{}
	assert.Equal(t, dateFormat, s.ParsedHumanDate())
	s.Restaurants = RestaurantMap{"1": &r}
	assert.Equal(t, r.ParsedHumanDate(), s.ParsedHumanDate())
}

func TestSite_Add(t *testing.T) {
	assert.Nil(t, (*Site)(nil).Add(&Restaurant{}))

	id := "id"
	r := Restaurant{ID: id}
	s := Site{}
	ret := s.Add(nil)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 0)
	ret = s.Add(&r)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 1)
}

func TestSite_Delete(t *testing.T) {
	assert.Nil(t, (*Site)(nil).Delete(""))

	s := Site{}
	// Test delete on nil map
	assert.NotPanics(t, func() { s.Delete("") })

	s.Restaurants = RestaurantMap{"id": {}}
	s.Delete("id")
	assert.Len(t, s.Restaurants, 0)
}

func TestSite_Set(t *testing.T) {
	assert.Nil(t, (*Site)(nil).Set(nil))

	id := "id"
	s := Site{
		Restaurants: RestaurantMap{id: {}},
	}
	ret := s.Set(nil)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 0)

	s.Restaurants[id] = &Restaurant{}
	rm := RestaurantMap{
		"1": {},
		"2": {},
	}
	s.Set(rm)
	assert.Len(t, s.Restaurants, 2)
	_, found := s.Restaurants[id]
	assert.False(t, found)
}

func TestSite_Get(t *testing.T) {
	id := "id"
	r := (*Site)(nil).Get(id)
	assert.Nil(t, r)

	s := Site{}
	r = s.Get(id)
	assert.Nil(t, r)

	restaurant := Restaurant{}
	s.Restaurants = RestaurantMap{id: &restaurant}
	r = s.Get(id)
	assert.NotNil(t, r)
	assert.Same(t, &restaurant, r)
}

func TestSite_NumDishes(t *testing.T) {
	assert.Equal(t, 0, (*Site)(nil).NumDishes())

	s := Site{
		Restaurants: RestaurantMap{
			"1": {Dishes: Dishes{{}, {}}},
			"2": {Dishes: Dishes{{}, {}}},
		},
	}
	assert.Equal(t, 4, s.NumDishes())
}

func TestSite_SetScraper(t *testing.T) {
	assert.Nil(t, (*Site)(nil).SetScraper(nil))

	s := Site{}
	scraper := &mockSiteScraper{}
	ret := s.SetScraper(scraper)
	assert.NotNil(t, ret)
	assert.Same(t, &s, ret)
	assert.NotNil(t, s.Scraper)
	assert.Same(t, scraper, s.Scraper)

	assert.NotNil(t, s.SetScraper(nil))
	assert.Nil(t, s.Scraper)
}

func TestSite_RunScraper(t *testing.T) {
	err := (*Site)(nil).RunScraper()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilSite)

	s := Site{}
	err = s.RunScraper()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errNilScraper)

	scrapeErr := errors.New("scrape error")
	s.Scraper = &mockSiteScraper{err: scrapeErr}
	err = s.RunScraper()
	assert.Error(t, err)
	assert.ErrorIs(t, err, scrapeErr)
	assert.Nil(t, s.Restaurants)

	rm := RestaurantMap{
		"1": {},
		"2": {},
	}
	s.Scraper = &mockSiteScraper{restaurants: rm}
	err = s.RunScraper()
	assert.NoError(t, err)
	assert.NotNil(t, s.Restaurants)
	assert.Len(t, s.Restaurants, 2)
}

func Test_Site_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(*Site)(nil).setIDIfEmpty()
	})
	s := Site{}
	s.setIDIfEmpty()
	assert.NotEmpty(t, s.ID)
}

// func Test_Site_UnmarshalJSON(t *testing.T) {
// 	data := []byte(`
// 		{
// 			"name": "sName",
// 			"id": "sID",
// 			"comment": "sComment",
// 			"restaurants": {
// 				"mapKey": {
// 					"name": "rName",
// 					"id": "rID",
// 					"url": "rURL",
// 					"address": "rAddr",
// 					"map_url": "rMapURL",
// 					"parsed_at": "2022-11-14T12:30:03.465655196+01:00",
// 					"dishes": [
// 						{
// 							"id": "dID",
// 							"name": "dName",
// 							"desc": "dDesc",
// 							"price": 1
// 						}
// 					]
// 				}
// 			}
// 		}
// 	`)
// 	var s *Site
// 	assert.NoError(t, json.Unmarshal(data, &s))
// 	assert.NotNil(t, s)
// 	assert.IsType(t, (*Site)(nil), s)
// 	assert.Equal(t, "sName", s.Name)
// 	assert.Equal(t, "sID", s.ID)
// 	assert.Equal(t, "sComment", s.Comment)
// 	assert.NotNil(t, s.Restaurants)
// 	assert.Len(t, s.Restaurants, 1)
// 	r := s.Restaurants["mapKey"]
// 	if assert.NotNil(t, r) {
// 		assert.IsType(t, (*Restaurant)(nil), r)
// 		assert.Equal(t, "rName", r.Name)
// 		assert.Equal(t, "rID", r.ID)
// 		assert.Equal(t, "rURL", r.URL)
// 		assert.Equal(t, "rAddr", r.Address)
// 		assert.Equal(t, "rMapURL", r.MapURL)
// 		assert.NotNil(t, r.Dishes)
// 		assert.Len(t, r.Dishes, 1)
// 		assert.Equal(t, "dID", r.Dishes[0].ID)
// 		assert.Equal(t, "dName", r.Dishes[0].Name)
// 		assert.Equal(t, "dDesc", r.Dishes[0].Desc)
// 		assert.Equal(t, 1, r.Dishes[0].Price)
// 	}
// }
