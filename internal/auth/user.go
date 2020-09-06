package auth

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 10
)

// User is a user of our application.
// It contains its team ID.
type User struct {
	Name         string `gorm:"primaryKey"`
	PasswordHash string `json:"-"`
	TeamID       uuid.UUID
}

var UserNotFound = errors.New("user not found")

// UserRepository is a service which will handle the storage of the users.
type UserRepository interface {
	Create(*User) error
	GetUserByName(name string) (*User, error)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
