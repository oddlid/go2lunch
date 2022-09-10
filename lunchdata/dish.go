package lunchdata

type Dish struct {
	Name  string `json:"dish_name"`
	Desc  string `json:"dish_desc"`
	GTag  string `json:"-"`
	Price int    `json:"dish_price"`
}

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

func (d Dish) String() string {
	return d.Name + " " + d.Desc
}
