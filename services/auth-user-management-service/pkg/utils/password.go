package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword encrypts password using bcrypt
func HashPassword(password string) (string, error) {
	// Cost = 14 is the complexity of the encryption algorithm
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plaintext one
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
