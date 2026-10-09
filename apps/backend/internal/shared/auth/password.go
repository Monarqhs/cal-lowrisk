// Package auth holds cross-module authentication primitives: password hashing,
// JWT issue/parse, and the Gin middleware that enforces authentication and roles.
// Kept in shared/ because every module that needs an authenticated caller depends on it.
package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a bcrypt hash of the plaintext password.
// bcrypt silently truncates input beyond 72 bytes; the register DTO caps length at 72.
func HashPassword(plaintext string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// VerifyPassword reports whether plaintext matches the stored bcrypt hash.
func VerifyPassword(hash, plaintext string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext)) == nil
}
