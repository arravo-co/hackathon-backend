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

type AuthParticipantInfoUpdateDTO struct {
	FirstName string `validate:"min=2,omitempty" json:"first_name"`
	LastName  string `validate:"min=2,omitempty" json:"last_name"`
	Email     string `validate:"email,omitempty" json:"email"`
	Gender    string `validate:"oneof=MALE FEMALE,omitempty" json:"gender"`
	State     string `json:"state,omitempty"`
}

type RegisterNewParticipantDTO struct {
	FirstName           string   `validate:"min=2" json:"first_name"`
	LastName            string   `validate:"min=2" json:"last_name"`
	Email               string   `validate:"email" json:"email"`
	Password            string   `validate:"min=7" json:"password"`
	PhoneNumber         string   `validate:"e164" json:"phone_number"`
	ConfirmPassword     string   `validate:"eqfield=Password" json:"confirm_password"`
	Gender              string   `validate:"oneof=MALE FEMALE" json:"gender"`
	Skillset            []string `validate:"min=1" json:"skillset"`
	State               string   `validate:"min=3" json:"state"`
	Type                string   `validate:"oneof=INDIVIDUAL TEAM" json:"type"`
	TeamSize            int      `json:"team_size"`
	DOB                 string   `validate:"required" json:"dob"`
	TeamName            string   `validate:"omitempty" json:"team_name"`
	EmploymentStatus    string   `validate:"oneof=STUDENT EMPLOYED UNEMPLOYED FREELANCER" json:"employment_status"`
	ExperienceLevel     string   `validate:"oneof=JUNIOR MID SENIOR" json:"experience_level"`
	Motivation          string   `validate:"min=100" json:"motivation"`
	HackathonExperience string   `json:"hackathon_experience"`
	YearsOfExperience   int      `json:"years_of_experience"`
	FieldOfStudy        string   `json:"field_of_study"`
	PreviousProjects    []string `json:"previous_projects"`
}

type AddToTeamInviteListData struct {
	HackathonId      string `bson:"hackathon_id"`
	ParticipantId    string `bson:"participant_id"`
	InviterEmail     string `bson:"inviter_email"`
	InviterFirstName string `bson:"inviter_first_name"`
	Email            string `bson:"email"`
	Role             string `bson:"role"`
}
type SelectTeamSolutionData struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
	SolutionId    string `bson:"solution_id"`
}
