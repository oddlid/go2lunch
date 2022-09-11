package lunchdata

type Dishes []*Dish

func (ds Dishes) Len() int {
	return len(ds)
}

func (ds Dishes) Empty() bool {
	return ds.Len() == 0
}

func (ds Dishes) SetGTag(tag string) {
	for _, d := range ds {
		if d != nil {
			d.GTag = tag
		}
	}
}
