package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_GetAllUsers(t *testing.T, s Store) {
	test_users := CreateTestUser(t, s, 5, false)

	users, err := s.User().GetAll()
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, test_users, users)
}

func TestStore_GetUserById(t *testing.T, s Store) {

	test_user := CreateTestUser(t, s, 3, false)[1]

	retrievedUser, err := s.User().GetById(test_user.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, test_user, retrievedUser)
}

func TestStore_GetUserByEmail(t *testing.T, s Store) {

	test_user := CreateTestUser(t, s, 3, false)[1]

	testCases := []struct {
		name     string
		email    string
		errorMsg string
	}{
		{
			name:     "success",
			email:    test_user.Email,
			errorMsg: "",
		},
		{
			name:     "empty email",
			email:    "",
			errorMsg: "mail: no address",
		},
		{
			name:     "invalid email",
			email:    "invalid",
			errorMsg: "mail: missing '@' or angle-addr",
		},
	}

	for _, tc := range testCases {
		retrievedUser, err := s.User().GetByEmail(tc.email)

		if tc.errorMsg == "" {
			assert.NoError(t, err)
			assert.EqualValues(t, test_user, retrievedUser)
		} else {
			assert.Equal(t, tc.errorMsg, err.Error())
		}
	}
}

func TestStore_CreateUser(t *testing.T, s Store) {

	testCases := []struct {
		name                string
		email               string
		password            string
		expectedErrorString string
	}{
		{
			name:                "valid",
			email:               "test1@test.test",
			password:            "test_password",
			expectedErrorString: "",
		},
		{
			name:                "missing email",
			email:               "",
			password:            "test_password",
			expectedErrorString: "mail: no address",
		},
		{
			name:                "missing password",
			email:               "test@test.test",
			password:            "",
			expectedErrorString: "a required field is empty",
		},
		{
			name:                "invalid email",
			email:               "invalid_email",
			password:            "test_password",
			expectedErrorString: "mail: missing '@' or angle-addr",
		},
	}
	for _, tc := range testCases {
		u, err := s.User().Insert(tc.email, tc.password, false, GetTestNow(t))
		if tc.expectedErrorString == "" {
			assert.Equal(t, tc.email, u.Email)
			assert.Equal(t, tc.password, u.Password)

		} else {
			assert.EqualError(t, err, tc.expectedErrorString)
		}
	}
}

func TestStore_UpdateUser(t *testing.T, s Store) {
	test_user := CreateTestUser(t, s, 1, false)[0]

	testCases := []struct {
		name                string
		userId              int64
		clauses             map[string]interface{}
		expectedErrorString string
	}{
		{
			name:   "update all fields",
			userId: test_user.ID,
			clauses: map[string]interface{}{
				"email":            "new@test.test",
				"password":         "newPassword",
				"active":           false,
				"email_verified":   true,
				"email_subscribed": false,
				"admin":            true,
			},
			expectedErrorString: "",
		},
		{
			name:   "invalid email",
			userId: test_user.ID,
			clauses: map[string]interface{}{
				"email": "invalid",
			},
			expectedErrorString: "mail: missing '@' or angle-addr",
		},
		{
			name:   "empty email",
			userId: test_user.ID,
			clauses: map[string]interface{}{
				"email": "",
			},
			expectedErrorString: "mail: no address",
		},
		{
			name:   "empty password",
			userId: test_user.ID,
			clauses: map[string]interface{}{
				"password": "",
			},
			expectedErrorString: "password is empty",
		},
		{
			name:   "user not found",
			userId: -1,
			clauses: map[string]interface{}{
				"email": "new@test.test",
			},
			expectedErrorString: "user not found",
		},
	}

	for _, tc := range testCases {
		err := s.User().Update(tc.userId, tc.clauses)
		if tc.expectedErrorString == "" {
			assert.NoError(t, err)
			u, err := s.User().GetById(test_user.ID)
			if err != nil {
				t.Fatal(err)
			}
			for key, value := range tc.clauses {
				switch key {
				case "email":
					assert.Equal(t, value, u.Email)
				case "password":
					assert.Equal(t, value, u.Password)
				case "active":
					assert.Equal(t, value, u.Active)
				case "email_verified":
					assert.Equal(t, value, u.EmailVerified)
				case "email_subscribed":
					assert.Equal(t, value, u.EmailSubscribed)
				case "admin":
					assert.Equal(t, value, u.Admin)
				}
			}

		} else {
			assert.EqualError(t, err, tc.expectedErrorString, tc.name)
		}
	}
}

func TestStore_DeleteUser(t *testing.T, s Store) {
	test_user := CreateTestUser(t, s, 1, false)[0]

	err := s.User().Delete(test_user.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.User().GetById(test_user.ID)
	assert.Error(t, err, "sql: no rows in result set")
}
