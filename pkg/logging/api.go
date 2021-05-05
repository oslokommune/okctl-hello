package logging

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLoggingMiddleware(logger *logrus.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		next.ServeHTTP(writer, request)

		duration := time.Since(start)

		logger.WithFields(logrus.Fields{
			"method":   request.Method,
			"route":    request.URL.Path,
			"duration": fmt.Sprintf("%vms", duration.Milliseconds()),
		}).Info("request")
	})
}
