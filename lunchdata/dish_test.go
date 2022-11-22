package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Dish_Clone(t *testing.T) {
	assert.Nil(t, (*Dish)(nil).Clone())

	d := Dish{
		ID:    "id",
		Name:  "name",
		Desc:  "desc",
		GTag:  "tag",
		Price: 1,
	}
	clone := d.Clone()
	assert.NotNil(t, clone)
	assert.Equal(t, &d, clone)
}

func Test_Dish_String(t *testing.T) {
	d := Dish{
		Name: "name",
		Desc: "desc",
	}
	assert.Equal(t, "name desc", d.String())
}

func Test_Dish_setIDIfEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		var d *Dish
		d.setIDIfEmpty()
	})
	d := Dish{}
	d.setIDIfEmpty()
	assert.NotEmpty(t, d.ID)
}

// func Test_Dish_UnmarshalJSON(t *testing.T) {
// 	data := []byte(`
// 		{
// 			"id": "dishID",
// 			"name": "dishName",
// 			"desc": "dishDesc",
// 			"price": 100
// 		}
// 	`)
// 	var d Dish
// 	assert.NoError(t, json.Unmarshal(data, &d))
// 	assert.NotNil(t, d)
// 	assert.Equal(t, "dishID", d.ID)
// 	assert.Equal(t, "dishName", d.Name)
// 	assert.Equal(t, "dishDesc", d.Desc)
// 	assert.Equal(t, 100, d.Price)
// }
