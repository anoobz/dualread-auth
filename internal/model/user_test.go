package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	type Credential struct {
		Email    string
		Password string
	}
	testCases := []struct {
		name          string
		credentials   Credential
		expectedError string
	}{
		{
			name: "valid",
			credentials: Credential{
				Email:    "test@test.test",
				Password: "test_password",
			},
			expectedError: "",
		},
		{
			name: "empty email",
			credentials: Credential{
				Email:    "",
				Password: "test_password",
			},
			expectedError: "mail: no address",
		},
		{
			name: "empty password",
			credentials: Credential{
				Email:    "test@test.test",
				Password: "",
			},
			expectedError: "a required field is empty",
		},
		{
			name: "invalid email",
			credentials: Credential{
				Email:    "invalid",
				Password: "test_password",
			},
			expectedError: "mail: missing '@' or angle-addr",
		},
	}

	for _, tc := range testCases {
		user, err := NewUser(
			tc.credentials.Email,
			tc.credentials.Password,
			false,
			time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local),
		)

		if tc.expectedError == "" {
			assert.NoError(t, err)
			assert.Equal(t, tc.credentials.Email, user.Email)
			assert.Equal(t, tc.credentials.Password, user.Password)
		} else {
			assert.EqualError(t, err, tc.expectedError)
		}
	}
}
