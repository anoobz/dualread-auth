package mockstore

import "github.com/anoobz/dualread/auth/internal/store"

type MockStore struct {
	userRepo      *MockUserRepo
	authTokenRepo *MockAuthTokenRepo
}

func NewMockStore() *MockStore {
	return &MockStore{
		userRepo:      &MockUserRepo{},
		authTokenRepo: &MockAuthTokenRepo{},
	}
}

func (s *MockStore) User() store.UserRepo {
	return s.userRepo
}

func (s *MockStore) AuthToken() store.AuthTokenRepo {
	return s.authTokenRepo
}
