package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httputil"
)

func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("/home/mattb/Dev/work/biarri/wond/frontend/app/"))

	r.PathPrefix("/").Handler(fs)

	http.Handle("/", r)

	http.ListenAndServe(":8123", nil)
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//		log.Println(r.URL)
		p.ServeHTTP(w, r)
	}
}
