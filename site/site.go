package site

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"time"
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

type Restaurants []Restaurant

type Site struct {
	Name        string      `json:"name"`
	ID          string      `json:"siteid"` // eg. se/gbg/lindholmen or the url it came from
	Comment     string      `json:"comment,omitempty"`
	Restaurants Restaurants `json:"restaurants"`
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
	defer f.Close()
	w := bufio.NewWriter(f)
	err = s.Encode(w)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
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
	defer f.Close()
	r := bufio.NewReader(f)
	return NewFromJSON(r)
}
