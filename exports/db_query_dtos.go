package exports

import (
	"time"
)

type UpdateAccountDocumentFilter struct {
	Email       string `bson:"email"`
	PhoneNumber string `bson:"phone_number"`
	Status      string `bson:"status"`
}

type UpdateAccountDocument struct {
	FirstName         string    `bson:"first_name,omitempty"`
	LastName          string    `bson:"last_name,omitempty"`
	Gender            string    `bson:"gender,omitempty"`
	LinkedInAddress   string    `bson:"linkedIn_address,omitempty"`
	GithubAddress     string    `bson:"github_address,omitempty"`
	State             string    `bson:"state,omitempty"`
	Bio               string    `bson:"bio,omitempty"`
	Status            string    `bson:"status,omitempty"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
	ProfilePictureUrl string    `bson:"profile_picture_url"`
}

type CreateAdminAccountData struct {
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	FirstName    string    `bson:"first_name"`
	LastName     string    `bson:"last_name"`
	Gender       string    `bson:"gender"`
	HackathonId  string    `bson:"hackathon_id"`
	Role         string    `bson:"role"`
	PhoneNumber  string    `bson:"phone_number"`
	Status       string    `bson:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
	CreatedAt         time.Time `bson:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at"`
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
	HackathonId      string `bson:"hackathon_id"`
	ParticipantId    string `bson:"participant_id"`
	InviterEmail     string `bson:"inviter_email"`
	InviterFirstName string `bson:"inviter_first_name"`
	Email            string `bson:"email"`
	Role             string `bson:"role"`
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
	DOB                 string   ` json:"dob"`
	TeamName            string   `validate:"omitempty" json:"team_name"`
	EmploymentStatus    string   `validate:"oneof=STUDENT EMPLOYED UNEMPLOYED FREELANCER" json:"employment_status"`
	ExperienceLevel     string   `validate:"oneof=JUNIOR MID SENIOR" json:"experience_level"`
	Motivation          string   `validate:"min=100" json:"motivation"`
	HackathonExperience string   `json:"hackathon_experience"`
	YearsOfExperience   int      `json:"years_of_experience"`
	FieldOfStudy        string   `json:"field_of_study"`
	PreviousProjects    []string `json:"previous_projects"`
}

type GetParticipantsFilterOpts struct {
	ParticipantId            *string
	ParticipantStatus        *string `validate:"omitempty, oneof UNREVIEWED REVIEWED AI_RANKED "`
	ReviewRanking_Eq         *int
	ReviewRanking_Top        *int
	Solution_Like            *string
	Limit                    *int
	SortByReviewRanking_Asc  *bool
	SortByReviewRanking_Desc *bool
}

type UpdateAdminAccountDocument struct {
	FirstName         string    `bson:"first_name,omitempty"`
	LastName          string    `bson:"last_name,omitempty"`
	Gender            string    `bson:"gender,omitempty"`
	State             string    `bson:"state,omitempty"`
	Bio               string    `bson:"bio,omitempty"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
	ProfilePictureUrl string    `bson:"profile_picture_url"`
}

type FilterGetManyAccountDocuments struct {
	Email_eq string
}

type UpdateSingleParticipantRecordFilter struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
}

type UpdateManyParticipantRecordFilter struct {
	HackathonId   string `bson:"hackathon_id"`
	ParticipantId string `bson:"participant_id"`
	Status        string `bson:"status"`
	Role          string `bson:"role"`
}
type UpdateParticipantRecordData struct {
	Status        string `bson:"status"`
	ReviewRanking int    `bson:"review_ranking"`
	TeamName      string `bson:"team_name"`
}

type AdminParticipantInfoUpdateDTO struct {
	Status        string `json:"status,omitempty"`
	ReviewRanking int    `json:"review_rank,omitempty"`
}
