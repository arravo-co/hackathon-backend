package repository

import (
	"time"

	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddMemberToParticipatingTeam
type AccountRepository struct {
	DB                  *query.Query
	FirstName           string               `json:"first_name"`
	LastName            string               `json:"last_name"`
	Email               string               `json:"email"`
	Gender              string               `json:"gender"`
	State               string               `json:"state"`
	Age                 int                  `json:"age"`
	DOB                 time.Time            `json:"dob"`
	AccountRole         string               `json:"role"`
	ParticipantId       string               `json:"participant_id"`
	TeamLeadEmail       string               `json:"team_lead_email"`
	TeamName            string               `json:"team_name"`
	TeamRole            string               `json:"team_role"`
	HackathonId         string               `json:"hackathon_id"`
	ParticipantType     string               `json:"type"`
	CoParticipants      []CoParticipantInfo  `json:"co_participants"`
	ParticipantEmail    string               `json:"participant_email"`
	InviteList          []exports.InviteInfo `json:"invite_list"`
	AccountStatus       string               `json:"account_status"`
	ParticipationStatus string               `json:"participation_status"`
	Skillset            []string             `json:"skillset"`
	PhoneNumber         string               `json:"phone_number"`
	EmploymentStatus    string               `json:"employment_status"`
	ExperienceLevel     string               `json:"experience_level"`
	Motivation          string               `json:"motivation"`
	Solution            *SolutionRepository  `json:"solution"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
}

func NewAccountRepository(q *query.Query) *AccountRepository {
	return &AccountRepository{}
}

func (acc *AccountRepository) CreateTeamMemberAccount(dataToSave *exports.CreateTeamMemberAccountData) (*entity.TeamMemberAccount, error) {
	acDoc, err := acc.DB.CreateTeamMemberAccount(dataToSave)
	if err != nil {
		return nil, err
	}
	accId := acDoc.Id.(primitive.ObjectID).Hex()
	return &entity.TeamMemberAccount{
		Email:             acDoc.Email,
		AccountId:         accId,
		FirstName:         acDoc.FirstName,
		LastName:          acDoc.LastName,
		Gender:            acDoc.Gender,
		PhoneNumber:       acDoc.PhoneNumber,
		ParticipantId:     acDoc.ParticipantId,
		Skillset:          acDoc.Skillset,
		HackathonId:       acDoc.HackathonId,
		State:             acDoc.State,
		Status:            acDoc.Status,
		DOB:               acDoc.DOB,
		TeamRole:          dataToSave.TeamRole,
		AccountRole:       dataToSave.Role,
		IsEmailVerified:   acDoc.IsEmailVerified,
		IsEmailVerifiedAt: acDoc.IsEmailVerifiedAt,
		LinkedInAddress:   acDoc.LinkedInAddress,
		CreatedAt:         acDoc.CreatedAt,
		UpdatedAt:         acDoc.UpdatedAt,
	}, nil
}

func (acc *AccountRepository) CreateParticipantAccount(dataToSave *exports.CreateParticipantAccountData) (*entity.Participant, error) {
	accountCol, err := acc.DB.CreateParticipantAccount(dataToSave)
	return &entity.Participant{
		FirstName:       accountCol.FirstName,
		LastName:        accountCol.LastName,
		DOB:             accountCol.DOB,
		Gender:          accountCol.Gender,
		State:           accountCol.State,
		ParticipantId:   accountCol.ParticipantId,
		PhoneNumber:     accountCol.PhoneNumber,
		Email:           accountCol.Email,
		IsEmailVerified: accountCol.IsEmailVerified,
		CreatedAt:       accountCol.CreatedAt,
		UpdatedAt:       accountCol.UpdatedAt,
	}, err
}
