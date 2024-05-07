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
	Skillset          []string  `bson:"skillset,omitempty"`
	ParticipantId     string    `bson:"participant_id,omitempty"`
	HackathonId       string    `bson:"hackathon_id"`
	State             string    `bson:"state,omitempty"`
	Role              string    `bson:"role,omitempty"`
	DOB               time.Time `bson:"dob,omitempty"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
	Status            string    `bson:"status"`
	CreatedAt         time.Time `bson:"created_at,omitempty"`
	UpdatedAt         time.Time `bson:"updated_at,omitempty"`
}

type ParticipantDocument struct {
	Id                  interface{}
	ParticipantId       string           `bson:"participant_id"`
	HackathonId         string           `bson:"hackathon_id"`
	Type                string           `bson:"type,omitempty"`
	TeamLeadEmail       string           `bson:"team_lead_email,omitempty"`
	TeamName            string           `bson:"team_name,omitempty"`
	CoParticipantEmails []string         `bson:"co_participant_emails,omitempty"`
	ParticipantEmail    string           `bson:"participant_email,omitempty"`
	GithubAddress       string           `bson:"github_address,omitempty"`
	InviteList          []InviteInfo     `bson:"invite_list,omitempty"`
	Status              string           `bson:"status,omitempty"`
	Solution            SolutionDocument `bson:"solution_document"`
	CreatedAt           time.Time        `bson:"created_at,omitempty"`
	UpdatedAt           time.Time        `bson:"updated_at,omitempty"`
}

type SolutionDocument struct {
	Title       string `bson:"title,omitempty"`
	Description string `bson:"description"`
	Unique      string `bson:"unique"`
	Url         string `bson:"url"`
}

type InviteInfo struct {
	Email     string    `bson:"email,omitempty"`
	Time      time.Time `bson:"time,omitempty"`
	InviterId string    `bson:"inviter_id,omitempty"`
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

type CreateAdminAccountData struct {
	Email        string `bson:"email"`
	PasswordHash string `bson:"password_hash"`
	FirstName    string `bson:"first_name"`
	LastName     string `bson:"last_name"`
	Gender       string `bson:"gender"`
	HackathonId  string `bson:"hackathon_id"`
	Role         string `bson:"role"`
	PhoneNumber  string `bson:"phone_number"`
	Status       string `bson:"status"`
}

type CreateAccountData struct {
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	FirstName    string    `bson:"first_name"`
	LastName     string    `bson:"last_name"`
	Gender       string    `bson:"gender"`
	State        string    `bson:"state"`
	Role         string    `bson:"role"`
	PhoneNumber  string    `bson:"phone_number"`
	DOB          time.Time `bson:"dob"`
	HackathonId  string    `bson:"hackathon_id"`
	Status       string    `bson:"status"`
}

type CreateParticipantAccountData struct {
	CreateAccountData
	ParticipantId string   `bson:"participant_id"`
	Skillset      []string `bson:"skillset"`
}

type CreateTeamMemberAccountData struct {
	CreateAccountData
	ParticipantId string   `bson:"participant_id"`
	Skillset      []string `bson:"skillset"`
}

type RemoveTeamMemberAccountData struct {
	Email         string `bson:"email"`
	ParticipantId string `bson:"participant_id"`
	HackathonId   string `bson:"hackathon_id"`
}

type AddMemberToParticipatingTeamData struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
	Email         string `bson:"email"`
	Role          string `bson:"role"`
	TeamRole      string `bson:"team_role"`
}

type RemoveMemberFromParticipatingTeamData struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
	MemberEmail   string `bson:"email"`
}

type AddToTeamInviteListData struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
	InviterEmail  string `bson:"inviter_email"`
	Email         string `bson:"email"`
	Role          string `bson:"role"`
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
	ParticipantId       string   `bson:"participant_id"`
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
