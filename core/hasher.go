package core

import "golang.org/x/crypto/bcrypt"

type Hasher struct{}

// Hashes the given password and returns the hashed password as a string.
// Returns an error if hashing fails.
func (h *Hasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hashedPassword := string(hashedBytes)
	return hashedPassword, nil
}

// Compares a hashed password with a plain password.
// Returns true if they match, false otherwise.
func (h *Hasher) Compare(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
