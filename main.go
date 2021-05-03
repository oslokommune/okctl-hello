package main

import (
	_ "embed"
	"github.com/oslokommune/okctl-hello/pkg/content"
	"github.com/oslokommune/okctl-hello/pkg/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

//go:embed public/index.html
var indexHtml []byte

//go:embed public/logo.png
var logo []byte

func main() {
	server := http.NewServeMux()

	server.Handle("/metrics", promhttp.Handler())

	server.Handle("/logo.png", content.LogoHandler(logo))
	server.Handle("/", monitoring.NewMonitoringMiddleware(
		content.StaticHtmlHandler(indexHtml),
	))

	log.Fatal(http.ListenAndServe(":3000", server))
}
