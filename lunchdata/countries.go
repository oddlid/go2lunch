package lunchdata

type Countries []Country
type CountryMatch func(c Country) bool

func (cs Countries) Len() int {
	return len(cs)
}

func (cs Countries) NumCities() int {
	total := 0
	for i := range cs {
		total += cs[i].NumCities()
	}
	return total
}

func (cs Countries) NumSites() int {
	total := 0
	for i := range cs {
		total += cs[i].NumSites()
	}
	return total
}

func (cs Countries) NumRestaurants() int {
	total := 0
	for i := range cs {
		total += cs[i].NumRestaurants()
	}
	return total
}

func (cs Countries) NumDishes() int {
	total := 0
	for i := range cs {
		total += cs[i].NumDishes()
	}
	return total
}

func (cs Countries) Total() int {
	total := 0
	for i := range cs {
		total += cs[i].Cities.Total()
	}
	return total + cs.Len()
}

func (cs Countries) Get(f CountryMatch) *Country {
	if idx := sliceIndex(cs, f); idx > -1 {
		return &cs[idx]
	}
	return nil
}

func (cs Countries) GetByID(id string) *Country {
	return cs.Get(func(c Country) bool { return c.ID == id })
}
