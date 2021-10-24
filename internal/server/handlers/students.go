// Package handlers provide handlers for vk-stats application
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GRbit/shkoding-rest/internal/storage"
	"github.com/go-chi/chi"
	"golang.org/x/xerrors"
)

func GetStudent(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ret interface{}
		ok := true

		studentID := chi.URLParam(r, "studentID")
		id, err := strconv.Atoi(studentID)
		if err != nil {
			ret = s.GetStudents()
		} else {
			ret, ok = s.GetStudent(int64(id))
		}

		if !ok {
			writeResp(w, rID(r), ret, http.StatusNotFound,
				xerrors.Errorf("can't find student with id='%s'", studentID))
		}

		writeResp(w, rID(r), ret, http.StatusOK, nil)
	}
}

func NewStudent(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m struct {
			Name string
			Telegram   string
		}

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			writeResp(w, rID(r), nil, http.StatusBadRequest, xerrors.Errorf("request json decoding: %w", err))
		}

		var ret interface{}

		ret = s.NewStudent(m.Name, m.Telegram)

		writeResp(w, rID(r), ret, http.StatusCreated, nil)
	}
}

func UpdateStudent(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m struct {
			ID int64
			Name string
			Telegram   string
		}

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			writeResp(w, rID(r), nil, http.StatusBadRequest, xerrors.Errorf("request json decoding: %w", err))
		}

		var ret interface{}

		ret, ok := s.UpdateStudent(m.ID, m.Name, m.Telegram)
		if !ok {
			writeResp(w, rID(r), ret, http.StatusNotFound,
				xerrors.Errorf("can't find student with id='%s'", m.ID))
		}

		writeResp(w, rID(r), ret, http.StatusOK, nil)
	}
}

func DeleteStudent(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentID := chi.URLParam(r, "studentID")
		id, err := strconv.Atoi(studentID)
		if err != nil {
			writeResp(w, rID(r), nil, http.StatusBadRequest,
				xerrors.Errorf("can't recognize student id='%s'", studentID))
		} else {
			s.DeleteStudent(int64(id))
		}


		writeResp(w, rID(r), "success", http.StatusOK, nil)
	}
}
