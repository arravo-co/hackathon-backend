package seeders

import (
	"context"
	"fmt"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateFakeAdminAccountOpts struct {
	HackathonId string
	Email       string
}

func CreateFakeAdminAccount(dbInstance *mongo.Database, opts *CreateFakeAdminAccountOpts) (*exports.AccountDocument, string, error) {
	accountCol := dbInstance.Collection("accounts")
	ctx := context.Context(context.Background())

	fake := faker.New()
	email := fake.Internet().Email()
	if opts.Email != "" {
		email = opts.Email
	}
	person := fake.Person()
	var gender string = fake.RandomStringElement([]string{"MALE", "FEMALE"})
	phone_number := fake.Phone().E164Number()
	password := fake.Internet().Password()
	password_hash, _ := exports.GenerateHashPassword(password)
	var hackathon_id string = "HACKATHON_ID_001"
	if opts.HackathonId != "" {
		hackathon_id = opts.HackathonId
	}
	acc := &exports.AccountDocument{
		Email:           email,
		PasswordHash:    password_hash,
		FirstName:       person.FirstName(),
		LastName:        person.LastName(),
		Gender:          gender,
		HackathonId:     hackathon_id,
		Role:            "ADMIN",
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
