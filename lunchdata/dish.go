package lunchdata

type Dish struct {
	Name  string `json:"dish_name"`
	Desc  string `json:"dish_desc"`
	GTag  string `json:"-"`
	Price int    `json:"dish_price"`
}

func (d Dish) String() string {
	return d.Name + " " + d.Desc
}
