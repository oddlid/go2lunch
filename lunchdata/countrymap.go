package lunchdata

type CountryMap map[string]*Country

func (cm CountryMap) Len() int {
	return len(cm)
}

func (cm CountryMap) Empty() bool {
	return cm.Len() == 0
}

func (cm CountryMap) NumCities() int {
	total := 0
	for _, c := range cm {
		total += c.NumCities()
	}
	return total
}

func (cm CountryMap) NumSites() int {
	total := 0
	for _, c := range cm {
		total += c.NumSites()
	}
	return total
}

func (cm CountryMap) NumRestaurants() int {
	total := 0
	for _, c := range cm {
		total += c.NumRestaurants()
	}
	return total
}

func (cm CountryMap) NumDishes() int {
	total := 0
	for _, c := range cm {
		total += c.NumDishes()
	}
	return total
}

func (cm CountryMap) Total() int {
	total := 0
	for _, c := range cm {
		total += c.Cities.Total()
	}
	return total + cm.Len()
}

func (cm CountryMap) Add(countries ...*Country) {
	if cm == nil {
		return
	}
	for _, c := range countries {
		if c != nil {
			cm[c.ID] = c
		}
	}
}

func (cm CountryMap) Delete(ids ...string) {
	for _, id := range ids {
		delete(cm, id)
	}
}

func (cm CountryMap) Get(id string) *Country {
	return cm[id]
}

func (cm CountryMap) SetGTag(tag string) {
	for _, c := range cm {
		c.SetGTag(tag)
	}
}
