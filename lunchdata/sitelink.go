package lunchdata

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
