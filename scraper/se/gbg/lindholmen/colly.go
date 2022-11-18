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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/google/uuid"
	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/rs/zerolog"

	"github.com/oddlid/go2lunch/lunchdata"
)

const (
	DefaultScrapeURL    = `https://lindholmen.uit.se/omradet/dagens-lunch?embed-mode=iframe`
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
	// userAgent           = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36`
)

type Scraper struct {
	URL    string
	Logger zerolog.Logger
}

// Encode ID field. Might find a better strategy for this later
// func getRestaurantID(name string) string {
// 	return url.PathEscape(strings.ToLower(name))
// }

// func byURL(url string) lunchdata.RestaurantMatch {
// 	return func(r lunchdata.Restaurant) bool { return r.URL == url }
// }

func getRestaurantNameLinkName(name string) string {
	return urlPrefix + hyphenRX.ReplaceAllString(
		restaurantNameReplacer.Replace(
			strings.ToLower(name),
		),
		"-",
	)
}

func (Scraper) CountryID() string {
	return countryID
}

func (Scraper) CityID() string {
	return cityID
}

func (Scraper) SiteID() string {
	return siteID
}

func (lhs *Scraper) Scrape() (lunchdata.Restaurants, error) {
	// lindholmen.se has changed the whole way they present menus. The menu is not available anymore on each restaurant page,
	// so we need to parse the single page with all restaurants and menus instead. This is not even hosted on lindholmen.se anymore,
	// but on https://lindholmen.uit.se/omradet/dagens-lunch?embed-mode=iframe (important to have the embed-mode in the url, or the site will be blocked with http auth)

	if lhs.URL == "" {
		lhs.URL = DefaultScrapeURL
	}

	restaurantMap := make(lunchdata.RestaurantMap)
	menuCollector := colly.NewCollector()
	extensions.RandomUserAgent(menuCollector)

	addrCollector := menuCollector.Clone()
	extensions.RandomUserAgent(addrCollector)
	addrCollector.Async = true
	if err := addrCollector.Limit(&colly.LimitRule{DomainGlob: "*.lindholmen.se", Parallelism: 32}); err != nil {
		return nil, err
	}

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

	menuCollector.OnHTML(selectorViewContent, func(e *colly.HTMLElement) {
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

			linkName := getRestaurantNameLinkName(name)
			restaurant := &lunchdata.Restaurant{
				Name:   name,
				ID:     linkName,
				URL:    linkName,
				Parsed: time.Now(),
			}

			h.DOM.NextFilteredUntil(selectorDishRow, selectorTitle).Each(func(_ int, s *goquery.Selection) {
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
				restaurant.Add(
					&lunchdata.Dish{
						ID:    uuid.NewString(),
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

			lhs.Logger.Trace().Str(keyURL, restaurant.URL).Msg("Starting scrape for maps link")
			if err := addrCollector.Visit(restaurant.URL); err != nil {
				lhs.Logger.Error().Err(err).Send()
			}
		})
	})

	lhs.Logger.Debug().
		Str(keyURL, lhs.URL).
		Msg("Start scraping menus")

	startTimeMenus := time.Now()
	startTimeAddrs := time.Now()
	if err := menuCollector.Visit(lhs.URL); err != nil {
		return nil, err
	}
	endTimeMenus := time.Now()
	addrCollector.Wait()
	endTimeAddrs := time.Now()

	lhs.Logger.Debug().
		Str(keyURL, lhs.URL).
		Dur(keyParsedTime, endTimeMenus.Sub(startTimeMenus)).
		Int(keyNumRestaurants, restaurantMap.Len()).
		Int(keyNumDishes, restaurantMap.NumDishes()).
		Msg("Restaurants and menus parsed")

	lhs.Logger.Debug().
		Str(keyURL, lhs.URL).
		Dur(keyParsedTime, endTimeAddrs.Sub(startTimeAddrs)).
		Msg("Map links parsed")

	lhs.Logger.Debug().
		Str(keyURL, lhs.URL).
		Dur(keyParsedTime, time.Since(startTimeMenus)).
		Msg("Site parsed")

	restaurants := make(lunchdata.Restaurants, 0, restaurantMap.Len())
	for _, restaurant := range restaurantMap {
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}
