package auth

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12 // same default cost as Laravel

// HashPassword — equivalent to Laravel's Hash::make($password)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

// CheckPassword — equivalent to Laravel's Hash::check($password, $hash)
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
