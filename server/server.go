package main

import (
	"encoding/json"
	"fmt"
	htmpl "html/template"
	"net/http"
	"strings"
	ttmpl "text/template"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/oddlid/go2lunch/lunchdata"
	"github.com/rs/zerolog/log"
)

const (
	headerKeyCT     = "Content-Type"
	headerKeyAccept = "Accept"
	headerValJSON   = "application/json; charset=UTF-8"
)

var (
	htmlTemplates *htmpl.Template
	textTemplates *ttmpl.Template
	htmlFiles     = []string{"lunchlist.html", "country.html", "city.html", "site.html", "default.html"} // virtual files
	textFiles     = []string{"lunchlist.txt", "country.txt", "city.txt", "site.txt"}                     // virtual files
	urlIds        = []string{"country_id", "city_id", "site_id", "restaurant_id"}
)

//type WebUser struct {
//    Username string `json:"username"`
//    Password string `json:"password"`
//}

//type JwtToken struct {
//    Token string `json:"token"`
//}

//type Exception struct {
//    Message string `json:"message"`
//}

func initTmpl() error {
	log.Debug().Msg("Looking for template folder...")
	tBox, err := rice.FindBox("tmpl")
	if err != nil {
		return err
	}

	htmplStr, err := tBox.String("allhtml.tpl")
	if err != nil {
		return err
	}
	htmlTemplates, err = htmpl.New("html").Parse(htmplStr)
	if err != nil {
		return err
	}

	ttmplStr, err := tBox.String("alltext.tpl")
	if err != nil {
		return err
	}
	textTemplates, err = ttmpl.New("text").Parse(ttmplStr)
	if err != nil {
		return err
	}

	log.Debug().Msg("Templates loaded and parsed")
	return nil
}

