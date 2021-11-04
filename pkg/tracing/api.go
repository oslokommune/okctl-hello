package tracing

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

const trackerName = "mytracker"

// HandlerFunc returns the health handler
func HandlerFunc(logger *logrus.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := handler(logger, writer)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func handler(logger *logrus.Logger, writer io.Writer) error {
	f, err := os.Create("traces.txt")
	if err != nil {
		return fmt.Errorf("creating trace file: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		return fmt.Errorf("creating exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Errorf("error shutting down tracer provider: %s", err.Error())
		}
	}()
	otel.SetTracerProvider(tp)

	err = work(writer)
	if err != nil {
		return fmt.Errorf("working: %w", err)
	}

	return nil
}

func work(writer io.Writer) error {
	var err error
	response := strings.Builder{}

	ctx, span := otel.Tracer(trackerName).Start(context.Background(), "Run")
	defer span.End()

	response.WriteString(fmt.Sprintln("Doing some work..."))

	digForGold(ctx)
	goToTheMoon(ctx)

	response.WriteString(fmt.Sprintln("Done working."))

	_, err = writer.Write([]byte(response.String()))
	if err != nil {
		return fmt.Errorf("creating response string: %w", err)
	}

	return nil
}

func digForGold(ctx context.Context) {
	_, span := otel.Tracer(trackerName).Start(ctx, "Moon")
	defer span.End()

	time.Sleep(time.Millisecond * 700)
}

func goToTheMoon(ctx context.Context) {
	_, span := otel.Tracer(trackerName).Start(ctx, "Moon")
	defer span.End()

	rndInt := rand.Intn(10) //nolint:gosec
	span.SetAttributes(attribute.String("someImportantValue", strconv.Itoa(rndInt)))

	time.Sleep(time.Millisecond * 200)

}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (*stdouttrace.Exporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("okctl-hello"),
			semconv.ServiceVersionKey.String("v0.0.10"),
			attribute.String("environment", "test"),
		),
	)
	return r
}
