package main

import (
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type configJSON struct {
	ProxyPort int `json:"proxyPort"`
	ProxyIp string `json:"proxyIp"`
	ProxyUrl string `json:"proxyUrl"`
	ServerPort int `json:"serverPort"`
	StaticDirectory string `json:"staticDirectory"`
}

func (c *configJSON) MakeProxyServerAddress() string {
	// TODO: how to handle https if that should be a thing in this app
	return fmt.Sprintf("http://%s:%d", c.ProxyIp, c.ProxyPort)
}

func (c *configJSON) MakeProxyUrl() string {
	if c.ProxyUrl[0:1] != "/" || c.ProxyUrl[len(c.ProxyUrl)-1:len(c.ProxyUrl)] != "/" {
		log.Println("WARNING: proxy url should start and end with a / character")
	}
	return fmt.Sprintf("%s{_dummy:.*}/", c.ProxyUrl)
}

func makeProxyHandler(serverAddress string) *httputil.ReverseProxy {
	app_remote, err := url.Parse(serverAddress)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(app_remote)
}

func Init() *configJSON {
	config := &configJSON{}

	// handle config file flag if set
	var configFileFlag string
	flag.StringVar(&configFileFlag, "config", "", "location of JSON config file")
	flag.Parse()

	if configFileFlag == "" {
		// set the sane defaults
		config.ProxyPort = 8001
		config.ProxyIp = "localhost"
		config.ProxyUrl = "/api/"
		config.ServerPort = 8000
		config.StaticDirectory = "."
	} else {
		log.Printf("Using config file %s\n", configFileFlag)

		// load config from file
		configRaw, err := ioutil.ReadFile(configFileFlag)
		if err != nil {
			panic("Error reading config file")
		}

		json.Unmarshal(configRaw, config)
		log.Println(string(configRaw), config)
	}

	if _, err := os.Stat(config.StaticDirectory); os.IsNotExist(err) {
		log.Println("static directory error;", err)
		os.Exit(0)
	}

	return config
}

func main() {
	r := mux.NewRouter()

	config := Init()

	fs := http.FileServer(http.Dir(config.StaticDirectory))

	proxy := makeProxyHandler(config.MakeProxyServerAddress())
	r.HandleFunc(config.MakeProxyUrl(), reverseProxyHandler(proxy))

	// To use my router the behavior is different to `http.Handle`
	// matching will only occur on a fixed path or using an expression to
	// handle depth
	// Below I use `PathPrefix` with a `Handler` to fix this.
	r.PathPrefix("/").Handler(handlers.CombinedLoggingHandler(os.Stdout, fs))
	http.Handle("/", r)

	log.Println("Serving directory", config.StaticDirectory)
	log.Printf("http://%s:%d\t:Web server\n", "localhost", config.ServerPort)
	log.Printf("http://%s:%d\t:Application server\n", config.ProxyIp, config.ProxyPort)
	log.Printf("http://%s:%d%s\t:Reverse Proxy to application server\n", "localhost", config.ServerPort, config.ProxyUrl)

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), nil)
	if err != nil {
		panic(err)
	}
}

func reverseProxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Reverse Proxy:\t", r.URL)
		p.ServeHTTP(w, r)
	}
}