func setupRouter() (pubR, admR *mux.Router) {
	const (
		mGET   = "GET"
		mPOST  = "POST"
		mDEL   = "DELETE"
		ppStat = "/static/"
		ppJSON = "/json/"
		ppHTML = "/html/"
		ppText = "/text/"
		ppUpd  = "/update/"
		ppAdm  = "/adm/"
		ppAdd  = "/add"
		ppDel  = "/del"
		slash  = "/"
	)

	ppath := func(upto int) string {
		var b strings.Builder
		for i := 0; i <= upto; i++ {
			fmt.Fprintf(&b, "/{%s}", urlIds[i])
		}
		b.WriteString(slash) // add final slash
		return b.String()
	}

	box := rice.MustFindBox("static")

	pubR = mux.NewRouter()
	pubR.PathPrefix(ppStat).Handler(http.StripPrefix(ppStat, http.FileServer(box.HTTPBox())))
	pubR.HandleFunc(slash, htmlIndexHandler).Methods(mGET)

	// json/api GET routes
	jsubr := pubR.PathPrefix(ppJSON).Subrouter().StrictSlash(true)
	jsubr.HandleFunc(slash, jsonAPIHandler).Methods(mGET)
	jsubr.HandleFunc(ppath(0), jsonAPIHandler).Methods(mGET)
	jsubr.HandleFunc(ppath(1), jsonAPIHandler).Methods(mGET)
	jsubr.HandleFunc(ppath(2), jsonAPIHandler).Methods(mGET)
	//jsubr.HandleFunc(ppath(3), jsonApiHandler).Methods(MGET)

	// regular HTML GET routes
	hsubr := pubR.PathPrefix(ppHTML).Subrouter().StrictSlash(true)
	hsubr.HandleFunc(slash, htmlTmplHandler).Methods(mGET)
	hsubr.HandleFunc(ppath(0), htmlTmplHandler).Methods(mGET)
	hsubr.HandleFunc(ppath(1), htmlTmplHandler).Methods(mGET)
	hsubr.HandleFunc(ppath(2), htmlTmplHandler).Methods(mGET)

	// text/plain GET routes
	tsubr := pubR.PathPrefix(ppText).Subrouter().StrictSlash(true)
	tsubr.HandleFunc(slash, textTmplHandler).Methods(mGET)
	tsubr.HandleFunc(ppath(0), textTmplHandler).Methods(mGET)
	tsubr.HandleFunc(ppath(1), textTmplHandler).Methods(mGET)
	tsubr.HandleFunc(ppath(2), textTmplHandler).Methods(mGET)

	// 2019-10-29 16:50: Disabling registration of update routes temporarily,
	// so I can build a docker image and try to run this newer version in prod
	// for a while, without risking unauthorized updates.

	// POST routes, for receiving updates in public router
	//usubr := pubR.PathPrefix(ppUpd).Subrouter() // .StrictSlash(true) // seemed to not be desirable
	//usubr.HandleFunc(slash, setGtagMW(updateHandler)).Methods(MPOST)
	//usubr.HandleFunc(ppath(2), setGtagMW(updateHandler)).Methods(MPOST)

	// Redirects just to not break old urls
	pubR.HandleFunc("/lindholmen.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "html/se/gbg/lindholmen", http.StatusMovedPermanently)
	}).Methods(mGET)
	pubR.HandleFunc("/lindholmen", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "html/se/gbg/lindholmen", http.StatusMovedPermanently)
	}).Methods(mGET)
	pubR.HandleFunc("/lindholmen.json", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "json/se/gbg/lindholmen", http.StatusMovedPermanently)
	}).Methods(mGET)
	pubR.HandleFunc("/lindholmen.txt", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "text/se/gbg/lindholmen", http.StatusMovedPermanently)
	}).Methods(mGET)

	// admin POST interface
	admR = mux.NewRouter()
	admSubr := admR.PathPrefix(ppAdm).Subrouter().StrictSlash(false)
	admSubr.HandleFunc(ppAdd+ppath(0), logInventoryMW(setGtagMW(addCountryHandler))).Methods(mPOST)
	admSubr.HandleFunc(ppDel+ppath(0), logInventoryMW(delCountryHandler)).Methods(mDEL)
	admSubr.HandleFunc(ppAdd+ppath(1), logInventoryMW(setGtagMW(addCityHandler))).Methods(mPOST)
	admSubr.HandleFunc(ppDel+ppath(1), logInventoryMW(delCityHandler)).Methods(mDEL)
	admSubr.HandleFunc(ppAdd+ppath(2), logInventoryMW(setGtagMW(addSiteHandler))).Methods(mPOST)
	admSubr.HandleFunc(ppDel+ppath(2), logInventoryMW(delSiteHandler)).Methods(mDEL)

	//	r.HandleFunc("/getauth", createTokenHandler).Methods(MPOST)
	//	r.HandleFunc("/testauth", authMiddleWare(testCreatedTokenHandler)).Methods(MGET)

	return
}

func htmlIndexHandler(w http.ResponseWriter, _ *http.Request) {
	if err := htmlTemplates.ExecuteTemplate(w, htmlFiles[4], getLunchList()); err != nil {
		log.Error().Err(err).Str("func", "htmlIndexHandler").Send()
	}
}

func genericTmplHandler(
	w http.ResponseWriter,
	r *http.Request,
	f func(tmplIndex int, wr http.ResponseWriter, obj interface{})) {
	vars := mux.Vars(r)
	index := 0
	var country *lunchdata.Country
	var city *lunchdata.City
	var site *lunchdata.Site

	countryID, found := vars[urlIds[index]]
	if !found {
		// show list of countries and return
		f(index, w, getLunchList())
		return
	}
	country = getLunchList().GetCountryByID(countryID)
	if country == nil {
		http.NotFound(w, r) // 404
		return
	}
	index++
	cityID, found := vars[urlIds[index]]
	if !found {
		// show list of cities below above country and return
		f(index, w, country)
		return
	}
	city = country.GetCityByID(cityID)
	if city == nil {
		http.NotFound(w, r) // 404
		return
	}
	index++
	siteID, found := vars[urlIds[index]]
	if !found {
		// show list of sites below above city and return
		f(index, w, city)
		return
	}
	site = city.GetSiteByID(siteID)
	if site == nil {
		http.NotFound(w, r) // 404
		return
	}

	// at this point, we might need to do a scrape

	index++
	f(index, w, site)
}

