// Package monitoring knows how to deliver metrics
package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

var (
	HitCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "okctl_hello_hits",
		Help: "Counts the number of hits to okctl-hello",
	})
)

func NewMonitoringMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		HitCounter.Inc()

		next.ServeHTTP(writer, request)
	})
}

func init() {
	err := prometheus.Register(HitCounter)
	if err != nil {
		log.Fatalln(err)
	}
}
