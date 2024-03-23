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
	AccessToken  string
	RefreshToken string
	Exp          time.Time
}

type jwtCustomClaims struct {
	Email     string `json:"email"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

type AuthUtilsPayload struct {
	Email     string
	LastName  string
	FirstName string
	Role      string
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