func textTmplHandler(w http.ResponseWriter, r *http.Request) {
	genericTmplHandler(w, r, func(tmplIdx int, wr http.ResponseWriter, obj interface{}) {
		if err := textTemplates.ExecuteTemplate(wr, textFiles[tmplIdx], obj); err != nil {
			log.Error().Err(err).Str("func", "textTmplHandler").Send()
		}
	})
}

func htmlTmplHandler(w http.ResponseWriter, r *http.Request) {
	genericTmplHandler(w, r, func(tmplIdx int, wr http.ResponseWriter, obj interface{}) {
		if err := htmlTemplates.ExecuteTemplate(wr, htmlFiles[tmplIdx], obj); err != nil {
			log.Error().Err(err).Str("func", "htmlTmplHandler").Send()
		}
	})
}

func jsonAPIHandler(w http.ResponseWriter, r *http.Request) {
	// I think maybe it could be a good idea to add gzip to this reply
	genericTmplHandler(w, r, func(_ int, wr http.ResponseWriter, obj interface{}) {
		wr.Header().Set(headerKeyCT, headerValJSON)
		if err := json.NewEncoder(wr).Encode(obj); err != nil {
			log.Error().Err(err).Str("func", "jsonApiHandler").Send()
		}
	})
}

// TODO: implement this in a way that makes it easy to render html to stdout etc...
//func dumpHtml(w io.Writer) {
//}

