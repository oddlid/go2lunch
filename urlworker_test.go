package urlworker

import (
	"testing"
	//"encoding/json"
	//"strings"
)

type httpbinGet struct {
	Args []struct{} `json:"args"`
}

var tests = []Request{
	{"http://httpbin.org/ip", "LHMS"},
	{"http://httpbin.org/ip", "Nemos"},
}

func TestSeed(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	in := Seed(done, tests...)
	lhms := <-in
	if lhms.Tag != tests[0].Tag {
		t.Errorf("Expected %q, got %q", tests[0].Tag, lhms.Tag)
	}
	nemos := <-in
	if nemos.Tag != tests[1].Tag {
		t.Errorf("Expected: %q, got %q", tests[1].Tag, nemos.Tag)
	}
}

// Call httpbinorg/get?somekey=somevalue and compare
// Requires JSON parsing. Too tired right now for that
func TestInitWithWorkers(t *testing.T) {
}
