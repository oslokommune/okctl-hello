package main

import (
	_ "embed"
	"log"
	"net/http"

	"github.com/oslokommune/okctl-hello/pkg/content"
	"github.com/oslokommune/okctl-hello/pkg/monitoring"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed public/index.html
var indexHtml []byte

//go:embed public/logo.png
var logo []byte

func main() {
	server := http.NewServeMux()

	server.Handle("/metrics", promhttp.Handler())

	server.Handle("/logo.png", monitoring.NewOlliCounterMiddleware(
		content.LogoHandler(logo),
	))
	server.Handle("/", monitoring.NewHitCounterMiddleware(
		content.StaticHtmlHandler(indexHtml),
	))

	log.Fatal(http.ListenAndServe(":3000", server))
}
