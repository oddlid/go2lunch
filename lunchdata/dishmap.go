package lunchdata

type DishMap map[string]*Dish

func (dm DishMap) Len() int {
	return len(dm)
}

func (dm DishMap) Dishes() Dishes {
	dishes := make(Dishes, 0, dm.Len())
	for _, dish := range dm {
		dishes = append(dishes, *dish)
	}
	return dishes
}
