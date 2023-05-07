package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUser_VerifyPassword(t *testing.T) {
	// Create a user with a known password hash
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &User{
		HashedPassword: string(hashedPassword),
	}

	// Test correct password
	err := user.VerifyPassword("password123")
	assert.NoError(t, err, "VerifyPassword returned an error for correct password")

	// Test incorrect password
	err = user.VerifyPassword("wrongpassword")
	assert.Error(t, err, "VerifyPassword did not return an error for incorrect password")
}

func TestUser_HashPassword(t *testing.T) {
	// Create a user
	user := &User{}

	// Test hashing a password
	password := "password123"
	err := user.HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword returned an error: %v", err)
	}

	// Verify that the hashed password is not empty
	assert.NotEmpty(t, user.HashedPassword, "HashPassword did not generate a hashed password")

	// Verify that the hashed password is different from the original password
	assert.NotEqual(t, user.HashedPassword, "HashPassword did not hash the password")
}
