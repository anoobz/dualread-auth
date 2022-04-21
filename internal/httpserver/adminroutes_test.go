package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anoobz/dualread/auth/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestServer_GetUser(t *testing.T) {
	s := NewTestServer(t)

	user := s.CreateTestUser(t, 1, false)[0]
	s.CreateTestUser(t, 1, true)

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		requestedUserId  int64
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name: "success",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			requestedUserId:  user.ID,
			expectedStatus:   http.StatusOK,
			expectedErrorMsg: "",
		},
		{
			name: "non-admin",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			requestedUserId:  user.ID,
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "unauthorized",
		},
		{
			name: "user not found",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			requestedUserId:  9999,
			expectedStatus:   http.StatusNotFound,
			expectedErrorMsg: "sql: no rows in result set",
		},
		{
			name: "invalid authorization token",
			loginPayload: map[string]string{
				"email":    "unauthorized@test.test",
				"password": "unauthorized",
			},
			requestedUserId:  user.ID,
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
			fmt.Sprintf("/auth/admin/user/%d", tc.requestedUserId),
			nil,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg == "" {
			u := &model.User{}
			json.NewDecoder(rec.Body).Decode(&u)
			assert.Equal(t, user, u, tc.name)
		} else {
			res := struct {
				ErrorMsg string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, tc.expectedErrorMsg, res.ErrorMsg, tc.name)
		}
	}
}
