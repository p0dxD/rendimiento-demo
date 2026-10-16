package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// tracer nombra los momentos propios de la app en sus trazas.
var tracer = otel.Tracer("rendimiento-demo")

// startTracing manda las trazas al colector de rendimiento, para la página
// del universo. Con observe: true la plataforma pone OTEL_EXPORTER_OTLP_ENDPOINT,
// OTEL_SERVICE_NAME y el muestreo; sin ellas no se manda nada, y los tramos
// que hace la app casi no cuestan.
func startTracing(ctx context.Context) func(context.Context) error {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		return func(context.Context) error { return nil }
	}
	exp, err := otlptracehttp.New(ctx)
	if err != nil {
		return func(context.Context) error { return nil }
	}
	res, _ := resource.New(ctx, resource.WithFromEnv(), resource.WithTelemetrySDK())
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp), sdktrace.WithResource(res))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp.Shutdown
}

// traced mide cada solicitud, con el nombre de su ruta ("GET /api/info").
func traced(h http.Handler) http.Handler {
	return otelhttp.NewHandler(h, "http", otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
		if r.Pattern != "" {
			return strings.TrimSuffix(r.Pattern, "{$}")
		}
		return r.Method
	}))
}
