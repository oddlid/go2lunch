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

func (s *mockSiteScraper) GetCountryID() string {
	return s.countryID
}

func (s *mockSiteScraper) GetCityID() string {
	return s.cityID
}

func (s *mockSiteScraper) GetSiteID() string {
	return s.siteID
}
