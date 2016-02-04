package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func makeProxyHandler(url string) httputil.NewSingleHostReverseProxy {
	app_remote, err := url.Parse(url)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(app_remote)
}

func main() {
	r := mux.NewRouter()

	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("filename not specified")
	}
	dirname := flag.Args()[0]

	fs := http.FileServer(http.Dir("/home/mattb/Dev/work/biarri/wond/frontend/app/"))

	log.Println("Serving", dirname)

	proxy1 := makeProxyHandler("http://localhost:6543")
	proxy2 := makeProxyHandler("http://localhost:6543/static")

	// To use my router the behavior is different to `http.Handle`
	// matching will only occur on a fixed path or using an expression to
	// handle depth
	// Below I use `PathPrefix` with a `Handler` to fix this.
	r.HandleFunc("/api/{_dummy:.*}/", handler(proxy1))
	r.HandleFunc("/static/{_dummy:.*}/", handler(proxy2))
	r.PathPrefix("/").Handler(fs)
	http.Handle("/", r)

	log.Println("Listening...on port :8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		p.ServeHTTP(w, r)
	}
}
