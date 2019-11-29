package engine

import (
	"testing"
	"time"
	log "github.com/sirupsen/logrus"
	"github.com/oddlid/go2lunch/site"
)

var (
	rpcclient = Setup()
	testSite = &site.Site{Name: "Lindholmen", ID: "se/gbg/lindholmen"}
)

func Setup() *RPCClient {
	//log.SetLevel(log.DebugLevel)
	go ListenAndServe(DEFAULT_DSN_PORT)
	time.Sleep(500 * time.Millisecond) // without this, the client fails to connect
	rc, err := NewRPCClient(DEFAULT_DSN_HOST + DEFAULT_DSN_PORT, time.Millisecond * 500)
	if err != nil {
		log.Fatal(err)
	}
	return rc
}

func TestColdGet(t *testing.T) {
	_, err := rpcclient.Get(testSite.ID)
	if err == nil {
		t.Error("Should get error trying to get key from empty storage\n")
	}
	//if s != nil {
	//	t.Errorf("No site instances should be in the storage yet: %q\n", s.ID)
	//}
}

func TestPut(t *testing.T) {
	lhms := &site.Restaurant{
		Name: "LHMS",
		Url: "http://www.lindholmen.se/restauranger/lindholmens-matsal",
		Parsed: time.Now(),
		Dishes: []site.Dish{
			site.Dish{Name: "Meatballs", Desc: "with mashed potatoes", Price: "85"},
			site.Dish{Name: "Lasagna", Desc: "with badger", Price: "85"},
		},

	}
	testSite.Add(lhms)
	ok, err := rpcclient.Put(testSite)
	if err != nil {
		t.Error(err)
	}
	t.Logf("client Put: %v\n", ok)
}

func TestWarmGet(t *testing.T) {
	s, err := rpcclient.Get(testSite.ID)
	if err != nil {
		t.Error(err)
	}
	if s == nil {
		t.Errorf("Site should exist: %s\n", testSite.ID)
	}
	if s.ID != testSite.ID {
		t.Errorf("Expected ID %q, got %q\n", testSite.ID, s.ID)
	}
	t.Logf("client got: %#v\n", s)
}

func TestGetKeys(t *testing.T) {
	keys, err := rpcclient.GetKeys()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Keys: %#v\n", keys)
}

func TestDelete(t *testing.T) {
	_, err := rpcclient.Delete(testSite.ID)
	if err != nil {
		t.Error(err)
	}

	site, err := rpcclient.Get(testSite.ID)
	if err != nil {
		t.Logf("Try Get() after Delete(), expect not found: %s\n", err)
	}
	if site != nil && site.ID == testSite.ID {
		t.Errorf("Site ID %q should not exist\n", site.ID)
	}
}

func TestClear(t *testing.T) {
	_, err := rpcclient.Clear()
	if err != nil {
		t.Errorf("Error clearing sites: %s\n", err)
	}
	t.Log("Sites cleared\n")
}
