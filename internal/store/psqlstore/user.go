package psqlstore

import (
	"database/sql"
	"errors"
	"net/mail"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/anoobz/dualread/auth/internal/model"
	"github.com/anoobz/dualread/auth/internal/store"
)

type SqlUserRepo struct {
	db   *sql.DB
	psql squirrel.StatementBuilderType
}

func NewSqlUserRepo(db *sql.DB, psql squirrel.StatementBuilderType) *SqlUserRepo {
	return &SqlUserRepo{
		db:   db,
		psql: psql,
	}
}

func (r *SqlUserRepo) GetById(id int64) (*model.User, error) {
	row := r.psql.Select("*").From("users").
		Where("id = ?", id).QueryRow()
	u, err := userFromRow(row)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *SqlUserRepo) GetAll() ([]*model.User, error) {
	rows, err := r.psql.Select("*").From("users").Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*model.User{}
	for rows.Next() {
		u, err := userFromRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *SqlUserRepo) GetByEmail(email string) (*model.User, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}

	row := r.psql.Select("*").From("users").
		Where("email = ?", email).QueryRow()
	u, err := userFromRow(row)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *SqlUserRepo) GetPage(page uint64) ([]*model.User, error) {
	rows, err := r.psql.Select("*").From("users").Limit(store.PAGE_COUNT).Offset(page * store.PAGE_COUNT).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*model.User{}
	for rows.Next() {
		u, err := userFromRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		return nil, errors.New("insufficient user count")
	}

	return users, nil
}

func (r *SqlUserRepo) Insert(
	email string,
	password string,
	admin bool,
	now time.Time,
) (*model.User, error) {
	u, err := model.NewUser(email, password, admin, now)
	if err != nil {
		return nil, err
	}

	if err := r.psql.Insert("users").
		Columns("email",
			"password",
			"active",
			"email_verified",
			"email_subscribed",
			"admin",
			"created",
			"last_login",
			"last_action").
		Values(u.Email,
			u.Password,
			u.Active,
			u.EmailVerified,
			u.EmailSubscribed,
			u.Admin,
			u.Created,
			u.LastLogin,
			u.LastAction).
		Suffix("RETURNING ID").
		QueryRow().
		Scan(&u.ID); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *SqlUserRepo) Update(id int64, clauses map[string]interface{}) error {
	email, ok := clauses["email"]
	if ok {
		_, err := mail.ParseAddress(email.(string))
		if err != nil {
			return err
		}
	}
	password, ok := clauses["password"]
	if ok {
		if password == "" {
			return errors.New("password is empty")
		}
	}

	res, err := r.psql.Update("users").
		SetMap(clauses).
		Where("id = ?", id).
		Exec()
	if err != nil {
		return err
	}

	updatedRowCount, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if updatedRowCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *SqlUserRepo) Delete(id int64) error {
	_, err := r.psql.Delete("users").Where("id = ?", id).Exec()
	if err != nil {
		return err
	}

	return nil
}

func userFromRow(row store.Row) (*model.User, error) {
	u := &model.User{}
	created := &time.Time{}
	lastLogin := &time.Time{}
	lastAction := &time.Time{}
	if err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Active,
		&u.EmailVerified,
		&u.EmailSubscribed,
		&u.Admin,
		&created,
		&lastLogin,
		&lastAction,
	); err != nil {
		return nil, err
	}

	u.Created = created.Local()
	u.LastLogin = lastLogin.Local()
	u.LastAction = lastAction.Local()

	return u, nil
}
