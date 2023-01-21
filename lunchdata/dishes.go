package lunchdata

type Dishes []Dish
type DishMatch func(d Dish) bool

func (ds Dishes) Len() int {
	return len(ds)
}

func (ds Dishes) Empty() bool {
	return ds.Len() == 0
}

func (ds Dishes) get(f DishMatch) *Dish {
	if idx := sliceIndex(ds, f); idx > -1 {
		return &ds[idx]
	}
	return nil
}

func (ds Dishes) getByID(id string) *Dish {
	return ds.get(func(d Dish) bool { return d.ID == id })
}

func (ds Dishes) setGTag(tag string) {
	for i := range ds {
		ds[i].setGTag(tag)
	}
}

func (ds Dishes) setIDIfEmpty() {
	for i := range ds {
		ds[i].setIDIfEmpty()
	}
}
