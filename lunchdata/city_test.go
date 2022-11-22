package lunchdata

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCity(t *testing.T) {
	id := "id"
	name := "name"
	c := NewCity(name, id)
	assert.NotNil(t, c)
	assert.IsType(t, (*City)(nil), c)
	assert.Equal(t, id, c.ID)
	assert.Equal(t, name, c.Name)
	assert.NotNil(t, c.Sites)
}

func Test_City_NumSites(t *testing.T) {
	assert.Zero(t, (*City)(nil).NumSites())

	c := City{Sites: SiteMap{"1": {}}}
	assert.Equal(t, 1, c.NumSites())
}

func Test_City_NumRestaurants(t *testing.T) {
	assert.Zero(t, (*City)(nil).NumRestaurants())

	c := City{
		Sites: SiteMap{
			"1": {
				Restaurants: RestaurantMap{
					"1": {},
				},
			},
		},
	}
	assert.Equal(t, 1, c.NumRestaurants())
}

func Test_City_NumDishes(t *testing.T) {
	assert.Zero(t, (*City)(nil).NumDishes())

	c := City{
		Sites: SiteMap{
			"1": {
				Restaurants: RestaurantMap{
					"1": {
						Dishes: Dishes{{}, {}},
					},
				},
			},
			"2": {
				Restaurants: RestaurantMap{
					"1": {
						Dishes: Dishes{{}, {}},
					},
				},
			},
		},
	}
	assert.Equal(t, 4, c.NumDishes())
}

func Test_City_setGTag(t *testing.T) {
	assert.NotPanics(t, func() { (*City)(nil).setGTag("") })

	c := City{
		Sites: SiteMap{
			"1": {
				Restaurants: RestaurantMap{
					"1": {
						Dishes: Dishes{{}, {}},
					},
				},
			},
			"2": {
				Restaurants: RestaurantMap{
					"1": {
						Dishes: Dishes{{}, {}},
					},
				},
			},
		},
	}
	tag := "sometag"
	ret := c.setGTag(tag)
	assert.Same(t, &c, ret)
	assert.Equal(t, tag, c.GTag)
	for _, s := range c.Sites {
		assert.Equal(t, tag, s.GTag)
		for _, r := range s.Restaurants {
			assert.Equal(t, tag, r.GTag)
			for _, d := range r.Dishes {
				assert.Equal(t, tag, d.GTag)
			}
		}
	}
}

func Test_City_Add(t *testing.T) {
	assert.Nil(t, (*City)(nil).Add(&Site{}))

	c := City{}
	assert.Nil(t, c.Sites)
	ret := c.Add(nil)
	assert.Same(t, &c, ret)
	assert.NotNil(t, c.Sites)

	s := Site{ID: "1"}
	c.Add(&s, nil)
	assert.Len(t, c.Sites, 1)
	assert.Same(t, &s, c.Sites["1"])
}

func Test_City_Delete(t *testing.T) {
	assert.Nil(t, (*City)(nil).Delete(""))

	c := City{
		Sites: SiteMap{
			"1": {
				Restaurants: RestaurantMap{
					"1": {
						Dishes: Dishes{{}, {}},
					},
				},
			},
			"2": {
				Restaurants: RestaurantMap{
					"1": {
						Dishes: Dishes{{}, {}},
					},
				},
			},
		},
	}
	c.Delete("3")
	assert.Len(t, c.Sites, 2)
	c.Delete("2")
	assert.Len(t, c.Sites, 1)
}

func Test_City_Get(t *testing.T) {
	assert.Nil(t, (*City)(nil).Get(""))

	c := City{
		Sites: SiteMap{
			"1": {
				Restaurants: RestaurantMap{
					"1": {
						Dishes: Dishes{{}, {}},
					},
				},
			},
			"2": {
				Restaurants: RestaurantMap{
					"1": {
						Dishes: Dishes{{}, {}},
					},
				},
			},
		},
	}
	s := c.Get("1")
	assert.NotNil(t, s)
	assert.Same(t, c.Sites["1"], s)

	s = c.Get("blah")
	assert.Nil(t, s)
}

func Test_City_RunSiteScrapers(t *testing.T) {
	assert.NotPanics(t, func() {
		(*City)(nil).RunSiteScrapers(nil, nil)
	})

	c := City{
		Sites: SiteMap{
			"1": {
				Scraper: &mockSiteScraper{
					err: errors.New("scrape error"),
				},
			},
		},
	}
	errChan := make(chan error, c.NumSites())
	wg := sync.WaitGroup{}
	c.RunSiteScrapers(&wg, errChan)
	wg.Wait()
	close(errChan)
	for err := range errChan {
		t.Log(err)
	}
}

func Test_City_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(*City)(nil).setIDIfEmpty()
	})
	c := City{}
	c.setIDIfEmpty()
	assert.NotEmpty(t, c.ID)
}
