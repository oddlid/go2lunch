package lindholmen

/*
2019-10-24 22:52
Turns out this solution takes many times as long as the original single thread/page scraper.
No doubt due to the many more http requests here. But I like the concept of colly, so I'm
keeping this code only for reference, but excluding it from builds.

2021-05-20 20:23
This is now the scraper used, as the Google Maps link feature demands we scrape each
restaurant page to get it, and thus this strategy fits the bill, even though it's
slower, which doesn't really matter, as it happens in the background.
*/

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/rs/zerolog"
)

const (
	userAgent           = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`
	countryID           = `se`
	cityID              = `gbg`
	siteID              = `lindholmen`
	selectorViewContent = `div.view-content`
	selectorContent     = `div.content`
	selectorTitle       = `h3.title`
	selectorDishRow     = `div.table-list__row`
	selectorDish        = `span.dish-name`
	selectorDishName    = `strong`
	selectorPrice       = `div.table-list__column--price`
	keyURL              = `url`
	keyRestaurant       = `restaurant`
	keyNumRestaurants   = `numRestaurants`
	keyDish             = `dish`
	keyNumDishes        = `numDishes`
	keyParsedTime       = `parsedTimeMS`
	keyMapURL           = `mapURL`
	keyAddr             = `addr`
	keyLink             = `link`
)

type LHScraper struct {
	URL    string
	Logger zerolog.Logger
}

// Encode ID field. Might find a better strategy for this later
func getRestaurantID(name string) string {
	return url.PathEscape(strings.ToLower(name))
}

func (LHScraper) GetCountryID() string {
	return countryID
}

func (LHScraper) GetCityID() string {
	return cityID
}

func (LHScraper) GetSiteID() string {
	return siteID
}

func (lhs *LHScraper) Scrape() (lunchdata.Restaurants, error) {
	// lindholmen.se has changed the whole way they present menus. The menu is not available anymore on each restaurant page,
	// so we need to parse the single page with all restaurants and menus instead. This is not even hosted on lindholmen.se anymore,
	// but on https://lindholmen.uit.se/omradet/dagens-lunch?embed-mode=iframe (important to have the embed-mode in the url, or the site will be blocked with http auth)

	restaurantMap := make(map[string]*lunchdata.Restaurant)
	collector := colly.NewCollector(colly.UserAgent(userAgent))

	collector.OnHTML(selectorViewContent, func(e *colly.HTMLElement) {
		e.ForEach(selectorTitle, func(i int, h *colly.HTMLElement) {
			name := strings.TrimSpace(h.ChildText("a"))
			// we only want the last part of the link, since the links on this page are not correct,
			// so we need to reconstruct them ourselves later
			link := strings.Replace(
				h.ChildAttr("a", "href"),
				"/restauranger/", "", 1,
			)

			lhs.Logger.Trace().
				Str(keyRestaurant, name).
				Msg("Adding restaurant")

			restaurant := lunchdata.NewRestaurant(
				name,
				getRestaurantID(name),
				"https://www.lindholmen.se/sv/"+link, // fill in the correct prefix for the link
				time.Now(),
			)

			h.DOM.NextFilteredUntil(selectorDishRow, selectorTitle).Each(func(i int, s *goquery.Selection) {
				dishSelection := s.Find(selectorDish)
				dishName := strings.TrimSpace(dishSelection.Find(selectorDishName).Text())
				dishDesc := dishSelection.Text()
				dishDesc = strings.TrimSpace(strings.Replace(dishDesc, dishName, "", 1))
				dishDesc = strings.Join(strings.Fields(dishDesc), " ")
				dishPrice := strings.TrimSpace(
					// we only want the value, so it can be stored as an int
					strings.Replace(
						s.Find(selectorPrice).Text(),
						"kr", "", 1,
					),
				)
				price, err := strconv.Atoi(dishPrice)
				if err != nil {
					lhs.Logger.Error().
						Err(err).
						Msg("Failed to parse dish price")
					price = -1
				}
				restaurant.AddDish(
					&lunchdata.Dish{
						Name:  dishName,
						Desc:  dishDesc,
						Price: price,
					},
				)
				lhs.Logger.Trace().
					Str(keyRestaurant, name).
					Str(keyDish, dishName).
					Send()
			})
			restaurantMap[restaurant.URL] = restaurant
		})
	})

	lhs.Logger.Debug().
		Str(keyURL, lhs.URL).
		Msg("Start scraping menus")

	startTimeMenus := time.Now()
	if err := collector.Visit(lhs.URL); err != nil {
		return nil, err
	}

	restaurants := make(lunchdata.Restaurants, 0, len(restaurantMap))
	for _, restaurant := range restaurantMap {
		restaurants = append(restaurants, restaurant)
	}

	lhs.Logger.Debug().
		Str(keyURL, lhs.URL).
		Dur(keyParsedTime, time.Since(startTimeMenus)).
		Int(keyNumRestaurants, restaurants.Len()).
		Int(keyNumDishes, restaurants.NumDishes()).
		Msg("Restaurants and menus parsed")

	addrCollector := collector.Clone()
	addrCollector.Async = true

	addrCollector.OnHTML(selectorContent, func(e *colly.HTMLElement) {
		lhs.Logger.Trace().
			Str(keyURL, e.Request.URL.String()).
			Msg("Looking for map link...")

		restaurant, found := restaurantMap[e.Request.URL.String()]
		if !found {
			lhs.Logger.Error().
				Str(keyURL, e.Request.URL.String()).
				Msg("No restaurant entry for URL")
			return
		}
		e.ForEachWithBreak("p > a", func(i int, h *colly.HTMLElement) bool {
			link := h.Attr("href")
			if strings.Contains(link, "maps.google.com") {
				mapURL, err := url.Parse(link)
				if err != nil {
					lhs.Logger.Error().Err(err).Send()
					return true
				}
				restaurant.MapURL = mapURL.String()

				query := mapURL.Query().Get("q")
				address, err := url.QueryUnescape(query)
				if err != nil {
					lhs.Logger.Error().Err(err).Send()
					return true
				}
				restaurant.Address = address

				lhs.Logger.Trace().
					Str(keyURL, e.Request.URL.String()).
					Str(keyMapURL, restaurant.MapURL).
					Str(keyAddr, address).
					Str(keyRestaurant, restaurant.Name).
					Msg("Parsed map URL and address")

				return false
			}
			lhs.Logger.Trace().
				Str(keyURL, e.Request.URL.String()).
				Str(keyLink, link).
				Msg("Not a map link")
			return true
		})
	})

	startTimeMapLinks := time.Now()
	for url := range restaurantMap {
		if err := addrCollector.Visit(url); err != nil {
			lhs.Logger.Error().Err(err).Send()
		}
	}
	addrCollector.Wait()

	lhs.Logger.Debug().
		Str(keyURL, lhs.URL).
		Dur(keyParsedTime, time.Since(startTimeMapLinks)).
		Msg("Map links parsed")

	lhs.Logger.Debug().
		Str(keyURL, lhs.URL).
		Dur(keyParsedTime, time.Since(startTimeMenus)).
		Msg("Site parsed")

	return restaurants, nil
}
