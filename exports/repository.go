package exports

import (
	"time"
)

type JudgeAccountRepository struct {
	DB                *JudgeRepositoryInterface
	Id                string
	FirstName         string
	LastName          string
	Email             string
	passwordHash      string
	Gender            string
	Role              string
	HackathonId       string
	Status            string
	State             string
	PhoneNumber       string
	Bio               string
	IsEmailVerified   bool
	ProfilePictureUrl string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type JudgeRepositoryInterface interface {
	GetJudgeAccountByEmail(email string) (*JudgeAccountRepository, error)
	CreateJudgeAccount(input *RegisterNewJudgeDTO) (*JudgeAccountRepository, error)
	UpdateJudgeAccount(email string, input *UpdateJudgeDTO) error
	GetJudges() ([]*JudgeAccountRepository, error)
	DeleteJudgeAccount(identifier string) (*JudgeAccountRepository, error)
	//GetJudgeAccountByEmail(email string) (*JudgeAccountRepository, error)
	//UpdateJudgeAccount(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*JudgeAccountRepository, error)
	UpdateJudgePassword(filter *UpdateAccountFilter, newPasswordHash string) (*JudgeAccountRepository, error)
}
