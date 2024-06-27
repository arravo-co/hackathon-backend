package data

import (
	"context"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeleteAccount(identifier string) (*exports.AccountDocument, error) {
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
	err = accountCol.FindOneAndDelete(context.TODO(), filter).Decode(&dataFromCol)
	return &dataFromCol, err
}

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

func GetAccountsByEmails(emails []string) ([]exports.AccountDocument, error) {
	var accounts []exports.AccountDocument
	accountCol, err := Datasource.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	var filterStructs bson.A
	for _, email := range emails {
		filterStructs = append(filterStructs, bson.M{"email": email})
	}
	cursor, err := accountCol.Find(ctx, bson.D{{"$or", filterStructs}} /*bson.D{{"$or", filterStructs}}*/)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &accounts)
	fmt.Println(accounts)
	if err != nil {
		return nil, err
	}
	return accounts, err
}

func GetAccountsOfJudges() ([]exports.AccountDocument, error) {
	var accounts []exports.AccountDocument
	accountCol, err := Datasource.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	var filterStructs bson.D = bson.D{{
		"role", "JUDGE"}}
	cursor, err := accountCol.Find(ctx, filterStructs /*bson.D{{"$or", filterStructs}}*/)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &accounts)
	fmt.Println(accounts)
	if err != nil {
		return nil, err
	}
	return accounts, err
}

func GetAccountsByParticipantIds(participantIds []string) ([]exports.AccountDocument, error) {
	var accounts []exports.AccountDocument
	accountCol, err := Datasource.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	var filterStructs bson.A
	for _, email := range participantIds {
		filterStructs = append(filterStructs, bson.M{"participant_id": email})
	}
	cursor, err := accountCol.Find(ctx, bson.D{{Key: "$or", Value: filterStructs}} /*bson.D{{"$or", filterStructs}}*/)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &accounts)
	fmt.Println(accounts)
	if err != nil {
		return nil, err
	}
	return accounts, err
}

func CreateTeamMemberAccount(dataToSave *exports.CreateTeamMemberAccountData) (*exports.AccountDocument, error) {
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
		HackathonId:     dataToSave.HackathonId,
		Role:            dataToSave.Role,
		PhoneNumber:     dataToSave.PhoneNumber,
		IsEmailVerified: true,
		Skillset:        dataToSave.Skillset,
		State:           dataToSave.State,
		ParticipantId:   dataToSave.ParticipantId,
		Status:          dataToSave.Status,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	fmt.Printf("Account\n")
	result, err := accountCol.InsertOne(ctx, acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.InsertedID)
	return &acc, err
}

func RemoveTeamMemberAccount(dataToSave *exports.RemoveTeamMemberAccountData) (*exports.AccountDocument, error) {
	accountCol, err := Datasource.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	acc := exports.AccountDocument{}
	fmt.Printf("Account\n")
	result := accountCol.FindOneAndDelete(ctx, dataToSave)
	err = result.Decode(&acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result)
	return &acc, err
}

/*
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
*/

func CreateAdminAccount(dataToSave *exports.CreateAdminAccountData) (*exports.AccountDocument, error) {
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
	acc := &exports.AccountDocument{
		Email:           dataToSave.Email,
		FirstName:       dataToSave.FirstName,
		LastName:        dataToSave.LastName,
		PhoneNumber:     dataToSave.PhoneNumber,
		Gender:          dataToSave.Gender,
		Role:            dataToSave.Role,
		IsEmailVerified: false,
		HackathonId:     dataToSave.HackathonId,
	}
	return acc, nil
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

func CreateParticipantAccount(dataToSave *exports.CreateParticipantAccountData) (*exports.AccountDocument, error) {
	accountCol, err := Datasource.GetAccountCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	acc := exports.AccountDocument{
		Email:            dataToSave.Email,
		PasswordHash:     dataToSave.PasswordHash,
		FirstName:        dataToSave.FirstName,
		LastName:         dataToSave.LastName,
		Gender:           dataToSave.Gender,
		HackathonId:      dataToSave.HackathonId,
		Role:             dataToSave.Role,
		PhoneNumber:      dataToSave.PhoneNumber,
		IsEmailVerified:  false,
		ParticipantId:    dataToSave.ParticipantId,
		Status:           dataToSave.Status,
		Skillset:         dataToSave.Skillset,
		State:            dataToSave.State,
		DOB:              dataToSave.DOB,
		Motivation:       dataToSave.Motivation,
		ExperienceLevel:  dataToSave.ExperienceLevel,
		EmploymentStatus: dataToSave.EmploymentStatus,
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
		HackathonId:     dataToSave.HackathonId,
		Role:            dataToSave.Role,
		PhoneNumber:     dataToSave.PhoneNumber,
		IsEmailVerified: false,
		Status:          dataToSave.Status,
		Bio:             dataToSave.Bio,
	}
	result, err := accountCol.InsertOne(ctx, &acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.InsertedID)
	return dataToSave, err
}

func UpdateAccountInfoByEmail(filter *exports.UpdateAccountFilter, dataInput *exports.UpdateAccountDocument) (*exports.AccountDocument, error) {
	accountCol, err := Datasource.GetAccountCollection()
	fmt.Printf("%+v", filter)
	accountDoc := exports.AccountDocument{}
	ctx := context.Context(context.Background())
	defer ctx.Done()
	if err != nil {
		return nil, err
	}
	updateDoc := bson.M{"$set": dataInput}
	after := options.After
	opts := []*options.FindOneAndUpdateOptions{{
		ReturnDocument: &after,
	}}
	result := accountCol.FindOneAndUpdate(ctx, bson.M{"email": filter.Email}, &updateDoc, opts...)
	err = result.Decode(&accountDoc)
	return &accountDoc, err
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
