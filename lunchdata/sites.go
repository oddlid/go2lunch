package lunchdata

type Sites []Site
type SiteMatch func(s Site) bool

func (ss Sites) Len() int {
	return len(ss)
}

func (ss Sites) NumRestaurants() int {
	total := 0
	for i := range ss {
		total += ss[i].NumRestaurants()
	}
	return total
}

func (ss Sites) NumDishes() int {
	total := 0
	for i := range ss {
		total += ss[i].NumDishes()
	}
	return total
}

func (ss Sites) Total() int {
	total := 0
	for i := range ss {
		total += ss[i].Restaurants.Total()
	}
	return total + ss.Len()
}

func (ss Sites) Get(f SiteMatch) *Site {
	if idx := sliceIndex(ss, f); idx > -1 {
		return &ss[idx]
	}
	return nil
}

func (ss Sites) GetByID(id string) *Site {
	return ss.Get(func(s Site) bool { return s.ID == id })
}
