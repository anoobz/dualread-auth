package psqlstore

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/anoobz/dualread/auth/internal/store"
)

type SqlStore struct {
	userRepo      *SqlUserRepo
	authTokenRepo *SqlAuthTokenRepo
}

func NewSqlStore(
	db *sql.DB,
	psql squirrel.StatementBuilderType,
) *SqlStore {
	return &SqlStore{
		userRepo:      NewSqlUserRepo(db, psql),
		authTokenRepo: NewSqlAuthTokenRepo(db, psql),
	}
}

func (s *SqlStore) User() store.UserRepo {
	return s.userRepo
}

func (s *SqlStore) AuthToken() store.AuthTokenRepo {
	return s.authTokenRepo
}
