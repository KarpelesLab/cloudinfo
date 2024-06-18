package cloudinfo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type reqKey struct {
	method string
	url    string
}

// cachedHttp is an object that can perform only GET HTTP queries and will always perform a given request only once.
//
// All requests have a timeout of 2 seconds.
type cachedHttp struct {
	reqs map[reqKey]*cachedRequest
	lk   sync.Mutex
}

type cachedRequest struct {
	req  *http.Request
	body []byte         // response body
	resp *http.Response // response
	err  error          // error
	run  sync.Once
}

func (c *cachedRequest) key() reqKey {
	return reqKey{method: c.req.Method, url: c.req.URL.String()}
}

func newCachedHttp() *cachedHttp {
	return &cachedHttp{
		reqs: make(map[reqKey]*cachedRequest),
	}
}

// getReq returns or create the given request
func (c *cachedHttp) getReq(req *cachedRequest) *cachedRequest {
	c.lk.Lock()
	defer c.lk.Unlock()

	key := req.key()
	if oldreq, ok := c.reqs[key]; ok {
		return oldreq
	}
	return req
}

// Get performs the HTTP request and return the body, response and any error
func (c *cachedHttp) Get(u string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	return c.Do(req)
}

// GetWithHeaders performs the HTTP request with the specified headers. The headers aren't used in the caching key
func (c *cachedHttp) GetWithHeaders(u string, hdrs map[string]string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range hdrs {
		req.Header.Set(k, v)
	}
	return c.Do(req)
}

// PutWithHeaders performs a HTTP PUT
func (c *cachedHttp) PutWithHeaders(u string, hdrs map[string]string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("PUT", u, nil)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range hdrs {
		req.Header.Set(k, v)
	}
	return c.Do(req)
}

// GetWithHost performs the HTTP request while connecting to the specified host
func (c *cachedHttp) GetWithHost(u, host string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	req.URL.Host = host
	return c.Do(req)
}

func (c *cachedHttp) Do(req *http.Request) ([]byte, *http.Response, error) {
	return c.getReq(&cachedRequest{req: req}).Do()
}

func (r *cachedRequest) Do() ([]byte, *http.Response, error) {
	r.run.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		resp, err := http.DefaultClient.Do(r.req.WithContext(ctx))
		if err != nil {
			r.err = err
			return
		}
		defer resp.Body.Close()

		dat, err := io.ReadAll(&io.LimitedReader{R: resp.Body, N: 1024 * 1024})
		if err != nil {
			r.err = err
			return
		}

		r.resp = resp
		r.body = dat

		// set r.err if status code isn't 200
		if resp.StatusCode != 200 {
			r.err = fmt.Errorf("HTTP Status: %s", resp.Status)
		}
	})

	return r.body, r.resp, r.err
}
