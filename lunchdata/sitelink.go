package lunchdata

import (
	"encoding/json"
	"io"
)

// Structs for linking into the larger structs
type SiteLink struct {
	CountryName string `json:"country_name"`
	CityName    string `json:"city_name"`
	SiteName    string `json:"site_name"`
	Comment     string `json:"comment,omitempty"`
	Url         string `json:"url"`
}

type SiteLinks []SiteLink

func (sls SiteLinks) Add(sl SiteLink) {
	sls = append(sls, sl)
}

func (sls SiteLinks) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(&sls)
}

func (sls SiteLinks) Decode(r io.Reader) error {
	return json.NewDecoder(r).Decode(&sls)
}
