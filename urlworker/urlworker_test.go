package urlworker

import (
	"testing"
	//"encoding/json"
	//"strings"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

//type httpbinGet struct {
//	Args []struct{} `json:"args"`
//}

//var tests = []Request{
//	{"http://httpbin.org/ip", "LHMS"},
//	{"http://httpbin.org/ip", "Nemos"},
//}

//func TestSeed(t *testing.T) {
//	done := make(chan struct{})
//	defer close(done)
//	in := Seed(done, tests...)
//	lhms := <-in
//	if lhms.Tag != tests[0].Tag {
//		t.Errorf("Expected %q, got %q", tests[0].Tag, lhms.Tag)
//	}
//	nemos := <-in
//	if nemos.Tag != tests[1].Tag {
//		t.Errorf("Expected: %q, got %q", tests[1].Tag, nemos.Tag)
//	}
//}

// Call httpbinorg/get?somekey=somevalue and compare
// Requires JSON parsing. Too tired right now for that
//func TestInitWithWorkers(t *testing.T) {
//}

//func TestInit(t *testing.T) {
//}

func TestSyntax(t *testing.T) {
	New(1)
}

func TestBlah(t *testing.T) {
	urls := []string{
		"https://www.lindholmen.se/restauranger/bistrot",
		"https://www.lindholmen.se/restauranger/cuckoos-nest",
		"https://www.lindholmen.se/restauranger/dirty-dough",
		"https://www.lindholmen.se/restauranger/district-one",
		"https://www.lindholmen.se/restauranger/kooperativet",
		"https://www.lindholmen.se/restauranger/ls-kitchen",
		"https://www.lindholmen.se/restauranger/ls-resto",
		"https://www.lindholmen.se/restauranger/ls-express",
		"https://www.lindholmen.se/restauranger/lindholmens-matsal",
		"https://www.lindholmen.se/restauranger/matminnen",
		"https://www.lindholmen.se/restauranger/mimolett",
		"https://www.lindholmen.se/restauranger/pier-eleven",
		"https://www.lindholmen.se/restauranger/restaurang-gothia",
		"https://www.lindholmen.se/restauranger/restaurang-aran",
	}
	uw := New(4)
	for _, url := range urls {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			continue
		}
		uw.AddRequest(req)
	}

	for resp := range uw.Get() {
		//t.Logf("%q: Status: %q, Length: %d", resp.Request.URL, resp.Status, resp.ContentLength)
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(doc.Find("head > title").Text())
		resp.Body.Close()
	}
}

