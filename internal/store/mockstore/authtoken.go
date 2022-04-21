package mockstore

import (
	"errors"

	"github.com/anoobz/dualread/auth/internal/model"
)

type MockAuthTokenRepo struct {
	authTokens []*model.AuthToken
}

func (r *MockAuthTokenRepo) Insert(
	uuid string,
	tokenString string,
	expires int64,
) error {
	t := &model.AuthToken{
		Uuid:        uuid,
		TokenString: tokenString,
		Expires:     expires,
	}
	r.authTokens = append(r.authTokens, t)
	return nil
}

func (r *MockAuthTokenRepo) GetAll() ([]*model.AuthToken, error) {
	return r.authTokens, nil
}

func (r *MockAuthTokenRepo) GetById(id string) (*model.AuthToken, error) {
	for _, t := range r.authTokens {
		if t.Uuid == id {
			return t, nil
		}
	}

	return nil, errors.New("sql: no rows in result set")
}

func (r *MockAuthTokenRepo) Delete(id string) error {
	for i, t := range r.authTokens {
		if t.Uuid == id {
			r.authTokens = append(r.authTokens[:i], r.authTokens[i+1:]...)
			return nil
		}
	}

	return errors.New("")
}
