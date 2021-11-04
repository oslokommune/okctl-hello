package main

import (
	_ "embed"
	"fmt"
	"github.com/oslokommune/okctl-hello/pkg/tracing"
	"net/http"
	"os"
	"strconv"

	"github.com/oslokommune/okctl-hello/pkg/communicationtest"

	"github.com/oslokommune/okctl-hello/pkg/loadtest"

	"github.com/oslokommune/okctl-hello/pkg/health"

	"github.com/oslokommune/okctl-hello/pkg/logging"
	"github.com/oslokommune/okctl-hello/pkg/postgres"
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

	if debug := os.Getenv("DEBUG"); debug != "true" {
		server.Handle("/health", health.HandlerFunc())
	} else {
		server.Handle("/health", health.DebugHandler(logger))
	}

	server.Handle("/metrics", promhttp.Handler())

	server.Handle("/logo.png", monitoring.NewOlliCounterMiddleware(
		content.LogoHandler(logo),
	))

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

	server.Handle("/burn-cpu", loadtest.HandlerFunc())

	server.Handle("/commtest", communicationtest.HandlerFunc())

	server.Handle("/tracing", tracing.HandlerFunc(logger))

	server.Handle("/", monitoring.NewHitCounterMiddleware(
		logging.NewLoggingMiddleware(
			logger,
			content.StaticHtmlHandler(indexHtml),
		),
	))

	logger.Fatal(http.ListenAndServe(":3000", server))
}

func enablePostgres(server *http.ServeMux, logger *logrus.Logger, rawDSN string) error {
	dsn := postgres.ParseDSN(rawDSN)

	if err := dsn.Validate(); err != nil {
		return fmt.Errorf("validating DSN: %w", err)
	}

	pgClient := postgres.Client{DSN: dsn}
	dbErrorFields := logrus.Fields{
		"database-host": pgClient.DSN.URI,
		"database-port": pgClient.DSN.Port,
		"database-user": pgClient.DSN.Username,
	}

	server.HandleFunc("/postgres/write", func(w http.ResponseWriter, r *http.Request) {
		err := pgClient.Open()
		if err != nil {
			logger.WithFields(dbErrorFields).Errorf("opening database: %s", err.Error())

			return
		}

		defer func() {
			_ = pgClient.Close()
		}()

		err = pgClient.Write()
		if err != nil {
			logger.WithFields(dbErrorFields).Errorf("writing to database: %s", err.Error())

			return
		}

		w.WriteHeader(http.StatusOK)
	})

	server.HandleFunc("/postgres/read", func(w http.ResponseWriter, r *http.Request) {
		err := pgClient.Open()
		if err != nil {
			logger.WithFields(dbErrorFields).Errorf("opening database: %s", err.Error())

			return
		}

		defer func() {
			_ = pgClient.Close()
		}()

		currentHits, err := pgClient.Read()
		if err != nil {
			logger.WithFields(dbErrorFields).Errorf("reading from database: %s", err.Error())

			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("hits: %s", strconv.Itoa(currentHits))))
	})

	return nil
}
