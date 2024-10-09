package exports

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParticipantTeamMembersWithAccountsAggregateDocument struct {
	Id               primitive.ObjectID                             `bson:"_id"`
	ParticipantId    string                                         `bson:"participant_id"`
	HackathonId      string                                         `bson:"hackathon_id"`
	Type             string                                         `bson:"type,omitempty"`
	TeamLeadEmail    string                                         `bson:"team_lead_email,omitempty"`
	SolutionId       string                                         `bson:"solution_id,omitempty"`
	Solution         ParticipantDocumentParticipantSelectedSolution `bson:"solution_document"`
	TeamName         string                                         `bson:"team_name,omitempty"`
	ParticipantEmail string                                         `bson:"participant_email,omitempty"`
	GithubAddress    string                                         `bson:"github_address,omitempty"`
	InviteList       []ParticipantDocumentTeamInviteInfo            `bson:"invite_list,omitempty"`
	ReviewRanking    int                                            `bson:"review_ranking,omitempty"`
	Status           string                                         `bson:"status,omitempty"`
	CreatedAt        time.Time                                      `bson:"created_at,omitempty"`
	UpdatedAt        time.Time                                      `bson:"updated_at,omitempty"`

	EmploymentStatus    string    `bson:"employment_status,omitempty"`
	ExperienceLevel     string    `bson:"experience_level,omitempty"`
	Motivation          string    `bson:"motivation,omitempty"`
	HackathonExperience string    `bson:"hackathon_experience"`
	YearsOfExperience   int       `bson:"years_of_experience"`
	FieldOfStudy        string    `bson:"field_of_study"`
	PreviousProjects    []string  `bson:"previous_projects"`
	IsEmailVerified     bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt   time.Time `bson:"is_email_verified_at,omitempty"`
	ProfilePictureUrl   string    `bson:"profile_picture_url,omitempty"`
	LinkedInAddress     string    `bson:"linkedIn_address,omitempty"`
	TeamLeadFirstName   string    `bson:"team_lead_first_name,omitempty"`
	TeamLeadLastName    string    `bson:"team_lead_last_name,omitempty"`
	TeamLeadGender      string    `bson:"team_lead_gender,omitempty"`
	TeamLeadAccountId   string    `bson:"team_lead_account_id,omitempty"`

	TeamLeadInfo   TeamLeadAndAccountAggregateData          `bson:"team_lead_info,omitempty"`
	CoParticipants []CoParticipantsAndAccountsAggregateData `bson:"co_participants,omitempty"`
}

type TeamLeadAndAccountAggregateData struct {
	HackathonId         string    `bson:"team_lead_hackathon_id"`
	AccountId           string    `bson:"id,omitempty"`
	Email               string    `bson:"email"`
	FirstName           string    `bson:"first_name"`
	LastName            string    `bson:"last_name"`
	Gender              string    `bson:"gender"`
	PhoneNumber         string    `bson:"phone_number"`
	Skillset            []string  `bson:"skillset"`
	AccountStatus       string    `bson:"status"`
	AccountRole         string    `bson:"role"`
	TeamRole            string    `bson:"team_role"`
	PasswordHash        string    `bson:"password_hash"`
	State               string    `bson:"state"`
	EmploymentStatus    string    `bson:"employment_status,omitempty"`
	ExperienceLevel     string    `bson:"experience_level,omitempty"`
	Motivation          string    `bson:"motivation,omitempty"`
	HackathonExperience string    `bson:"hackathon_experience"`
	YearsOfExperience   int       `bson:"years_of_experience"`
	FieldOfStudy        string    `bson:"field_of_study"`
	PreviousProjects    []string  `bson:"previous_projects"`
	IsEmailVerified     bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt   time.Time `bson:"is_email_verified_at,omitempty"`
	//Status              string    `bson:"status,omitempty"`
	ProfilePictureUrl string    `bson:"profile_picture_url,omitempty"`
	LinkedInAddress   string    `bson:"linkedIn_address,omitempty"`
	DOB               time.Time `bson:"dob,omitempty"`
	CreatedAt         time.Time `bson:"created_at"`
	UpdateAt          time.Time `bson:"update_at"`
}

type CoParticipantsAndAccountsAggregateData struct {
	HackathonId         string    `bson:"hackathon_id"`
	AccountId           string    `bson:"id,omitempty"`
	Email               string    `bson:"email"`
	FirstName           string    `bson:"first_name"`
	LastName            string    `bson:"last_name"`
	Gender              string    `bson:"gender"`
	PhoneNumber         string    `bson:"phone_number"`
	Skillset            []string  `bson:"skillset"`
	AccountStatus       string    `bson:"status"`
	AccountRole         string    `bson:"account_role"`
	TeamRole            string    `bson:"team_role"`
	PasswordHash        string    `bson:"password_hash"`
	ParticipantId       string    `bson:"participant_id"`
	PreviousProjects    []string  `bson:"previous_projects"`
	EmploymentStatus    string    `bson:"employment_status,omitempty"`
	ExperienceLevel     string    `bson:"experience_level,omitempty"`
	Motivation          string    `bson:"motivation,omitempty"`
	HackathonExperience string    `bson:"hackathon_experience"`
	YearsOfExperience   int       `bson:"years_of_experience"`
	FieldOfStudy        string    `bson:"field_of_study"`
	IsEmailVerified     bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt   time.Time `bson:"is_email_verified_at,omitempty"`
	//Status              string    `bson:"status,omitempty"`
	ProfilePictureUrl string    `bson:"profile_picture_url,omitempty"`
	LinkedInAddress   string    `bson:"linkedIn_address,omitempty"`
	State             string    `bson:"state"`
	DOB               time.Time `bson:"dob,omitempty"`
	CreatedAt         time.Time `bson:"created_at"`
	UpdateAt          time.Time `bson:"update_at"`
}
type GetParticipantsWithAccountsAggregateFilterOpts struct {
	ParticipantId            *string
	ParticipantStatus        *string `validate:"omitempty, oneof UNREVIEWED REVIEWED AI_RANKED REVIEW_DISQUALIFIED TEAM_ONBOARDING SOLUTION_SELECTION SOLUTION_IMPLEMENTATION SHORTLISTED COMPETITION_WON"`
	ParticipantType          *string `validate:"omitempty, oneof TEAM "`
	ReviewRanking_Eq         *int
	ReviewRanking_Top        *int
	Solution_Like            *string
	Limit                    *int
	SortByReviewRanking_Asc  *bool
	SortByReviewRanking_Desc *bool
}
