package auth

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 12 // mesmo custo padrão do Laravel

// HashSenha — equivalente ao Hash::make($password) do Laravel
func HashSenha(senha string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senha), bcryptCost)
	return string(bytes), err
}

// VerificarSenha — equivalente ao Hash::check($password, $hash) do Laravel
func VerificarSenha(senha, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha))
	return err == nil
}
