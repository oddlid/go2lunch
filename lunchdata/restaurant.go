package lunchdata

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	Name     string    `json:"name"`
	ID       string    `json:"id"`
	URL      string    `json:"url,omitempty"`
	GTag     string    `json:"-"`
	Address  string    `json:"address"`
	MapURL   string    `json:"map_url"`
	ParsedAt time.Time `json:"parsed_at"`
	Dishes   Dishes    `json:"dishes"`
}

func (r *Restaurant) NumDishes() int {
	if r == nil {
		return 0
	}
	return r.Dishes.Len()
}

// ParsedRFC3339 returns the date in RFC3339 format
func (r *Restaurant) ParsedRFC3339() string {
	if r == nil {
		return time.Time{}.Format(time.RFC3339)
	}
	return r.ParsedAt.Format(time.RFC3339)
}

// ParsedHumanDate returns a more human readable date/time format, without too much detail
func (r *Restaurant) ParsedHumanDate() string {
	if r == nil {
		return time.Time{}.Format(dateFormat)
	}
	return r.ParsedAt.Format(dateFormat)
}

func (r *Restaurant) Get(f DishMatch) *Dish {
	if r == nil {
		return nil
	}
	return r.Dishes.Get(f)
}

func (r *Restaurant) GetByID(id string) *Dish {
	if r == nil {
		return nil
	}
	return r.Dishes.GetByID(id)
}

func (r *Restaurant) setIDIfEmpty() {
	if r == nil {
		return
	}
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	r.Dishes.setIDIfEmpty()
}
