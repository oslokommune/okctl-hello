package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"

	"github.com/oslokommune/okctl-hello/pkg/logging"
	"github.com/sirupsen/logrus"

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

	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{},
		Level:     logrus.InfoLevel,
	}

	server.Handle("/metrics", promhttp.Handler())

	server.Handle("/logo.png", monitoring.NewOlliCounterMiddleware(
		content.LogoHandler(logo),
	))
	server.Handle("/", monitoring.NewHitCounterMiddleware(
		logging.NewLoggingMiddleware(
			logger,
			content.StaticHtmlHandler(indexHtml),
		),
	))

	log.Fatal(http.ListenAndServe(":3000", server))
}
