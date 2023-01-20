package lunchdata

type RestaurantMap map[string]*Restaurant

// func (rm RestaurantMap) Clone() RestaurantMap {
// 	clone := make(RestaurantMap)
// 	for _, r := range rm {
// 		clone[r.ID] = r.Clone()
// 	}
// 	return clone
// }

func (rm RestaurantMap) Len() int {
	return len(rm)
}

func (rm RestaurantMap) Empty() bool {
	return rm.Len() == 0
}

func (rm RestaurantMap) NumDishes() int {
	total := 0
	for _, r := range rm {
		total += r.NumDishes()
	}
	return total
}

func (rm RestaurantMap) Total() int {
	return rm.Len() + rm.NumDishes()
}

func (rm RestaurantMap) Add(restaurants ...*Restaurant) {
	if rm == nil {
		return
	}
	for _, r := range restaurants {
		if r != nil {
			rm[r.ID] = r
		}
	}
}

func (rm RestaurantMap) Delete(ids ...string) {
	for _, id := range ids {
		delete(rm, id)
	}
}

func (rm RestaurantMap) Get(id string) *Restaurant {
	return rm[id]
}

func (rm RestaurantMap) setGTag(tag string) {
	for _, r := range rm {
		r.setGTag(tag)
	}
}

func (rm RestaurantMap) setIDIfEmpty() {
	for _, r := range rm {
		r.setIDIfEmpty()
	}
}
