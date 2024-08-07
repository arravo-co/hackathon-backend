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
	Role                string       `json:"role"`
	HackathonId         string       `json:"hackathon_id"`
	Status              string       `json:"status"`
	PhoneNumber         string       `json:"phone_number"`
	Age                 int          `json:"age"`
	DOB                 time.Time    `json:"dob"`
	ParticipantId       string       `json:"participant_id"`
	TeamLeadEmail       string       `json:"team_lead_email"`
	TeamName            string       `json:"team_name"`
	TeamRole            string       `json:"team_role"`
	Type                string       `json:"type"`
	CoParticipantEmails []string     `json:"co_participant_emails"`
	ParticipantEmail    string       `json:"participant_email"`
	InviteList          []InviteInfo `json:"invite_list"`
	Skillset            []string     `json:"skillset"`
	Solution            *struct {
		Id               string `json:"id"`
		Title            string `json:"title"`
		Description      string `json:"description"`
		SolutionImageUrl string `json:"solution_image_url"`
		CreatorId        string `json:"creator_id"`
	} `json:"solution"`
	Bio                 string   `bson:"bio,omitempty"`
	EmploymentStatus    string   `bson:"employment_status,omitempty"`
	ExperienceLevel     string   `bson:"experience_level,omitempty"`
	Motivation          string   `bson:"motivation,omitempty"`
	HackathonExperience string   `bson:"hackathon_experience,omitempty"`
	YearsOfExperience   int      `bson:"years_of_experience,omitempty"`
	FieldOfStudy        string   `bson:"field_of_study,omitempty"`
	PreviousProjects    []string `bson:"previous_projects,omitempty"`
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
