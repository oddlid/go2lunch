package urlworker


import (
	//"crypto/tls"
	//log "github.com/Sirupsen/logrus"
	"net/http"
	"time"
)

const (
	DEF_UA    string  = "oddware.urlworker/0.1"
	DEF_TMOUT float64 = 3.0
)

type HttpClient struct {
	Client *http.Client
	UA     string
}

func NewHttpClient() *HttpClient {
	client := &HttpClient{
		Client: http.DefaulClient,
		UA:     DEF_UA,
	}

	return client
}

func (c *HttpClient) SkipCertCheck(skip bool) {
	c.Transport.TLSClientConfig.InsecureSkipVerify = skip
}

func (c *HttpClient) SetKeepAlives(val bool) {
	c.Transport.DisableKeepAlives = val
}

func (c *HttpClient) SetTimeout(sec float64) {
	c.Client.Timeout = time.Second * time.Duration(sec)
}

func (c *HttpClient) Get(url) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UA)

	return c.Do(req)
}
