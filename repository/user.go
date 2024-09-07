package repository

import (
	"errors"

	"github.com/arravoco/hackathon_backend/exports"
)

func (accRepo *AccountRepository) ChangePassword(dataInput *exports.PasswordChangeData) (*exports.AccountRepository, error) {
	acc, err := accRepo.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return nil, err
	}
	valid, _ := exports.ComparePasswordAndHash(dataInput.OldPassword, acc.PasswordHash)
	if !valid {
		return nil, errors.New("email/password does not match")
	}
	passwordHash, err := exports.GenerateHashPassword(dataInput.NewPassword)
	if err != nil {
		return nil, err
	}
	acc, err = accRepo.UpdatePasswordByEmail(&exports.UpdateAccountDocumentFilter{Email: dataInput.Email}, passwordHash)
	return acc, err
}
