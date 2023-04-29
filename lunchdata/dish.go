package lunchdata

import "fmt"

type Dish struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Desc  string  `json:"desc"`
	Price float32 `json:"price"`
}

// String implements fmt.Stringer
func (d Dish) String() string {
	return fmt.Sprintf("%s %s :: %.2f", d.Name, d.Desc, d.Price)
}
