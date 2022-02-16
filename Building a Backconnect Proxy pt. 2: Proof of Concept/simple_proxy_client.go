package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {
	address := "http://ip-api.com/json"

	// proxyURL, err := url.Parse("<proxy-address>")
	// if err != nil {
	// 	panic(err)
	// }

	// Alternatively client := http.DefaultClient
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment, // Or specify http.ProxyURL(proxyURL)
		},
	}

	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", dump)
}
