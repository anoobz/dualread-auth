package store

import (
	"testing"

	"github.com/anoobz/dualread/auth/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestStore_InsertRefreshToken(t *testing.T, s Store) {
	testUser := CreateTestUser(t, s, 1, false)[0]
	testToken, err := model.NewRefreshToken(testUser)
	if err != nil {
		t.Fatal(err)
	}
	err = s.AuthToken().Insert(testToken.Uuid, testToken.TokenString, testToken.Expires)
	if err != nil {
		t.Fatal(err)
	}

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

func TestStore_GetTokenPage(t *testing.T, s Store) {
	testUser := CreateTestUser(t, s, 1, false)[0]
	testTokens := CreateTestToken(t, s, 25, testUser)

	testCases := []struct {
		name             string
		pageId           uint64
		expectedErrorMsg string
		expectedTokens   []*model.AuthToken
	}{
		{
			name:             "first page",
			pageId:           0,
			expectedErrorMsg: "",
			expectedTokens:   testTokens[:1*PAGE_COUNT],
		},
		{
			name:             "not full last page",
			pageId:           1,
			expectedErrorMsg: "",
			expectedTokens:   testTokens[1*PAGE_COUNT:],
		},
		{
			name:             "invalid page",
			pageId:           999,
			expectedErrorMsg: "insufficient token count",
			expectedTokens:   nil,
		},
	}

	for _, tc := range testCases {
		tokens, err := s.AuthToken().GetPage(tc.pageId)

		if tc.expectedErrorMsg == "" {
			assert.Equal(t, tc.expectedTokens, tokens, tc.name)
		} else {
			assert.Equal(t, tc.expectedErrorMsg, err.Error(), tc.name)
		}
	}
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
