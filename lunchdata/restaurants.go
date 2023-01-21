package lunchdata

type Restaurants []Restaurant
type RestaurantMatch func(r Restaurant) bool

func (rs Restaurants) Len() int {
	return len(rs)
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

func (rs Restaurants) GetByID(id string) *Restaurant {
	return rs.Get(func(r Restaurant) bool { return r.ID == id })
}

func (rs Restaurants) first() *Restaurant {
	if len(rs) > 0 {
		return &rs[0]
	}
	return nil
}

func (rs Restaurants) setGTag(tag string) {
	for i := range rs {
		rs[i].setGTag(tag)
	}
}

func (rs Restaurants) setIDIfEmpty() {
	for i := range rs {
		rs[i].setIDIfEmpty()
	}
}
