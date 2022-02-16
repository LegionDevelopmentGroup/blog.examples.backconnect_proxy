package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func CustomDialer(ctx context.Context, network string, addr string) (net.Conn, error) {
	altIP := "1.1.1.1" // Custom IP
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

func main() {
	url := "http://ip-api.com/json"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	transport := &http.Transport{
		Proxy:       http.ProxyFromEnvironment,
		DialContext: CustomDialer,
	}

	resp, err := transport.RoundTrip(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
