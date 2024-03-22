package utils

import (
	"errors"
	"fmt"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword() string {
	passwordHash := password.MustGenerate(10, 2, 1, false, false)
	return passwordHash
}

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
		fmt.Printf("%v\n", err)
		return false, errors.New("authentication failed with identifier/password mismatch")
	}
	return true, nil
}
