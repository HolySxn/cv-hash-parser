package httpHandler

import (
	"log/slog"
	"net/http"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) ParseHash(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Parsing hash request received")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hash parsed successfully"))
}
