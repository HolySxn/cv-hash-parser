package httpHandler

import (
	"net/http"
)

func NewServer(
	handler *Handler,
) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("POST /parse-hash", handler.ParseHash)

	return router
}
