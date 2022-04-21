package httpserver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/anoobz/dualread/auth/internal/httpserver"
	"github.com/anoobz/dualread/auth/internal/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	godotenv.Load("../../.env")

	os.Exit(m.Run())
}

func TestServer_Register(t *testing.T) {
	s := httpserver.NewTestServer(t)

	testCases := []struct {
		name                string
		payload             map[string]string
		expectedStatus      int
		expectedErrorString string
	}{
		{
			name: "success",
			payload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			expectedStatus:      http.StatusCreated,
			expectedErrorString: "",
		},
		{
			name: "empty email",
			payload: map[string]string{
				"email":    "",
				"password": "test_password",
			},
			expectedStatus:      http.StatusBadRequest,
			expectedErrorString: "mail: no address",
		},
		{
			name: "empty password",
			payload: map[string]string{
				"email":    "test@test.test",
				"password": "",
			},
			expectedStatus:      http.StatusBadRequest,
			expectedErrorString: "a required field is empty",
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": "test_password",
			},
			expectedStatus:      http.StatusBadRequest,
			expectedErrorString: "mail: missing '@' or angle-addr",
		},
	}

	for _, tc := range testCases {
		rec := httptest.NewRecorder()

		b := &bytes.Buffer{}
		if err := json.NewEncoder(b).Encode(tc.payload); err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/auth/register", b)
		if err != nil {
			t.Fatal(err)
		}

		s.ServeHTTP(rec, req)

		assert.Equal(
			t, tc.expectedStatus,
			rec.Code, fmt.Sprintf("Case name: %s", tc.name),
		)

		if rec.Code == http.StatusCreated {
			u := &model.User{}
			if err := json.NewDecoder(rec.Body).Decode(&u); err != nil {
				t.Fatal(err)
			}
			assert.NotEmpty(t, u.ID, fmt.Sprintf("Case name: %s", tc.name))
		} else {
			res := &struct {
				Error string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)

			assert.Equal(
				t, tc.expectedErrorString, res.Error,
				fmt.Sprintf("Case name: %s", tc.name),
			)
		}
	}
}

func TestServer_LoginSuccess(t *testing.T) {
	s := httpserver.NewTestServer(t)

	s.CreateTestUser(t, 1, false)
	s.CreateTestUser(t, 1, true)

	testCases := []struct {
		name           string
		payload        map[string]string
		admin          bool
		expectedStatus int
	}{
		{
			name: "non-admin",
			payload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			admin:          false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "non-admin",
			payload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			admin:          true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		rec := httptest.NewRecorder()

		req := s.CreateTestRequest(t, http.MethodPost, "/auth/login", tc.payload)
		s.ServeHTTP(rec, req)

		assert.Equal(
			t, tc.expectedStatus,
			rec.Code, fmt.Sprintf("Case name: %s", tc.name),
		)

		//Assert if the access token is created
		token := &model.AuthToken{}
		if err := json.NewDecoder(rec.Body).Decode(&token); err != nil {
			t.Fatal(err)
		}
		assert.NotEmpty(t, token.TokenString, fmt.Sprintf("Case name: %s", tc.name))

		//Assert if the refresh token http cookie is properly set
		refreshCookie := rec.Result().Cookies()
		assert.Condition(t, func() (success bool) {
			for _, c := range refreshCookie {
				if c.Name == "refresh_token" {
					return true
				}
			}
			return false
		})
	}
}

func TestServer_LoginError(t *testing.T) {
	s := httpserver.NewTestServer(t)

	s.CreateTestUser(t, 1, false)

	testCases := []struct {
		name             string
		payload          map[string]string
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name: "empty email",
			payload: map[string]string{
				"email":    "",
				"password": "test_password0",
			},
			expectedStatus:   http.StatusNotFound,
			expectedErrorMsg: "mail: no address",
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": "test_password0",
			},
			expectedStatus:   http.StatusNotFound,
			expectedErrorMsg: "mail: missing '@' or angle-addr",
		},
		{
			name: "empty password",
			payload: map[string]string{
				"email":    "test0@test.test",
				"password": "",
			},
			expectedStatus:   http.StatusNotFound,
			expectedErrorMsg: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			name: "user not found",
			payload: map[string]string{
				"email":    "notFound@test.test",
				"password": "test_password0",
			},
			expectedStatus:   http.StatusNotFound,
			expectedErrorMsg: "sql: no rows in result set",
		},
	}

	for _, tc := range testCases {
		rec := httptest.NewRecorder()

		req := s.CreateTestRequest(t, http.MethodPost, "/auth/login", tc.payload)

		s.ServeHTTP(rec, req)

		assert.Equal(
			t, tc.expectedStatus,
			rec.Code, fmt.Sprintf("Case name: %s", tc.name),
		)

		errorMsg := struct {
			Error string
		}{}
		if err := json.NewDecoder(rec.Body).Decode(&errorMsg); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tc.expectedErrorMsg, errorMsg.Error)
	}
}

