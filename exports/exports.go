package exports

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
)

type AuthUtilsInterface interface {
	BasicLogin(dataInput *AuthUtilsBasicLoginData) (*AuthUtilsBasicLoginSuccessData, error)
	GenerateAccessToken(payload *AuthUtilsPayload) (string, error)
	GetJWTConfig() echojwt.Config
	VerifyToken(dataInput *AuthUtilsVerifyTokenData) error
	InitiateEmailVerification(dataInput *AuthUtilsConfigTokenData) (*TokenData, error)
	CompleteEmailVerification(dataInput *AuthUtilsCompleteEmailVerificationData) error
	ChangePassword(dataInput *AuthUtilsChangePasswordData) error
	InitiatePasswordRecovery(dataInput *AuthUtilsConfigTokenData) (*TokenData, error)
	CompletePasswordRecovery(dataInput *AuthUtilsCompletePasswordRecoveryData) (interface{}, error)
}
type AuthUtilsBasicLoginData struct {
	Identifier string
	Password   string
	Role       string
}

type AuthUtilsBasicLoginSuccessData struct {
	AccessToken         string    `json:"access_token"`
	RefreshToken        string    `json:"refresh_token"`
	Exp                 time.Time `json:"exp"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Email               string    `json:"email"`
	Gender              string    `json:"gender"`
	State               string    `json:"state"`
	passwordHash        string
	AccountRole         string                                  `json:"role"`
	HackathonId         string                                  `json:"hackathon_id"`
	AccountStatus       string                                  `json:"account_status"`
	ParticipantStatus   string                                  `json:"participant_status"`
	PhoneNumber         string                                  `json:"phone_number"`
	Age                 int                                     `json:"age"`
	DOB                 time.Time                               `json:"dob"`
	ParticipantId       string                                  `json:"participant_id"`
	TeamLeadEmail       string                                  `json:"team_lead_email"`
	TeamName            string                                  `json:"team_name"`
	TeamRole            string                                  `json:"team_role"`
	Type                string                                  `json:"type"`
	CoParticipants      []AuthUtilsParticipantCoParticipantInfo `json:"co_participants"`
	ParticipantEmail    string                                  `json:"participant_email"`
	InviteList          []ParticipantDocumentTeamInviteInfo     `json:"invite_list"`
	Skillset            []string                                `json:"skillset"`
	Solution            *AuthUtilsParticipantSolutionInfo       `json:"solution"`
	Bio                 string                                  `json:"bio,omitempty"`
	EmploymentStatus    string                                  `json:"employment_status,omitempty"`
	ExperienceLevel     string                                  `json:"experience_level,omitempty"`
	Motivation          string                                  `json:"motivation,omitempty"`
	HackathonExperience string                                  `json:"hackathon_experience,omitempty"`
	YearsOfExperience   int                                     `json:"years_of_experience,omitempty"`
	FieldOfStudy        string                                  `json:"field_of_study,omitempty"`
	PreviousProjects    []string                                `json:"previous_projects,omitempty"`
	ProfilePictureUrl   string                                  `json:"profile_picture_url,omitempty"`
	CreatedAt           time.Time                               `json:"created_at,omitempty"`
	UpdatedAt           time.Time                               `json:"updated_at,omitempty"`
}

type AuthUtilsParticipantCoParticipantInfo struct {
	AccountId           string    `json:"account_id"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Email               string    `json:"email"`
	Gender              string    `json:"gender"`
	State               string    `json:"state"`
	Age                 int       `json:"age"`
	DOB                 time.Time `json:"dob"`
	AccountStatus       string    `json:"account_status"`
	AccountRole         string    `json:"account_role"`
	TeamRole            string    `json:"team_role"`
	ParticipantId       string    `json:"participant_id"`
	HackathonId         string    `json:"hackathon_id"`
	Skillset            []string  `json:"skillset"`
	PhoneNumber         string    `json:"phone_number"`
	EmploymentStatus    string    `json:"employment_status"`
	ExperienceLevel     string    `json:"experience_level"`
	Motivation          string    `json:"motivation"`
	HackathonExperience string    `json:"hackathon_experience"`
	YearsOfExperience   int       `json:"years_of_experience"`
	FieldOfStudy        string    `json:"field_of_study"`
	PreviousProjects    []string  `json:"previous_projects"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type AuthUtilsParticipantSolutionInfo struct {
	Id               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Objective        string `json:"objective"`
	SolutionImageUrl string `json:"solution_image_url"`
	CreatorId        string `json:"creator_id"`
}

type jwtCustomClaims struct {
	Email     string `json:"email"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

type AuthUtilsPayload struct {
	AccountId       string
	HackathonId     string
	Email           string
	LastName        string
	FirstName       string
	Role            string
	IsParticipant   bool
	ParticipantType string
	ParticipantId   string
}

type AuthUtilsVerifyTokenData struct {
	Email          string
	Token          string
	Scope          string
	TokenType      string
	TokenTypeValue string
}

type AuthUtilsConfigTokenData struct {
	Email string
	TTL   time.Time
}

type AuthUtilsCompleteEmailVerificationData struct {
	Email string
	Token string
}

type AuthUtilsCompletePasswordRecoveryData struct {
	Email       string
	Token       string
	NewPassword string
}

type AuthUtilsChangePasswordData struct {
	Email       string
	OldPassword string
	NewPassword string
}
