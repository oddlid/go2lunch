package lunchdata

type Sites []*Site

func (ss Sites) Len() int {
	return len(ss)
}

func (ss Sites) Empty() bool {
	return ss.Len() == 0
}

func (ss Sites) NumRestaurants() int {
	total := 0
	for _, site := range ss {
		total += site.NumRestaurants()
	}
	return total
}

func (ss Sites) NumDishes() int {
	total := 0
	for _, s := range ss {
		total += s.NumDishes()
	}
	return total
}

func (ss Sites) Total() int {
	total := 0
	for _, s := range ss {
		total += s.Restaurants.Total()
	}
	return total + ss.Len()
}

func (ss Sites) setGTag(tag string) {
	for _, s := range ss {
		s.setGTag(tag)
	}
}

func (ss Sites) AsMap() SiteMap {
	sMap := make(SiteMap)
	sMap.Add(ss...)
	return sMap
}
