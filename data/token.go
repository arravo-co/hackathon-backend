package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpsertToken(dataInput *exports.UpsertTokenData) (*exports.TokenData, error) {
	fmt.Printf("\nuuuuuuuuuuuuuuuuuuuuuuuuuu\n%+v\npppppppppppppppppppppppppp\n", dataInput)
	tokenCol, err := Datasource.GetTokenCollection()
	if err != nil {
		return nil, err
	}
	tokenInfo := &exports.TokenData{}
	filter := bson.M{
		"token_type":       dataInput.TokenType,
		"token_type_value": dataInput.TokenTypeValue,
	}
	dataInput.UpdatedAt = time.Now()
	var upsert bool = true
	returnDoc := options.After
	updateDoc := bson.M{"$set": dataInput, "$setOnInsert": bson.M{"created_at": time.Now()}}
	err = tokenCol.FindOneAndUpdate(context.TODO(), filter, updateDoc, &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnDoc,
	}).Decode(tokenInfo)
	return tokenInfo, err
}

func VerifyToken(dataInput *exports.VerifyTokenData) error {
	fmt.Printf("\n%+v\n", dataInput)
	var tokenInfo exports.TokenData
	tokenCol, err := Datasource.GetTokenCollection()
	if err != nil {
		return err
	}
	result := tokenCol.FindOne(context.TODO(), dataInput)
	err = result.Decode(&tokenInfo)
	if err != nil {
		exports.MySugarLogger.Error(err)
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
