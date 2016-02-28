package urlworker

// Inspiration for this was gotten from: https://blog.golang.org/pipelines

import (
	"crypto/tls"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

const UA string = "go2lunch.urlworker/0.1"

type Request struct {
	Url string `json:"url"`
	Tag string `json:"name,omitempty"`
}

type Requests []Request

type Response struct {
	Res *http.Response
	Req Request
	//Err error
}

//type Responses []Response

func (reqs *Requests) NewFromJSON(rdr io.Reader) error {
	dec := json.NewDecoder(rdr)
	var r Requests
	err := dec.Decode(&r)
	if err != nil && err != io.EOF {
		return err
	}
	*reqs = r
	return nil
}

func Seed(done <-chan struct{}, reqs ...Request) <-chan Request {
	out := make(chan Request, len(reqs))
	defer close(out)
	for _, req := range reqs {
		select {
		case out <- req:
		case <-done:
		}
	}
	return out
}

func Get(done <-chan struct{}, in <-chan Request) <-chan Response {
	out := make(chan Response)

	go func() {
		defer close(out)

		for req := range in {

			hreq, err := http.NewRequest("GET", req.Url, nil)
			if err != nil {
				log.Errorf("Unable to setup request for URL %q: %s", req.Url, err)
				return
			}

			hreq.Header.Set("User-Agent", UA)
			tr := &http.Transport{DisableKeepAlives: true} // we're not reusing the connection, so don't let it hang open

			if strings.Index(req.Url, "https") >= 0 {
				// Not checking certs right now, but should do that when shit gets more real
				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			}

			client := &http.Client{Transport: tr}
			res, err := client.Do(hreq)
			// I could maybe return the err in the Response struct instead of bailing here
			if err != nil {
				log.Errorf("Error fetching URL: %s", err)
				return
			}

			select {
			case out <- Response{res, req}:
			case <-done:
				return
			}
		}
	}()

	return out
}

func Merge(done <-chan struct{}, cs ...<-chan Response) <-chan Response {
	var wg sync.WaitGroup
	out := make(chan Response)

	output := func(c <-chan Response) {
		defer wg.Done()
		for v := range c {
			select {
			case out <- v:
			case <-done:
				return
			}
		}
		//wg.Done()
	}
	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func InitWithWorkers(workers int, done <-chan struct{}, reqs ...Request) <-chan Response {
	in := Seed(done, reqs...)
	minions := make([]<-chan Response, workers)
	for i := 0; i < workers; i++ {
		minions[i] = Get(done, in)
	}
	return Merge(done, minions...)
}

func Init(done <-chan struct{}, reqs ...Request) <-chan Response {
	numCPU := runtime.NumCPU()
	return InitWithWorkers(numCPU, done, reqs...)
}
