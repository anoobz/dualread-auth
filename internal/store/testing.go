package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/anoobz/dualread/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func CreateTestUser(t *testing.T, s Store, count int, admin bool) []*model.User {
	t.Helper()

	testTime := GetTestNow(t)
	users, err := s.User().GetAll()
	if err != nil {
		t.Fatal(err)
	}
	userCount := len(users)

	newUsers := []*model.User{}
	for i := userCount; i < count+userCount; i++ {
		encryptedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(fmt.Sprintf("test_password%d", i)),
			10,
		)
		if err != nil {
			t.Fatal(err)
		}
		u, err := s.User().Insert(
			fmt.Sprintf("test%d@test.test", i),
			string(encryptedPassword),
			admin,
			testTime,
		)
		if err != nil {
			t.Fatal(err)
		}
		newUsers = append(newUsers, u)
	}

	return newUsers
}

func CreateTestToken(
	t *testing.T,
	s Store,
	count int,
	user *model.User,
) []*model.AuthToken {
	tokens := []*model.AuthToken{}
	for i := 0; i < count; i++ {
		token, err := model.NewRefreshToken(user)
		if err != nil {
			t.Fatal(err)
		}

		err = s.AuthToken().Insert(token.Uuid, token.TokenString, token.Expires)
		if err != nil {
			t.Fatal(err)
		}
		tokens = append(tokens, token)
	}

	return tokens
}

func GetTestNow(t *testing.T) time.Time {
	t.Helper()
	return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local)
}
