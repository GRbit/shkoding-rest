// Package server provides go-chi server for vk-stats application.
package server

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/GRbit/shkoding-rest/internal/server/handlers"
	"github.com/GRbit/shkoding-rest/internal/storage"
)

// Serve runs http.ListenAndServe on go-chi router with address specified in Config
func Serve() error {
	var cfg serviceConfig

	if _, err := flags.Parse(&cfg); err != nil {
		return err
	}

	s, err := initService(&cfg)
	if err != nil {
		return err
	}

	log.Info().Msg("Creating router...")

	r := chi.NewRouter()
	r.Use(tokenChecker(cfg.Service.SystemToken))
	r.Use(reqIDSetter)
	r.Use(logger)
	r.Get("/health", handlers.Health)
	r.Post("/logger_level", handlers.SetLoggerLevel())
	r.Get("/students", handlers.GetStudent(s))
	r.Post("/students", handlers.NewStudent(s))
	r.Put("/students", handlers.UpdateStudent(s))
	r.Patch("/students", handlers.UpdateStudent(s))
	r.Delete("/students", handlers.UpdateStudent(s))

	log.Info().Msg("Everything configured. ListenAndServe.")

	return http.ListenAndServe(cfg.Service.Addr, r)
}

func initService(cfg *serviceConfig) (s *storage.Storage, err error) {
	initLogger(cfg)

	s, err = storage.New(storage.Config{
		Debug:           cfg.Service.Debug,
		Logger:          &log.Logger,
	})

	if err != nil {
		return s, err
	}

	return s, err
}

func initLogger(cfg *serviceConfig) {
	log.Logger = log.
		Level(zerolog.DebugLevel).
		With().Timestamp().
		Str("service", cfg.Service.Name).
		Str("service_version", Version).
		Str("service_built", Built).
		Logger()

	switch cfg.Service.LogLevel {
	case "debug":
		log.Logger = log.Level(zerolog.DebugLevel)
		cfg.Service.Debug = true
	case "info":
		log.Logger = log.Level(zerolog.InfoLevel)
	case "warn":
		log.Logger = log.Level(zerolog.WarnLevel)
	case "error":
		log.Logger = log.Level(zerolog.ErrorLevel)
	}

	if cfg.Service.Debug {
		log.Logger = log.Level(zerolog.DebugLevel)
	}

	if cfg.Service.Console {
		log.Logger = log.
			Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}).
			With().Caller().Timestamp().Logger()
		zerolog.TimeFieldFormat = time.StampMicro
	}

	log.Info().Msg("Logger configured")
}
