// Package handlers provide handlers for vk-stats application
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
			writeResp(w, ret, http.StatusNotFound, xerrors.Errorf("can't find student with id='%s'", studentID))
			return
		}

		writeResp(w, ret, http.StatusOK, nil)
	}
}

type NewStudentMessage struct {
	Name     string
	Telegram string
}

func NewStudentDocs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResp(w, NewStudentMessage{
			Name:     "here comes the name of a student",
			Telegram: "here you can add telegram nickname",
		}, http.StatusCreated, nil)
	}
}

func NewStudent(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m NewStudentMessage

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			writeResp(w, nil, http.StatusBadRequest, xerrors.Errorf("request json decoding: %w", err))
			return
		}

		if m.Name == "" && m.Telegram == "" {
			writeResp(w, nil, http.StatusBadRequest, xerrors.Errorf("no name or telegram parameters: %v", m))
			return
		}

		var ret interface{}

		ret = s.NewStudent(m.Name, m.Telegram)

		writeResp(w, ret, http.StatusCreated, nil)
	}
}

type UpdateStudentMessage struct {
	ID       int64
	Name     string
	Telegram string
}

func UpdateStudentDocs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeResp(w, struct{ ID, Name, Telegram string }{
			ID:       "Put there an existing ID of a student you want to change",
			Name:     "Put there new name of a student. Leave empty is you don't want to change it.",
			Telegram: "Put there new telegram nickname of a student. Leave empty is you don't want to change it.",
		}, http.StatusOK, nil)
	}
}

func UpdateStudent(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m UpdateStudentMessage

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			writeResp(w, nil, http.StatusBadRequest, xerrors.Errorf("request json decoding: %w", err))
			return
		}

		var ret interface{}

		ret, ok := s.UpdateStudent(m.ID, m.Name, m.Telegram)
		if !ok {
			writeResp(w, ret, http.StatusNotFound, xerrors.Errorf("can't find student with id='%d'", m.ID))
			return
		}

		writeResp(w, ret, http.StatusOK, nil)
	}
}

func DeleteStudent(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		studentID := chi.URLParam(r, "studentID")

		id, err := strconv.Atoi(studentID)
		if err != nil {
			writeResp(w, nil, http.StatusBadRequest, xerrors.Errorf("can't recognize student id='%s'", studentID))
			return
		}

		s.DeleteStudent(int64(id))

		writeResp(w, "success", http.StatusOK, nil)
	}
}
