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

// GetAllUsers lists all users.
func (u UserStore) GetAllUsers() ([]entities.User, error) {
	var users []entities.User

	if response := u.db.Find(&users); response.Error != nil {
		return users, response.Error
	}
	return users, nil
}

// GetUser returns a user from its ID.
func (u UserStore) GetUser(id string) (user entities.User, err error) {
	if result := u.db.Find(&user, "id = ?", id); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// DeleteUser deletes a user from database.
func (u UserStore) DeleteUser(id string) error {
	result := u.db.Delete(&entities.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateUser updates user information.
func (u UserStore) UpdateUser(id string, userForm *entities.UserForm) (user entities.User, err error) {
	// Hash password
	// -------------
	hashedPassword := sha512.Sum512([]byte(userForm.Password))

	result := u.db.Model(&entities.User{}).Where("id = ?", id).Select("lastname", "firstname", "username", "password").Updates(entities.User{
		Lastname:  userForm.Lastname,
		Firstname: userForm.Firstname,
		Username:  userForm.Username,
		Password:  hex.EncodeToString(hashedPassword[:]),
	})
	if result.Error != nil {
		return user, result.Error
	}

	user, err = u.GetUser(id)
	if err != nil {
		return user, err
	}
	return user, err
}
