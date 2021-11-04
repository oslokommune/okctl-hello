package tracingreceiver

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/oslokommune/okctl-hello/pkg/tracing"
	"io"
	"net/http"
)

// HandlerFunc returns the health handler
func HandlerFunc() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := handler(writer, request)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}
	}
}

func handler(writer io.Writer, request *http.Request) error {
	tracer, closer := tracing.InitJaeger("okctl-hello")
	//goland:noinspection GoUnhandledErrorResult
	defer closer.Close()

	// Extract the context from the headers
	spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))
	serverSpan := tracer.StartSpan("server", ext.RPCServerOption(spanCtx))
	defer serverSpan.Finish()

	_, err := writer.Write([]byte(fmt.Sprintln("Tracing receiver says hello")))
	if err != nil {
		return fmt.Errorf("writing: %w", err)
	}

	return nil
}
