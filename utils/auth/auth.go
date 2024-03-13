package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jaevor/go-nanoid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
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

type jwtCustomClaims struct {
	Email     string `json:"email"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

type Payload struct {
	Email     string
	LastName  string
	FirstName string
	Role      string
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

func InitiateEmailVerification(email string) (interface{}, error) {
	_, err := data.GetParticipantByEmail(email)
	if err != nil {
		return nil, errors.New("email not found in record")
	}
	tokenFunc, _ := nanoid.Custom("1234567890", 6)
	token := tokenFunc()
	ttl := time.Now().Add(time.Minute * 15)
	tokenData, err := data.CreateToken(&data.CreateTokenData{
		Token:          token,
		TokenType:      "EMAIL",
		TokenTypeValue: email,
		TTL:            ttl,
		Status:         "PENDING",
	})
	return tokenData, err
}
