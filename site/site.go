package site

import (
	"time"
	"io"
	"bufio"
	"os"
	"encoding/json"
)

type Dish struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Price string `json:"price"`
}

type Restaurant struct {
	Name   string    `json:"name"`
	Url    string    `json:"url"`
	Parsed time.Time `json:"date"`
	Dishes []Dish    `json:"dishes"`
}

type Site struct {
	Name        string       `json:"name"`
	ID          string       `json:"siteid"`
	Comment     string       `json:"comment,omitempty"`
	Restaurants []Restaurant `json:"restaurants"`
}

func (r Restaurant) ParsedRFC3339() string {
	return r.Parsed.Format(time.RFC3339)
}

func (r *Restaurant) Add(d Dish) *Restaurant {
	r.Dishes = append(r.Dishes, d)
	return r
}

func (s *Site) Add(r Restaurant) *Site {
	s.Restaurants = append(s.Restaurants, r)
	return s
}

func (s *Site) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(s)
}

func (s *Site) SaveJSON(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	return s.Encode(w)
}

func NewFromJSON(r io.Reader) (*Site, error) {
	dec := json.NewDecoder(r)
	var s Site
	if err := dec.Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func NewFromFile(fileName string) (*Site, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	return NewFromJSON(r)
}

