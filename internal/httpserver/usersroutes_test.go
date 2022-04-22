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

func TestServer_GetAllUser(t *testing.T) {
	s := NewTestServer(t)

	user := s.CreateTestUser(t, 1, false)[0]
	admin := s.CreateTestUser(t, 1, true)[0]

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name: "success",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			expectedStatus:   http.StatusOK,
			expectedErrorMsg: "",
		},
		{
			name: "non-admin",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
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

		rec := httptest.NewRecorder()
		req := s.CreateTestRequest(
			t, http.MethodGet,
			"/auth/admin/user",
			nil,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg == "" {
			u := []*model.User{}
			json.NewDecoder(rec.Body).Decode(&u)
			assert.EqualValues(t, []*model.User{user, admin}, u, tc.name)
		} else {
			res := struct {
				ErrorMsg string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, tc.expectedErrorMsg, res.ErrorMsg, tc.name)
		}
	}
}

func TestServer_GetUserPage(t *testing.T) {
	s := NewTestServer(t)

	users := s.CreateTestUser(t, 25, false)
	admin := s.CreateTestUser(t, 1, true)

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		page             uint64
		expectedStatus   int
		expectedUsers    []*model.User
		expectedErrorMsg string
	}{
		{
			name: "get page 1",
			loginPayload: map[string]string{
				"email":    "test25@test.test",
				"password": "test_password25",
			},
			page:             0,
			expectedStatus:   http.StatusOK,
			expectedUsers:    users[:store.PAGE_COUNT],
			expectedErrorMsg: "",
		},
		{
			name: "not full last page",
			loginPayload: map[string]string{
				"email":    "test25@test.test",
				"password": "test_password25",
			},
			page:             1,
			expectedStatus:   http.StatusOK,
			expectedUsers:    append(users[store.PAGE_COUNT:], admin[0]),
			expectedErrorMsg: "",
		},
		{
			name: "non-admin",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			page:             0,
			expectedStatus:   http.StatusUnauthorized,
			expectedUsers:    nil,
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
			expectedUsers:    nil,
			expectedErrorMsg: "token contains an invalid number of segments",
		},
		{
			name: "page does not exist",
			loginPayload: map[string]string{
				"email":    "test25@test.test",
				"password": "test_password25",
			},
			page:             999,
			expectedStatus:   http.StatusNotFound,
			expectedUsers:    nil,
			expectedErrorMsg: "insufficient user count",
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
			fmt.Sprintf("/auth/admin/user-page/%d", tc.page),
			nil,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg == "" {
			u := []*model.User{}
			json.NewDecoder(rec.Body).Decode(&u)
			assert.EqualValues(t, tc.expectedUsers, u, tc.name)
		} else {
			res := struct {
				ErrorMsg string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, tc.expectedErrorMsg, res.ErrorMsg, tc.name)
		}
	}
}

