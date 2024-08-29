package entity

import (
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	valueobjects "github.com/arravoco/hackathon_backend/value_objects"
)

// AddMemberToParticipatingTeam
type Participant struct {
	FirstName           string                           `json:"first_name"`
	LastName            string                           `json:"last_name"`
	Email               string                           `json:"email"`
	Gender              string                           `json:"gender"`
	State               string                           `json:"state"`
	Age                 int                              `json:"age"`
	DOB                 time.Time                        `json:"dob"`
	AccountRole         string                           `json:"role"`
	ParticipantId       string                           `json:"participant_id"`
	TeamLeadEmail       string                           `json:"team_lead_email"`
	TeamName            string                           `json:"team_name"`
	TeamRole            string                           `json:"team_role"`
	HackathonId         string                           `json:"hackathon_id"`
	ParticipantType     string                           `json:"type"`
	CoParticipants      []valueobjects.CoParticipantInfo `json:"co_participants"`
	ParticipantEmail    string                           `json:"participant_email"`
	InviteList          []exports.InviteInfo             `json:"invite_list"`
	AccountStatus       string                           `json:"account_status"`
	ParticipationStatus string                           `json:"participation_status"`
	ReviewRanking       int                              `json:"review_ranking"`
	Skillset            []string                         `json:"skillset"`
	PhoneNumber         string                           `json:"phone_number"`
	EmploymentStatus    string                           `json:"employment_status"`
	ExperienceLevel     string                           `json:"experience_level"`
	Motivation          string                           `json:"motivation"`
	HackathonExperience string                           `json:"hackathon_experience"`
	YearsOfExperience   int                              `json:"years_of_experience"`
	FieldOfStudy        string                           `json:"field_of_study"`
	PreviousProjects    []string                         `json:"previous_projects"`
	IsEmailVerified     bool                             `json:"is_email_verified"`
	EmailVerifiedAt     time.Time                        `json:"email_verified_at"`
	Solution            *Solution                        `json:"solution"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           time.Time                        `json:"updated_at"`
}

type InviteInfo struct {
}
