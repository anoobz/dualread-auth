package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anoobz/dualread/auth/internal/model"
	"github.com/anoobz/dualread/auth/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestServer_GetAuthToken(t *testing.T) {
	s := NewTestServer(t)

	s.CreateTestUser(t, 1, true)
	user := s.CreateTestUser(t, 1, false)[0]

	token := s.CreateTestToken(t, 1, user)[0]

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		requestedTokenId string
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name: "success",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			requestedTokenId: token.Uuid,
			expectedStatus:   http.StatusOK,
			expectedErrorMsg: "",
		},
		{
			name: "non-admin",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			requestedTokenId: token.Uuid,
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "unauthorized",
		},
		{
			name: "user not found",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			requestedTokenId: "user-not-found",
			expectedStatus:   http.StatusNotFound,
			expectedErrorMsg: "sql: no rows in result set",
		},
		{
			name: "invalid authorization token",
			loginPayload: map[string]string{
				"email":    "unauthorized@test.test",
				"password": "unauthorized",
			},
			requestedTokenId: token.Uuid,
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "token contains an invalid number of segments",
		},
	}

	for _, tc := range testCases {
		accessToken := s.LoginTestUser(
			t,
			tc.loginPayload["email"],
			tc.loginPayload["password"],
		)

		rec := httptest.NewRecorder()
		req := s.CreateTestRequest(
			t, http.MethodGet,
			fmt.Sprintf("/auth/admin/auth-token/%s", tc.requestedTokenId),
			nil,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg == "" {
			token := &model.AuthToken{}
			json.NewDecoder(rec.Body).Decode(&token)
			assert.Equal(t, token, token, tc.name)
		} else {
			res := struct {
				ErrorMsg string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, tc.expectedErrorMsg, res.ErrorMsg, tc.name)
		}
	}
}

func TestServer_GetAllAuthToken(t *testing.T) {
	s := NewTestServer(t)

	s.CreateTestUser(t, 1, true)
	user := s.CreateTestUser(t, 1, false)[0]

	testTokens := s.CreateTestToken(t, 15, user)

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name: "success",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			expectedStatus:   http.StatusOK,
			expectedErrorMsg: "",
		},
		{
			name: "non-admin",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "unauthorized",
		},
		{
			name: "invalid authorization token",
			loginPayload: map[string]string{
				"email":    "unauthorized@test.test",
				"password": "unauthorized",
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "token contains an invalid number of segments",
		},
	}

	for _, tc := range testCases {
		accessToken := s.LoginTestUser(
			t,
			tc.loginPayload["email"],
			tc.loginPayload["password"],
		)

		//Delete the token that was created during login request
		authTokens, err := s.store.AuthToken().GetAll()
		if err != nil {
			t.Fatal(err)
		}
		err = s.store.AuthToken().Delete(authTokens[len(authTokens)-1].Uuid)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		req := s.CreateTestRequest(
			t, http.MethodGet,
			"/auth/admin/auth-token",
			nil,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg == "" {
			tokens := []*model.AuthToken{}
			json.NewDecoder(rec.Body).Decode(&tokens)
			assert.EqualValues(t, testTokens, tokens, tc.name)
		} else {
			res := struct {
				ErrorMsg string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, tc.expectedErrorMsg, res.ErrorMsg, tc.name)
		}
	}
}

func TestServer_GetAuthTokenPage(t *testing.T) {
	s := NewTestServer(t)

	s.CreateTestUser(t, 1, true)
	user := s.CreateTestUser(t, 1, false)[0]

	testTokens := s.CreateTestToken(t, 25, user)

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		page             uint64
		expectedStatus   int
		expectedTokens   []*model.AuthToken
		expectedErrorMsg string
	}{
		{
			name: "get page 1",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			page:             0,
			expectedStatus:   http.StatusOK,
			expectedTokens:   testTokens[:store.PAGE_COUNT],
			expectedErrorMsg: "",
		},
		{
			name: "not full last page",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			page:             1,
			expectedStatus:   http.StatusOK,
			expectedTokens:   testTokens[store.PAGE_COUNT:],
			expectedErrorMsg: "",
		},
		{
			name: "non-admin",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			page:             0,
			expectedStatus:   http.StatusUnauthorized,
			expectedTokens:   nil,
			expectedErrorMsg: "unauthorized",
		},
		{
			name: "invalid authorization token",
			loginPayload: map[string]string{
				"email":    "unauthorized@test.test",
				"password": "unauthorized",
			},
			page:             0,
			expectedStatus:   http.StatusUnauthorized,
			expectedTokens:   nil,
			expectedErrorMsg: "token contains an invalid number of segments",
		},
		{
			name: "page does not exist",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			page:             999,
			expectedStatus:   http.StatusNotFound,
			expectedTokens:   nil,
			expectedErrorMsg: "insufficient token count",
		},
	}

	for _, tc := range testCases {
		accessToken := s.LoginTestUser(
			t,
			tc.loginPayload["email"],
			tc.loginPayload["password"],
		)
		//Delete the token that was created during login request
		authTokens, err := s.store.AuthToken().GetAll()
		if err != nil {
			t.Fatal(err)
		}
		err = s.store.AuthToken().Delete(authTokens[len(authTokens)-1].Uuid)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		req := s.CreateTestRequest(
			t, http.MethodGet,
			fmt.Sprintf("/auth/admin/auth-token-page/%d", tc.page),
			nil,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg == "" {
			tokens := []*model.AuthToken{}
			json.NewDecoder(rec.Body).Decode(&tokens)
			assert.Equal(t, tc.expectedTokens, tokens, tc.name)
		} else {
			res := struct {
				ErrorMsg string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, tc.expectedErrorMsg, res.ErrorMsg, tc.name)
		}
	}
}

func TestServer_DeleteAuthToken(t *testing.T) {
	s := NewTestServer(t)

	user := s.CreateTestUser(t, 1, false)
	s.CreateTestUser(t, 1, true)

	tokens := s.CreateTestToken(t, 2, user[0])

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		deleteTokenId    string
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name: "success",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			deleteTokenId:    tokens[0].Uuid,
			expectedStatus:   http.StatusOK,
			expectedErrorMsg: "",
		},
		{
			name: "non-admin",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			deleteTokenId:    tokens[1].Uuid,
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "unauthorized",
		},
		{
			name: "invalid authorization token",
			loginPayload: map[string]string{
				"email":    "unauthorized@test.test",
				"password": "unauthorized",
			},
			deleteTokenId:    tokens[1].Uuid,
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "token contains an invalid number of segments",
		},
	}

	for _, tc := range testCases {
		accessToken := s.LoginTestUser(
			t,
			tc.loginPayload["email"],
			tc.loginPayload["password"],
		)

		rec := httptest.NewRecorder()
		req := s.CreateTestRequest(
			t, http.MethodDelete,
			fmt.Sprintf("/auth/admin/auth-token/%s", tc.deleteTokenId),
			nil,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg == "" {
			_, err := s.store.AuthToken().GetById(tc.deleteTokenId)
			assert.EqualError(t, err, "sql: no rows in result set", tc.name)
		} else {
			res := struct {
				ErrorMsg string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, tc.expectedErrorMsg, res.ErrorMsg, tc.name)
		}
	}
}
