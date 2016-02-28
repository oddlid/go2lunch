package main

import (
	"bufio"
	//"github.com/oddlid/go2lunch/site"
	"github.com/oddlid/go2lunch/urlworker"
	"testing"
	"os"
)

func TestJSONLoad(t *testing.T) {
	reqs := &urlworker.Requests{}
	f, err := os.Open("lh_restaurant_urls.json")
	if err != nil {
		t.Error(err)
	}
	r := bufio.NewReader(f)
	err = reqs.NewFromJSON(r)
	if err != nil {
		t.Error(err)
	}
	t.Logf("URLS: %#v\n", reqs)
}
