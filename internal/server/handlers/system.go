// Package handlers provide handlers for vk-stats application
package handlers

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
)

type loggerLevel struct {
	Level string `json:"logger_level"`
}

func (m *loggerLevel) validate() error {
	m.Level = strings.ToLower(m.Level)

	switch m.Level {
	case "debug", "info", "warn", "error":
	default:
		return xerrors.Errorf("loggerLevel: unknown logger level: '%s'", m.Level)
	}

	return nil
}

// SetLoggerLevel set global logging level
func SetLoggerLevel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
			lvl zerolog.Level
			msg loggerLevel
		)

		if err = bind(r, &msg); err != nil {
			writeResp(w, rID(r), nil, http.StatusBadRequest,
				xerrors.Errorf("msg '%T' bind err: %w", msg, err))

			return
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

		writeResp(w, rID(r), nil, http.StatusOK, nil)
	}
}
