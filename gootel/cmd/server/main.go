package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isaacdsc/gootel/internal/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	otelapi "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer         trace.Tracer
	meter          metric.Meter
	requestCounter metric.Int64Counter
	requestLatency metric.Float64Histogram
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configuration from environment variables
	cfg := otel.Config{
		ServiceName:    getEnv("SERVICE_NAME", "gootel-service"),
		ServiceVersion: getEnv("SERVICE_VERSION", "1.0.0"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		OTLPEndpoint:   getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
	}

	// Setup OpenTelemetry
	provider, err := otel.Setup(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to setup OpenTelemetry: %v", err)
	}
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := provider.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error shutting down OpenTelemetry: %v", err)
		}
	}()

	// Initialize tracer and meter
	tracer = otelapi.Tracer(cfg.ServiceName)
	meter = otelapi.Meter(cfg.ServiceName)

	// Create metrics
	if err := setupMetrics(); err != nil {
		log.Fatalf("Failed to setup metrics: %v", err)
	}

	// Setup HTTP server with OpenTelemetry instrumentation
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/users", usersHandler)
	mux.HandleFunc("/api/process", processHandler)

	// Wrap the mux with OpenTelemetry HTTP instrumentation
	handler := otelhttp.NewHandler(mux, "http-server",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
		cancel()
	}()

	log.Printf("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func setupMetrics() error {
	var err error

	requestCounter, err = meter.Int64Counter("http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return fmt.Errorf("failed to create request counter: %w", err)
	}

	requestLatency, err = meter.Float64Histogram("http_request_duration_seconds",
		metric.WithDescription("HTTP request latency in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return fmt.Errorf("failed to create request latency histogram: %w", err)
	}

	return nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	start := time.Now()

	// Record metrics
	defer func() {
		duration := time.Since(start).Seconds()
		requestCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("endpoint", "/"),
			attribute.String("method", r.Method),
		))
		requestLatency.Record(ctx, duration, metric.WithAttributes(
			attribute.String("endpoint", "/"),
			attribute.String("method", r.Method),
		))
	}()

	response := map[string]string{
		"message": "Welcome to Gootel - Go OpenTelemetry Demo",
		"version": "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	start := time.Now()

	defer func() {
		duration := time.Since(start).Seconds()
		requestCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("endpoint", "/api/users"),
			attribute.String("method", r.Method),
		))
		requestLatency.Record(ctx, duration, metric.WithAttributes(
			attribute.String("endpoint", "/api/users"),
			attribute.String("method", r.Method),
		))
	}()

	// Create a child span for database operation simulation
	ctx, span := tracer.Start(ctx, "fetch-users-from-db",
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", "SELECT"),
		),
	)
	defer span.End()

	// Simulate database query
	users, err := fetchUsersFromDB(ctx)
	if err != nil {
		span.RecordError(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	span.SetAttributes(attribute.Int("users.count", len(users)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func fetchUsersFromDB(ctx context.Context) ([]map[string]interface{}, error) {
	// Simulate database latency
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	users := []map[string]interface{}{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
		{"id": 2, "name": "Bob", "email": "bob@example.com"},
		{"id": 3, "name": "Charlie", "email": "charlie@example.com"},
	}

	return users, nil
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	start := time.Now()

	defer func() {
		duration := time.Since(start).Seconds()
		requestCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("endpoint", "/api/process"),
			attribute.String("method", r.Method),
		))
		requestLatency.Record(ctx, duration, metric.WithAttributes(
			attribute.String("endpoint", "/api/process"),
			attribute.String("method", r.Method),
		))
	}()

	// Demonstrate nested spans
	result, err := processData(ctx)
	if err != nil {
		http.Error(w, "Processing failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func processData(ctx context.Context) (map[string]interface{}, error) {
	ctx, span := tracer.Start(ctx, "process-data")
	defer span.End()

	// Step 1: Validate
	if err := validateData(ctx); err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Step 2: Transform
	transformed, err := transformData(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Step 3: Store
	if err := storeData(ctx, transformed); err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.Bool("processing.success", true))

	return map[string]interface{}{
		"status":  "processed",
		"data":    transformed,
		"traceId": span.SpanContext().TraceID().String(),
	}, nil
}

func validateData(ctx context.Context) error {
	_, span := tracer.Start(ctx, "validate-data")
	defer span.End()

	time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
	span.SetAttributes(attribute.Bool("validation.passed", true))
	return nil
}

func transformData(ctx context.Context) (map[string]interface{}, error) {
	_, span := tracer.Start(ctx, "transform-data")
	defer span.End()

	time.Sleep(time.Duration(20+rand.Intn(30)) * time.Millisecond)

	data := map[string]interface{}{
		"transformed": true,
		"timestamp":   time.Now().Unix(),
		"value":       rand.Float64() * 100,
	}

	span.SetAttributes(
		attribute.Bool("transform.success", true),
		attribute.Float64("transform.value", data["value"].(float64)),
	)

	return data, nil
}

func storeData(ctx context.Context, data map[string]interface{}) error {
	_, span := tracer.Start(ctx, "store-data",
		trace.WithAttributes(
			attribute.String("db.system", "redis"),
			attribute.String("db.operation", "SET"),
		),
	)
	defer span.End()

	time.Sleep(time.Duration(15+rand.Intn(25)) * time.Millisecond)
	span.SetAttributes(attribute.Bool("store.success", true))
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
