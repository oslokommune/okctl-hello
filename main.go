package main

import (
	_ "embed"
	"net/http"
	"os"

	"github.com/oslokommune/okctl-hello/pkg/communicationtest"
	"github.com/oslokommune/okctl-hello/pkg/killswitch"

	"github.com/oslokommune/okctl-hello/pkg/loadtest"

	"github.com/oslokommune/okctl-hello/pkg/health"

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

var logger = &logrus.Logger{
	Out:       os.Stdout,
	Formatter: &logrus.JSONFormatter{},
	Level:     logrus.InfoLevel,
}

type route struct {
	path    string
	handler http.Handler
}

var routes = []route{
	{
		path: "/",
		handler: monitoring.NewHitCounterMiddleware(logging.NewLoggingMiddleware(
			logger,
			content.StaticHtmlHandler(indexHtml),
		)),
	},
	{
		path:    "/logo.png",
		handler: monitoring.NewOlliCounterMiddleware(content.LogoHandler(logo)),
	},
	{
		path:    "/metrics",
		handler: promhttp.Handler(),
	},
	{
		path: "/health",
		handler: func() http.Handler {
			if debug := os.Getenv("DEBUG"); debug == "true" {
				return health.DebugHandler(logger)
			}

			return health.HandlerFunc()
		}(),
	},
	{
		path:    "/burn-cpu",
		handler: loadtest.HandlerFunc(),
	},
	{
		path:    "/commtest",
		handler: communicationtest.HandlerFunc(),
	},
	{
		path:    "/kill",
		handler: killswitch.HandlerFunc(logger),
	},
}

func main() {
	server := http.NewServeMux()

	for _, route := range routes {
		server.Handle(route.path, route.handler)
	}

	if rawDSN := os.Getenv("DSN"); rawDSN != "" {
		logger.Info("Found DSN. Enabling Postgres integration")

		err := enablePostgres(server, logger, rawDSN)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Info("No DSN found. Ignoring Postgres integration")

		server.HandleFunc("/postgres/read", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("Postgres integration is disabled. Use the DSN environment variable to activate."))
		})

		server.HandleFunc("/postgres/write", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("Postgres integration is disabled. Use the DSN environment variable to activate."))
		})
	}

	logger.Fatal(http.ListenAndServe(":3000", server))
}
