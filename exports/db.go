package exports

import (
	"time"
)

type MongoDBConnConfig struct {
	Url    string
	DBName string
}
type AccountDocument struct {
	Id                  interface{} `json:"_id"`
	Email               string      `bson:"email,omitempty"`
	FirstName           string      `bson:"first_name,omitempty"`
	LastName            string      `bson:"last_name,omitempty"`
	Gender              string      `bson:"gender,omitempty"`
	LinkedInAddress     string      `bson:"linkedIn_address,omitempty"`
	PasswordHash        string      `bson:"password_hash,omitempty"`
	PhoneNumber         string      `bson:"phone_number,omitempty"`
	Skillset            []string    `bson:"skillset,omitempty"`
	ParticipantId       string      `bson:"participant_id,omitempty"`
	HackathonId         string      `bson:"hackathon_id,omitempty"`
	State               string      `bson:"state,omitempty"`
	Role                string      `bson:"role,omitempty"`
	DOB                 time.Time   `bson:"dob,omitempty"`
	Bio                 string      `bson:"bio,omitempty"`
	EmploymentStatus    string      `bson:"employment_status,omitempty"`
	ExperienceLevel     string      `bson:"experience_level,omitempty"`
	Motivation          string      `bson:"motivation,omitempty"`
	HackathonExperience string      `bson:"hackathon_experience,omitempty"`
	YearsOfExperience   int         `bson:"years_of_experience,omitempty"`
	FieldOfStudy        string      `bson:"field_of_study,omitempty"`
	PreviousProjects    []string    `bson:"previous_projects,omitempty"`
	IsEmailVerified     bool        `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt   time.Time   `bson:"is_email_verified_at,omitempty"`
	Status              string      `bson:"status,omitempty"`
	ProfilePictureUrl   string      `bson:"profile_picture_url,omitempty"`
	CreatedAt           time.Time   `bson:"created_at,omitempty"`
	UpdatedAt           time.Time   `bson:"updated_at,omitempty"`
}

type ParticipantDocument struct {
	Id               interface{}
	ParticipantId    string           `bson:"participant_id"`
	HackathonId      string           `bson:"hackathon_id"`
	Type             string           `bson:"type,omitempty"`
	TeamLeadEmail    string           `bson:"team_lead_email,omitempty"`
	SolutionId       string           `bson:"solution_id,omitempty"`
	Solution         SolutionDocument `bson:"solution"`
	TeamName         string           `bson:"team_name,omitempty"`
	CoParticipants   []CoParticipant  `bson:"co_participants,omitempty"`
	ParticipantEmail string           `bson:"participant_email,omitempty"`
	GithubAddress    string           `bson:"github_address,omitempty"`
	InviteList       []InviteInfo     `bson:"invite_list,omitempty"`
	ReviewRanking    int              `bson:"review_ranking,omitempty"`
	Status           string           `bson:"status,omitempty"`
	CreatedAt        time.Time        `bson:"created_at,omitempty"`
	UpdatedAt        time.Time        `bson:"updated_at,omitempty"`
}

type CoParticipant struct {
	Email string `bson:"email,omitempty"`
	Role  string `bson:"role,omitempty"`
}

type InviteInfo struct {
	Email     string    `bson:"email,omitempty" json:"email,omitempty"`
	Time      time.Time `bson:"time,omitempty" json:"time,omitempty"`
	InviterId string    `bson:"inviter_id,omitempty" json:"inviter_id,omitempty"`
}

type ParticipantScoreDocument struct {
	Id            interface{}
	HackathonId   string      `validate:"required" json:"hackathon_id"`
	ParticipantId string      `validate:"required" json:"participant_id"`
	Stage         string      `validate:"required" json:"stage"`
	TotalScore    float32     `validate:"required" json:"score"`
	ScoresInfo    []ScoreInfo `validate:"required" json:"scores_info"`
}

type SolutionDocument struct {
	Id               interface{} `bson:"_id,omitempty" `
	Title            string      `bson:"title,omitempty"`
	Objective        string      `bson:"objective,omitempty"`
	Description      string      `bson:"description,omitempty"`
	HackathonId      string      `bson:"hackathon_id,omitempty"`
	CreatorId        string      `bson:"creator_id,omitempty"`
	SolutionImageUrl string      `bson:"solution_image_url,omitempty"`
	CreatedAt        time.Time   `bson:"created_at,omitempty"`
	UpdatedAt        time.Time   `bson:"updated_at,omitempty"`
}

type ParticipantAccountWithCoParticipantsDocument struct {
	Id                  interface{} `bson:"_id"`
	Email               string      `bson:"email,omitempty"`
	FirstName           string      `bson:"first_name,omitempty"`
	LastName            string      `bson:"last_name,omitempty"`
	Gender              string      `bson:"gender,omitempty"`
	PhoneNumber         string      `bson:"phone_number,omitempty"`
	Skillset            []string    `bson:"skillset,omitempty"`
	ParticipantId       string      `bson:"participant_id,omitempty"`
	HackathonId         string      `bson:"hackathon_id"`
	State               string      `bson:"state,omitempty"`
	Role                string      `bson:"role,omitempty"`
	DOB                 time.Time   `bson:"dob,omitempty"`
	Bio                 string      `bson:"bio,omitempty"`
	EmploymentStatus    string      `bson:"employment_status,omitempty"`
	ExperienceLevel     string      `bson:"experience_level,omitempty"`
	Motivation          string      `bson:"motivation,omitempty"`
	HackathonExperience string      `bson:"hackathon_experience"`
	YearsOfExperience   int         `bson:"years_of_experience"`
	FieldOfStudy        string      `bson:"field_of_study"`
	PreviousProjects    []string    `bson:"previous_projects"`
	IsEmailVerified     bool        `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt   time.Time   `bson:"is_email_verified_at,omitempty"`
	Status              string      `bson:"status,omitempty"`
	ProfilePictureUrl   string      `bson:"profile_picture_url,omitempty"`
	LinkedInAddress     string      `bson:"linkedIn_address,omitempty"`
	CreatedAt           time.Time   `bson:"created_at,omitempty"`
	UpdatedAt           time.Time   `bson:"updated_at,omitempty"`

	Type             string                           `bson:"type,omitempty"`
	TeamName         string                           `bson:"team_name,omitempty"`
	SolutionId       string                           `bson:"solution_id,omitempty"`
	CoParticipants   []CoParticipantAggregateDocument `bson:"co_participants,omitempty"`
	ParticipantEmail string                           `bson:"participant_email,omitempty"`
	GithubAddress    string                           `bson:"github_address,omitempty"`
	InviteList       []InviteInfo                     `bson:"invite_list,omitempty"`
}

