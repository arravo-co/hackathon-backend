package exports

import (
	"time"

	"github.com/arravoco/hackathon_backend/dtos"
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
	GetJudgeByEmail(email string) (*JudgeAccountRepository, error)
	CreateJudgeAccount(input *RegisterNewJudgeDTO) (*JudgeAccountRepository, error)
	UpdateJudgeProfile(email string, input dtos.UpdateJudgeDTO) error
	GetJudges() ([]*JudgeAccountRepository, error)
	DeleteJudgeAccount(identifier string) (*JudgeAccountRepository, error)
	GetJudgeAccountByEmail(email string) (*JudgeAccountRepository, error)
	UpdateJudgeAccountInfoByEmail(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*JudgeAccountRepository, error)
	UpdateJudgePasswordByEmail(filter *UpdateAccountFilter, newPasswordHash string) (*JudgeAccountRepository, error)
}
