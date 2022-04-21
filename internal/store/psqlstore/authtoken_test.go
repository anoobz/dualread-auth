package psqlstore

import (
	"testing"

	"github.com/anoobz/dualread/auth/internal/store"
)

func TestRepository_InsertRefreshToken(t *testing.T) {
	s, dbTearUp := CreateTestStore(t)
	defer dbTearUp("refresh_token", "users")

	store.TestStore_InsertRefreshToken(t, s)
}

func TestStore_GetAllToken(t *testing.T) {
	s, dbTearUp := CreateTestStore(t)
	defer dbTearUp("refresh_token", "users")

	store.TestStore_GetAllToken(t, s)
}

func TestStore_DeleteToken(t *testing.T) {
	s, dbTearUp := CreateTestStore(t)
	defer dbTearUp("refresh_token", "users")

	store.TestStore_DeleteToken(t, s)
}
