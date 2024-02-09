package helper

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(plainPassword string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	return string(result), err
}

func CompareHashAndPlainPassword(hashPassword string, plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword)) == nil
}
