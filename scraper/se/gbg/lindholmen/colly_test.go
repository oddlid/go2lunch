package lindholmen

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_Scrape(t *testing.T) {
	// Content is commented out, since this should only be tested manually against a local webserver

	l := zerolog.New(zerolog.NewTestWriter(t))
	lhs := Scraper{
		Logger: l,
		URL:    "http://localhost:8080", // start a local webserver to test this
	}
	rs, err := lhs.Scrape()
	assert.NoError(t, err)
	assert.NotNil(t, rs)

	l.Debug().Int("numRestaurants", rs.Len()).Send()
	for _, restaurant := range rs {
		l.Debug().
			Str("Restaurant", restaurant.Name).
			Str("ID", restaurant.ID).
			Str("URL", restaurant.URL).
			Str("GTag", restaurant.GTag).
			Str("Addr", restaurant.Address).
			Time("Parsed", restaurant.Parsed).
			Send()
		for _, dish := range restaurant.Dishes {
			l.Debug().
				Str("Dish", dish.Name).
				Str("Desc", dish.Desc).
				Str("GTag", dish.GTag).
				Int("Price", dish.Price).
				Send()
		}
	}
}
