package loadtest

import (
	"net/http"
	"runtime"
	"time"
)

const burnSeconds = 30

// HandlerFunc returns the health handler
func HandlerFunc() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)

		BurnCPUFor30seconds()

		_, _ = writer.Write([]byte("Done burning CPU\n"))
	}
}

// BurnCPUFor30seconds uses max CPU usage for limited time
func BurnCPUFor30seconds() {
	done := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				//nolint:staticcheck
				default:
				}
			}
		}()
	}

	time.Sleep(time.Second * burnSeconds)
	close(done)
}
