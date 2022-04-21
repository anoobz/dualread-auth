package mockstore

import (
	"errors"
	"net/mail"
	"time"

	"github.com/anoobz/dualread/auth/internal/model"
	"github.com/anoobz/dualread/auth/internal/store"
)

type MockUserRepo struct {
	users []*model.User
}

func (r *MockUserRepo) GetById(id int64) (*model.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, errors.New("sql: no rows in result set")
}

func (r *MockUserRepo) GetByEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("mail: no address")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, errors.New("sql: no rows in result set")
}

func (r *MockUserRepo) GetAll() ([]*model.User, error) {
	return r.users, nil
}

func (r *MockUserRepo) GetPage(page uint64) ([]*model.User, error) {
	if len(r.users) <= int(page*store.PAGE_COUNT) {
		return nil, errors.New("insufficient user count")
	}

	users := []*model.User{}
	for i := page * store.PAGE_COUNT; i < (page+1)*store.PAGE_COUNT; i++ {
		if i >= uint64(len(r.users)) {
			break
		}
		users = append(users, r.users[i])
	}

	return users, nil
}

func (r *MockUserRepo) Insert(
	email string,
	password string,
	admin bool,
	now time.Time,
) (*model.User, error) {
	u, err := model.NewUser(email, password, admin, now)
	if err != nil {
		return nil, err
	}

	u.ID = int64(len(r.users) + 1)
	r.users = append(r.users, u)

	return u, nil
}

func (r *MockUserRepo) Update(id int64, clauses map[string]interface{}) error {
	for i, u := range r.users {
		if u.ID == id {
			updatedUser := r.users[i]
			for key, value := range clauses {
				switch key {
				case "email":
					_, err := mail.ParseAddress(value.(string))
					if err != nil {
						return err
					}
					updatedUser.Email = value.(string)
				case "password":
					if value == "" {
						return errors.New("password is empty")
					}
					updatedUser.Password = value.(string)
				case "active":
					updatedUser.Active = value.(bool)
				case "email_verified":
					updatedUser.EmailVerified = value.(bool)
				case "email_subscribed":
					updatedUser.EmailSubscribed = value.(bool)
				case "admin":
					updatedUser.Admin = value.(bool)
				}
			}
			return nil
		}
	}

	return errors.New("user not found")
}

func (r *MockUserRepo) Delete(id int64) error {
	for i, u := range r.users {
		if u.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}

	return errors.New("")
}
