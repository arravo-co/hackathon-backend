package exports

import (
	"time"
)

type AccountDocument struct {
	Email             string    `bson:"email,omitempty"`
	FirstName         string    `bson:"first_name,omitempty"`
	LastName          string    `bson:"last_name,omitempty"`
	Gender            string    `bson:"gender,omitempty"`
	LinkedInAddress   string    `bson:"linkedIn_address,omitempty"`
	PasswordHash      string    `bson:"password_hash,omitempty"`
	PhoneNumber       string    `bson:"phone_number,omitempty"`
	ParticipantId     string    `bson:"participant_id,omitempty"`
	HackathonId       string    `bson:"hackathon_id"`
	State             string    `bson:"state,omitempty"`
	Role              string    `bson:"role,omitempty"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
}

type ParticipantDocument struct {
	Id                  interface{}
	ParticipantId       string   `bson:"participant_id"`
	HackathonId         string   `bson:"hackathon_id"`
	Type                string   `bson:"type,omitempty"`
	TeamLeadEmail       string   `bson:"team_lead_email,omitempty"`
	TeamName            string   `bson:"team_name,omitempty"`
	CoParticipantEmails []string `bson:"co_participant_emails,omitempty"`
	ParticipantEmail    string   `bson:"participant_email,omitempty"`
	GithubAddress       string   `bson:"github_address,omitempty"`
}

type ParticipantScoreDocument struct {
	Id            interface{}
	HackathonId   string      `validate:"required" json:"hackathon_id"`
	ParticipantId string      `validate:"required" json:"participant_id"`
	Stage         string      `validate:"required" json:"stage"`
	TotalScore    float32     `validate:"required" json:"score"`
	ScoresInfo    []ScoreInfo `validate:"required" json:"scores_info"`
}

type UpdateAccountFilter struct {
	Email       string `bson:"email"`
	PhoneNumber string `bson:"phone_number"`
}

type UpdateAccountDocument struct {
	FirstName         string    `bson:"first_name,omitempty"`
	LastName          string    `bson:"last_name,omitempty"`
	Gender            string    `bson:"gender,omitempty"`
	LinkedInAddress   string    `bson:"linkedIn_address,omitempty"`
	GithubAddress     string    `bson:"github_address,omitempty"`
	State             string    `bson:"state,omitempty"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
}

type CreateAccountData struct {
	Email        string `bson:"email"`
	PasswordHash string `bson:"password_hash"`
	FirstName    string `bson:"first_name"`
	LastName     string `bson:"last_name"`
	Gender       string `bson:"gender"`
	State        string `bson:"state"`
	Role         string `bson:"role"`
	PhoneNumber  string `bson:"phone_number"`
}

type CreateIndividualParticipantAccountData struct {
	CreateAccountData
	LinkedInAddress string `bson:"linkedIn_address"`
}

type CreateTeamParticipantAccountData struct {
	Email               string   `bson:"email"`
	ParticipantId       string   `bson:"participant_id"`
	Role                string   `bson:"role"`
	HackathonId         string   `bson:"hackathon_id"`
	Type                string   `bson:"type,omitempty"`
	TeamLeadEmail       string   `bson:"team_lead_email,omitempty"`
	TeamName            string   `bson:"team_lead_email,omitempty"`
	CoParticipantEmails []string `bson:"co_participant_emails,omitempty"`
	GithubAddress       string   `bson:"github_address,omitempty"`
}

type CreateParticipantRecordData struct {
	ParticipantEmail    string   `bson:"participant_email,omitempty"`
	TeamLeadEmail       string   `bson:"team_lead_email,omitempty"`
	TeamName            string   `bson:"team_name,omitempty"`
	CoParticipantEmails []string `bson:"co_participant_emails,omitempty"`
	Type                string   `bson:"type"`
	HackathonId         string   `bson:"hackathon_id"`
	GithubAddress       string   `bson:"github_address,omitempty"`
}

type TeamParticipantRecordCreatedData struct {
	ParticipantId       string   `bson:"participant_id"`
	Role                string   `bson:"role"`
	HackathonId         string   `bson:"hackathon_id"`
	Type                string   `bson:"type,omitempty"`
	TeamLeadEmail       string   `bson:"team_lead_email,omitempty"`
	TeamName            string   `bson:"team_lead_email,omitempty"`
	CoParticipantEmails []string `bson:"co_participant_emails,omitempty"`
	GithubAddress       string   `bson:"github_address,omitempty"`
}

type CreateJudgeAccountData struct {
	CreateAccountData
}
