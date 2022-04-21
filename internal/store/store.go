package store

import (
	"time"

	"github.com/anoobz/dualread/auth/internal/model"
)

const (
	PAGE_COUNT = 20
)

type UserRepo interface {
	GetById(id int64) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetAll() ([]*model.User, error)
	GetPage(page uint64) ([]*model.User, error)
	Insert(
		email string,
		password string,
		admin bool,
		now time.Time,
	) (*model.User, error)
	Update(id int64, clauses map[string]interface{}) error
	Delete(id int64) error
}

type AuthTokenRepo interface {
	GetById(id string) (*model.AuthToken, error)
	GetAll() ([]*model.AuthToken, error)
	GetPage(pageId uint64) ([]*model.AuthToken, error)
	Insert(
		uuid string,
		tokenString string,
		expires int64,
	) error
	Delete(id string) error
}

type Store interface {
	User() UserRepo
	AuthToken() AuthTokenRepo
}
