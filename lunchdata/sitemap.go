package lunchdata

type SiteMap map[string]*Site

func (sm SiteMap) Len() int {
	return len(sm)
}

func (sm SiteMap) Empty() bool {
	return sm.Len() == 0
}

func (sm SiteMap) NumRestaurants() int {
	total := 0
	for _, s := range sm {
		total += s.NumRestaurants()
	}
	return total
}

func (sm SiteMap) NumDishes() int {
	total := 0
	for _, s := range sm {
		total += s.NumDishes()
	}
	return total
}

func (sm SiteMap) Total() int {
	total := 0
	for _, s := range sm {
		total += s.Restaurants.Total()
	}
	return total + sm.Len()
}

func (sm SiteMap) Add(sites ...*Site) {
	for _, site := range sites {
		if site != nil {
			sm[site.ID] = site
		}
	}
}
func (sm SiteMap) Delete(ids ...string) {
	for _, id := range ids {
		delete(sm, id)
	}
}

func (sm SiteMap) Get(id string) *Site {
	if sm == nil {
		return nil
	}
	s, found := sm[id]
	if !found {
		return nil
	}
	return s
}

func (sm SiteMap) SetGTag(tag string) {
	for _, s := range sm {
		s.SetGTag(tag)
	}
}
