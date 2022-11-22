package lunchdata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_RestaurantMap_Clone(t *testing.T) {
	assert.Empty(t, (RestaurantMap)(nil).Clone())

	rm := RestaurantMap{
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
	}
	clone := rm.Clone()
	assert.Equal(t, rm, clone)
}

func Test_RestaurantMap_Len_whenNil(t *testing.T) {
	var rm RestaurantMap
	assert.Equal(t, 0, rm.Len())
}

func Test_RestaurantMap_Len(t *testing.T) {
	rm := make(RestaurantMap)
	rm["one"] = &Restaurant{}
	assert.Equal(t, 1, rm.Len())
}

func Test_RestaurantMap_Empty(t *testing.T) {
	assert.True(t, (RestaurantMap)(nil).Empty())

	rm := RestaurantMap{"1": {}}
	assert.False(t, rm.Empty())
}

func Test_RestaurantMap_NumDishes(t *testing.T) {
	assert.Zero(t, (RestaurantMap)(nil).NumDishes())

	rm := RestaurantMap{
		"1": {Dishes: Dishes{{}, {}}},
		"2": {Dishes: Dishes{{}, {}}},
		"3": {Dishes: Dishes{{}}},
	}
	assert.Equal(t, 5, rm.NumDishes())
}

func Test_RestaurantMap_Total(t *testing.T) {
	assert.Zero(t, (RestaurantMap)(nil).Total())

	rm := RestaurantMap{
		"1": {Dishes: Dishes{{}, {}}},
		"2": {Dishes: Dishes{{}, {}}},
		"3": {Dishes: Dishes{{}}},
	}
	assert.Equal(t, 8, rm.Total())
}

func Test_RestaurantMap_Add(t *testing.T) {
	assert.NotPanics(t, func() { (RestaurantMap)(nil).Add(&Restaurant{}) })

	rm := make(RestaurantMap)
	rm.Add(nil)
	assert.Equal(t, 0, rm.Len())

	rm.Add(&Restaurant{})
	assert.Equal(t, 1, rm.Len())
}

func Test_ResturantMap_Delete(t *testing.T) {
	assert.NotPanics(t, func() {
		(RestaurantMap)(nil).Delete("")
	})

	r := Restaurant{
		ID: "test",
	}
	rm := make(RestaurantMap)
	rm[r.ID] = &r
	assert.Equal(t, 1, len(rm))

	rm.Delete(r.ID)
	assert.Equal(t, 0, len(rm))
}

func Test_RestaurantMap_Get(t *testing.T) {
	assert.Nil(t, (RestaurantMap)(nil).Get(""))

	id := "id"
	r := Restaurant{}
	rm := RestaurantMap{id: &r}
	got := rm.Get(id)
	assert.NotNil(t, got)
	assert.Same(t, &r, got)

	assert.Nil(t, rm.Get("otherid"))
}

func Test_RestaurantMap_setGTag(t *testing.T) {
	assert.NotPanics(t, func() { (RestaurantMap)(nil).setGTag("") })

	rm := RestaurantMap{
		"1": {Dishes: Dishes{{}, {}}},
		"2": {Dishes: Dishes{{}, {}}},
		"3": {Dishes: Dishes{{}}},
	}
	tag := "sometag"
	rm.setGTag(tag)
	for _, r := range rm {
		assert.Equal(t, tag, r.GTag)
		for _, d := range r.Dishes {
			assert.Equal(t, tag, d.GTag)
		}
	}
}

func Test_RestaurantMap_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		(RestaurantMap)(nil).setIDIfEmpty()
	})
	rm := RestaurantMap{
		"1": {},
	}
	rm.setIDIfEmpty()
	assert.NotEmpty(t, rm["1"].ID)
}

// func Test_RestaurantMap_UnmarshalJSON(t *testing.T) {
// 	data := []byte(`
// 		{
// 			"mapKey": {
// 				"name": "rName",
// 				"id": "rID",
// 				"url": "rURL",
// 				"address": "rAddr",
// 				"map_url": "rMapURL",
// 				"parsed_at": "2022-11-14T12:30:03.465655196+01:00",
// 				"dishes": [
// 					{
// 						"id": "dID",
// 						"name": "dName",
// 						"desc": "dDesc",
// 						"price": 1
// 					}
// 				]
// 			}
// 		}
// 	`)
// 	var rm RestaurantMap
// 	assert.NoError(t, json.Unmarshal(data, &rm))
// 	assert.NotNil(t, rm)
// 	assert.IsType(t, (RestaurantMap)(nil), rm)
// 	r := rm["mapKey"]
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
