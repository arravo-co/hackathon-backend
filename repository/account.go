package repository

import (

	//"gitee.com/golang-module/carbon"

	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
	//"github.com/golang-module/carbon"
)

type AccountRepository struct {
	DB *query.Query
}

func NewAccountRepository(q *query.Query) *AccountRepository {
	return &AccountRepository{
		DB: q,
	}
}

func (p *AccountRepository) CreateAccount(input *exports.CreateAccountData) (*exports.AccountRepository, error) {
	accDoc, err := p.DB.CreateAccount(input)
	if err != nil {
		return nil, err
	}

	// emit created event

	return &exports.AccountRepository{
		Id:                  accDoc.Id.Hex(),
		Email:               accDoc.Email,
		LastName:            accDoc.LastName,
		FirstName:           accDoc.FirstName,
		Gender:              accDoc.Gender,
		PasswordHash:        accDoc.PasswordHash,
		Status:              accDoc.Status,
		State:               accDoc.State,
		LinkedInAddress:     accDoc.LinkedInAddress,
		DOB:                 accDoc.DOB,
		EmploymentStatus:    accDoc.EmploymentStatus,
		IsEmailVerified:     accDoc.IsEmailVerified,
		IsEmailVerifiedAt:   accDoc.IsEmailVerifiedAt,
		ExperienceLevel:     accDoc.ExperienceLevel,
		HackathonExperience: accDoc.HackathonExperience,
		YearsOfExperience:   accDoc.YearsOfExperience,
		Motivation:          accDoc.Motivation,
		CreatedAt:           accDoc.CreatedAt,
		UpdatedAt:           accDoc.UpdatedAt,
	}, nil
}

func (p *AccountRepository) FindAccountIdentifier(identifier string) (*exports.AccountRepository, error) {
	accDoc, err := p.DB.FindAccountIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	return &exports.AccountRepository{
		Id:                  accDoc.Id.Hex(),
		Email:               accDoc.Email,
		LastName:            accDoc.LastName,
		FirstName:           accDoc.FirstName,
		Gender:              accDoc.Gender,
		PasswordHash:        accDoc.PasswordHash,
		Status:              accDoc.Status,
		State:               accDoc.State,
		LinkedInAddress:     accDoc.LinkedInAddress,
		DOB:                 accDoc.DOB,
		EmploymentStatus:    accDoc.EmploymentStatus,
		IsEmailVerified:     accDoc.IsEmailVerified,
		IsEmailVerifiedAt:   accDoc.IsEmailVerifiedAt,
		ExperienceLevel:     accDoc.ExperienceLevel,
		HackathonExperience: accDoc.HackathonExperience,
		YearsOfExperience:   accDoc.YearsOfExperience,
		Motivation:          accDoc.Motivation,
		CreatedAt:           accDoc.CreatedAt,
		UpdatedAt:           accDoc.UpdatedAt,
		Role:                accDoc.Role,
		ParticipantId:       accDoc.ParticipantId,
		PhoneNumber:         accDoc.PhoneNumber,
		ProfilePictureUrl:   accDoc.ProfilePictureUrl,
		HackathonId:         accDoc.HackathonId,
	}, nil
}

func MarkAccountAsDeleted(identifier string) (*exports.AccountRepository, error) {
	return nil, nil
}

func (p *AccountRepository) GetAccountByEmail(email string) (*exports.AccountRepository, error) {
	accDoc, err := p.DB.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	return &exports.AccountRepository{
		Id:                  accDoc.Id.Hex(),
		Email:               accDoc.Email,
		LastName:            accDoc.LastName,
		FirstName:           accDoc.FirstName,
		Gender:              accDoc.Gender,
		PasswordHash:        accDoc.PasswordHash,
		Status:              accDoc.Status,
		State:               accDoc.State,
		LinkedInAddress:     accDoc.LinkedInAddress,
		DOB:                 accDoc.DOB,
		EmploymentStatus:    accDoc.EmploymentStatus,
		IsEmailVerified:     accDoc.IsEmailVerified,
		IsEmailVerifiedAt:   accDoc.IsEmailVerifiedAt,
		ExperienceLevel:     accDoc.ExperienceLevel,
		HackathonExperience: accDoc.HackathonExperience,
		YearsOfExperience:   accDoc.YearsOfExperience,
		Motivation:          accDoc.Motivation,
		CreatedAt:           accDoc.CreatedAt,
		UpdatedAt:           accDoc.UpdatedAt,
	}, nil
}
func (p *AccountRepository) GetAccountsByEmail(emails []string) ([]*exports.AccountRepository, error) {
	return nil, nil
}
func (p *AccountRepository) UpdateAccount(filter *exports.UpdateAccountFilter, input *exports.UpdateAccountDTO) error {
	p.DB.UpdateAccountInfoByEmail(filter, &exports.UpdateAccountDocument{
		IsEmailVerified:   input.IsEmailVerified,
		IsEmailVerifiedAt: input.IsEmailVerifiedAt,
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Bio:               input.Bio,
		State:             input.State,
		Gender:            input.Gender,
	})
	return nil
}

func (p *AccountRepository) UpdatePasswordByEmail(filter *exports.UpdateAccountFilter, newPasswordHash string) (*exports.AccountRepository, error) {
	accDoc, err := p.DB.UpdatePasswordByEmail(filter, newPasswordHash)
	if err != nil {
		return nil, err
	}

	return &exports.AccountRepository{
		Id:                  accDoc.Id.Hex(),
		Email:               accDoc.Email,
		LastName:            accDoc.LastName,
		FirstName:           accDoc.FirstName,
		Gender:              accDoc.Gender,
		PasswordHash:        accDoc.PasswordHash,
		Status:              accDoc.Status,
		State:               accDoc.State,
		LinkedInAddress:     accDoc.LinkedInAddress,
		DOB:                 accDoc.DOB,
		EmploymentStatus:    accDoc.EmploymentStatus,
		IsEmailVerified:     accDoc.IsEmailVerified,
		IsEmailVerifiedAt:   accDoc.IsEmailVerifiedAt,
		ExperienceLevel:     accDoc.ExperienceLevel,
		HackathonExperience: accDoc.HackathonExperience,
		YearsOfExperience:   accDoc.YearsOfExperience,
		Motivation:          accDoc.Motivation,
		CreatedAt:           accDoc.CreatedAt,
		UpdatedAt:           accDoc.UpdatedAt,
	}, nil
}

func (p *AccountRepository) DeleteAccount(identifier string) (*exports.AccountRepository, error) {
	return nil, nil
}
