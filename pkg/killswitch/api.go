// Package killswitch provides a simple kill switch for the application
package killswitch

import (
	"fmt"
	"net/http"
	"os"
)

type logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

const (
	killModeGraceful = "graceful"
	killModePanic    = "panic"
)

// HandlerFunc returns the health handler
func HandlerFunc(log logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Debug("Got kill request")

		mode := request.URL.Query().Get("mode")
		if mode == "" || mode == killModeGraceful {
			log.Debug("Graceful shutdown")

			os.Exit(0)
		}

		if mode == killModePanic {
			log.Debug("Panic shutdown")

			panic("killswitch")
		}

		log.Debugf("Unknown mode %s", mode)

		writer.WriteHeader(http.StatusBadRequest)

		_, err := writer.Write([]byte(fmt.Sprintf("400 BAD REQUEST: Unknown mode %s", mode)))
		if err != nil {
			log.Errorf("Failed to write response: %s", err.Error())
		}
	}
}
