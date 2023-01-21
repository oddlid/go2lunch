package lunchdata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// func Test_NewRestaurant(t *testing.T) {
// 	name := "Bistrot"
// 	id := "id"
// 	url := "url"
// 	parsed := time.Now()

// 	r := NewRestaurant(name, id, url, parsed)
// 	assert.NotNil(t, r)
// 	assert.IsType(t, (*Restaurant)(nil), r)
// 	assert.Equal(t, name, r.Name)
// 	assert.Equal(t, id, r.ID)
// 	assert.Equal(t, url, r.URL)
// 	assert.NotNil(t, r.Dishes)
// 	assert.Len(t, r.Dishes, 0)
// }

// func Test_Restaurant_Clone(t *testing.T) {
// 	assert.Nil(t, (*Restaurant)(nil).Clone())

// 	r := Restaurant{
// 		Name:     "rName",
// 		ID:       "rID",
// 		URL:      "rURL",
// 		GTag:     "rTAG",
// 		Address:  "rAddr",
// 		MapURL:   "rMapUrl",
// 		ParsedAt: time.Now(),
// 		Dishes: Dishes{
// 			{
// 				Name:  "dName",
// 				ID:    "dID",
// 				Desc:  "dDesc",
// 				Price: 1,
// 				GTag:  "dTAG",
// 			},
// 		},
// 	}
// 	clone := r.Clone()
// 	assert.NotNil(t, clone)
// 	assert.Equal(t, &r, clone)
// }

func Test_Restaurant_NumDishes(t *testing.T) {
	assert.Equal(t, 0, (*Restaurant)(nil).NumDishes())

	r := Restaurant{
		Dishes: Dishes{{}, {}},
	}
	assert.Equal(t, 2, r.NumDishes())
}

func Test_Restaurant_Empty(t *testing.T) {
	assert.True(t, (*Restaurant)(nil).Empty())

	r := Restaurant{
		Dishes: Dishes{{}, {}},
	}
	assert.False(t, r.Empty())
}

func Test_Restaurant_Set(t *testing.T) {
	assert.NotPanics(t, func() {
		(*Restaurant)(nil).Set(nil)
	})

	r := Restaurant{}
	assert.Nil(t, r.Dishes)

	ds := Dishes{}
	ret := r.Set(ds)
	assert.Same(t, &r, ret)
	assert.NotNil(t, r.Dishes)
	assert.Equal(t, ds, r.Dishes)

	r.Set(nil)
	assert.Nil(t, r.Dishes)
}

func Test_Restaurant_Add(t *testing.T) {
	assert.NotPanics(t, func() {
		(*Restaurant)(nil).Add()
	})

	r := Restaurant{}
	ret := r.Add()
	assert.Same(t, &r, ret)
	assert.Nil(t, r.Dishes)
	ds := Dishes{{}, {}}
	ret = r.Add(ds...)
	assert.Same(t, &r, ret)
	assert.NotNil(t, r.Dishes)
	assert.Equal(t, len(ds), len(r.Dishes))
}

func Test_Restaurant_ParsedRFC3339(t *testing.T) {
	assert.NotEmpty(t, (*Restaurant)(nil).ParsedRFC3339())

	now := time.Now()
	r := Restaurant{ParsedAt: now}
	assert.Equal(t, now.Format(time.RFC3339), r.ParsedRFC3339())
}

func Test_Restaurant_ParsedHumanDate(t *testing.T) {
	assert.NotEmpty(t, (*Restaurant)(nil).ParsedHumanDate())

	now := time.Now()
	r := Restaurant{ParsedAt: now}
	assert.Equal(t, now.Format(dateFormat), r.ParsedHumanDate())
}

func Test_Restaurant_setGTag(t *testing.T) {
	assert.Nil(t, (*Restaurant)(nil).setGTag(""))

	gtag := "sometag"
	r := Restaurant{
		Dishes: Dishes{
			{Name: "Middag"},
			{Name: "Lunch"},
		},
	}
	ret := r.setGTag(gtag)
	assert.Same(t, &r, ret)
	for _, dish := range r.Dishes {
		assert.Equal(t, gtag, dish.GTag)
	}
	assert.Equal(t, gtag, r.GTag)
}

func Test_Restaurant_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		var r *Restaurant
		r.setIDIfEmpty()
	})
	r := Restaurant{}
	r.setIDIfEmpty()
	assert.NotEmpty(t, r.ID)
}

// func Test_Restaurant_UnmarshalJSON(t *testing.T) {
// 	data := []byte(`
// 		{
// 			"name": "rName",
// 			"id": "rID",
// 			"url": "rURL",
// 			"address": "rAddr",
// 			"map_url": "rMapURL",
// 			"parsed_at": "2022-11-14T12:30:03.465655196+01:00",
// 			"dishes": [
// 				{
// 					"id": "dID",
// 					"name": "dName",
// 					"desc": "dDesc",
// 					"price": 1
// 				}
// 			]
// 		}
// 	`)
// 	var r *Restaurant
// 	assert.NoError(t, json.Unmarshal(data, &r))
// 	assert.NotNil(t, r)
// 	assert.IsType(t, (*Restaurant)(nil), r)
// 	assert.Equal(t, "rName", r.Name)
// 	assert.Equal(t, "rID", r.ID)
// 	assert.Equal(t, "rURL", r.URL)
// 	assert.Equal(t, "rAddr", r.Address)
// 	assert.Equal(t, "rMapURL", r.MapURL)
// 	assert.NotNil(t, r.Dishes)
// 	assert.Len(t, r.Dishes, 1)
// 	assert.Equal(t, "dID", r.Dishes[0].ID)
// 	assert.Equal(t, "dName", r.Dishes[0].Name)
// 	assert.Equal(t, "dDesc", r.Dishes[0].Desc)
// 	assert.Equal(t, 1, r.Dishes[0].Price)
// }
