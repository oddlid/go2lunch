package lunchdata

type RestaurantMap map[string]*Restaurant

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

func (rm RestaurantMap) SetGTag(tag string) {
	for _, r := range rm {
		r.SetGTag(tag)
	}
}
