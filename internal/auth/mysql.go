package auth

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// UserRepositoryMysl is an implementation of UserRepository using a Mysl database as storage.
type UserRepositoryMysl struct {
	db *gorm.DB
}

func NewUserRepositoryMysl(db *gorm.DB) (*UserRepositoryMysl, error) {
	if db == nil {
		return nil, errors.New("cannot instantiate new user repository Mysl: no db provided")
	}

	return &UserRepositoryMysl{db: db}, nil
}

func (r *UserRepositoryMysl) Create(user *User) error {
	return errors.Wrap(r.db.Create(user).Error, "cannot create user")
}

func (r *UserRepositoryMysl) GetUserByName(name string) (*User, error) {
	user := &User{}
	err := r.db.Where("name = ?", name).First(user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Wrap(err, "cannot get user")
	}

	return user, nil
}
