package services

type RegisterNewJudgeDTO struct {
	FirstName       string `validate:"min=2,required" json:"first_name"`
	LastName        string `validate:"min=2,required" json:"last_name"`
	Email           string `validate:"email,required" json:"email"`
	Password        string `validate:"min=7" json:"password"`
	ConfirmPassword string `validate:"eqfield=Password" json:"confirm_password"`
	Gender          string `validate:"oneof=MALE FEMALE" json:"gender"`
	State           string `json:"state"`
	Bio             string `json:"bio"`
}

type UpdateJudgeDTO struct {
	FirstName         string `validate:"omitempty,min=2"`
	LastName          string `validate:"omitempty,min=2"`
	Gender            string `validate:"omitempty,oneof=MALE FEMALE"`
	State             string
	Bio               string
	ProfilePictureUrl string
}

type CompleteNewTeamMemberRegistrationDTO struct {
	FirstName        string   `validate:"required"`
	LastName         string   `validate:"required"`
	Email            string   `validate:"email"`
	Password         string   `validate:"min=5"`
	PhoneNumber      string   `validate:"e164"`
	ConfirmPassword  string   `validate:"eqfield=Password"`
	Gender           string   `validate:"oneof= MALE FEMALE"`
	Skillset         []string `validate:"min=1"`
	State            string   `validate:"required"`
	DOB              string   `validate:"required"`
	ParticipantId    string   `validate:"required"`
	TeamLeadEmail    string   `validate:"email"`
	HackathonId      string   `validate:"required"`
	TeamRole         string   `validate:"oneof= TEAM_MEMBER"`
	EmploymentStatus string   `validate:"required"`
	ExperienceLevel  string   `validate:"required"`
	Motivation       string   `validate:"required"`
}
