// Package monitoring knows how to deliver metrics
package monitoring

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	HitCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "okctl_route_hits",
		Help: "Counts the number of hits to a route",
	})
	ImageLoads = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "okctl_olli_counter",
		Help: "Counts the number of hits to the Olli logo",
	})
)

func NewHitCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		HitCounter.Inc()

		next.ServeHTTP(writer, request)
	})
}

func NewOlliCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ImageLoads.Inc()

		next.ServeHTTP(writer, request)
	})
}

func init() {
	err := prometheus.Register(HitCounter)
	if err != nil {
		log.Fatalln(err)
	}

	err = prometheus.Register(ImageLoads)
	if err != nil {
		log.Fatalln(err)
	}
}
