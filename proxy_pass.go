package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	app_remote, err := url.Parse("http://localhost:6543")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(app_remote)
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		r.URL.Path = "/"
		p.ServeHTTP(w, r)
	}
}
