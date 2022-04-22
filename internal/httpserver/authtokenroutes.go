package httpserver

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *server) getAuthToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		t, err := s.store.AuthToken().GetById(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, t)
	}
}

func (s *server) getAuthTokenPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.ParseUint(mux.Vars(r)["page"], 10, 64)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		t, err := s.store.AuthToken().GetPage(page)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, t)
	}
}

func (s *server) getAllAuthTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := s.store.AuthToken().GetAll()
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, t)
	}
}

func (s *server) deleteAuthToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		err := s.store.AuthToken().Delete(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}
