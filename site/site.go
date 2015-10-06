package site

//import (
//	"log"
//	"io/ioutil"
//	"encoding/json"
//  "os"
//  "github.com/codegangsta/cli"
//)

type Dish struct {
	Name  string
	Desc  string
	Price string
}

type Restaurant struct {
	Name   string
	Url    string
	Dishes []Dish
}

type Site struct {
	Name        string
	Comment     string
	Restaurants []Restaurant
}

type LunchParser interface {
	Scrape() (Restaurant, error)
}


func (r *Restaurant) Add(d *Dish) *Restaurant {
	r.Dishes = append(r.Dishes, *d)
	return r
}

func (s *Site) Add(r *Restaurant) *Site {
	s.Restaurants = append(s.Restaurants, *r)
	return s
}


