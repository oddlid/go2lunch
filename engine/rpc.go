package engine

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/oddlid/go2lunch/site"
	"net"
	"net/rpc"
	"sync"
	"time"
)

const (
	DEFAULT_DSN_HOST string = "127.0.0.1"
	DEFAULT_DSN_PORT string = ":10666"
)

var NotFoundError = errors.New("Site not found")

type RPC struct {
	sites map[string]site.Site
	mx *sync.RWMutex
}

type RPCClient struct {
	conn *rpc.Client
}

//type Sites []site.Site
type Keys []string

func NewRPC() *RPC {
	return &RPC{
		sites: make(map[string]site.Site),
		mx:    &sync.RWMutex{},
	}
}

func NewRPCClient(dsn string, timeout time.Duration) (*RPCClient, error) {
	c, err := net.DialTimeout("tcp", dsn, timeout)
	if err != nil {
		return nil, err
	}
	return &RPCClient{conn: rpc.NewClient(c)}, nil
}

func ListenAndServe(dsn string) {
	rpc.Register(NewRPC())
	l, err := net.Listen("tcp", dsn)
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	rpc.Accept(l)
}

func (r *RPC) GetKeys(_ bool, keys *Keys) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	k := make(Keys, 0, len(r.sites))
	for _, s := range r.sites {
		k = append(k, s.ID)
	}
	*keys = k
	return nil
}

func (r *RPC) Get(key string, site *site.Site) error {
	r.mx.RLock()
	defer r.mx.RUnlock()
	s, found := r.sites[key]
	if !found {
		return NotFoundError
	}
	*site = s
	return nil
}

func (r *RPC) Put(site *site.Site, ack *bool) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.sites[site.ID] = *site
	*ack = true
	return nil
}

func (r *RPC) Delete(key string, ack *bool) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	_, found := r.sites[key]
	if !found {
		return NotFoundError
	}

	delete(r.sites, key)
	*ack = true
	return nil
}

func (r *RPC) Clear(_ bool, ack *bool) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.sites = make(map[string]site.Site)
	*ack = true
	return nil
}

func (c *RPCClient) Get(key string) (*site.Site, error) {
	var site site.Site
	err := c.conn.Call("RPC.Get", key, &site)
	return &site, err
}

func (c *RPCClient) GetKeys() (*Keys, error) {
	var keys Keys
	err := c.conn.Call("RPC.GetKeys", true, &keys)
	return &keys, err
}

func (c *RPCClient) Put(site *site.Site) (bool, error) {
	var added bool
	err := c.conn.Call("RPC.Put", site, &added)
	return added, err
}

func (c *RPCClient) Delete(key string) (bool, error) {
	var deleted bool
	err := c.conn.Call("RPC.Delete", key, &deleted)
	return deleted, err
}

func (c *RPCClient) Clear() (bool, error) {
	var cleared bool
	err := c.conn.Call("RPC.Clear", true, &cleared)
	return cleared, err
}
