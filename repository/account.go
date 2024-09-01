package repository

import (
	"time"

	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
)

/*
type Datasource interface {
	DeleteAccount(identifier string) (*exports.AccountDocument, error)
	FindAccountIdentifier(identifier string) (*exports.AccountDocument, error)
	GetAccountByEmail(email string) (*exports.AccountDocument, error)
	GetAccountsByEmails(emails []string) ([]exports.AccountDocument, error)
	CreateJudgeAccount(dataToSave *exports.CreateJudgeAccountData) (*exports.CreateJudgeAccountData, error)
}
*/
// AddMemberToParticipatingTeam
type AccountRepository struct {
	DB                  *query.Query
	FirstName           string                                      `json:"first_name"`
	LastName            string                                      `json:"last_name"`
	Email               string                                      `json:"email"`
	Gender              string                                      `json:"gender"`
	State               string                                      `json:"state"`
	Age                 int                                         `json:"age"`
	DOB                 time.Time                                   `json:"dob"`
	AccountRole         string                                      `json:"role"`
	ParticipantId       string                                      `json:"participant_id"`
	TeamLeadEmail       string                                      `json:"team_lead_email"`
	TeamName            string                                      `json:"team_name"`
	TeamRole            string                                      `json:"team_role"`
	HackathonId         string                                      `json:"hackathon_id"`
	ParticipantType     string                                      `json:"type"`
	CoParticipants      []CoParticipantInfo                         `json:"co_participants"`
	ParticipantEmail    string                                      `json:"participant_email"`
	InviteList          []exports.ParticipantDocumentTeamInviteInfo `json:"invite_list"`
	AccountStatus       string                                      `json:"account_status"`
	ParticipationStatus string                                      `json:"participation_status"`
	Skillset            []string                                    `json:"skillset"`
	PhoneNumber         string                                      `json:"phone_number"`
	EmploymentStatus    string                                      `json:"employment_status"`
	ExperienceLevel     string                                      `json:"experience_level"`
	Motivation          string                                      `json:"motivation"`
	Solution            *SolutionRepository                         `json:"solution"`
	CreatedAt           time.Time                                   `json:"created_at"`
	UpdatedAt           time.Time                                   `json:"updated_at"`
}

func NewAccountRepository(q *query.Query) *AccountRepository {
	return &AccountRepository{
		DB: q,
	}
}
