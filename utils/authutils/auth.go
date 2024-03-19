package authutils

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/arravoco/hackathon_backend/utils/email"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jaevor/go-nanoid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func BasicLogin(dataInput *BasicLoginData) (*BasicLoginSuccessData, error) {
	accountCol, err := db.GetAccountCollection()
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}
	filter := bson.D{{
		"$or", bson.A{
			bson.D{{"email", dataInput.Identifier}},
			bson.D{{"username", dataInput.Identifier}},
		}},
	}
	dataFromCol := data.AccountDocument{}
	err = accountCol.FindOne(context.TODO(), filter).Decode(&dataFromCol)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, errors.New("no email or username provided that matches record")
	}
	fmt.Printf("%+v\n", dataFromCol)
	_, err = utils.ComparePasswordAndHash(dataInput.Password, dataFromCol.PasswordHash)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	accessToken, err := GenerateAccessToken(&Payload{
		Email:     dataFromCol.Email,
		LastName:  dataFromCol.LastName,
		FirstName: dataFromCol.FirstName,
		Role:      dataFromCol.Role,
	})
	return &BasicLoginSuccessData{
		AccessToken: accessToken,
	}, err
}

func GenerateAccessToken(payload *Payload) (string, error) {
	claims := &jwtCustomClaims{
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
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(config.GetSecretKey()),
	}
	return config
}

func VerifyToken(dataInput *VerifyTokenData) error {
	_, err := data.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return errors.New("email not found in record")
	}
	err = data.VerifyToken(&data.VerifyTokenData{
		Token:          dataInput.Token,
		TokenType:      dataInput.TokenType,
		TokenTypeValue: dataInput.TokenTypeValue,
		Scope:          dataInput.Scope,
	})
	if err != nil {
		utils.MySugarLogger.Error(err)
		return err
	}
	return nil
}

func InitiateEmailVerification(dataInput *ConfigTokenData) error {
	_, err := data.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return errors.New("email not found in record")
	}
	tokenFunc, _ := nanoid.Custom("1234567890", 6)
	token := tokenFunc()
	tokenData, err := data.CreateToken(&data.CreateTokenData{
		Token:          token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		TTL:            dataInput.TTL,
		Status:         "PENDING",
		Scope:          "EMAIL_VERIFICATION",
	})
	if err != nil {
		utils.MySugarLogger.Error(err)
		return errors.New("failed to generate token: ")
	}
	email.SendEmailVerificationEmail(&email.SendEmailVerificationEmailData{
		Email:    dataInput.Email,
		Token:    tokenData.Token,
		TokenTTL: dataInput.TTL,
		Subject:  "Email Verification",
	})
	return nil
}

func CompleteEmailVerification(dataInput *CompleteEmailVerificationData) error {
	err := VerifyToken(&VerifyTokenData{
		Token:          dataInput.Token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		Scope:          "EMAIL_VERIFICATION",
	})
	if err != nil {
		utils.MySugarLogger.Error(err)
		return err
	}
	_, err = data.UpdateParticipantInfoByEmail(&data.UpdateAccountFilter{
		Email: dataInput.Email,
	}, &data.UpdateAccountDocument{
		IsEmailVerified:   true,
		IsEmailVerifiedAt: time.Now(),
	})
	if err != nil {
		utils.MySugarLogger.Error(err)
		return errors.New("unable to complete token verification")
	}

	email.SendEmailVerificationCompleteEmail(&email.SendEmailVerificationCompleteEmailData{
		Email:   dataInput.Email,
		Subject: "Email Verification Success",
	})
	return nil
}

func ChangePassword(dataInput *ChangePasswordData) error {
	accountDoc, err := data.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return errors.New("user info not found in record")
	}
	_, err = utils.ComparePasswordAndHash(dataInput.OldPassword, accountDoc.PasswordHash)
	if err != nil {
		return err
	}
	hash, err := utils.GenerateHashPassword(dataInput.NewPassword)
	if err != nil {
		return err
	}
	accountDoc, err = data.UpdatePasswordByEmail(&data.UpdateAccountFilter{Email: dataInput.Email}, hash)

	// emit an emit here
	return nil
}

func InitiatePasswordRecovery(dataInput *ConfigTokenData) error {
	_, err := data.GetAccountByEmail(dataInput.Email)
	if err != nil {
		return errors.New("email not found in record")
	}
	tokenFunc, _ := nanoid.Custom("1234567890", 6)
	token := tokenFunc()
	tokenData, err := data.CreateToken(&data.CreateTokenData{
		Token:          token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		TTL:            dataInput.TTL,
		Status:         "PENDING",
		Scope:          "PASSWORD_RECOVERY",
	})
	if err != nil {
		utils.MySugarLogger.Error(err)
		return errors.New("failed to generate token: ")
	}
	email.SendPasswordRecoveryEmail(&email.SendPasswordRecoveryEmailData{
		Email:    dataInput.Email,
		Token:    tokenData.Token,
		TokenTTL: dataInput.TTL,
		Subject:  "Password Recovery",
	})
	return nil
}

func CompletePasswordRecovery(dataInput *CompletePasswordRecoveryData) error {
	err := VerifyToken(&VerifyTokenData{
		Token:          dataInput.Token,
		TokenType:      "EMAIL",
		TokenTypeValue: dataInput.Email,
		Scope:          "PASSWORD_RECOVERY",
	})
	if err != nil {
		utils.MySugarLogger.Error(err)
		return err
	}
	newPasswordHash, err := utils.GenerateHashPassword(dataInput.NewPassword)
	if err != nil {
		utils.MySugarLogger.Error(err)
		return err
	}
	_, err = data.UpdatePasswordByEmail(&data.UpdateAccountFilter{Email: dataInput.Email}, newPasswordHash)

	if err != nil {
		utils.MySugarLogger.Error(err)
		return errors.New("unable to complete password recovery")
	}

	email.SendPasswordRecoveryCompleteEmail(&email.SendPasswordRecoveryCompleteEmailData{
		Email:   dataInput.Email,
		Subject: "Password Recovery Success",
	})
	return nil
}
