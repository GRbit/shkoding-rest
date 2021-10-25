// Package handlers provide handlers for vk-stats application
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
)

type loggerLevel struct {
	Level string `json:"logger_level"`
}

// SetLoggerLevel set global logging level
func SetLoggerLevel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
			lvl zerolog.Level
			msg loggerLevel
		)

		if err = json.NewDecoder(r.Body).Decode(&msg); err != nil {
			writeResp(w, nil, http.StatusBadRequest,
				xerrors.Errorf("msg '%T' bind err: %w", msg, err))
		}

		switch msg.Level {
		case "debug", "dbg":
			lvl = zerolog.DebugLevel
		case "info":
			lvl = zerolog.InfoLevel
		case "warn":
			lvl = zerolog.WarnLevel
		case "error", "err":
			lvl = zerolog.ErrorLevel
		case "fatal":
			lvl = zerolog.FatalLevel
		}

		log.Logger = log.Logger.Level(lvl)

		writeResp(w, nil, http.StatusOK, nil)
	}
}
