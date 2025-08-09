package httpHandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/HolySxn/cv-hash-parser/internal/service"
)

type Handler struct {
	logger  *slog.Logger
	service *service.Service
}

func NewHandler(logger *slog.Logger, service *service.Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

type ParseHashRequest struct {
	CvURL string `json:"cv_url"`
}

func (h *Handler) ParseHash(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Parsing hash request received")

	var req ParseHashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validateRequest(req); err != nil {
		h.logger.Error("validation failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.ProcessCV(r.Context(), req.CvURL); err != nil {
		h.logger.Error("failed to process cv", "error", err)
		http.Error(w, "failed to process cv", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hash parsed and email sent successfully"))
}

func (h *Handler) validateRequest(req ParseHashRequest) error {
	if req.CvURL == "" {
		return errors.New("cv_url is required")
	}
	if _, err := url.Parse(req.CvURL); err != nil {
		return errors.New("invalid cv_url format")
	}
	return nil
}
