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

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/oddlid/go2lunch/lunchdata"
	log "github.com/sirupsen/logrus"
)

const (
	UA         = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	TAG        = "LHScraper"
	COUNTRY_ID = "se"
	CITY_ID    = "gbg"
	SITE_ID    = "lindholmen"
)

var (
	logger = log.WithField("ScraperName", TAG)
)

type LHScraper struct{}

// Encode ID field. Might find a better strategy for this later
func getRestaurantID(name string) string {
	return url.PathEscape(strings.ToLower(name))
}

func (lhs LHScraper) GetCountryID() string {
	return COUNTRY_ID
}

func (lhs LHScraper) GetCityID() string {
	return CITY_ID
}

func (lhs LHScraper) GetSiteID() string {
	return SITE_ID
}

// helper func for parsing out address from maps url found on lindholmen.se
func getAddress(mapUrl string) string {
	// remove map url prefix
	query := strings.Replace(
		mapUrl,
		"http://maps.google.com/?q=",
		"",
		1,
	)
	// remove encoded newlines
	query = strings.Replace(
		query,
		"\n",
		"",
		-1,
	)
	// Remove leading and trailing '+'
	query = strings.Trim(query, "+")
	// Create plain text
	tmp, err := url.QueryUnescape(query)
	if nil == err {
		query = tmp
	}
	// Might be empty
	return query
}

/*
2019-10-23 21:37:
This is very strange... I basically copied this code from collytest.go
When I run collytest.go directly, it all works fine, but when I run this code included in
the larger program, run via a goroutine, I only get the first restaurant with it's dishes,
repeated for the number of restaurants that were actually parsed.
This sounds just like the problems I've had with (not) using pointers in the lunchdata structures...
Need to find out why this happens.
2019-10-24 22:51
Update: Turns out it was due to some obscure pointer dereferencing in Go I wasn't aware of. Not related
to the code in this module at all.
*/
func (lhs LHScraper) Scrape() (lunchdata.Restaurants, error) {
	rmap := make(map[string]*lunchdata.Restaurant)
	// create a queue for holding links to each restaurant to parse
	q, _ := queue.New(
		4, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	// Create collector with callbacks for picking out links to each restaurant page from the overview page
	lc := colly.NewCollector(colly.UserAgent(UA))
	lc.OnHTML("h3.restaurant-name > a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		q.AddURL(link) // add to queue

		restaurantName := strings.TrimSpace(e.Text)
		_, found := rmap[restaurantName]
		if !found {
			rmap[restaurantName] = lunchdata.NewRestaurant(
				restaurantName,
				getRestaurantID(restaurantName),
				link,
				time.Now(),
			)
		}
	})

	// This will block until done parsing, so we're sure to have all links before continuing
	t_start := time.Now()
	lc.Visit("https://www.lindholmen.se/omradet/restauranger")
	numLinks, err := q.Size()
	if err != nil {
		numLinks = -1
	}
	logger.WithFields(log.Fields{
		"Seconds":            time.Duration(time.Now().Sub(t_start)).Seconds(),
		"NumRestaurantLinks": numLinks,
	}).Debug("Time to parse overview page")

	// Create collector with callbacks for parsing a restaurant detail page, and picking out dishes from it
	rc := lc.Clone()
	rc.OnHTML("div.node.node-restaurant.node-full", func(e *colly.HTMLElement) {
		restaurantName := strings.TrimSpace(e.ChildText("h1.content__title.page-title.restaurant-name"))
		restaurant, found := rmap[restaurantName]
		if !found {
			restaurant = lunchdata.NewRestaurant(
				restaurantName,
				getRestaurantID(restaurantName),
				e.Request.URL.String(),
				time.Now(),
			)
			rmap[restaurantName] = restaurant
		}

		// Update @ 2021-05-19
		// Babak Eghbali suggested it would be good to have a Google Maps link to each restaurant,
		// so of course he will get that!
		e.ForEachWithBreak("div.restaurant-facts > p > a.right-arrow", func(_ int, el *colly.HTMLElement) bool {
			restaurant.Address = getAddress(el.Attr("href")) // might be empty
			return false
		})


		// Fetch dishes
		e.ForEach("div.node.node-dish", func(_ int, el *colly.HTMLElement) {
			dishName := strings.TrimSpace(el.ChildText("span.dish-name > strong"))
			// remove dish name from desc, so we don't get double up
			dishDesc := strings.TrimSpace(
				strings.Replace(
					el.ChildText("div.table-list__column.table-list__column--dish"),
					dishName,
					"",
					1,
				),
			)
			// replace redundant whitespace in desc, as we often get that from lindholmen.se
			dishDesc = strings.Join(
				strings.Fields(dishDesc), // split on any whitespace
				" ",                      // replace with just one space
			)
			// remove "kr" from price, so we get a pure int
			dishPrice := strings.TrimSpace(
				strings.Replace(
					el.ChildText("div.table-list__column.table-list__column--price"),
					"kr",
					"",
					1,
				),
			)

			price, err := strconv.Atoi(dishPrice)
			if err != nil {
				price = -1
			}

			restaurant.AddDish(
				lunchdata.Dish{
					Name:  dishName,
					Desc:  dishDesc,
					Price: price,
				})
		})
	})

	// start parsing all urls in the queue
	t_start = time.Now()
	q.Run(rc)
	t_end := time.Now()

	// get a []Restaurant slice from our map, so we can use it for setting restaurants for a site
	rs := make(lunchdata.Restaurants, 0, len(rmap))
	for _, r := range rmap {
		if r.HasDishes() {
			rs.Add(*r)
		}
	}

	logger.WithFields(log.Fields{
		"Seconds":     time.Duration(t_end.Sub(t_start)).Seconds(),
		"Restaurants": rs.Len(),
		"Dishes":      rs.NumDishes(),
	}).Debug("Parse results")

	return rs, nil
}
