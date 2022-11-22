package lunchdata

type PropertyField int

const (
	ID PropertyField = iota
	GTag
	Name
	Desc
	Price
	URL
	Address
	MapURL
	ParsedAt
	Comment
	DishList
	RestaurantList
	SiteList
	CityList
	CountryList
)
