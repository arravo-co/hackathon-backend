package data

import (
	"context"
	"fmt"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/events"
	eventsdtos "github.com/arravoco/hackathon_backend/events_dtos"
	"go.mongodb.org/mongo-driver/bson"
)

type AccountDocument struct {
	Email           string `bson:"email"`
	PasswordHash    string `bson:"password_hash"`
	FirstName       string `bson:"first_name"`
	LastName        string `bson:"last_name"`
	Gender          string `bson:"gender"`
	LinkedInAddress string `bson:"linkedIn_address"`
	GithubAddress   string `bson:"github_address"`
	State           string `bson:"state"`
}

type CreateParticipantAccountData struct {
	Email           string `bson:"email"`
	PasswordHash    string `bson:"password_hash"`
	FirstName       string `bson:"first_name"`
	LastName        string `bson:"last_name"`
	Gender          string `bson:"gender"`
	LinkedInAddress string `bson:"linkedIn_address"`
	GithubAddress   string `bson:"github_address"`
	State           string `bson:"state"`
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
	fmt.Printf("Account\n")
	fmt.Printf("%#v\n", dataToSave)
	result, err := accountCol.InsertOne(ctx, dataToSave)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.InsertedID)
	events.EmitParticipantAccountCreated(&eventsdtos.ParticipantAccountCreatedEventData{
		ParticipantEmail: dataToSave.Email,
		LastName:         dataToSave.LastName,
		FirstName:        dataToSave.FirstName,
		EventData:        eventsdtos.EventData{EventName: "ParticipantAccountCreated"},
	})
	return dataToSave, err
}
