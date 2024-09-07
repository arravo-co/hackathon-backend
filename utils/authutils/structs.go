package authutils

import (
	"time"
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

type Payload struct {
	Email     string
	LastName  string
	FirstName string
	Role      string
}

type VerifyTokenData struct {
	Email          string
	Token          string
	Scope          string
	TokenType      string
	TokenTypeValue string
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

type GenerateTokenArgs struct {
	Token          string
	TokenType      string
	TokenTypeValue string
	Scope          string
	TTL            time.Time
	Status         string
}

type TokenData struct {
	Id             string
	Token          string
	TokenType      string
	TokenTypeValue string
	Scope          string
	TTL            time.Time
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
