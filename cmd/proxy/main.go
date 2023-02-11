package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	host1 = "http://localhost:8081"
	host2 = "http://localhost:8082"
	port  = ":8080"
)

var (
	count int
)

func main() {
	http.HandleFunc("/", loadBalacer)
	log.Fatal(http.ListenAndServe(port, nil))
}

func loadBalacer(res http.ResponseWriter, req *http.Request) {
	url := getProxyURL()
	logRequestPayload(url)
	serveReversProxy(url, res, req)
}

func getProxyURL() string {
	var servers = []string{host1, host2}
	server := servers[count]
	count++
	count = count % len(servers)
	return server
}

func logRequestPayload(proxyUrl string) {
	log.Printf("proxy_url: %s\n", proxyUrl)

}
func serveReversProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(res, req)
}
