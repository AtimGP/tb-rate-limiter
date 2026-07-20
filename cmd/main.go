package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"hl-rate-limiter/internal/handlers"
	"hl-rate-limiter/internal/middleware"
	"hl-rate-limiter/internal/models"
	"hl-rate-limiter/internal/storage"
)

type mockService struct{}

func (s *mockService) Verify(ctx context.Context, req models.LimiterRequest) (models.LimiterResponse, error) {
	return models.LimiterResponse{
		Allowed: true,
		Remaining: req.Limit-1,
		ResetAfter: "1s",
	}, nil
}

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	service := storage.NewMemoryLimiter()
	limiterHandler := handlers.NewLimiterHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/limiter/check", limiterHandler.CreateResponse)
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