package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/oslokommune/okctl-hello/pkg/postgres"
	"github.com/sirupsen/logrus"
)

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
