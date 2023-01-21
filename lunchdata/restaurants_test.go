package lunchdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestaurants_Len_whenNil(t *testing.T) {
	var rs Restaurants
	assert.Zero(t, rs.Len())
}

func TestRestaurants_Len(t *testing.T) {
	rs := Restaurants{{}, {}}
	assert.Equal(t, 2, rs.Len())
}

func TestRestaurants_Empty(t *testing.T) {
	assert.True(t, (Restaurants)(nil).Empty())

	rs := Restaurants{{}}
	assert.False(t, rs.Empty())
}

func TestRestaurants_NumDishes(t *testing.T) {
	rs := Restaurants{
		{Dishes: Dishes{{}, {}}},
		{Dishes: Dishes{{}}},
		{},
	}
	assert.Equal(t, 3, rs.NumDishes())
}

func TestRestaurants_Total(t *testing.T) {
	assert.Zero(t, (Restaurants)(nil).Total())

	rs := Restaurants{
		{Dishes: Dishes{{}, {}}},
		{Dishes: Dishes{{}}},
		{},
	}
	assert.Equal(t, 6, rs.Total())
}

func Test_Restaurants_Get(t *testing.T) {
	rs := Restaurants{
		{URL: "a"},
		{URL: "b"},
	}
	f := func(url string) RestaurantMatch {
		return func(r Restaurant) bool {
			return r.URL == url
		}
	}
	assert.Nil(t, rs.Get(f("c")))
	assert.Same(t, &rs[0], rs.Get(f("a")))
	assert.Same(t, &rs[1], rs.Get(f("b")))
}

func Benchmark_Restaurants_Get(b *testing.B) {
	// It can bee seen when running this benchmark that the time per operation is
	// multiplied by the index of the element to be found when we search slices.
	rs := Restaurants{
		{URL: "a"},
		{URL: "b"},
		{URL: "c"},
	}
	f := func(url string) RestaurantMatch {
		return func(r Restaurant) bool {
			return r.URL == url
		}
	}
	for i := 0; i < b.N; i++ {
		_ = rs.Get(f("c"))
	}
}

func Benchmark_Restaurant_GetFromMap(b *testing.B) {
	rm := map[string]Restaurant{
		"a": {},
		"b": {},
		"c": {},
	}
	for i := 0; i < b.N; i++ {
		_ = rm["c"]
	}
}

func Test_Restaurants_Delete(t *testing.T) {
	rs := Restaurants{
		{URL: "a"},
		{URL: "b"},
	}
	f := func(url string) RestaurantMatch {
		return func(r Restaurant) bool {
			return r.URL == url
		}
	}
	var nilRS Restaurants
	assert.False(t, nilRS.Delete(f("")))
	assert.False(t, rs.Delete(f("c")))
	assert.Len(t, rs, 2)
	assert.True(t, rs.Delete(f("a")))
	assert.Len(t, rs, 1)
	assert.True(t, rs.Delete(f("b")))
	assert.Len(t, rs, 0)
}

func Benchmark_Restaurants_Delete(b *testing.B) {
	const id = `blah`
	f := func(r Restaurant) bool { return r.ID == id }
	r := Restaurant{ID: id}
	rs := Restaurants{
		{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
		r,
	}
	for i := 0; i < b.N; i++ {
		rs.Delete(f)
		rs = append(rs, r)
	}
}

func TestRestaurants_setGTag(t *testing.T) {
	gtag := "sometag"
	rs := Restaurants{
		{Dishes: Dishes{{}, {}}},
		{Dishes: Dishes{{}, {}}},
	}
	rs.setGTag(gtag)
	for i := range rs {
		assert.Equal(t, gtag, rs[i].GTag)
		for j := range rs[i].Dishes {
			assert.Equal(t, gtag, rs[i].Dishes[j].GTag)
		}
	}
}
