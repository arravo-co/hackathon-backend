package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/utils"
)

type TokenData struct {
	Id             interface{} `bson:"_id"`
	Token          string      `bson:"token"`
	TokenType      string      `bson:"token_type"`
	TokenTypeValue string      `bson:"token_type_value"`
	TTL            time.Time   `bson:"ttl"`
	Status         string      `bson:"status"`
}

type CreateTokenData struct {
	Token          string    `bson:"token"`
	TokenType      string    `bson:"token_type"`
	TokenTypeValue string    `bson:"token_type_value"`
	TTL            time.Time `bson:"ttl"`
	Status         string    `bson:"status"`
}

type VerifyTokenData struct {
	Token          string `bson:"token"`
	TokenType      string `bson:"token_type"`
	TokenTypeValue string `bson:"token_type_value"`
}

func CreateToken(dataInput *CreateTokenData) (*TokenData, error) {
	tokenCol, err := db.GetTokenCollection()
	if err != nil {
		return nil, err
	}
	result, err := tokenCol.InsertOne(context.TODO(), dataInput)
	if err != nil {
		return nil, err
	}
	tokenInfo := &TokenData{
		Id:             result.InsertedID,
		Token:          dataInput.Token,
		TokenType:      dataInput.TokenType,
		TokenTypeValue: dataInput.TokenTypeValue,
		TTL:            dataInput.TTL,
	}
	return tokenInfo, err
}

func VerifyToken(dataInput *VerifyTokenData) (bool, error) {
	fmt.Printf("\n%+v\n", dataInput)
	var tokenInfo TokenData
	tokenCol, err := db.GetTokenCollection()
	if err != nil {
		return false, err
	}
	result := tokenCol.FindOne(context.TODO(), dataInput)
	err = result.Decode(&tokenInfo)
	if err != nil {
		utils.MySugarLogger.Error(err)
		return false, errors.New("unable to verify token")
	}
	if tokenInfo.Token == "" {
		return false, errors.New("token does not exist")
	}
	if tokenInfo.TTL.Before(time.Now()) {
		return false, errors.New("token has expired")
	}
	if tokenInfo.Status == "VERIFIED" {
		return false, errors.New("token has been used for verification of this entity in the past")
	}

	return true, err
}
