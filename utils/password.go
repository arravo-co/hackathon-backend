package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	passwordHashByte, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	passwordHash := string(passwordHashByte)
	return passwordHash, nil
}

func ComparePasswordAndHash(password, passwordHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
