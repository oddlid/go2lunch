package gld

/*
Playing with how to implement the data structures with generics
*/

type Identifiable interface {
	GetID() string
}

type LunchDataContainer[T Identifiable] struct {
	Items []T
}

type MatchFunc[T Identifiable] func(T) bool

func (ldc *LunchDataContainer[T]) Get(f MatchFunc[T]) T {
	var result T
	for _, item := range ldc.Items {
		if f(item) {
			result = item
			break
		}
	}
	return result
}

func (ldc *LunchDataContainer[T]) Add(item T) {
	ldc.Items = append(ldc.Items, item)
}

type ID struct {
	ID   string
	Name string
}

func (id ID) GetID() string {
	return id.ID
}

type Dish struct {
	ID
	Desc  string
	Price float32
}

type Restaurant struct {
	ID
	Dishes LunchDataContainer[*Dish]
}

type Site struct {
	ID
	Restaurants LunchDataContainer[*Restaurant]
}

// func (r *Restaurant) GetByID(id string) *Dish {
// 	if r == nil {
// 		return nil
// 	}
// 	return r.Dishes.Get(func(d *Dish) bool { return d.ID.ID == id })
// }

// func (s *Site) GetByID(id string) *Restaurant {
// 	if s == nil {
// 		return nil
// 	}
// 	return s.Restaurants.Get(func(r *Restaurant) bool { return r.ID.ID == id })
// }
