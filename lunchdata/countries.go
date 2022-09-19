package lunchdata

type Countries []*Country

func (cs Countries) Len() int {
	return len(cs)
}

func (cs Countries) Empty() bool {
	return cs.Len() == 0
}

func (cs Countries) NumCities() int {
	total := 0
	for _, c := range cs {
		total += c.NumCities()
	}
	return total
}

func (cs Countries) NumSites() int {
	total := 0
	for _, c := range cs {
		total += c.NumSites()
	}
	return total
}

func (cs Countries) NumRestaurants() int {
	total := 0
	for _, c := range cs {
		total += c.NumRestaurants()
	}
	return total
}

func (cs Countries) NumDishes() int {
	total := 0
	for _, c := range cs {
		total += c.NumDishes()
	}
	return total
}

func (cs Countries) Total() int {
	total := 0
	for _, c := range cs {
		total += c.Cities.Total()
	}
	return total + cs.Len()
}

func (cs Countries) SetGTag(tag string) {
	for _, c := range cs {
		c.SetGTag(tag)
	}
}

func (cs Countries) AsMap() CountryMap {
	cMap := make(CountryMap)
	cMap.Add(cs...)
	return cMap
}
