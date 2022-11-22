package lunchdata

type Cities []*City

func (cs Cities) Len() int {
	return len(cs)
}

func (cs Cities) Empty() bool {
	return cs.Len() == 0
}

func (cs Cities) NumSites() int {
	total := 0
	for _, c := range cs {
		total += c.NumSites()
	}
	return total
}

func (cs Cities) NumRestaurants() int {
	total := 0
	for _, c := range cs {
		total += c.NumRestaurants()
	}
	return total
}

func (cs Cities) NumDishes() int {
	total := 0
	for _, c := range cs {
		total += c.NumDishes()
	}
	return total
}

func (cs Cities) Total() int {
	total := 0
	for _, c := range cs {
		total += c.Sites.Total()
	}
	return total + cs.Len()
}

func (cs Cities) setGTag(tag string) {
	for _, c := range cs {
		c.setGTag(tag)
	}
}

func (cs Cities) AsMap() CityMap {
	cMap := make(CityMap)
	cMap.Add(cs...)
	return cMap
}
