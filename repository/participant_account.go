package repository

import (

	//"gitee.com/golang-module/carbon"

	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
	//"github.com/golang-module/carbon"
)

type ParticipantAccountRepository struct {
	DB *query.Query
}

func NewParticipantAccountRepository(q *query.Query) *ParticipantAccountRepository {
	return &ParticipantAccountRepository{
		DB: q,
	}
}

// GetParticipantAccountByEmail
func (p *ParticipantAccountRepository) CreateTeamMemberAccount(dataInput *exports.CreateTeamMemberAccountData) (*exports.ParticipantAccountRepository, error) {
	partDoc, err := p.DB.CreateTeamMemberAccount(dataInput)
	if err != nil {
		return nil, err
	}

	return &exports.ParticipantAccountRepository{
		Email:               partDoc.Email,
		LastName:            partDoc.LastName,
		FirstName:           partDoc.FirstName,
		Gender:              partDoc.Gender,
		PasswordHash:        partDoc.PasswordHash,
		Status:              partDoc.Status,
		State:               partDoc.State,
		LinkedInAddress:     partDoc.LinkedInAddress,
		DOB:                 partDoc.DOB,
		EmploymentStatus:    partDoc.EmploymentStatus,
		IsEmailVerified:     partDoc.IsEmailVerified,
		IsEmailVerifiedAt:   partDoc.IsEmailVerifiedAt,
		ExperienceLevel:     partDoc.ExperienceLevel,
		HackathonExperience: partDoc.HackathonExperience,
		YearsOfExperience:   partDoc.YearsOfExperience,
		Motivation:          partDoc.Motivation,
		CreatedAt:           partDoc.CreatedAt,
		UpdatedAt:           partDoc.UpdatedAt,
	}, nil
}

func (p *ParticipantAccountRepository) CreateParticipantAccount(input *exports.CreateParticipantAccountData) (*exports.ParticipantAccountRepository, error) {

	/*areEmailsInCache := cache.FindEmailInCache(input.Email)
	if !areEmailsInCache {
		//return nil, errors.New("email already exists")
	}*/

	acc, err := p.DB.CreateParticipantAccount(input)
	if err != nil {
		return nil, err
	}

	// emit created event

	return &exports.ParticipantAccountRepository{
		HackathonId:   acc.HackathonId,
		ParticipantId: acc.ParticipantId,
		CreatedAt:     acc.CreatedAt,
		UpdatedAt:     acc.UpdatedAt,
	}, err
}

func MarkParticipantAccountAsDeleted(identifier string) (*exports.ParticipantAccountRepository, error) {
	return nil, nil
}

func (p *ParticipantAccountRepository) GetParticipantAccountByEmail(email string) (*exports.ParticipantAccountRepository, error) {
	return nil, nil
}
func (p *ParticipantAccountRepository) GetParticipantAccountsByEmail(emails []string) ([]*exports.ParticipantAccountRepository, error) {
	return nil, nil
}
func (p *ParticipantAccountRepository) UpdateParticipantAccount(email string, input *exports.UpdateParticipantDTO) error {
	return nil
}
func (p *ParticipantAccountRepository) GetParticipantAccounts() ([]*exports.ParticipantAccountRepository, error) {
	return nil, nil
}
func (p *ParticipantAccountRepository) DeleteParticipantAccount(identifier string) (*exports.ParticipantAccountRepository, error) {
	return nil, nil
}
func (p *ParticipantAccountRepository) MarkParticipantAccountAsDeleted(identifier string) (*exports.ParticipantAccountRepository, error) {
	return nil, nil
}

// GetJudgeAccountByEmail(email string) (*JudgeAccountRepository, error)
// UpdateJudgeAccount(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*JudgeAccountRepository, error)
func (p *ParticipantAccountRepository) UpdateParticipantPassword(filter *exports.UpdateAccountFilter, newPasswordHash string) (*exports.ParticipantAccountRepository, error) {
	return nil, nil
}
