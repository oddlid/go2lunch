package gld

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_misc(t *testing.T) {
	t.Parallel()

	d := Dish{
		ID: ID{
			ID:   "d1",
			Name: "Meatballs",
		},
		Desc:  "with mashed potatoes",
		Price: 120.0,
	}
	r := Restaurant{
		ID: ID{
			ID:   "r1",
			Name: "Bistrot",
		},
	}
	s := Site{
		ID: ID{
			ID:   "LH",
			Name: "Lindholmen",
		},
	}

	r.Dishes.Add(&d)
	s.Restaurants.Add(&r)

	ret := s.Restaurants.Get(func(r *Restaurant) bool {
		return r.ID.ID == "r1"
	}).Dishes.Get(func(d *Dish) bool {
		return d.ID.ID == "d1"
	})
	assert.NotNil(t, ret)
	assert.Same(t, &d, ret)
}
