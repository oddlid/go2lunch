package lunchdata

type Dishes []*Dish

func (ds Dishes) Len() int {
	return len(ds)
}

func (ds Dishes) Empty() bool {
	return ds.Len() == 0
}

func (ds Dishes) Clone() Dishes {
	ret := make(Dishes, 0, ds.Len())
	for _, d := range ds {
		ret = append(ret, d.Clone())
	}
	return ret
}

func (ds Dishes) setGTag(tag string) {
	if ds == nil {
		return
	}
	for _, d := range ds {
		if d != nil {
			d.GTag = tag
		}
	}
}

func (ds Dishes) setIDIfEmpty() {
	for _, d := range ds {
		d.setIDIfEmpty()
	}
}
