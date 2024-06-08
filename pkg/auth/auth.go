package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password using bcrypt with a generated salt
func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

// CheckPasswordHash compares the hashed password with the plain password
func CheckPasswordHash(password, hashedPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}
