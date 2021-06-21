// Package content knows how to deliver content
package health

import (
	"github.com/oslokommune/okctl-hello/pkg/logging"
	"github.com/sirupsen/logrus"
	"net/http"
)

// HandlerFunc returns the health handler
func HandlerFunc() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("OK"))
	}
}

func DebugHandler(logger *logrus.Logger) http.Handler {
	return logging.NewLoggingMiddleware(
		logger,
		HandlerFunc(),
	)
}