type CoParticipantAggregateDocument struct {
	Email               string    `bson:"email,omitempty"`
	AccountRole         string    `bson:"account_role,omitempty"`
	FirstName           string    `bson:"first_name,omitempty"`
	LastName            string    `bson:"last_name,omitempty"`
	Gender              string    `bson:"gender,omitempty"`
	PhoneNumber         string    `bson:"phone_number,omitempty"`
	Skillset            []string  `bson:"skillset,omitempty"`
	ParticipantId       string    `bson:"participant_id,omitempty"`
	HackathonId         string    `bson:"hackathon_id,omitempty"`
	State               string    `bson:"state,omitempty"`
	TeamRole            string    `bson:"team_role,omitempty"`
	DOB                 time.Time `bson:"dob,omitempty"`
	Bio                 string    `bson:"bio,omitempty"`
	EmploymentStatus    string    `bson:"employment_status,omitempty"`
	ExperienceLevel     string    `bson:"experience_level,omitempty"`
	Motivation          string    `bson:"motivation,omitempty"`
	HackathonExperience string    `json:"hackathon_experience"`
	YearsOfExperience   int       `json:"years_of_experience"`
	FieldOfStudy        string    `json:"field_of_study"`
	PreviousProjects    []string  `json:"previous_projects"`
	IsEmailVerified     bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt   time.Time `bson:"is_email_verified_at,omitempty"`
	Status              string    `bson:"status,omitempty"`
	ProfilePictureUrl   string    `bson:"profile_picture_url,omitempty"`
	LinkedInAddress     string    `bson:"linkedIn_address,omitempty"`
	CreatedAt           time.Time `bson:"created_at,omitempty"`
	UpdatedAt           time.Time `bson:"updated_at,omitempty"`
}
