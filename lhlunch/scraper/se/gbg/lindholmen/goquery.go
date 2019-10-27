package lindholmen

/*
This is almost exactly the same scraper that has been for years, just adopted for
this new interface fitting the new code design, @2019-10-24

I have a version using colly as well, but no matter how much I like the code in
that one, it turns out you just can't beat the speed of this single thread solution,
parsing a single page and getting all info without more http requests.
*/

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/go2lunch/lunchdata"
)

const (
	URL = "https://www.lindholmen.se/pa-omradet/dagens-lunch"
	TAG = "LHScraper"
)

type LHScraper struct{} // only for having something to implement the SiteScraper interface

// getRestaurantUrl is a helper that tries to return the link to each restaurant
func getRestaurantUrl(name string) string {
	urlBase := "https://www.lindholmen.se/restauranger/"
	name = strings.Replace(name, "´", "", -1)  // remove forward apostrophe
	name = strings.Replace(name, "`", "", -1)  // remove backward apostrophe
	name = strings.Replace(name, "'", "", -1)  // remove regular apostrophe
	name = strings.Replace(name, " ", "-", -1) // replace space with hyphen
	name = strings.ToLower(name)
	name = strings.Replace(name, "ä", "a", -1)
	return urlBase + name
}

// Encode ID field. Might find a better strategy for this later
func getRestaurantID(name string) string {
	return url.PathEscape(strings.ToLower(name))
}

func (lhs *LHScraper) Scrape() (lunchdata.Restaurants, error) {
	csel := []string{
		"h3.title",
		"div.table-list__row",
		"span.dish-name",
		"strong",
		"div.table-list__column.table-list__column--price",
	}

	log.Infof(
		"%s: Starting scrape of %q",
		TAG,
		URL,
	)

	rs := make(lunchdata.Restaurants, 0)
	t_start := time.Now()

	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return nil, err
	}

	doc.Find(csel[0]).Each(
		func(i int, sel1 *goquery.Selection) {

			rname := strings.TrimSpace(sel1.Find("a").Text())

			r := lunchdata.NewRestaurant(
				rname,
				getRestaurantID(rname),
				getRestaurantUrl(rname),
				time.Now(),
			)

			sel1.NextFilteredUntil(csel[1], csel[0]).Each(
				func(j int, sel2 *goquery.Selection) {

					dname := strings.TrimSpace(sel2.Find(csel[2]).Find(csel[3]).Text())
					ddesc := strings.TrimSpace(strings.Replace(sel2.Find(csel[2]).Text(), dname, "", 1))
					ddesc = strings.Join(strings.Fields(ddesc), " ") // remove redundant WS inside string
					dprice := strings.TrimSpace(strings.Replace(sel2.Find(csel[4]).Text(), "kr", "", 1))

					price, err := strconv.Atoi(dprice)
					if err != nil {
						price = -1
					}

					dish := lunchdata.NewDish(
						dname,
						ddesc,
						price,
					)
					r.AddDish(*dish)
				},
			)
			rs.Add(*r)
		},
	)
	log.Infof("%s: Scrape done in %f seconds", TAG, time.Duration(time.Now().Sub(t_start)).Seconds())
	log.Debugf("%s: Parsed %d restaurants with %d dishes in total", TAG, rs.Len(), rs.NumDishes())

	return rs, nil
}
