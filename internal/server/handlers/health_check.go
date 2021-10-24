package handlers

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// Health just replies 200 OK
func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		log.Error().Err(err).Msg("HealthCheck: body write err")
	}
}
