package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/anoobz/dualread/auth/internal/model"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) register() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.Password == "" {
			s.error(
				w,
				r,
				http.StatusBadRequest,
				errors.New("a required field is empty"),
			)
			return
		}

		bcryptCost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcryptCost,
		)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		user, err := s.store.User().Insert(
			req.Email, string(encryptedPassword),
			false,
			time.Now(),
		)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusCreated, user)
	}
}

func (s *server) login() http.HandlerFunc {
	type payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		payload := payload{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := s.store.User().GetByEmail(payload.Email)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.Password), []byte(payload.Password),
		); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		at, err := model.NewAccessToken(user)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		rt, err := model.NewRefreshToken(user)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.store.AuthToken().Insert(rt.Uuid, rt.TokenString, rt.Expires)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    rt.TokenString,
			Expires:  time.Unix(rt.Expires, 0),
			HttpOnly: true,
		})

		s.respond(w, r, http.StatusOK, at)
	}
}

func (s *server) refreshAccessToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := s.getRefreshToken(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		refresh_uuid := fmt.Sprintf("%s", claims["refresh_uuid"])
		_, err = s.store.AuthToken().GetById(refresh_uuid)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		user_id, err := strconv.ParseInt(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		user, err := s.store.User().GetById(user_id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		accessToken, err := model.NewAccessToken(user)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, accessToken)
	}
}
