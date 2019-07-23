package main

/*
This file is hardcoded for scraping from lindholmen.se
Should separate this functionality into a separate module/binary or something,
and have an interface for delivering results from any scraper to a running server.
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

func scrape(url string) (lunchdata.Restaurants, error) {
	const logTag string = "scrape()"
	csel := []string{
		"h3.title",
		"div.table-list__row",
		"span.dish-name",
		"strong",
		"div.table-list__column.table-list__column--price",
	}
	var num_restaurants int
	var num_dishes int

	rs := make(lunchdata.Restaurants, 0)

	t_start := time.Now()
	log.Infof("%s: Starting scrape of %q @ %s", logTag, url, t_start.Format(time.RFC3339))
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return rs, err
	}

	doc.Find(csel[0]).Each(func(i int, sel1 *goquery.Selection) {
		rname := strings.TrimSpace(sel1.Find("a").Text())

		r := lunchdata.NewRestaurant(rname, getRestaurantID(rname), getRestaurantUrl(rname), GTAG_ID, time.Now())
		//log.Debugf("%s: Got Restaurant: %#v", logTag, r)
		num_restaurants++

		sel1.NextFilteredUntil(csel[1], csel[0]).Each(func(j int, sel2 *goquery.Selection) {
			dname := strings.TrimSpace(sel2.Find(csel[2]).Find(csel[3]).Text())
			ddesc := strings.TrimSpace(strings.Replace(sel2.Find(csel[2]).Text(), dname, "", 1))
			ddesc = strings.Join(strings.Fields(ddesc), " ") // remove redundant WS inside string
			dprice := strings.TrimSpace(strings.Replace(sel2.Find(csel[4]).Text(), "kr", "", 1))
			price, err := strconv.Atoi(dprice)
			if err != nil {
				price = -1
			}
			dish := lunchdata.NewDish(dname, ddesc, GTAG_ID, price)
			//log.Debugf("%s: Got Dish: %#v", logTag, dish)
			r.AddDish(*dish)
			num_dishes++

		})

		rs.Add(*r)
	})
	log.Infof("%s: Scrape done in %f seconds", logTag, time.Duration(time.Now().Sub(t_start)).Seconds())
	log.Infof("%s: Parsed %d restaurants with %d dishes in total", logTag, num_restaurants, num_dishes)

	return rs, nil
}

func update() error {
	rs, err := scrape(_site.url)
	if err != nil {
		return err
	}
	_site.Lock()
	//_site.s.Restaurants = rs
	_site.setLHRestaurants(rs)
	_site.Unlock()
	return nil
}
