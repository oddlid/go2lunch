package main

import (
	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/go2lunch/site"
	"strings"
	"time"
)

func scrape(url string) (site.Restaurants, error) {
	csel := []string{
		"h3.title",
		"div.table-list__row",
		"span.dish-name",
		"strong",
		"div.table-list__column.table-list__column--price",
	}
	var num_restaurants int
	var num_dishes int

	rs := site.Restaurants{}

	t_start := time.Now()
	log.Infof("Starting scrape of %q @ %s", url, t_start.Format(time.RFC3339))
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return rs, err
	}

	doc.Find(csel[0]).Each(func(i int, sel1 *goquery.Selection) {
		rname := sel1.Find("a").Text()
		log.Debugf("Found restaurant: %q", rname)

		r := &site.Restaurant{Name: rname, Parsed: time.Now(), Url: url}
		num_restaurants++

		sel1.NextFilteredUntil(csel[1], csel[0]).Each(func(j int, sel2 *goquery.Selection) {
			dname := strings.TrimSpace(sel2.Find(csel[2]).Find(csel[3]).Text())
			ddesc := strings.TrimSpace(strings.Replace(sel2.Find(csel[2]).Text(), dname, "", 1))
			dprice := strings.TrimSpace(strings.Replace(sel2.Find(csel[4]).Text(), "kr", "", 1))
			r.Add(site.Dish{Name: dname, Desc: ddesc, Price: dprice})
			num_dishes++

			log.Debugf("Found dish: %q", dname)
		})

		rs = append(rs, *r)
	})
	log.Infof("Scrape done in %f seconds", time.Duration(time.Now().Sub(t_start)).Seconds())
	log.Infof("Parsed %d restaurants with %d dishes in total", num_restaurants, num_dishes)

	return rs, nil
}

func update() error {
	rs, err := scrape(SRC_URL)
	if err != nil {
		return err
	}
	_site.Lock()
	_site.s.Restaurants = rs
	_site.Unlock()
	return nil
}
