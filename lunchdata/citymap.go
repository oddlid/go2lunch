package lunchdata

type CityMap map[string]*City

func (cm CityMap) Len() int {
	return len(cm)
}

func (cm CityMap) Empty() bool {
	return cm.Len() == 0
}

func (cm CityMap) NumSites() int {
	total := 0
	for _, c := range cm {
		total += c.NumSites()
	}
	return total
}

func (cm CityMap) NumRestaurants() int {
	total := 0
	for _, c := range cm {
		total += c.NumRestaurants()
	}
	return total
}

func (cm CityMap) NumDishes() int {
	total := 0
	for _, c := range cm {
		total += c.NumDishes()
	}
	return total
}

func (cm CityMap) Total() int {
	total := 0
	for _, c := range cm {
		total += c.Sites.Total()
	}
	return total + cm.Len()
}

func (cm CityMap) Add(cities ...*City) {
	if cm == nil {
		return
	}
	for _, c := range cities {
		if c != nil {
			cm[c.ID] = c
		}
	}
}

func (cm CityMap) Delete(ids ...string) {
	for _, id := range ids {
		delete(cm, id)
	}
}

func (cm CityMap) Get(id string) *City {
	return cm[id]
}

func (cm CityMap) SetGTag(tag string) {
	for _, c := range cm {
		c.SetGTag(tag)
	}
}
