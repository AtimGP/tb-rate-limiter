package main

import (
	"log/slog"
	"net/http"
	"os"

	"tb-rate-limiter/internal/handlers"
	"tb-rate-limiter/internal/middleware"
	"tb-rate-limiter/internal/storage"
)

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	service := storage.NewMemoryLimiter()
	limiterHandler := handlers.NewLimiterHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/limiter/verify", limiterHandler.CreateResponse)
	mux.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
		panic("Паника")
	})

	wrappedHandler := middleware.Logger(middleware.Recovery(mux))

	slog.Info("Starting server", "port", "8080")

	if err := http.ListenAndServe(":8080", wrappedHandler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}