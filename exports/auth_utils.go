package exports

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type BasicLoginData struct {
	Identifier string
	Password   string
	Role       string
}

type BasicLoginSuccessData struct {
	AccessToken  string
	RefreshToken string
	Exp          time.Time
}

type JwtCustomClaims struct {
	Email     string `json:"email"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

type Payload struct {
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

type ConfigTokenData struct {
	Email string
	TTL   time.Time
}

type CompleteEmailVerificationData struct {
	Email string
	Token string
}

type CompletePasswordRecoveryData struct {
	Email       string
	Token       string
	NewPassword string
}

type ChangePasswordData struct {
	Email       string
	OldPassword string
	NewPassword string
}

type MyJWTCustomClaims struct {
	Email           string `json:"email"`
	LastName        string `json:"last_name"`
	FirstName       string `json:"first_name"`
	Role            string `json:"role"`
	ParticipantId   string `json:"participant_id"`
	AccountId       string `json:"account_id"`
	HackathonId     string `json:"hackathon_id"`
	IsParticipant   bool   `json:"is_participant"`
	ParticipantType string `json:"participant_type"`
	jwt.RegisteredClaims
}

func GetPayload(c echo.Context) *Payload {
	jwtData := c.Get("user").(*jwt.Token)
	claims := jwtData.Claims.(*MyJWTCustomClaims)
	tokenData := Payload{
		Email:     claims.Email,
		LastName:  claims.LastName,
		FirstName: claims.FirstName,
		Role:      claims.Role,
	}
	return &tokenData
}
