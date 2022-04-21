package mockstore

import (
	"testing"

	"github.com/anoobz/dualread/auth/internal/store"
)

func TestStore_InsertRefreshToken(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_InsertRefreshToken(t, s)
}

func TestStore_GetAllToken(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_GetAllToken(t, s)
}

func TestStore_GetTokenPage(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_GetTokenPage(t, s)
}

func TestStore_DeleteToken(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_DeleteToken(t, s)
}
