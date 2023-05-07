package domain

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID         int
	UserName       string
	HashedPassword string
	CreatedAt      time.Time
	DeletedAt      *time.Time
}

func (u *User) VerifyPassword(pwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(pwd)); err != nil {
		return fmt.Errorf("password do not match: %w", err)
	}
	return nil
}

func (u *User) HashPassword(pwd string) error {
	// Hashing the password
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("GenerateFromPassword: %w", err)
	}
	u.HashedPassword = string(hash)

	return nil
}
