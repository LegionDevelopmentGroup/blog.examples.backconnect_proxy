package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./testclient <target-site> <proxy-address>")
		return
	}

	address := os.Args[1]
	proxyAddress := os.Args[2]

	proxyURL, err := url.Parse(proxyAddress)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		panic(err)
	}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("Request Took: ", time.Now().Sub(start))

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}


	fmt.Printf("%q\n", dump)
}
