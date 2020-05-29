package models

import "github.com/dgrijalva/jwt-go"

// User definition for the application, including
// a user id, email, password hash, and more.
type User struct {
	UserID       uint64 `json:"user_id,omitempty"`
	Email        string `json:"email,omitempty"`
	PasswordHash []byte `json:"password_hash,omitempty"`
}

// Claims defines the JWT claims for the application.
// These values will be encoded into issued JWT tokens.
type Claims struct {
	UserID uint64 `json:"user_id,omitempty"`
	jwt.StandardClaims
}
