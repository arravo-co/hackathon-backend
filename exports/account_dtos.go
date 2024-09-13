package exports

import "time"

type RegisterNewAccountDTO struct {
	FirstName           string    `validate:"min=2" json:"first_name"`
	LastName            string    `validate:"min=2" json:"last_name"`
	Email               string    `validate:"email" json:"email"`
	Password            string    `validate:"min=7" json:"password"`
	PhoneNumber         string    `validate:"e164" json:"phone_number"`
	ConfirmPassword     string    `validate:"eqfield=Password" json:"confirm_password"`
	Gender              string    `validate:"oneof=MALE FEMALE" json:"gender"`
	Skillset            []string  `validate:"min=1" json:"skillset"`
	State               string    `validate:"min=3" json:"state"`
	Type                string    `validate:"oneof=INDIVIDUAL TEAM" json:"type"`
	TeamSize            int       `json:"team_size"`
	TeamName            string    `validate:"omitempty" json:"team_name"`
	EmploymentStatus    string    `validate:"oneof=STUDENT EMPLOYED UNEMPLOYED FREELANCER" json:"employment_status"`
	ExperienceLevel     string    `validate:"oneof=JUNIOR MID SENIOR" json:"experience_level"`
	Motivation          string    `validate:"min=100" json:"motivation"`
	HackathonExperience string    `json:"hackathon_experience"`
	YearsOfExperience   int       `json:"years_of_experience"`
	FieldOfStudy        string    `json:"field_of_study"`
	PreviousProjects    []string  `json:"previous_projects"`
	HackathonId         string    `json:"hackathon_id"`
	Role                string    `json:"role"`
	Status              string    `json:"status"`
	ParticipantId       string    `json:"participant_id"`
	DOB                 time.Time `json:"dob"`
}

type UpdateAccountDTO struct {
	FirstName         string `validate:"omitempty,min=2"`
	LastName          string `validate:"omitempty,min=2"`
	Gender            string `validate:"omitempty, oneof=MALE FEMALE"`
	State             string
	Bio               string
	Status            string
	ProfilePictureUrl string
	IsEmailVerified   bool
	IsEmailVerifiedAt time.Time
}

type PasswordChangeData struct {
	Email       string
	OldPassword string
	NewPassword string
}
