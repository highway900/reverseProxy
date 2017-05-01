package main

import (
	"time"
	"net/http"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
  	w.Write([]byte("The time is: " + tm))
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", 8001), nil)
	if err != nil {
		panic(err)
	}
}