package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenData struct {
	Id             interface{} `bson:"_id"`
	Token          string      `bson:"token"`
	TokenType      string      `bson:"token_type"`
	TokenTypeValue string      `bson:"token_type_value"`
	Scope          string      `bson:"scope"`
	TTL            time.Time   `bson:"ttl"`
	Status         string      `bson:"status"`
}

type CreateTokenData struct {
	Token          string    `bson:"token"`
	TokenType      string    `bson:"token_type"`
	TokenTypeValue string    `bson:"token_type_value"`
	Scope          string    `bson:"scope"`
	TTL            time.Time `bson:"ttl"`
	Status         string    `bson:"status"`
}

type VerifyTokenData struct {
	Token          string `bson:"token"`
	TokenType      string `bson:"token_type"`
	TokenTypeValue string `bson:"token_type_value"`
	Scope          string `bson:"scope"`
}

func CreateToken(dataInput *CreateTokenData) (*TokenData, error) {
	tokenCol, err := db.GetTokenCollection()
	tokenInfo := &TokenData{}
	if err != nil {
		return nil, err
	}
	filter := struct {
		Token string
	}{}
	var upsert bool = true
	updateDoc := bson.M{"$set": dataInput}
	result := tokenCol.FindOneAndUpdate(context.TODO(), filter, updateDoc, &options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	})
	err = result.Decode(tokenInfo)
	return tokenInfo, err
}

func VerifyToken(dataInput *VerifyTokenData) error {
	fmt.Printf("\n%+v\n", dataInput)
	var tokenInfo TokenData
	tokenCol, err := db.GetTokenCollection()
	if err != nil {
		return err
	}
	result := tokenCol.FindOne(context.TODO(), dataInput)
	err = result.Decode(&tokenInfo)
	if err != nil {
		utils.MySugarLogger.Error(err)
		return errors.New("unable to verify token")
	}
	if tokenInfo.Token == "" {
		return errors.New("token does not exist")
	}
	if tokenInfo.TTL.Before(time.Now()) {
		return errors.New("token has expired")
	}
	if tokenInfo.Status == "VERIFIED" {
		return errors.New("token has been used for verification of this entity in the past")
	}

	return nil
}
