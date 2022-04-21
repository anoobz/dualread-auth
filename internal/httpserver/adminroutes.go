package httpserver

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func (s *server) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.validateAdmin(w, r) {
			s.error(w, r, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		u, err := s.store.User().GetById(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) getUserPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) insertUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) validateAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims := r.Context().Value(ctxKey("claims")).(jwt.MapClaims)
	return claims["admin"].(bool)
}
