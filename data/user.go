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
	Id                string    `bson:"_id"`
	Email             string    `bson:"email,omitempty"`
	PasswordHash      string    `bson:"password_hash,omitempty"`
	PhoneNumber       string    `bson:"phone_number,omitempty"`
	FirstName         string    `bson:"first_name,omitempty"`
	LastName          string    `bson:"last_name,omitempty"`
	Gender            string    `bson:"gender,omitempty"`
	LinkedInAddress   string    `bson:"linkedIn_address,omitempty"`
	GithubAddress     string    `bson:"github_address,omitempty"`
	State             string    `bson:"state,omitempty"`
	Role              string    `bson:"role,omitempty"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
}

type UpdateAccountFilter struct {
	Email       string `bson:"email"`
	PhoneNumber string `bson:"phone_number"`
}

type UpdateAccountDocument struct {
	FirstName         string    `bson:"first_name,omitempty"`
	LastName          string    `bson:"last_name,omitempty"`
	Gender            string    `bson:"gender,omitempty"`
	LinkedInAddress   string    `bson:"linkedIn_address,omitempty"`
	GithubAddress     string    `bson:"github_address,omitempty"`
	State             string    `bson:"state,omitempty"`
	Role              string    `bson:"role,omitempty"`
	IsEmailVerified   bool      `bson:"is_email_verified,omitempty"`
	IsEmailVerifiedAt time.Time `bson:"is_email_verified_at,omitempty"`
}

type CreateUserAccountData struct {
	Email        string `bson:"email"`
	PasswordHash string `bson:"password_hash"`
	FirstName    string `bson:"first_name"`
	LastName     string `bson:"last_name"`
	Gender       string `bson:"gender"`
	State        string `bson:"state"`
	Role         string `bson:"role"`
}

type CreateParticipantAccountData struct {
	CreateUserAccountData
	LinkedInAddress string `bson:"linkedIn_address"`
	GithubAddress   string `bson:"github_address"`
}

type CreateJudgeAccountData struct {
	CreateUserAccountData
}

func GetAccountByEmail(email string) (*AccountDocument, error) {
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

func CreateJudgeAccount(dataToSave *CreateJudgeAccountData) (*CreateJudgeAccountData, error) {
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
func UpdateParticipantInfoByEmail(filter *UpdateAccountFilter, dataInput *UpdateAccountDocument) (*AccountDocument, error) {
	accountCol, err := db.GetAccountCollection()
	fmt.Printf("%+v", filter)
	accountDoc := AccountDocument{}
	ctx := context.Context(context.Background())
	defer ctx.Done()
	if err != nil {
		return nil, err
	}
	updateDoc := bson.M{"$set": dataInput}
	result := accountCol.FindOneAndUpdate(ctx, bson.M{"email": filter.Email}, &updateDoc)
	err = result.Decode(&accountDoc)
	return &accountDoc, err
}

func UpdatePasswordByEmail(filter *UpdateAccountFilter, newPasswordHash string) (*AccountDocument, error) {
	accountCol, err := db.GetAccountCollection()
	fmt.Printf("%+v", filter)
	accountDoc := AccountDocument{}
	ctx := context.Context(context.Background())
	defer ctx.Done()
	if err != nil {
		return nil, err
	}
	updateDoc := bson.M{"$set": bson.M{"password_hash": newPasswordHash}}
	result := accountCol.FindOneAndUpdate(ctx, bson.M{"email": filter.Email}, &updateDoc)
	err = result.Decode(&accountDoc)
	return &accountDoc, err
}
