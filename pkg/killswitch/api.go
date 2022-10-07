// Package killswitch provides a simple kill switch for the application
package killswitch

import (
	"net/http"
	"os"
)

type logger interface {
	Debug(args ...interface{})
}

// HandlerFunc returns the health handler
func HandlerFunc(log logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Debug("Got kill request")

		mode := request.URL.Query().Get("mode")
		if mode == "" || mode == "graceful" {
			log.Debug("Graceful shutdown")

			os.Exit(0)
		}

		log.Debug("Disgraceful shutdown")

		panic("Killed brutally")
	}
}
