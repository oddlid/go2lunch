package lunchdata

type Dish struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	GTag  string `json:"-"`
	Price int    `json:"price"`
}

func (d Dish) String() string {
	return d.Name + " " + d.Desc
}
