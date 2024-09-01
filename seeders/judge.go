package seeders

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateFakeJudgeAccount(dbInstance *mongo.Database) (*exports.AccountDocument, string, error) {
	accountCol := dbInstance.Collection("accounts")
	ctx := context.Context(context.Background())

	fake := faker.New()
	email := fake.Internet().Email()
	person := fake.Person()
	var gender string
	rand.Shuffle(2, func(i, j int) {
		genderList := []string{"MALE", "FEMALE"}
		gender = genderList[i]
	})
	phone_number := fake.Phone().E164Number()
	password := fake.Internet().Password()
	password_hash, _ := exports.GenerateHashPassword(password)
	acc := &exports.AccountDocument{
		Email:           email,
		PasswordHash:    password_hash,
		FirstName:       person.FirstName(),
		LastName:        person.LastName(),
		Gender:          gender,
		HackathonId:     "HACKATHON_ID_001",
		Role:            "JUDGE",
		PhoneNumber:     phone_number,
		IsEmailVerified: false,
		Status:          "EMAIL_UNVERIFIED",
		Bio:             "Short bio",
	}
	result, err := accountCol.InsertOne(ctx, acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, "", err
	}
	fmt.Printf("%#v", result.InsertedID)
	acc.Id = result.InsertedID.(primitive.ObjectID)
	//fmt.Printf("%#v", acc)
	return acc, password, err
}
