package lunchdata

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

// Structs for linking into the larger structs
type SiteLink struct {
	CountryName string `json:"country_name"`
	CountryID   string `json:"country_id"`
	CityName    string `json:"city_name"`
	CityID      string `json:"city_id"`
	SiteName    string `json:"site_name"`
	SiteID      string `json:"site_id"`
	SiteKey     string `json:"site_key"`
	Comment     string `json:"comment,omitempty"`
	URL         string `json:"url"`
}

// Stripped down struct only for linking to the correct site, for saving/loading keys
type SiteKeyLink struct {
	CountryID string `json:"country_id"`
	CityID    string `json:"city_id"`
	SiteID    string `json:"site_id"`
	SiteKey   string `json:"site_key"`
}

type SiteLinks []*SiteLink
type SiteKeyLinks []*SiteKeyLink

func (sls *SiteLinks) Add(sl *SiteLink) {
	*sls = append(*sls, sl)
}

func (sks *SiteKeyLinks) Add(sk *SiteKeyLink) {
	*sks = append(*sks, sk)
}

func (sls *SiteLinks) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(sls)
}

func (sks *SiteKeyLinks) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(sks)
}

func (sls *SiteLinks) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(sls)
}

func (sks *SiteKeyLinks) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(sks)
}

func (sls *SiteLinks) SaveJSON(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = sls.Encode(w)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func (sks *SiteKeyLinks) SaveJSON(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = sks.Encode(w)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func SiteLinksFromJSON(r io.Reader) (SiteLinks, error) {
	sls := make(SiteLinks, 0)
	if err := sls.Decode(r); err != nil {
		return nil, err
	}
	return sls, nil
}

func SiteKeyLinksFromJSON(r io.Reader) (SiteKeyLinks, error) {
	sks := make(SiteKeyLinks, 0)
	if err := sks.Decode(r); err != nil {
		return nil, err
	}
	return sks, nil
}

func SiteLinksFromFile(fileName string) (SiteLinks, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	return SiteLinksFromJSON(r)
}

func SiteKeyLinksFromFile(fileName string) (SiteKeyLinks, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	return SiteKeyLinksFromJSON(r)
}
