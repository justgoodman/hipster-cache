package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var addr = flag.String("listen-address",":4000","The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(*addr, nil)
}
