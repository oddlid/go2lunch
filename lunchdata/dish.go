package lunchdata

import (
	"encoding/json"
	"io"
)

type Dish struct {
	Name  string `json:"dish_name"`
	Desc  string `json:"dish_desc"`
	Price int    `json:"dish_price"`
	Gtag  string `json:"-"`
}

type Dishes []Dish

func (ds *Dishes) Add(d *Dish) {
	*ds = append(*ds, *d)
}

func NewDish(name, desc, tag string, price int) *Dish {
	return &Dish{
		Name:  name,
		Desc:  desc,
		Price: price,
		Gtag:  tag,
	}
}

func (d *Dish) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(d)
}

func (d *Dish) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(d)
}

func DishFromJSON(r io.Reader) (*Dish, error) {
	d := &Dish{}
	if err := d.Decode(r); err != nil {
		return nil, err
	}
	return d, nil
}
