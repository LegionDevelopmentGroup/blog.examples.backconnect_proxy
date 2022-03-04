package main

import (
	"context"
	"github.com/elazarl/goproxy"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Rotater struct {
	availableIPs []string
	currentIndex int
	m            *sync.Mutex
}

func (r *Rotater) nextIP() string {
	r.m.Lock()
	defer r.m.Unlock()
	if r.currentIndex >= len(r.availableIPs) {
		r.currentIndex = 0
	}
	n := r.availableIPs[r.currentIndex]
	r.currentIndex += 1
	return n
}

func (r *Rotater) CustomDialer(ctx context.Context, network string, addr string) (net.Conn, error) {
	altIP := r.nextIP()
	ipAddress := net.ParseIP(altIP)
	d := net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP: ipAddress,
		},
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return d.Dial(network, addr)
}

type TransportWrapper struct {
	transport *http.Transport
}

func (t *TransportWrapper) RoundTrip(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Response, error) {
	return t.transport.RoundTrip(req)
}

func WrapTransport(t *http.Transport) *TransportWrapper {
	return &TransportWrapper{t}
}

func main() {
	rotater := &Rotater{
		availableIPs: []string{
			"x.x.x.1",
			"x.x.x.2",
		},
		m: &sync.Mutex{},
	}

	proxy := goproxy.NewProxyHttpServer()

	// For debugging
	proxy.Verbose = true

	proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		// This handles incoming https Requests, (they use CONNECT) to forward the encrypted request
		// For the basic version we reject with goproxy.RejectConnect
		// To handle https proxy requests we'd return gorpoxy.MitmConnect

		// Note there is a bug in elazarl/goproxy they fail to correctly close the connections (they use defer close in loops). This causes theye proxy (while mitm) where a connection may be re-used and therefore the IP used for a request is not rotated.
		// It does generally work for the simple IP retrieval example, however should be addressed before use in production scraping.
		return goproxy.MitmConnect, host
	})

	// An alternative would be to reject all connect requests.	
	/*proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return goproxy.RejectConnect, host
	})*/

	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

		ctx.RoundTripper = WrapTransport(
			&http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				DialContext:           rotater.CustomDialer,
				MaxIdleConns:          1,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     true,
			},
		)

		return r, nil
	})

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
