package handlers

import (
	"context"
	"net/http"
	"encoding/json"

	"tb-rate-limiter/internal/models"
)

type LimiterService interface {
	Verify(ctx context.Context, req models.LimiterRequest) (models.LimiterResponse, error)
}

type LimiterHandler struct {
	service LimiterService
}

func NewLimiterHandler(service LimiterService) *LimiterHandler {
	return &LimiterHandler{service: service}
}

func(h *LimiterHandler) CreateResponse(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newReq models.LimiterRequest
	if err := json.NewDecoder(r.Body).Decode(&newReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if newReq.Key == "" || newReq.Limit < 0 {
		http.Error(w, "Data is invalid", http.StatusBadRequest)
		return
	}

	newRes, err := h.service.Verify(r.Context(), newReq)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if newRes.Allowed != true {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	w.WriteHeader(http.StatusOK)
}