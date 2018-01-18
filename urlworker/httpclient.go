package urlworker


import (
	//"crypto/tls"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"time"
	"io"
	"io/ioutil"
	"bytes"
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
	return client.SetTimeout(DEF_TMOUT)
}

func (c *HttpClient) SkipCertCheck(skip bool) *HttpClient {
	c.Transport.TLSClientConfig.InsecureSkipVerify = skip
	return c
}

func (c *HttpClient) SetKeepAlives(val bool) *HttpClient {
	c.Transport.DisableKeepAlives = val
	return c
}

func (c *HttpClient) SetTimeout(sec float64) *HttpClient {
	c.Client.Timeout = time.Second * time.Duration(sec)
	return c
}

func (c *HttpClient) Setup(ua string, timeout float64, keepAlive, skipcertcheck bool) *HttpClient {
	if ua != "" {
		c.UA = ua
	}
	return c.SetTimeout(timeout).SetKeepAlives(keepAlive).SkipCertCheck(skipcertcheck)
}

func (c *HttpClient) Get(url) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UA)

	return c.Do(req)
}

// GetResponseBody reads http.Response.Body, closes it, and returns the contents as a byte slice
func GetResponseBody(res *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetResponseReader(res *http.Response) bytes.Reader {
	b, err := GetResponseBody(res)
	if err != nil {
		log.Error(err)
		return nil
	}
	return bytes.NewReader(b)
}

func CloseResponse(res *http.Response) error {
	// See https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
	io.Copy(ioutil.Discard, res.Body)
	return res.Body.Close()
}

