package psqlstore

import (
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/anoobz/dualread/auth/internal/model"
	"github.com/anoobz/dualread/auth/internal/store"
)

type SqlAuthTokenRepo struct {
	db   *sql.DB
	psql squirrel.StatementBuilderType
}

func NewSqlAuthTokenRepo(
	db *sql.DB,
	psql squirrel.StatementBuilderType,
) *SqlAuthTokenRepo {
	return &SqlAuthTokenRepo{
		db:   db,
		psql: psql,
	}
}

func (r *SqlAuthTokenRepo) Insert(
	uuid string,
	tokenString string,
	expires int64,
) error {
	_, err := r.psql.Insert("refresh_token").
		Columns("id", "token_string", "expires").
		Values(uuid, tokenString, expires).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

func (r *SqlAuthTokenRepo) GetAll() ([]*model.AuthToken, error) {
	rows, err := r.psql.Select("*").From("refresh_token").Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := []*model.AuthToken{}
	for rows.Next() {
		t, err := tokenFromRow(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}

	return tokens, nil
}

func (r *SqlAuthTokenRepo) GetPage(page uint64) ([]*model.AuthToken, error) {
	rows, err := r.psql.Select("*").
		From("refresh_token").
		Limit(store.PAGE_COUNT).
		Offset(page * store.PAGE_COUNT).
		Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := []*model.AuthToken{}
	for rows.Next() {
		t, err := tokenFromRow(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}

	if len(tokens) == 0 {
		return nil, errors.New("insufficient token count")
	}

	return tokens, nil
}

func (r *SqlAuthTokenRepo) GetById(id string) (*model.AuthToken, error) {
	row := r.psql.Select("*").From("refresh_token").Where("id = ?", id).QueryRow()
	t, err := tokenFromRow(row)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *SqlAuthTokenRepo) Delete(id string) error {
	_, err := r.psql.Delete("refresh_token").Where("id = ?", id).Exec()
	if err != nil {
		return err
	}

	return nil
}

func tokenFromRow(row store.Row) (*model.AuthToken, error) {
	t := &model.AuthToken{}
	if err := row.Scan(&t.Uuid, &t.TokenString, &t.Expires); err != nil {
		return nil, err
	}

	return t, nil
}
