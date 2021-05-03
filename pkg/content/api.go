// Package content knows how to deliver content
package content

import "net/http"

// StaticHandler returns the provided bytes array on each call
func StaticHandler(content []byte) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write(content)
	}
}