func addCountryHandler(w http.ResponseWriter, r *http.Request) {
	// sLog.Debug("Entering addCountryHandler...")

	w.Header().Set(headerKeyAccept, headerValJSON)

	//	vars := mux.Vars(r)
	//
	//	countryID, found := vars[urlIds[0]]
	//	if !found {
	//		http.NotFound(w, r) // 404
	//		return
	//	}
	//	country := _site.ll.GetCountryById(countryID)
	//	if country == nil {
	//		http.Error(w, "Non-existent country code", http.StatusInternalServerError)
	//		return
	//	}

	newCountry, err := lunchdata.CountryFromJSON(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getLunchList().AddCountry(newCountry)
}

func delCountryHandler(w http.ResponseWriter, r *http.Request) {
	// sLog.WithField("func", "delCountryHandler").Debug("Entering delCountryHandler...")
	w.Header().Set(headerKeyAccept, headerValJSON)

	vars := mux.Vars(r)

	countryID, found := vars[urlIds[0]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}

	if !getLunchList().HasCountry(countryID) {
		http.Error(w, "Non-existent country code", http.StatusInternalServerError)
		return
	}

	getLunchList().DeleteCountry(countryID)
}

func addCityHandler(w http.ResponseWriter, r *http.Request) {
	// acLog := sLog.WithField("func", "addCityHandler")
	// acLog.Debug("Entering addCityHandler...")
	w.Header().Set(headerKeyAccept, headerValJSON)

	vars := mux.Vars(r)

	countryID, found := vars[urlIds[0]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}
	cityID, found := vars[urlIds[1]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}

	country := getLunchList().GetCountryByID(countryID)
	if country == nil {
		http.Error(w, "Non-existent country code", http.StatusInternalServerError)
		return
	}

	if country.HasCity(cityID) {
		log.Debug().Str("cityID", cityID).Msg("City already exists, overwriting")
	}

	city, err := lunchdata.CityFromJSON(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	country.AddCity(city)
}

func delCityHandler(w http.ResponseWriter, r *http.Request) {
	// sLog.WithField("func", "delCityHandler").Debug("Entering delCityHandler...")

	vars := mux.Vars(r)

	countryID, found := vars[urlIds[0]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}
	cityID, found := vars[urlIds[1]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}

	country := getLunchList().GetCountryByID(countryID)
	if country == nil {
		http.Error(w, "Non-existent country code", http.StatusInternalServerError)
		return
	}

	if !country.HasCity(cityID) {
		http.Error(w, "Non-existent city code", http.StatusInternalServerError)
		return
	}

	country.DeleteCity(cityID)
}

func addSiteHandler(w http.ResponseWriter, r *http.Request) {
	// asLog := sLog.WithField("func", "addSiteHandler")
	// asLog.Debug("Entering addSiteHandler...")
	w.Header().Set(headerKeyAccept, headerValJSON)

	vars := mux.Vars(r)

	countryID, found := vars[urlIds[0]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}
	cityID, found := vars[urlIds[1]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}
	siteID, found := vars[urlIds[2]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}

	country := getLunchList().GetCountryByID(countryID)
	if country == nil {
		http.Error(w, "Non-existent country code", http.StatusInternalServerError)
		return
	}

	city := country.GetCityByID(cityID)
	if city == nil {
		http.Error(w, "Non-existent city code", http.StatusInternalServerError)
		return
	}

	if city.HasSite(siteID) {
		log.Debug().Str("siteID", siteID).Msg("Site already exists, overwriting")
	}

	site, err := lunchdata.SiteFromJSON(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	city.AddSite(site)

	//	sl := _site.ll.GetSiteLinkById(countryID, cityID, siteID)
	//	if sl != nil {
	//		token, err := getTokenForSiteLink(*sl)
	//		if err == nil {
	//			// We need to get a new reference to the site here
	//			site = _site.ll.GetSiteById(countryID, cityID, siteID)
	//			site.Key = token
	//			sLog.Debugf("addSiteHandler: Got key: %q", token)
	//		}
	//	} else {
	//		sLog.Debug("addSiteHandler: got no sitelink to generate key for")
	//	}
}

func delSiteHandler(w http.ResponseWriter, r *http.Request) {
	// sLog.WithField("func", "delSiteHandler").Debug("Entering delSiteHandler...")

	vars := mux.Vars(r)

	countryID, found := vars[urlIds[0]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}
	cityID, found := vars[urlIds[1]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}
	siteID, found := vars[urlIds[2]]
	if !found {
		http.NotFound(w, r) // 404
		return
	}

	country := getLunchList().GetCountryByID(countryID)
	if country == nil {
		http.Error(w, "Non-existent country code", http.StatusInternalServerError)
		return
	}

	city := country.GetCityByID(cityID)
	if city == nil {
		http.Error(w, "Non-existent city code", http.StatusInternalServerError)
		return
	}

	if !city.HasSite(siteID) {
		http.Error(w, "Non-existent site code", http.StatusInternalServerError)
		return
	}

	city.DeleteSite(siteID)
}

//func updateHandler(w http.ResponseWriter, r *http.Request) {
//	// Thoughts...
//	// We should probably accept something like this:
//	// A full site posted to a /country/city/
//	// A restaurant posted to a /country/city/site/
//	// A dish posted to a /country/city/site/restaurant/ ? No, this is one step too far
//
//	// 2019-08-07 21:02:
//	// We should only accept a json encoded instance of Restaurants that we will add to an existing Site.
//	// This since the site needs to have a key for verification that the scraper is authorized.
//	// Reading a whole site and replace it would make the key useless.
//
//	w.Header().Set(HDR_KEY_ACCEPT, HDR_VAL_JSON)
//
//	vars := mux.Vars(r)
//
//	countryID, found := vars[urlIds[0]]
//	if !found {
//		http.NotFound(w, r) // 404
//		return
//	}
//	cityID, found := vars[urlIds[1]]
//	if !found {
//		http.NotFound(w, r) // 404
//		return
//	}
//	siteID, found := vars[urlIds[2]]
//	if !found {
//		http.NotFound(w, r) // 404
//		return
//	}
//
//	site := getLunchList().GetSiteById(countryID, cityID, siteID)
//	if site == nil {
//		http.NotFound(w, r) // 404
//		return
//	}
//
//	// Get header with auth token here, and compare it to the key set for the site.
//	// If there is no key for the site, return error.
//	// If the keys don't match, return error.
//	// If the keys match, continue.
//
//	// eg:
//	//	if r.Header.Get("x-auth-token") != "admin" {
//	//		w.WriteHeader(http.StatusUnauthorized)
//	//		return
//	//	}
//
//	rs, err := lunchdata.RestaurantsFromJSON(r.Body)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	site.SetRestaurants(rs)
//
//	// return 201 created on success
//}

//func getTokenForSiteLink(sl lunchdata.SiteLink, secret string) (string, error) {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"country_name": sl.CountryName,
//		"country_id":   sl.CountryID,
//		"city_name":    sl.CityName,
//		"city_id":      sl.CityID,
//		"site_name":    sl.SiteName,
//		"site_id":      sl.SiteID,
//		"comment":      sl.Comment,
//		"url":          sl.Url,
//	})
//	tokenString, err := token.SignedString([]byte(secret))
//	if err != nil {
//		sLog.Error(err.Error())
//		return "", err
//	}
//	return tokenString, nil
//}

// JWT stuff modified from: https://www.thepolyglotdeveloper.com/2017/03/authenticate-a-golang-api-with-json-web-tokens/

// simplistic and naive first version just to get going
//func createTokenHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set(HDR_KEY_ACCEPT, HDR_VAL_JSON)
//
//	// Idea @ 2019-08-07 17:52:
//	// We could maybe create the token based on an instance of SiteLink.
//	// This would make more sense than user/pass as the sitelink would be unique for each site,
//	// which is probably what we need.
//	// The site would then have a field with the token generated from the sitelink to it and the app secret.
//	// As long as the secret is actually secret, that should be enough to prevent unauthorized updates.
//	// But, I might very well be wrong...
//
//	var wu WebUser
//	err := json.NewDecoder(r.Body).Decode(&wu)
//	if err != nil {
//		sLog.Error(err.Error())
//	}
//
//	// at this point we should validate credentials before proceeding
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"username": wu.Username,
//		"password": wu.Password,
//	})
//	tokenString, err := token.SignedString([]byte("secret"))
//	if err != nil {
//		sLog.Error(err.Error())
//	}
//
//	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
//}
//
//func testCreatedTokenHandler(w http.ResponseWriter, r *http.Request) {
//	decoded := context.Get(r, "decoded")
//	var wu WebUser
//	mapstructure.Decode(decoded.(jwt.MapClaims), &wu)
//	json.NewEncoder(w).Encode(wu)
//}
//
//func authMiddleWare(next http.HandlerFunc) http.HandlerFunc {
//	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//		authHdr := req.Header.Get("Authorization")
//		if authHdr != "" {
//			bearerToken := strings.Split(authHdr, " ")
//			if len(bearerToken) == 2 {
//				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
//					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//						return nil, fmt.Errorf("Error parsing JWT token")
//					}
//					return []byte("secret"), nil
//				})
//				if err != nil {
//					json.NewEncoder(w).Encode(Exception{Message: err.Error()})
//					return
//				}
//				if token.Valid {
//					context.Set(req, "decoded", token.Claims)
//					next(w, req)
//				} else {
//					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
//				}
//			}
//		} else {
//			json.NewEncoder(w).Encode(Exception{Message: "Authorization header required"})
//		}
//	})
//}

func setGtagMW(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next(w, req)
		getLunchList().PropagateGtag(_gtag)
	})
}

func logInventoryMW(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next(w, req)
		logInventory()
	})
}