func TestServer_RefreshAccessToken(t *testing.T) {
	s := httpserver.NewTestServer(t)

	s.CreateTestUser(t, 1, false)

	rec := httptest.NewRecorder()

	// Login with the test user credential to recieve a refresh token
	payload := map[string]string{
		"email":    "test0@test.test",
		"password": "test_password0",
	}
	req := s.CreateTestRequest(t, http.MethodPost, "/auth/login", payload)
	s.ServeHTTP(rec, req)
	res := struct {
		model.AuthToken
		ErrorMsg string `json:"error"`
	}{}
	assert.Equal(t, http.StatusOK, rec.Code)

	// Send access token refresh request
	req = s.CreateTestRequest(t, http.MethodPost, "/auth/refresh-access-token", nil)
	// Append refresh token cookie from the response of the login request
	for _, c := range rec.Result().Cookies() {
		if c.Name == "refresh_token" {
			req.AddCookie(c)
		}
	}

	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	json.NewDecoder(rec.Body).Decode(&res)
	assert.NotEmpty(t, res.TokenString)
	assert.Empty(t, res.ErrorMsg)
}

func TestServer_RefreshAccessToken_RefreshTokenDoesNotExist(t *testing.T) {
	s := httpserver.NewTestServer(t)

	s.CreateTestUser(t, 1, false)

	// Login with the test user credential to recieve a refresh token
	payload := map[string]string{
		"email":    "test0@test.test",
		"password": "test_password0",
	}
	req := s.CreateTestRequest(t, http.MethodPost, "/auth/login", payload)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	res := struct {
		model.AuthToken
		ErrorMsg string `json:"error"`
	}{}
	assert.Equal(t, http.StatusOK, rec.Code)

	// Send access token refresh request
	req = s.CreateTestRequest(t, http.MethodPost, "/auth/refresh-access-token", nil)
	// Append refresh token cookie from the response of the login request
	for _, c := range rec.Result().Cookies() {
		if c.Name == "refresh_token" {
			req.AddCookie(c)
			// Delete refresh token from database to make the validation fail
			s.DeleteTestRefreshToken(t, c)
		}
	}
	rec = httptest.NewRecorder()

	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	json.NewDecoder(rec.Body).Decode(&res)
	assert.Equal(t, "sql: no rows in result set", res.ErrorMsg)
}

func TestServer_RefreshAccessToken_RefreshTokenUserDoesNotExist(t *testing.T) {
	s := httpserver.NewTestServer(t)

	s.CreateTestUser(t, 1, false)

	// Login with the test user credential to recieve a refresh token
	payload := map[string]string{
		"email":    "test0@test.test",
		"password": "test_password0",
	}
	req := s.CreateTestRequest(t, http.MethodPost, "/auth/login", payload)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	res := struct {
		model.AuthToken
		ErrorMsg string `json:"error"`
	}{}
	assert.Equal(t, http.StatusOK, rec.Code)

	// Send access token refresh request
	req = s.CreateTestRequest(t, http.MethodPost, "/auth/refresh-access-token", nil)
	// Append refresh token cookie from the response of the login request
	for _, c := range rec.Result().Cookies() {
		if c.Name == "refresh_token" {
			req.AddCookie(c)
			// Delete refresh token from database to make the validation fail
			s.DeleteTestRefreshTokenUser(t, c)
		}
	}
	rec = httptest.NewRecorder()

	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	json.NewDecoder(rec.Body).Decode(&res)
	assert.Equal(t, "sql: no rows in result set", res.ErrorMsg)
}
