package entity

import "time"

type TeamMemberWithParticipantRecord struct {
	Email           string   `json:"email"`
	AccountId       string   `json:"account_id,omitempty"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Gender          string   `json:"gender,omitempty"`
	LinkedInAddress string   `json:"linkedIn_address,omitempty"`
	PhoneNumber     string   `json:"phone_number,omitempty"`
	Skillset        []string `json:"skillset,omitempty"`
	ParticipantId   string   `json:"participant_id,omitempty"`
	HackathonId     string   `json:"hackathon_id"`

	State               string    `json:"state,omitempty"`
	AccountRole         string    `json:"account_role,omitempty"`
	TeamRole            string    `json:"team_role,omitempty"`
	DOB                 time.Time `json:"dob,omitempty"`
	Motivation          string    `json:"motivation,omitempty"`
	HackathonExperience string    `json:"hackathon_experience"`
	YearsOfExperience   int       `json:"years_of_experience"`
	FieldOfStudy        string    `json:"field_of_study"`
	PreviousProjects    []string  `json:"previous_projects"`
	IsEmailVerified     bool      `json:"is_email_verified,omitempty"`
	IsEmailVerifiedAt   time.Time `json:"is_email_verified_at,omitempty"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
