package lunchdata

type mockSiteScraper struct {
	err         error
	countryID   string
	cityID      string
	siteID      string
	restaurants Restaurants
}

func (s *mockSiteScraper) Scrape() (Restaurants, error) {
	return s.restaurants, s.err
}

func (s *mockSiteScraper) CountryID() string {
	return s.countryID
}

func (s *mockSiteScraper) CityID() string {
	return s.cityID
}

func (s *mockSiteScraper) SiteID() string {
	return s.siteID
}
