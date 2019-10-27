package main

import (
	"strings"
	"github.com/oddlid/go2lunch/lunchdata"
)

var lldata string = `
{
  "countries": {
    "se": {
      "country_name": "Sweden",
      "country_id": "se",
      "cities": {
        "gbg": {
          "city_name": "Gothenburg",
          "city_id": "gbg",
          "sites": {
            "lindholmen": {
              "site_name": "Lindholmen",
              "site_id": "lindholmen",
              "site_comment": "Gruvan",
              "restaurants": {}
            }
          }
        }
      }
    }
  }
}
`

func getLunchList() *lunchdata.LunchList {
	if nil == _lunchList {
		//_lunchList = lunchdata.NewLunchList()
		ll, err := lunchdata.LunchListFromJSON(strings.NewReader(lldata))
		if err != nil {
			panic(err)
		}
		_lunchList = ll
	}
	return _lunchList
}


