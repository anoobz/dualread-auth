package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_InsertRefreshToken(t *testing.T, s Store) {

	testUser := CreateTestUser(t, s, 1, false)[0]
	testToken := CreateTestToken(t, s, 1, testUser)[0]

	insertedToken, err := s.AuthToken().GetById(testToken.Uuid)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testToken, insertedToken)
}

func TestStore_GetAllToken(t *testing.T, s Store) {
	t.Helper()

	testUser := CreateTestUser(t, s, 1, false)[0]
	testTokens := CreateTestToken(t, s, 5, testUser)

	insertedTokens, err := s.AuthToken().GetAll()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testTokens, insertedTokens)
}

func TestStore_DeleteToken(t *testing.T, s Store) {
	t.Helper()

	testUser := CreateTestUser(t, s, 1, false)[0]
	testTokenToRemove := CreateTestToken(t, s, 5, testUser)[1]

	err := s.AuthToken().Delete(testTokenToRemove.Uuid)
	if err != nil {
		t.Fatal(err)
	}

	tokens, err := s.AuthToken().GetAll()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotContains(t, tokens, testTokenToRemove)
}
