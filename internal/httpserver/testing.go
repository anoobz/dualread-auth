package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/anoobz/dualread/auth/internal/model"
	"github.com/anoobz/dualread/auth/internal/store"
	"github.com/anoobz/dualread/auth/internal/store/mockstore"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

func NewTestLogger(t *testing.T) *log.Logger {
	t.Helper()

	return log.New(ioutil.Discard, "", 0)
}

func NewTestServer(t *testing.T) *server {
	t.Helper()

	store := mockstore.CreateTestStore(t)
	logger := NewTestLogger(t)

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		t.Fatal(err)
	}

	return NewServer(store, logger, port)
}

func (s *server) CreateTestUser(t *testing.T, count int, admin bool) []*model.User {
	t.Helper()

	return store.CreateTestUser(t, s.store, count, admin)
}

func (s *server) CreateTestRequest(
	t *testing.T,
	method string,
	url string,
	payload map[string]string,
) *http.Request {
	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(payload)

	req, err := http.NewRequest(method, url, b)

	if err != nil {
		t.Fatal(err)
	}

	return req
}

func (s *server) DeleteTestRefreshTokenUser(t *testing.T, cookie *http.Cookie) {
	t.Helper()

	jwtToken, err := jwt.Parse(
		cookie.Value,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}
			return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if !jwtToken.Valid {
		t.Fatal(errors.New("invalid token"))
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal(errors.New("invalid token"))
	}

	user_id, err := strconv.ParseInt(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	s.store.User().Delete(user_id)
}

func (s *server) DeleteTestRefreshToken(t *testing.T, cookie *http.Cookie) {
	t.Helper()

	jwtToken, err := jwt.Parse(
		cookie.Value,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method: %v",
					token.Header["alg"],
				)
			}
			return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if !jwtToken.Valid {
		t.Fatal(errors.New("invalid token"))
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal(errors.New("invalid token"))
	}

	refresh_uuid := fmt.Sprintf("%s", claims["refresh_uuid"])

	err = s.store.AuthToken().Delete(refresh_uuid)
	if err != nil {
		t.Fatal(err)
	}
}

func (s *server) LoginTestUser(t *testing.T, email string, password string) string {
	t.Helper()

	rec := httptest.NewRecorder()
	req := s.CreateTestRequest(
		t, http.MethodPost, "/auth/login",
		map[string]string{"email": email, "password": password},
	)
	s.ServeHTTP(rec, req)

	res := struct {
		model.AuthToken
		ErrorMsg string `json:"error"`
	}{}
	err := json.NewDecoder(rec.Body).Decode(&res)
	if err != nil {
		t.Fatal(err)
	}

	return res.TokenString
}
