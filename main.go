package main

import (
	_ "embed"
	"github.com/oslokommune/okctl-hello/pkg/content"
	"github.com/oslokommune/okctl-hello/pkg/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

// go:embed public/index.html
var indexHtml []byte

func main() {
	server := http.NewServeMux()

	server.Handle("/metrics", promhttp.Handler())

	server.Handle("/", monitoring.NewMonitoringMiddleware(
		content.StaticHandler(indexHtml),
	))

	log.Fatal(http.ListenAndServe(":3000", server))
}
