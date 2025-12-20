package authService

import "golang.org/x/crypto/bcrypt"

// test hooks
var bcryptGenerate = bcrypt.GenerateFromPassword

func hashPassword(password string) (string, error) {
	b, err := bcryptGenerate([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func comparePassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
