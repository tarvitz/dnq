package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/http2"
)

var (
	httpServer  *http.Server
	http2Server = &http2.Server{}
	cmd         *Command
)

func init() {
	http.Handle("/", Default())
	http.Handle("/echo", Echo())
	http.Handle("/inline", Inline())
	http.Handle("/reload", Reload())
	http.Handle("/metrics", promhttp.Handler())
}
