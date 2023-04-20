package lunchdata

type Dish struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price int    `json:"price"`
}

// String implements fmt.Stringer
func (d *Dish) String() string {
	if d == nil {
		return ""
	}
	return d.Name + " " + d.Desc
}
