package server

// import (
// 	"strings"

// 	"github.com/oddlid/go2lunch/lunchdata"
// )

// var lldata = `
// {
// 	"gtag": "UA-126840341-2",
//   "countries": {
//     "se": {
//       "country_name": "Sweden",
//       "country_id": "se",
//       "cities": {
//         "gbg": {
//           "city_name": "Gothenburg",
//           "city_id": "gbg",
//           "sites": {
//             "lindholmen": {
//               "site_name": "Lindholmen",
//               "site_id": "lindholmen",
//               "site_comment": "Gruvan",
//               "restaurants": {}
//             }
//           }
//         }
//       }
//     }
//   }
// }
// `

// func getLunchList() *lunchdata.LunchList {
// 	if _lunchList == nil {
// 		ll, err := lunchdata.LunchListFromJSON(strings.NewReader(lldata))
// 		if err != nil {
// 			// Ok to panic here, because this should not be possible to get wrong
// 			panic(err)
// 		}
// 		// If we load the LunchList from JSON, and we only have a top-level
// 		// Gtag, we need to propagate it now after load
// 		if ll.Gtag != "" {
// 			ll.PropagateGtag(ll.Gtag)
// 		}
// 		_lunchList = ll
// 	}
// 	return _lunchList
// }
