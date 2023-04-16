package lunchdata

type Country struct {
	Name   string `json:"name"`
	ID     string `json:"id"` // preferably international country code, like "se", "no", and so on
	Cities Cities `json:"cities"`
}

func (c *Country) NumCities() int {
	if c == nil {
		return 0
	}
	return c.Cities.Len()
}

func (c *Country) NumSites() int {
	if c == nil {
		return 0
	}
	return c.Cities.NumSites()
}

func (c *Country) NumRestaurants() int {
	if c == nil {
		return 0
	}
	return c.Cities.NumRestaurants()
}

func (c *Country) NumDishes() int {
	if c == nil {
		return 0
	}
	return c.Cities.NumDishes()
}

func (c *Country) Get(f CityMatch) *City {
	if c == nil {
		return nil
	}
	return c.Cities.Get(f)
}

func (c *Country) GetByID(id string) *City {
	if c == nil {
		return nil
	}
	return c.Cities.GetByID(id)
}
