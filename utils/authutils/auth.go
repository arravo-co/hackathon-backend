package authutils

import (
	"errors"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jaevor/go-nanoid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthUtils struct {
}

func BasicLogin(dataInput *exports.AuthUtilsBasicLoginData) (*exports.AuthUtilsBasicLoginSuccessData, error) {
	accountDoc, err := data.FindAccountIdentifier(dataInput.Identifier)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, errors.New("no email or username provided that matches record")
	}
	_, err = exports.ComparePasswordAndHash(dataInput.Password, accountDoc.PasswordHash)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	accessToken, err := GenerateAccessToken(&exports.AuthUtilsPayload{
		Email:     accountDoc.Email,
		LastName:  accountDoc.LastName,
		FirstName: accountDoc.FirstName,
		Role:      accountDoc.Role,
	})
	return &exports.AuthUtilsBasicLoginSuccessData{
		AccessToken: accessToken,
	}, err
}

func GenerateAccessToken(payload *exports.AuthUtilsPayload) (string, error) {
	claims := &MyJWTCustomClaims{
		payload.Email,
		payload.FirstName,
		payload.LastName,
		payload.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.GetSecretKey()))
	return t, err
}

func GetJWTConfig() echojwt.Config {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(MyJWTCustomClaims)
		},
		SigningKey: []byte(config.GetSecretKey()),
	}
	return config
}

func VerifyToken(dataInput *exports.AuthUtilsVerifyTokenData) error {
	_, err := data.GetAccountByEmail(dataInput.TokenTypeValue)
	if err != nil {
		return errors.New("email not found in record")
	}
	err = data.VerifyToken(&exports.VerifyTokenData{
		Token:          dataInput.Token,
		TokenType:      dataInput.TokenType,
		TokenTypeValue: dataInput.TokenTypeValue,
		Scope:          dataInput.Scope,
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	return nil
}

func InitiateEmailVerification(dataInput *exports.AuthUtilsConfigTokenData) (*exports.TokenData, error) {
	_, err := data.GetAccountByEmail(dataInput.Email)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, errors.New("email not found in record")
	}
	tokenFunc, _ := nanoid.Custom("1234567890", 6)
	token := tokenFunc()
	tokenData, err := data.UpsertToken(&exports.UpsertTokenData{
		Token:          token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		TTL:            dataInput.TTL,
		Status:         "PENDING",
		Scope:          "EMAIL_VERIFICATION",
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, errors.New("failed to generate token: ")
	}
	return tokenData, nil
}

func CompleteEmailVerification(dataInput *exports.AuthUtilsCompleteEmailVerificationData) error {
	err := VerifyToken(&exports.AuthUtilsVerifyTokenData{
		Email:          dataInput.Email,
		Token:          dataInput.Token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		Scope:          "EMAIL_VERIFICATION",
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return err
	}
	_, err = data.UpdateParticipantInfoByEmail(&exports.UpdateAccountFilter{
		Email: dataInput.Email,
	}, &exports.UpdateAccountDocument{
		IsEmailVerified:   true,
		IsEmailVerifiedAt: time.Now(),
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return errors.New("unable to complete token verification")
	}

	return nil
}

func ChangePassword(dataInput *exports.AuthUtilsChangePasswordData) error {
	accountDoc, err := data.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return errors.New("user info not found in record")
	}
	_, err = exports.ComparePasswordAndHash(dataInput.OldPassword, accountDoc.PasswordHash)
	if err != nil {
		return err
	}
	hash, err := exports.GenerateHashPassword(dataInput.NewPassword)
	if err != nil {
		return err
	}
	accountDoc, err = data.UpdatePasswordByEmail(&exports.UpdateAccountFilter{Email: dataInput.Email}, hash)

	// emit an emit here
	return nil
}

func InitiatePasswordRecovery(dataInput *exports.AuthUtilsConfigTokenData) (*exports.TokenData, error) {
	_, err := data.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return nil, errors.New("email not found in record")
	}
	tokenFunc, _ := nanoid.Custom("1234567890", 6)
	token := tokenFunc()
	tokenData, err := data.UpsertToken(&exports.UpsertTokenData{
		Token:          token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		TTL:            dataInput.TTL,
		Status:         "PENDING",
		Scope:          "PASSWORD_RECOVERY",
	})
	return tokenData, err
}

func CompletePasswordRecovery(dataInput *exports.AuthUtilsCompletePasswordRecoveryData) (interface{}, error) {
	err := VerifyToken(&exports.AuthUtilsVerifyTokenData{
		Token:          dataInput.Token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		Scope:          "PASSWORD_RECOVERY",
	})
	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, err
	}
	newPasswordHash, err := exports.GenerateHashPassword(dataInput.NewPassword)
	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, err
	}
	_, err = data.UpdatePasswordByEmail(&exports.UpdateAccountFilter{Email: dataInput.Email}, newPasswordHash)

	if err != nil {
		exports.MySugarLogger.Error(err)
		return nil, errors.New("unable to complete password recovery")
	}

	return nil, nil
}
