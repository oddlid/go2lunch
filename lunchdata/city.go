package lunchdata

type City struct {
	Name  string `json:"name"`
	ID    string `json:"id"` // e.g. osl, gbg or something like the airlines use
	Sites Sites  `json:"sites"`
}

func (c *City) NumSites() int {
	if c == nil {
		return 0
	}
	return c.Sites.Len()
}

func (c *City) NumRestaurants() int {
	if c == nil {
		return 0
	}
	return c.Sites.NumRestaurants()
}

func (c *City) NumDishes() int {
	if c == nil {
		return 0
	}
	return c.Sites.NumDishes()
}

func (c *City) Get(f SiteMatch) *Site {
	if c == nil {
		return nil
	}
	return c.Sites.Get(f)
}

func (c *City) GetByID(id string) *Site {
	if c == nil {
		return nil
	}
	return c.Sites.GetByID(id)
}
