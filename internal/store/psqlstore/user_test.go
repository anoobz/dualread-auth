package psqlstore

import (
	"testing"

	"github.com/anoobz/dualread/auth/internal/store"
)

func TestStore_GetAllUsers(t *testing.T) {
	s, dbTearUp := CreateTestStore(t)
	defer dbTearUp("users")

	store.TestStore_GetAllUsers(t, s)
}

func TestStore_GetUserById(t *testing.T) {
	s, dbTearUp := CreateTestStore(t)
	defer dbTearUp("users")

	store.TestStore_GetUserById(t, s)
}

func TestStore_GetUserByEmail(t *testing.T) {
	s, dbTearUp := CreateTestStore(t)
	defer dbTearUp("users")

	store.TestStore_GetUserByEmail(t, s)
}

func TestStore_CreateUser(t *testing.T) {
	s, dbTearUp := CreateTestStore(t)
	defer dbTearUp("users")

	store.TestStore_CreateUser(t, s)
}

func TestStore_UpdateUser(t *testing.T) {
	s, dbTearUp := CreateTestStore(t)
	defer dbTearUp("users")

	store.TestStore_UpdateUser(t, s)
}

func TestStore_DeleteUser(t *testing.T) {
	s, dbTearUp := CreateTestStore(t)
	defer dbTearUp("users")

	store.TestStore_DeleteUser(t, s)
}
