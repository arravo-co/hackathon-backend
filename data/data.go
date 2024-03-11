package data

import (
	"context"
	"fmt"

	"github.com/arravoco/hackathon_backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

type AccountDocument struct {
	Email        string `bson:"email"`
	PasswordHash string `bson:"password_hash"`
	FirstName    string `bson:"first_name"`
	LastName     string `bson:"last_name"`
}

type CreateParticipantAccountData struct {
	Email        string `bson:"email"`
	PasswordHash string `bson:"password_hash"`
	FirstName    string `bson:"first_name"`
	LastName     string `bson:"last_name"`
}

func GetParticipantByEmail(email string) (*AccountDocument, error) {
	accountCol, err := db.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filterStruct := bson.D{
		{Key: "email", Value: email},
	}
	result := accountCol.FindOne(ctx, filterStruct)
	accountDoc := AccountDocument{}
	err = result.Decode(&accountDoc)
	return &accountDoc, err
}

func CreateParticipantAccount(dataToSave *CreateParticipantAccountData) (*CreateParticipantAccountData, error) {
	accountCol, err := db.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	result, err := accountCol.InsertOne(ctx, dataToSave)
	fmt.Println(result.InsertedID)
	return dataToSave, err
}
