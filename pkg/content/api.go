// Package content knows how to deliver content
package content

import "net/http"

// StaticHtmlHandler returns the provided bytes array on each call
func StaticHtmlHandler(content []byte) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		writer.WriteHeader(http.StatusOK)

		_, _ = writer.Write(content)
	}
}

// LogoHandler returns the set logo
func LogoHandler(logo []byte) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "image/png")
		writer.WriteHeader(http.StatusOK)

		_, _ = writer.Write(logo)
	}
}
