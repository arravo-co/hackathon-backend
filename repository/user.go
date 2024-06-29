package repository

import (
	"errors"
	"fmt"

	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/exports"
)

type User interface {
	Register()
	InitiatePasswordChange()
	CompletePasswordChange()
	InitiatePasswordRecovery()
	CompletePasswordRecovery()
}

type PasswordChangeData struct {
	Email       string
	OldPassword string
	NewPassword string
}

func ChangePassword(dataInput *PasswordChangeData) (*exports.AccountDocument, error) {
	acc, err := data.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return nil, err
	}
	fmt.Println(acc)
	valid, _ := exports.ComparePasswordAndHash(dataInput.OldPassword, acc.PasswordHash)
	if !valid {
		return nil, errors.New("email/password does not match")
	}
	passwordHash, err := exports.GenerateHashPassword(dataInput.NewPassword)
	if err != nil {
		return nil, err
	}
	acc, err = data.UpdatePasswordByEmail(&exports.UpdateAccountFilter{Email: dataInput.Email}, passwordHash)
	return acc, err
}
