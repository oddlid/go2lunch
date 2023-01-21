package lunchdata

type Restaurants []Restaurant
type RestaurantMatch func(r Restaurant) bool

func (rs Restaurants) Len() int {
	return len(rs)
}

func (rs Restaurants) Empty() bool {
	return rs.Len() == 0
}

func (rs Restaurants) NumDishes() int {
	total := 0
	for i := range rs {
		total += rs[i].NumDishes()
	}
	return total
}

func (rs Restaurants) Total() int {
	return rs.Len() + rs.NumDishes()
}

func (rs Restaurants) Get(f RestaurantMatch) *Restaurant {
	if idx := sliceIndex(rs, f); idx > -1 {
		return &rs[idx]
	}
	return nil
}

func (rs *Restaurants) Delete(f RestaurantMatch) bool {
	if idx := sliceIndex(*rs, f); idx > -1 {
		*rs = deleteByIndex(*rs, idx)
		return true
	}
	return false
}

func (rs Restaurants) setGTag(tag string) {
	for i := range rs {
		rs[i].setGTag(tag)
	}
}
