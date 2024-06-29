package valueobjects

import "time"

type CoParticipantInfo struct {
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Email            string    `json:"email"`
	Gender           string    `json:"gender"`
	State            string    `json:"state"`
	Age              int       `json:"age"`
	DOB              time.Time `json:"dob"`
	AccountStatus    string    `json:"account_status"`
	AccountRole      string    `json:"account_role"`
	TeamRole         string    `json:"team_role"`
	ParticipantId    string    `json:"participant_id"`
	HackathonId      string    `json:"hackathon_id"`
	Skillset         []string  `json:"skillset"`
	PhoneNumber      string    `json:"phone_number"`
	EmploymentStatus string    `json:"employment_status"`
	ExperienceLevel  string    `json:"experience_level"`
	Motivation       string    `json:"motivation"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
