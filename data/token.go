package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateToken(dataInput *exports.CreateTokenData) (*exports.TokenData, error) {
	tokenCol, err := db.GetTokenCollection()
	tokenInfo := &exports.TokenData{}
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

func VerifyToken(dataInput *exports.VerifyTokenData) error {
	fmt.Printf("\n%+v\n", dataInput)
	var tokenInfo exports.TokenData
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
