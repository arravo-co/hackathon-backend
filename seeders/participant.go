package seeders

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/jaevor/go-nanoid"
	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateParticipantAccountOpts struct {
	Status        *string
	HackathonId   *string
	Email         *string
	ParticipantId *string
}

func CreateFakeParticipantAccount(dbInstance *mongo.Database, opts *CreateParticipantAccountOpts) (*exports.AccountDocument, string, error) {
	accountCol := dbInstance.Collection("accounts")
	ctx := context.Context(context.Background())

	fake := faker.New()
	email := fake.Internet().Email()
	person := fake.Person()
	var hackathon_id string
	var account_status string = fake.RandomStringElement([]string{"EMAIL_UNVERIFIED", "EMAIL_VERIFIED"})
	if opts != nil {
		if opts.Email != nil {
			email = *opts.Email
		}
		if opts.Status != nil {
			account_status = *opts.Status
		}
		if opts.HackathonId != nil {
			hackathon_id = *opts.HackathonId
		}
	}
	var gender string = fake.RandomStringElement([]string{"MALE", "FEMALE"})
	var employment_status string = fake.RandomStringElement([]string{"STUDENT",
		"EMPLOYED",
		"UNEMPLOYED",
		"FREELANCER"})

	var state string = fake.RandomStringElement([]string{"STUDENT",
		"LAGOS",
		"ABUJA",
		"IBADAN"})
	var skillset []string = fake.Lorem().Sentences(fake.RandomDigit())
	var course string = fake.Lorem().Word()
	var experience_level string = fake.RandomStringElement([]string{"JUNIOR",
		"MID",
		"SENIOR"})
	var IsEmailVerified bool = fake.BoolWithChance(50)
	var participant_id string
	if opts.ParticipantId != nil {
		participant_id = *opts.ParticipantId
	} else {
		gen, _ := nanoid.CustomASCII("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 11)
		participant_id = strings.Join([]string{"PARTICIPANT_ID", gen()}, "_") // "PART_ID"
	}
	phone_number := fake.Phone().E164Number()
	password := fake.Internet().Password()
	password_hash, _ := exports.GenerateHashPassword(password)
	years_of_experience := rand.Intn(7)
	dob := fake.Time().TimeBetween(time.Date(1980, time.December, 31, 0, 0, 0, 0, time.Local), time.Date(2004, time.December, 31, 0, 0, 0, 0, time.Local))
	acc := &exports.AccountDocument{
		Email:               email,
		PasswordHash:        password_hash,
		FirstName:           person.FirstName(),
		LastName:            person.LastName(),
		Gender:              gender,
		HackathonId:         hackathon_id,
		Role:                "PARTICIPANT",
		PhoneNumber:         phone_number,
		IsEmailVerified:     IsEmailVerified,
		IsEmailVerifiedAt:   fake.Time().TimeBetween(time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local), time.Now()),
		Status:              account_status,
		ParticipantId:       participant_id,
		EmploymentStatus:    employment_status,
		ExperienceLevel:     experience_level,
		HackathonExperience: fake.Lorem().Sentence(300),
		Motivation:          fake.Lorem().Sentence(50),
		PreviousProjects:    fake.Lorem().Sentences(40),
		DOB:                 dob,
		YearsOfExperience:   years_of_experience,
		State:               state,
		FieldOfStudy:        course,
		LinkedInAddress:     fake.Internet().URL(),
		Skillset:            skillset,
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
