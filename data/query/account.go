package query

import (
	"context"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (q *Query) DeleteAccount(identifier string) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) FindAccountIdentifier(identifier string) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) GetAccountByEmail(email string) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) GetAccountsByEmails(emails []string) ([]exports.AccountDocument, error) {
	var accounts []exports.AccountDocument
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) GetAccountsOfJudges() ([]exports.AccountDocument, error) {
	var accounts []exports.AccountDocument
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) GetAccountsByParticipantIds(participantIds []string) ([]exports.AccountDocument, error) {
	var accounts []exports.AccountDocument
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) CreateTeamMemberAccount(dataToSave *exports.CreateTeamMemberAccountData) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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
	acc.Id = result.InsertedID
	return &acc, err
}

func (q *Query) RemoveTeamMemberAccount(dataToSave *exports.RemoveTeamMemberAccountData) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) CreateAdminAccount(dataToSave *exports.CreateAdminAccountData) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) CreateAccount(dataToSave *exports.CreateAccountData) (interface{}, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) CreateParticipantAccount(dataToSave *exports.CreateParticipantAccountData) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Println("Create")
	acc := exports.AccountDocument{
		Email:               dataToSave.Email,
		PasswordHash:        dataToSave.PasswordHash,
		FirstName:           dataToSave.FirstName,
		LastName:            dataToSave.LastName,
		Gender:              dataToSave.Gender,
		HackathonId:         dataToSave.HackathonId,
		Role:                dataToSave.Role,
		PhoneNumber:         dataToSave.PhoneNumber,
		IsEmailVerified:     false,
		ParticipantId:       dataToSave.ParticipantId,
		Status:              dataToSave.Status,
		Skillset:            dataToSave.Skillset,
		State:               dataToSave.State,
		DOB:                 dataToSave.DOB,
		Motivation:          dataToSave.Motivation,
		ExperienceLevel:     dataToSave.ExperienceLevel,
		EmploymentStatus:    dataToSave.EmploymentStatus,
		YearsOfExperience:   dataToSave.YearsOfExperience,
		FieldOfStudy:        dataToSave.FieldOfStudy,
		HackathonExperience: dataToSave.HackathonExperience,
		PreviousProjects:    dataToSave.PreviousProjects,
	}
	fmt.Println("acciiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii")
	ctx := context.Context(context.Background())
	result, err := accountCol.InsertOne(ctx, acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Println(result.InsertedID)
	return &acc, err
}

func (q *Query) CreateJudgeAccount(dataToSave *exports.CreateJudgeAccountData) (*exports.CreateJudgeAccountData, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) UpdateAccountInfoByEmail(filter *exports.UpdateAccountFilter, dataInput *exports.UpdateAccountDocument) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) UpdateParticipantInfoByEmail(filter *exports.UpdateAccountFilter, dataInput *exports.UpdateAccountDocument) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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

func (q *Query) UpdatePasswordByEmail(filter *exports.UpdateAccountFilter, newPasswordHash string) (*exports.AccountDocument, error) {
	accountCol, err := q.Datasource.GetAccountCollection()
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
