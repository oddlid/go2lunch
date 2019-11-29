package urlworker

// Inspiration for this was gotten from: https://blog.golang.org/pipelines

import (
	//"crypto/tls"
	//"encoding/json"
	//"io"
	"net/http"
	//"runtime"
	//"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	PKG_NAME = "urlworker"
)

//type Request struct {
//	Url string `json:"url"`
//	Tag string `json:"tag,omitempty"`
//}

//type Requests []Request

//type Response struct {
//	Res *http.Response
//	Req Request
//	//Err error
//}

//type Responses []Response

type UrlWorker struct {
	Done       chan struct{}
	client     *HttpClient
	reqs       []*http.Request
	wg         sync.WaitGroup
	numWorkers int
}

func New(workers int) *UrlWorker {
	return &UrlWorker{
		Done:       make(chan struct{}),
		reqs:       make([]*http.Request, 0),
		client:     NewHttpClient(),
		numWorkers: workers,
	}
}

//func NewRequest(url string) *Request {
//	return &Request{
//		Url: url,
//	}
//}

func (uw *UrlWorker) AddRequest(r *http.Request) {
	uw.reqs = append(uw.reqs, r)
}

func (uw *UrlWorker) SetRequests(rs []*http.Request) {
	uw.reqs = rs
}

func (uw *UrlWorker) GetClient() *HttpClient {
	return uw.client
}

func (uw *UrlWorker) GetRequests() []*http.Request {
	return uw.reqs
}

//func (reqs *Requests) NewFromJSON(rdr io.Reader) error {
//	dec := json.NewDecoder(rdr)
//	var r Requests
//	err := dec.Decode(&r)
//	if err != nil && err != io.EOF {
//		return err
//	}
//	*reqs = r
//	return nil
//}

//func (uw *UrlWorker) seed(done <-chan struct{}, reqs ...Request) <-chan Request {
//	out := make(chan Request, len(reqs))
//	defer close(out)
//	for _, req := range reqs {
//		select {
//		case out <- req:
//		case <-done:
//		}
//	}
//	return out
//}

func (uw *UrlWorker) seed() <-chan *http.Request {
	out := make(chan *http.Request, len(uw.reqs))
	defer close(out)
	for _, req := range uw.reqs {
		select {
		case out <- req:
		case <-uw.Done:
		}
	}
	return out
}

// should rewrite this to use Get from httpclient.go
//func Get(done <-chan struct{}, in <-chan Request) <-chan Response {
//	out := make(chan Response)
//
//	go func() {
//		defer close(out)
//
//		for req := range in {
//			hreq, err := http.NewRequest("GET", req.Url, nil)
//			if err != nil {
//				log.Errorf("Unable to setup request for URL %q: %s", req.Url, err)
//				return
//			}
//
//			hreq.Header.Set("User-Agent", UA)
//			tr := &http.Transport{DisableKeepAlives: true} // we're not reusing the connection, so don't let it hang open
//
//			if strings.Index(req.Url, "https") >= 0 {
//				// Not checking certs right now, but should do that when shit gets more real
//				tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
//			}
//
//			client := &http.Client{Transport: tr}
//			res, err := client.Do(hreq)
//			// I could maybe return the err in the Response struct instead of bailing here
//			if err != nil {
//				log.Errorf("Error fetching URL: %s", err)
//				return
//			}
//
//			select {
//			case out <- Response{res, req}:
//			case <-done:
//				return
//			}
//		}
//	}()
//
//	return out
//}

//func (uw *UrlWorker) get(done <-chan struct{}, in <-chan Request) <-chan Response {
//	out := make(chan Response)
//
//	go func() {
//		defer close(out)
//		//hc := NewHttpClient()
//		// configure the client here if needed
//		for req := range in {
//			//res, err := hc.Get(req.Url)
//			res, err := uw.client.Get(req.Url)
//			if err != nil {
//				log.Errorf("Error fetching URL: %s", err)
//				return
//			}
//
//			select {
//			case out <- Response{res, req}:
//			case <-done:
//				return
//			}
//		}
//	}()
//
//	return out
//}

func (uw *UrlWorker) get(in <-chan *http.Request) <-chan *http.Response {
	out := make(chan *http.Response)

	go func() {
		defer close(out)
		for req := range in {
			res, err := uw.client.Get(req)
			if err != nil {
				log.WithFields(log.Fields{
					"pkg":    PKG_NAME,
					"func":   "get",
					"URL":    req.URL,
					"ErrMSG": err.Error(),
				}).Error("Fetch error")
				return
			}
			select {
			case out <- res:
			case <-uw.Done:
				return
			}
		}
	}()

	return out
}

//func (uw *UrlWorker) merge(done <-chan struct{}, cs ...<-chan Response) <-chan Response {
//	var wg sync.WaitGroup
//	out := make(chan Response)
//
//	output := func(c <-chan Response) {
//		defer wg.Done()
//		for v := range c {
//			select {
//			case out <- v:
//			case <-done:
//				return
//			}
//		}
//	}
//	wg.Add(len(cs))
//
//	for _, c := range cs {
//		go output(c)
//	}
//
//	go func() {
//		wg.Wait()
//		close(out)
//	}()
//	return out
//}

func (uw *UrlWorker) merge(cs ...<-chan *http.Response) <-chan *http.Response {
	//var wg sync.WaitGroup
	out := make(chan *http.Response)

	output := func(c <-chan *http.Response) {
		defer uw.wg.Done()
		for v := range c {
			select {
			case out <- v:
			case <-uw.Done:
				return
			}
		}
	}
	uw.wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		uw.wg.Wait()
		close(out)
	}()

	return out
}

//func (uw *UrlWorker) InitWithWorkers(workers int, done <-chan struct{}, reqs ...Request) <-chan Response {
//	in := uw.seed(done, reqs...)
//	minions := make([]<-chan Response, workers)
//	for i := 0; i < workers; i++ {
//		minions[i] = uw.get(done, in)
//	}
//	return uw.merge(done, minions...)
//}

//func (uw *UrlWorker) InitWithWorkers(workers int) <-chan Response {
//	in := uw.seed()
//	minions := make([]<-chan Response, workers)
//	for i := 0; i < workers; i++ {
//		minions[i] = uw.get(in)
//	}
//	return uw.merge(minions...)
//}

func (uw *UrlWorker) Get() <-chan *http.Response {
	in := uw.seed()
	workers := make([]<-chan *http.Response, uw.numWorkers)
	for i := 0; i < uw.numWorkers; i++ {
		workers[i] = uw.get(in)
	}
	return uw.merge(workers...)
}


//func (uw *UrlWorker) Init(done <-chan struct{}, reqs ...Request) <-chan Response {
//	numCPU := runtime.NumCPU()
//	return uw.InitWithWorkers(numCPU, done, reqs...)
//}

//func (uw *UrlWorker) Init() <-chan Response {
//	numCPU := runtime.NumCPU()
//	return uw.InitWithWorkers(numCPU)
//}