func TestServer_InsertUser(t *testing.T) {
	s := NewTestServer(t)

	s.CreateTestUser(t, 1, false)
	s.CreateTestUser(t, 1, true)

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		insertPayload    map[string]interface{}
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name: "success",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			insertPayload: map[string]interface{}{
				"email":    "inserted0@test.test",
				"password": "inserted_password0",
				"admin":    false,
			},
			expectedStatus:   http.StatusCreated,
			expectedErrorMsg: "",
		},
		{
			name: "non-admin",
			loginPayload: map[string]string{
				"email":    "test0@test.test",
				"password": "test_password0",
			},
			insertPayload: map[string]interface{}{
				"email":    "inserted1@test.test",
				"password": "inserted_password1",
				"admin":    false,
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
			insertPayload: map[string]interface{}{
				"email":    "inserted1@test.test",
				"password": "inserted_password1",
				"admin":    false,
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "token contains an invalid number of segments",
		},
		{
			name: "empty email",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			insertPayload: map[string]interface{}{
				"email":    "",
				"password": "inserted_password1",
				"admin":    false,
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorMsg: "mail: no address",
		},
		{
			name: "invalid email",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			insertPayload: map[string]interface{}{
				"email":    "invalid",
				"password": "inserted_password1",
				"admin":    false,
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorMsg: "mail: missing '@' or angle-addr",
		},
		{
			name: "empty password",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			insertPayload: map[string]interface{}{
				"email":    "inserted1@test.test",
				"password": "",
				"admin":    false,
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorMsg: "a required field is empty",
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
			t, http.MethodPost,
			"/auth/admin/user",
			tc.insertPayload,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg == "" {
			u := &model.User{}
			json.NewDecoder(rec.Body).Decode(&u)
			assert.Equal(t, tc.insertPayload["email"], u.Email, tc.name)
			assert.Equal(t, tc.insertPayload["password"], u.Password, tc.name)
			assert.Equal(t, tc.insertPayload["admin"], u.Admin, tc.name)
		} else {
			res := struct {
				ErrorMsg string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, tc.expectedErrorMsg, res.ErrorMsg, tc.name)
		}
	}
}

func TestServer_UpdateUser(t *testing.T) {
	s := NewTestServer(t)

	user := s.CreateTestUser(t, 2, false)
	s.CreateTestUser(t, 1, true)

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		updateUSerId     int64
		clauses          map[string]interface{}
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name: "success",
			loginPayload: map[string]string{
				"email":    "test2@test.test",
				"password": "test_password2",
			},
			updateUSerId: user[0].ID,
			clauses: map[string]interface{}{
				"email":            "new@test.test",
				"password":         "newPassword",
				"active":           false,
				"email_verified":   true,
				"email_subscribed": false,
				"admin":            true,
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
			updateUSerId: user[1].ID,
			clauses: map[string]interface{}{
				"email":            "new@test.test",
				"password":         "newPassword",
				"active":           false,
				"email_verified":   true,
				"email_subscribed": false,
				"admin":            true,
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
			updateUSerId: user[1].ID,
			clauses: map[string]interface{}{
				"email":            "new@test.test",
				"password":         "newPassword",
				"active":           false,
				"email_verified":   true,
				"email_subscribed": false,
				"admin":            true,
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "token contains an invalid number of segments",
		},
		{
			name: "empty email",
			loginPayload: map[string]string{
				"email":    "test2@test.test",
				"password": "test_password2",
			},
			updateUSerId: user[1].ID,
			clauses: map[string]interface{}{
				"email":            "",
				"password":         "newPassword",
				"active":           false,
				"email_verified":   true,
				"email_subscribed": false,
				"admin":            true,
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorMsg: "mail: no address",
		},
		{
			name: "invalid email",
			loginPayload: map[string]string{
				"email":    "test2@test.test",
				"password": "test_password2",
			},
			updateUSerId: user[1].ID,
			clauses: map[string]interface{}{
				"email":            "invalid",
				"password":         "newPassword",
				"active":           false,
				"email_verified":   true,
				"email_subscribed": false,
				"admin":            true,
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorMsg: "mail: missing '@' or angle-addr",
		},
		{
			name: "empty password",
			loginPayload: map[string]string{
				"email":    "test2@test.test",
				"password": "test_password2",
			},
			updateUSerId: user[1].ID,
			clauses: map[string]interface{}{
				"email":            "new@test.test",
				"password":         "",
				"active":           false,
				"email_verified":   true,
				"email_subscribed": false,
				"admin":            true,
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorMsg: "password is empty",
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
			t, http.MethodPost,
			fmt.Sprintf("/auth/admin/user/%d", tc.updateUSerId),
			tc.clauses,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg != "" {
			res := struct {
				ErrorMsg string `json:"error"`
			}{}
			json.NewDecoder(rec.Body).Decode(&res)
			assert.Equal(t, tc.expectedErrorMsg, res.ErrorMsg, tc.name)
		}
	}
}

func TestServer_DeleteUser(t *testing.T) {
	s := NewTestServer(t)

	user := s.CreateTestUser(t, 3, false)
	s.CreateTestUser(t, 1, true)

	testCases := []struct {
		name             string
		loginPayload     map[string]string
		deleteUserId     int64
		expectedStatus   int
		expectedErrorMsg string
	}{
		{
			name: "success",
			loginPayload: map[string]string{
				"email":    "test3@test.test",
				"password": "test_password3",
			},
			deleteUserId:     user[2].ID,
			expectedStatus:   http.StatusOK,
			expectedErrorMsg: "",
		},
		{
			name: "non-admin",
			loginPayload: map[string]string{
				"email":    "test1@test.test",
				"password": "test_password1",
			},
			deleteUserId:     user[0].ID,
			expectedStatus:   http.StatusUnauthorized,
			expectedErrorMsg: "unauthorized",
		},
		{
			name: "invalid authorization token",
			loginPayload: map[string]string{
				"email":    "unauthorized@test.test",
				"password": "unauthorized",
			},
			deleteUserId:     user[0].ID,
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
			fmt.Sprintf("/auth/admin/user/%d", tc.deleteUserId),
			nil,
		)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		s.ServeHTTP(rec, req)

		assert.Equal(t, tc.expectedStatus, rec.Code, tc.name)
		if tc.expectedErrorMsg == "" {
			_, err := s.store.User().GetById(tc.deleteUserId)
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
