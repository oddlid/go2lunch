package urlworker

import (
	//"crypto/tls"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	HDR_KEY_UA = "User-Agent"
	DEF_UA     = "oddware.urlworker/0.2"
	DEF_TMOUT  = 3.0
)

type HttpClient struct {
	Client    *http.Client
	UA        string
	transport *http.Transport
}

func NewHttpClient() *HttpClient {
	client := &HttpClient{
		Client:    http.DefaultClient,
		UA:        DEF_UA,
		transport: http.DefaultTransport.(*http.Transport),
	}
	return client.SetTimeout(DEF_TMOUT)
}

func (c *HttpClient) SkipCertCheck(skip bool) *HttpClient {
	//c.Client.Transport.TLSClientConfig.InsecureSkipVerify = skip
	c.transport.TLSClientConfig.InsecureSkipVerify = skip
	c.Client.Transport = c.transport
	return c
}

func (c *HttpClient) SetKeepAlives(val bool) *HttpClient {
	//c.Client.Transport.DisableKeepAlives = val
	c.transport.DisableKeepAlives = val
	c.Client.Transport = c.transport
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

//func (c *HttpClient) Get(url string) (*http.Response, error) {
//	req, err := http.NewRequest(http.MethodGet, url, nil)
//	if err != nil {
//		return nil, err
//	}
//	req.Header.Set(HDR_KEY_UA, c.UA)
//
//	return c.Client.Do(req)
//}

func (c *HttpClient) Get(req *http.Request) (*http.Response, error) {
	req.Header.Set(HDR_KEY_UA, c.UA)
	return c.Client.Do(req)
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

func GetResponseReader(res *http.Response) *bytes.Reader {
	b, err := GetResponseBody(res)
	if err != nil {
		log.Errorf("%s: %s", PKG_NAME, err)
		return nil
	}
	return bytes.NewReader(b)
}

func CloseResponse(res *http.Response) error {
	// See https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
	io.Copy(ioutil.Discard, res.Body)
	return res.Body.Close()
}
