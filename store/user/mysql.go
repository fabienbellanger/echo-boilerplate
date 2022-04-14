package user

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/fabienbellanger/echo-boilerplate/db"
	"github.com/fabienbellanger/echo-boilerplate/entities"
	"github.com/google/uuid"
)

// UserStore ...
type UserStore struct {
	db *db.DB
}

// New returns a new UserStore
func New(db *db.DB) UserStore {
	return UserStore{db: db}
}

// Login authenticate a user
func (u UserStore) Login(username, password string) (user entities.User, err error) {
	// Hash password
	// -------------
	passwordBytes := sha512.Sum512([]byte(password))
	password = hex.EncodeToString(passwordBytes[:])

	if result := u.db.Where(&entities.User{Username: username, Password: password}).First(&user); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// Register creates a new user in database
func (u UserStore) Register(user *entities.User) error {
	// UUID
	// ----
	user.ID = uuid.New().String()

	// Hash password
	// -------------
	passwordBytes := sha512.Sum512([]byte(user.Password))
	user.Password = hex.EncodeToString(passwordBytes[:])

	if result := u.db.Create(&user); result.Error != nil {
		return result.Error
	}
	return nil
}
