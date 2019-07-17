package lunchdata

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

// A giant list of everything
type LunchList struct {
	Countries map[string]*Country `json:"countries"`
}

func NewLunchList() *LunchList {
	return &LunchList{
		Countries: make(map[string]*Country),
	}
}

func (ll *LunchList) AddCountry(c Country) *LunchList {
	ll.Countries[c.ID] = &c
	return ll
}

func (ll *LunchList) GetCountryById(id string) *Country {
	return ll.Countries[id]
}

func (ll *LunchList) NumCountries() int {
	return len(ll.Countries)
}

func (ll *LunchList) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(ll)
}

func (ll *LunchList) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(ll)
}

func (ll *LunchList) SaveJSON(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = ll.Encode(w)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func LunchListFromJSON(r io.Reader) (*LunchList, error) {
	ll := &LunchList{}
	if err := ll.Decode(r); err != nil {
		return nil, err
	}
	return ll, nil
}

func LunchListFromFile(fileName string) (*LunchList, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	return LunchListFromJSON(r)
}

func (ll *LunchList) GetSiteLinks() SiteLinks {
	sl := make(SiteLinks, 0)

	for _, country := range ll.Countries {
		for _, city := range country.Cities {
			for _, site := range city.Sites {
				sl = append(sl, SiteLink{
					CountryName: country.Name,
					CityName:    city.Name,
					SiteName:    site.Name,
					Url:         site.ID,
				})
			}
		}
	}

	return sl
}
