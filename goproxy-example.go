package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Set("X-GoProxy", "yxorPoG-X")
			return r, nil
		})

	go runEchoServer()

	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func runEchoServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := httputil.DumpRequest(r, false)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "Request Data: %s", string(data))
	})
	log.Fatal(http.ListenAndServe(":9090", mux))
}
