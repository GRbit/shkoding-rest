// Package storage provide methods for reading and writing ads data
package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type internalStorageMap struct {
	Students          map[int64]*Student
	studentsIncrement int64
	sync.RWMutex
}

// Config database connection
type Config struct {
	Debug  bool
	Logger Logger
}

// Logger interface implements just standard Printf function
type Logger interface {
	Printf(format string, v ...interface{})
}

// Storage execute query to database
type Storage struct {
	m internalStorageMap
}

// New makes a new Storage instance
func New(cfg Config) (*Storage, error) {
	log.Info().Msg("Creating new storage")
	defer log.Info().Msg(fmt.Sprintf("New storage created. cfg: %v", cfg))

	storage := new(Storage)
	storage.m.Students = make(map[int64]*Student)

	if cfg.Debug {
		go func() {
			time.Sleep(time.Second * 10)
			storage.dbStats()
		}()
	}

	return storage, nil
}

// dbStats is debug function to watch database state
func (s *Storage) dbStats() {
	s.m.RLock()
	defer s.m.RUnlock()

	log.Debug().
		Int("number_of_students", len(s.m.Students)).
		Msg("DB_STATS")
}
