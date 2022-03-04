/*
webfront is an HTTP server and reverse proxy.
It reads a JSON-formatted rule file like this:
	[
		{"Host": "example.com", "Serve": "/var/www"},
		{"Host": "example.org", "Forward": "localhost:8080"}
	]
For all requests to the host example.com (or any name ending in
".example.com") it serves files from the /var/www directory.
For requests to example.org, it forwards the request to the HTTP
server listening on localhost port 8080.
Usage of webfront:
  -http address
    	HTTP listen address (default ":http")
  -letsencrypt_cache directory
    	letsencrypt cache directory (default is to disable HTTPS)
  -poll interval
    	rule file poll interval (default 10s)
  -rules file
    	rule definition file
webfront was written by Andrew Gerrand <adg@golang.org>
*/
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/acme/autocert"
)

var (
	httpAddr     = flag.String("http", ":http", "HTTP listen `address`")
	metricsAddr  = flag.String("metrics", "", "metrics HTTP listen `address`")
	letsCacheDir = flag.String("letsencrypt_cache", "", "letsencrypt cache `directory` (default is to disable HTTPS)")
	ruleFile     = flag.String("rules", "", "rule definition `file`")
	pollInterval = flag.Duration("poll", time.Second*10, "rule file poll `interval`")
)
var hitCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "webfront_hits",
		Help: "Cumulative hits since startup.",
	},
	[]string{"host"},
)

func init() {
	prometheus.MustRegister(hitCounter)
}

func main() {
	flag.Parse()

	s, err := NewServer(*ruleFile, *pollInterval)
	if err != nil {
		log.Fatal(err)
	}
	if *metricsAddr != "" {
		go func() {
			log.Fatal(http.ListenAndServe(*metricsAddr, promhttp.Handler()))
		}()
	}
	if *letsCacheDir != "" {
		m := &autocert.Manager{
			Cache:      autocert.DirCache(*letsCacheDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: s.hostPolicy,
		}

		c := tls.Config{GetCertificate: m.GetCertificate}
		l, err := tls.Listen("tcp", ":https", &c)
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			log.Fatal(http.Serve(l, s))
		}()
		log.Fatal(http.ListenAndServe(*httpAddr, m.HTTPHandler(s)))
	}
}
