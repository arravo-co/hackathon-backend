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

func (p *ParticipantAccountRepository) FindAccountIdentifier(identifier string) (*exports.ParticipantAccountRepository, error) {
	partAccDoc, err := p.DB.FindAccountIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	return &exports.ParticipantAccountRepository{
		Email:               partAccDoc.Email,
		LastName:            partAccDoc.LastName,
		FirstName:           partAccDoc.FirstName,
		Gender:              partAccDoc.Gender,
		PasswordHash:        partAccDoc.PasswordHash,
		Status:              partAccDoc.Status,
		State:               partAccDoc.State,
		LinkedInAddress:     partAccDoc.LinkedInAddress,
		DOB:                 partAccDoc.DOB,
		EmploymentStatus:    partAccDoc.EmploymentStatus,
		IsEmailVerified:     partAccDoc.IsEmailVerified,
		IsEmailVerifiedAt:   partAccDoc.IsEmailVerifiedAt,
		ExperienceLevel:     partAccDoc.ExperienceLevel,
		HackathonExperience: partAccDoc.HackathonExperience,
		YearsOfExperience:   partAccDoc.YearsOfExperience,
		Motivation:          partAccDoc.Motivation,
		CreatedAt:           partAccDoc.CreatedAt,
		UpdatedAt:           partAccDoc.UpdatedAt,
	}, nil
}

func MarkParticipantAccountAsDeleted(identifier string) (*exports.ParticipantAccountRepository, error) {
	return nil, nil
}

func (p *ParticipantAccountRepository) GetParticipantAccountByEmail(email string) (*exports.ParticipantAccountRepository, error) {
	partAcc, err := p.DB.GetAccountByEmail(email)
	if err != nil {
		return nil, err
	}
	return &exports.ParticipantAccountRepository{
		FirstName:           partAcc.FirstName,
		LastName:            partAcc.LastName,
		ParticipantId:       partAcc.ParticipantId,
		HackathonId:         partAcc.HackathonId,
		Email:               partAcc.Email,
		ExperienceLevel:     partAcc.ExperienceLevel,
		EmploymentStatus:    partAcc.EmploymentStatus,
		IsEmailVerified:     partAcc.IsEmailVerified,
		HackathonExperience: partAcc.HackathonExperience,
		YearsOfExperience:   partAcc.YearsOfExperience,
		IsEmailVerifiedAt:   partAcc.IsEmailVerifiedAt,
		LinkedInAddress:     partAcc.LinkedInAddress,
		FieldOfStudy:        partAcc.FieldOfStudy,
		Id:                  partAcc.Id.Hex(),
		Role:                partAcc.Role,
		Status:              partAcc.Status,
		State:               partAcc.State,
		Skillset:            partAcc.Skillset,
		PasswordHash:        partAcc.PasswordHash,
		PhoneNumber:         partAcc.PhoneNumber,
		PreviousProjects:    partAcc.PreviousProjects,
		ProfilePictureUrl:   partAcc.ProfilePictureUrl,
	}, nil
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
func (p *ParticipantAccountRepository) UpdateParticipantPassword(filter *exports.UpdateAccountDocumentFilter, newPasswordHash string) (*exports.ParticipantAccountRepository, error) {
	return nil, nil
}

func (c *ParticipantAccountRepository) GetParticipantAccountWithParticipantRecord() {

}
