package mockstore

import (
	"testing"

	"github.com/anoobz/dualread/auth/internal/store"
)

func TestStore_GetAllUsers(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_GetAllUsers(t, s)
}

func TestStore_GetUserById(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_GetUserById(t, s)
}

func TestStore_GetUserByEmail(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_GetUserByEmail(t, s)
}

func TestStore_CreateUser(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_CreateUser(t, s)
}

func TestStore_UpdateUser(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_UpdateUser(t, s)
}

func TestStore_DeleteUser(t *testing.T) {
	s := CreateTestStore(t)

	store.TestStore_DeleteUser(t, s)
}
