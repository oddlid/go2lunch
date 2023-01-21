package lunchdata

import "github.com/google/uuid"

type Dish struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	GTag  string `json:"-"`
	Price int    `json:"price"`
}

func (d *Dish) String() string {
	if d == nil {
		return ""
	}
	return d.Name + " " + d.Desc
}

func (d *Dish) setGTag(tag string) {
	if d == nil {
		return
	}
	d.GTag = tag
}

func (d *Dish) setIDIfEmpty() {
	if d == nil {
		return
	}
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
}
