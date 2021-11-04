package tracing

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HandlerFunc returns the health handler
func HandlerFunc(logger *logrus.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := handler(logger, writer)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}
	}
}

func handler(logger *logrus.Logger, writer io.Writer) error {
	tracer, closer := InitJaeger("okctl-hello")
	//goland:noinspection GoUnhandledErrorResult
	defer closer.Close()

	span := tracer.StartSpan("work")
	defer span.Finish()

	err := work(logger, writer, tracer)
	if err != nil {
		return fmt.Errorf("working: %w", err)
	}

	return nil
}

func work(logger *logrus.Logger, writer io.Writer, tracer opentracing.Tracer) error {
	var err error
	response := strings.Builder{}

	response.WriteString(fmt.Sprintln("Doing some work..."))

	digForGold(tracer)
	err = goToTheMoon(logger, tracer)
	if err != nil {
		return fmt.Errorf("going to the moon: %w", err)
	}

	response.WriteString(fmt.Sprintln("Done working."))

	_, err = writer.Write([]byte(response.String()))
	if err != nil {
		return fmt.Errorf("creating response string: %w", err)
	}

	return nil
}

func digForGold(tracer opentracing.Tracer) {
	span := tracer.StartSpan("digForGold")
	defer span.Finish()

	time.Sleep(time.Millisecond * 700)
}

func goToTheMoon(logger *logrus.Logger, tracer opentracing.Tracer) error {
	span := tracer.StartSpan("goToTheMoon")
	defer span.Finish()

	url := "http://localhost:3000/tracing-receiver"
	req, _ := http.NewRequest("GET", url, nil)

	// Set some tags on the clientSpan to annotate that it's the client span. The additional HTTP tags are useful for debugging purposes.
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")

	// Inject the client span context into the headers
	err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		return fmt.Errorf("injecting: %w", err)
	}
	resp, _ := http.DefaultClient.Do(req)

	respString, err := ioutil.ReadAll(resp.Body)
	logger.Infof("Response: %s\n", string(respString))

	return nil

}
