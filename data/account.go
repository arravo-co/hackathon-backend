package data

import (
	"context"
	"fmt"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
)

func FindAccountIdentifier(identifier string) (*exports.AccountDocument, error) {
	accountCol, err := Datasource.GetAccountCollection()
	if err != nil {
		return nil, err
	}
	filter := bson.D{{
		"$or", bson.A{
			bson.D{{"email", identifier}},
			bson.D{{"username", identifier}},
		}},
	}
	dataFromCol := exports.AccountDocument{}
	err = accountCol.FindOne(context.TODO(), filter).Decode(&dataFromCol)
	return &dataFromCol, err
}

func GetAccountByEmail(email string) (*exports.AccountDocument, error) {
	accountCol, err := Datasource.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filterStruct := bson.D{
		{Key: "email", Value: email},
	}
	result := accountCol.FindOne(ctx, filterStruct)
	accountDoc := exports.AccountDocument{}
	err = result.Decode(&accountDoc)
	return &accountDoc, err
}

func CreateTeamParticipantAccount(dataToSave *exports.CreateTeamParticipantAccountData) (*exports.CreateTeamParticipantAccountData, error) {
	accountCol, err := Datasource.GetAccountCollection()
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
	return dataToSave, err
}

func CreateAccount(dataToSave *exports.CreateAccountData) (interface{}, error) {
	accountCol, err := Datasource.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	result, err := accountCol.InsertOne(ctx, dataToSave)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.InsertedID)
	return dataToSave, nil
}

func CreateIndividualParticipantAccount(dataToSave *exports.CreateIndividualParticipantAccountData) (*exports.AccountDocument, error) {
	accountCol, err := Datasource.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	acc := exports.AccountDocument{
		Email:           dataToSave.Email,
		PasswordHash:    dataToSave.PasswordHash,
		FirstName:       dataToSave.FirstName,
		LastName:        dataToSave.LastName,
		Gender:          dataToSave.Gender,
		HackathonId:     hackathonId,
		Role:            dataToSave.Role,
		PhoneNumber:     dataToSave.PhoneNumber,
		IsEmailVerified: false,
	}
	result, err := accountCol.InsertOne(ctx, acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Println(result.InsertedID)
	return &acc, err
}

func CreateJudgeAccount(dataToSave *exports.CreateJudgeAccountData) (*exports.CreateJudgeAccountData, error) {
	accountCol, err := Datasource.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	acc := exports.AccountDocument{
		Email:           dataToSave.Email,
		PasswordHash:    dataToSave.PasswordHash,
		FirstName:       dataToSave.FirstName,
		LastName:        dataToSave.LastName,
		Gender:          dataToSave.Gender,
		HackathonId:     hackathonId,
		Role:            dataToSave.Role,
		PhoneNumber:     dataToSave.PhoneNumber,
		IsEmailVerified: false,
	}
	result, err := accountCol.InsertOne(ctx, &acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.InsertedID)
	return dataToSave, err
}

func UpdateParticipantInfoByEmail(filter *exports.UpdateAccountFilter, dataInput *exports.UpdateAccountDocument) (*exports.AccountDocument, error) {
	accountCol, err := Datasource.GetAccountCollection()
	fmt.Printf("%+v", filter)
	accountDoc := exports.AccountDocument{}
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

func UpdatePasswordByEmail(filter *exports.UpdateAccountFilter, newPasswordHash string) (*exports.AccountDocument, error) {
	accountCol, err := Datasource.GetAccountCollection()
	fmt.Printf("%+v", filter)
	accountDoc := exports.AccountDocument{}
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
