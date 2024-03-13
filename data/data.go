package data

import (
	"context"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/events"
	eventsdtos "github.com/arravoco/hackathon_backend/events_dtos"
	"go.mongodb.org/mongo-driver/bson"
)

type AccountDocument struct {
	Id              string `bson:"_id"`
	Email           string `bson:"email"`
	PasswordHash    string `bson:"password_hash"`
	FirstName       string `bson:"first_name"`
	LastName        string `bson:"last_name"`
	Gender          string `bson:"gender"`
	LinkedInAddress string `bson:"linkedIn_address"`
	GithubAddress   string `bson:"github_address"`
	State           string `bson:"state"`
	Role            string `bson:"role"`
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

type TokenData struct {
	Id             string    `bson:"_id"`
	Token          string    `bson:"token"`
	TokenType      string    `bson:"token_type"`
	TokenTypeValue string    `bson:"token_type_value"`
	TTL            time.Time `bson:"ttl"`
	Status         string    `bson:"status"`
}

type CreateTokenData struct {
	Token          string    `bson:"token"`
	TokenType      string    `bson:"token_type"`
	TokenTypeValue string    `bson:"token_type_value"`
	TTL            time.Time `bson:"ttl"`
	Status         string    `bson:"status"`
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
		Id:             result.InsertedID.(string),
		Token:          dataInput.Token,
		TokenType:      dataInput.TokenType,
		TokenTypeValue: dataInput.TokenTypeValue,
		TTL:            dataInput.TTL,
	}
	return tokenInfo, err
}
