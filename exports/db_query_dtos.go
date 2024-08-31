package exports

import (
	"time"
)

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
	Bio               string    `bson:"bio,omitempty"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
	ProfilePictureUrl string    `bson:"profile_picture_url"`
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
	Email             string    `bson:"email"`
	PasswordHash      string    `bson:"password_hash"`
	FirstName         string    `bson:"first_name"`
	LastName          string    `bson:"last_name"`
	Gender            string    `bson:"gender"`
	State             string    `bson:"state"`
	Role              string    `bson:"role"`
	PhoneNumber       string    `bson:"phone_number"`
	HackathonId       string    `bson:"hackathon_id"`
	Status            string    `bson:"status"`
	ProfilePictureUrl string    `bson:"profile_picture_url"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateParticipantAccountData struct {
	CreateAccountData
	ParticipantId       string    `bson:"participant_id"`
	Skillset            []string  `bson:"skillset"`
	DOB                 time.Time `bson:"dob"`
	EmploymentStatus    string    `bson:"employment_status"`
	ExperienceLevel     string    `bson:"experience_level"`
	Motivation          string    `bson:"motivation"`
	HackathonExperience string    `bson:"hackathon_experience"`
	YearsOfExperience   int       `bson:"years_of_experience"`
	FieldOfStudy        string    `bson:"field_of_study"`
	PreviousProjects    []string  `bson:"previous_projects"`
}

type CreateTeamMemberAccountData struct {
	CreateAccountData
	ParticipantId    string    `bson:"participant_id"`
	TeamRole         string    `bson:"team_role"`
	Skillset         []string  `bson:"skillset"`
	DOB              time.Time `bson:"dob"`
	EmploymentStatus string    `bson:"employment_status"`
	ExperienceLevel  string    `bson:"experience_level"`
	Motivation       string    `bson:"motivation"`
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
	ParticipantId     string
	ParticipantEmail  string
	TeamLeadFirstName string
	TeamLeadLastName  string
	TeamLeadGender    string
	TeamLeadAccountId string
	TeamLeadEmail     string
	TeamName          string
	CoParticipants    []CoParticipant
	Type              string
	HackathonId       string
	GithubAddress     string
	ReviewRanking     int
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
	Bio string `bson:"bio,omitempty"`
}

type UpdateAccountData struct {
	FirstName         string `bson:"first_name"`
	LastName          string `bson:"last_name"`
	Gender            string `bson:"gender"`
	State             string `bson:"state"`
	Role              string `bson:"role"`
	PhoneNumber       string `bson:"phone_number"`
	Status            string `bson:"status"`
	ProfilePictureUrl string `bson:"profile_picture_url"`
	Bio               string `bson:"bio"`
}

type CreateSolutionData struct {
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Objective   string `bson:"objective,omitempty"`
	CreatorId   string `bson:"creator_id"`
	HackathonId string `bson:"hackathon_id"`
}

type GetSolutionsQueryData struct {
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Objective   string `bson:"objective,omitempty"`
	CreatorId   string `bson:"creator_id"`
	HackathonId string `bson:"hackathon_id"`
}

type UpdateSolutionData struct {
	Title            string `bson:"title,omitempty"`
	Description      string `bson:"description,omitempty"`
	CreatorId        string `bson:"creator_id"`
	Objective        string `bson:"objective,omitempty"`
	SolutionImageUrl string `bson:"solution_image_url,omitempty"`
}

type SelectTeamSolutionData struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
	SolutionId    string `bson:"solution_id"`
}
