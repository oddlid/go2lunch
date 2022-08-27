package lunchdata

type Dish struct {
	Name  string `json:"dish_name"`
	Desc  string `json:"dish_desc"`
	Price int    `json:"dish_price"`
	Gtag  string `json:"-"`
}

type Dishes []*Dish

func (ds Dishes) Len() int {
	return len(ds)
}

func (d Dish) String() string {
	return d.Name + " " + d.Desc
}
