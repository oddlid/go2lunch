package lunchdata

type Restaurants []*Restaurant

func (rs Restaurants) Len() int {
	return len(rs)
}

func (rs Restaurants) Empty() bool {
	return rs.Len() == 0
}

func (rs Restaurants) NumDishes() int {
	total := 0
	for _, r := range rs {
		total += r.NumDishes()
	}
	return total
}

func (rs Restaurants) Total() int {
	return rs.Len() + rs.NumDishes()
}

func (rs Restaurants) setGTag(tag string) {
	for _, r := range rs {
		r.setGTag(tag)
	}
}

func (rs Restaurants) AsMap() RestaurantMap {
	rMap := make(RestaurantMap)
	rMap.Add(rs...)
	return rMap
}
