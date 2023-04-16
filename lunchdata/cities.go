package lunchdata

type Cities []City
type CityMatch func(c City) bool

func (cs Cities) Len() int {
	return len(cs)
}

func (cs Cities) NumSites() int {
	total := 0
	for i := range cs {
		total += cs[i].NumSites()
	}
	return total
}

func (cs Cities) NumRestaurants() int {
	total := 0
	for i := range cs {
		total += cs[i].NumRestaurants()
	}
	return total
}

func (cs Cities) NumDishes() int {
	total := 0
	for i := range cs {
		total += cs[i].NumDishes()
	}
	return total
}

func (cs Cities) Total() int {
	total := 0
	for i := range cs {
		total += cs[i].Sites.Total()
	}
	return total + cs.Len()
}

func (cs Cities) Get(f CityMatch) *City {
	if idx := sliceIndex(cs, f); idx > -1 {
		return &cs[idx]
	}
	return nil
}

func (cs Cities) GetByID(id string) *City {
	return cs.Get(func(c City) bool { return c.ID == id })
}
