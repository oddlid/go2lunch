package lunchdata

// A giant list of everything
type LunchList struct {
	ID        string    `json:"id"`
	Countries Countries `json:"countries"`
}

func (l *LunchList) NumCountries() int {
	if l == nil {
		return 0
	}
	return l.Countries.Len()
}

func (l *LunchList) NumCities() int {
	if l == nil {
		return 0
	}
	return l.Countries.NumCities()
}

func (l *LunchList) NumSites() int {
	if l == nil {
		return 0
	}
	return l.Countries.NumSites()
}

func (l *LunchList) NumRestaurants() int {
	if l == nil {
		return 0
	}
	return l.Countries.NumRestaurants()
}

func (l *LunchList) NumDishes() int {
	if l == nil {
		return 0
	}
	return l.Countries.NumDishes()
}

func (l *LunchList) Get(f CountryMatch) *Country {
	if l == nil {
		return nil
	}
	return l.Countries.Get(f)
}

func (l *LunchList) GetByID(id string) *Country {
	if l == nil {
		return nil
	}
	return l.Countries.GetByID(id)
}
