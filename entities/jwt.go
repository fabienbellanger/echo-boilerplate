package entities

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// TODO: Add unit tests

// Claims are custom claims extending default ones
type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	jwt.StandardClaims
}

// NewClaims creates a new Claims
func NewClaims(id, username, firstname, lastname string, lifetime int) *Claims {
	return &Claims{
		id,
		username,
		lastname,
		firstname,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(lifetime)).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   "API authentication", // Subject of the JWT (the user)
			Issuer:    "API",                // Issuer of the JWT
			Audience:  "Client",             // Recipient for which the JWT is intended
		},
	}
}

// GenerateJWT generates token
func (c Claims) GenerateJWT(algo, secret string) (string, error) {
	if algo == "HS512" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
		return token.SignedString([]byte(secret))
	}
	return "", errors.New("unsupported JWT algo")
}
