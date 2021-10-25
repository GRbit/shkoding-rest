package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// xSystemTokenHeader is system token header
const (
	xSystemTokenHeader = "X-System-Token"
	xContentType       = "Content-Type"
	responseLoggingLen = 100
)

type statusRecorder struct {
	http.ResponseWriter
	status   int
	response string
}

// Write wraps ResponseWrite Write() method to save response.
func (rec *statusRecorder) Write(r []byte) (int, error) {
	l := len(r)
	if l > responseLoggingLen {
		l = responseLoggingLen
	}

	if rec.status == http.StatusOK { // if request served correctly, cut response for logs
		rec.response = string(r[:l])
	}

	return rec.ResponseWriter.Write(r)
}

// WriteHeader wraps ResponseWrite WriteHeader() method to save response code.
func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// maybe one day we add custom logger here
func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().
			Str("path", r.URL.Path).
			Str("raw_query", r.URL.RawQuery).
			Str("method", r.Method).
			Str("event_type", "request_received").
			Msg("request received")

		rec := statusRecorder{w, http.StatusOK, ""}

		next.ServeHTTP(&rec, r)

		log.Info().
			Str("response", rec.response).
			Int("response_code", rec.status).
			Str("path", r.URL.Path).
			Str("event_type", "request_served").
			Msg("request served")
	})
}

// tokenChecker grant access if system-token is valid
func tokenChecker(systemToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			st := r.Header.Get(xSystemTokenHeader)
			if st == "" || st != systemToken {
				w.WriteHeader(http.StatusUnauthorized)

				log.Info().Str("path", r.URL.Path).
					Str("system_token", r.Header.Get(xSystemTokenHeader)).
					Str("method", r.Method).
					Msg("invalid system token")

				if _, err := w.Write([]byte("invalid system token")); err != nil {
					log.Error().Err(err).Str("path", r.URL.Path).
						Str("method", r.Method).
						Msg("http.ResponseWriter Write() err")
				}

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// contentTypeChecker check that content-type header is application/json
func contentTypeChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t := r.Header.Get(xContentType); t != "application/json" {
			w.WriteHeader(http.StatusBadRequest)

			lg := log.With().Str("component", "middleware").Logger()

			b := []byte(fmt.Sprintf(
				`{"result":null,"error":{"message":"unknown Content-Type to handle", "details":"%s"}}`, t))

			if _, err := w.Write(b); err != nil {
				lg.Error().Bytes("body", b).Err(err).Msg("body write err")
			}
		}

		next.ServeHTTP(w, r)
	})
}
