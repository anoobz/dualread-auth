package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s *server) getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		page, err := strconv.ParseUint(mux.Vars(r)["page"], 10, 64)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		users, err := s.store.User().GetPage(page)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, users)
	}
}

func (s *server) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.User().GetAll()
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		s.respond(w, r, http.StatusOK, users)
	}
}

func (s *server) insertUser() http.HandlerFunc {
	type payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Admin    bool   `json:"admin"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		p := &payload{}
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		u, err := s.store.User().Insert(p.Email, p.Password, p.Admin, time.Now())
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		clauses := map[string]interface{}{}
		if err := json.NewDecoder(r.Body).Decode(&clauses); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.store.User().Update(id, clauses)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.store.User().Delete(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}
