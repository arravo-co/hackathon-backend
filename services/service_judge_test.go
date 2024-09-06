package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/seeders"
	testsetup "github.com/arravoco/hackathon_backend/testdbsetup"
	testrepos "github.com/arravoco/hackathon_backend/testdbsetup/test_repos"
	"go.mongodb.org/mongo-driver/bson"
)

func TestRegisterJudge(t *testing.T) {
	testsetup.SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := testsetup.GetMongoInstance(cfg)
	q := testsetup.GetQueryInstance(dbInstance)
	judgeAccountRepository := testrepos.GetJudgeAccountRepositoryWithQueryInstance(q)
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})

	/*
		rmq_url := os.Getenv("RABBITMQ_URL")
		ch, err := rabbitmq.GetRMQChannelWithURL(rabbitmq.SetupRMQConfig{
			Url: rmq_url,
		})
		if err != nil {
			panic(err)
		}
		publisher := publishers.NewRMQPublisherWithChannel(ch)*/

	service := NewService(&ServiceConfig{
		JudgeAccountRepository: judgeAccountRepository,
		//Publisher:              publisher,
	})
	dataInput := &RegisterNewJudgeByAdminDTO{
		FirstName:       "john",
		LastName:        "doe",
		Password:        "password",
		ConfirmPassword: "password",
		Gender:          "MALE",
		Email:           "trinitietp@gmail.com",
		Bio:             "Good bio",
		InviterEmail:    "temitope.alabi@arravo.co",
		InviterName:     "trinitietp",
	}
	judgeRepo, err := service.RegisterNewJudge(dataInput)
	if err != nil {
		t.Fatal(err)
	}
	if judgeRepo.Id == "" {
		t.Fatal(errors.New("id not returned"))
	}
	if judgeRepo.Email == "" {
		t.Fatal(errors.New("email not returned"))
	}
	if judgeRepo.FirstName == "" {
		t.Fatal(errors.New("firstname not returned"))
	}
	if judgeRepo.LastName == "" {
		t.Fatal(errors.New("lastname not returned"))
	}
	if judgeRepo.Bio == "" {
		t.Fatal(errors.New("bio not returned"))
	}

}

func TestUpdateJudgeInfo(t *testing.T) {
	testsetup.SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := testsetup.GetMongoInstance(cfg)
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	q := testsetup.GetQueryInstance(dbInstance)
	judgeAccountRepository := testrepos.GetJudgeAccountRepositoryWithQueryInstance(q)
	accInDB, _, err := seeders.CreateFakeJudgeAccount(dbInstance)
	if err != nil {
		panic(err)
	}

	service := NewService(&ServiceConfig{JudgeAccountRepository: judgeAccountRepository})
	t.Run("should update firstname", func(t *testing.T) {

		firstNameToUpdate := "Sam"
		judgeEntity, err := service.UpdateJudgeInfo(accInDB.Email, &UpdateJudgeDTO{
			FirstName: firstNameToUpdate,
		})
		if err != nil {
			t.Fatal(err)
		}

		accQuery := &exports.AccountDocument{}
		result := dbInstance.Collection("accounts").FindOne(context.Background(), bson.M{"email": accInDB.Email})
		err = result.Decode(accQuery)
		if err != nil {
			panic(err)
		}
		if accQuery.FirstName != firstNameToUpdate {
			t.Fatalf(fmt.Sprintf("Expected %s, got %s", firstNameToUpdate, accQuery.FirstName))
		}
		if judgeEntity.FirstName != firstNameToUpdate {
			t.Fatalf(fmt.Sprintf("Failed to update entity: Expected %s, got %s", firstNameToUpdate, judgeEntity.FirstName))
		}
	})

	t.Run("should update lastname", func(t *testing.T) {

		lastNameToUpdate := "Harry"
		judgeEntity, err := service.UpdateJudgeInfo(accInDB.Email, &UpdateJudgeDTO{
			LastName: lastNameToUpdate,
		})
		if err != nil {
			t.Fatal(err)
		}

		accQuery := &exports.AccountDocument{}
		result := dbInstance.Collection("accounts").FindOne(context.Background(), bson.M{"email": accInDB.Email})
		err = result.Decode(accQuery)
		if err != nil {
			panic(err)
		}
		if accQuery.LastName != lastNameToUpdate {
			t.Fatalf(fmt.Sprintf("Expected %s, got %s", lastNameToUpdate, accQuery.LastName))
		}
		if judgeEntity.LastName != lastNameToUpdate {
			t.Fatalf(fmt.Sprintf("Failed to update entity: Expected %s, got %s", lastNameToUpdate, judgeEntity.LastName))
		}
	})

	t.Run("should update bio", func(t *testing.T) {

		bioToUpdate := "updated bio"
		judgeEntity, err := service.UpdateJudgeInfo(accInDB.Email, &UpdateJudgeDTO{
			Bio: bioToUpdate,
		})
		if err != nil {
			t.Fatal(err)
		}

		accQuery := &exports.AccountDocument{}
		result := dbInstance.Collection("accounts").FindOne(context.Background(), bson.M{"email": accInDB.Email})
		err = result.Decode(accQuery)
		if err != nil {
			panic(err)
		}
		if accQuery.Bio != bioToUpdate {
			t.Fatalf(fmt.Sprintf("Expected %s, got %s", bioToUpdate, accQuery.Bio))
		}
		if judgeEntity.Bio != bioToUpdate {
			t.Fatalf(fmt.Sprintf("Failed to update entity: Expected %s, got %s", bioToUpdate, judgeEntity.Bio))
		}
	})

}

func TestGetJudgeByEmail(t *testing.T) {

	testsetup.SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := testsetup.GetMongoInstance(cfg)
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	q := testsetup.GetQueryInstance(dbInstance)
	judgeAccountRepository := testrepos.GetJudgeAccountRepositoryWithQueryInstance(q)
	accInDB, _, err := seeders.CreateFakeJudgeAccount(dbInstance)
	if err != nil {
		panic(err)
	}

	service := NewService(&ServiceConfig{JudgeAccountRepository: judgeAccountRepository})

	judgeEnt, err := service.GetJudgeByEmail(accInDB.Email)
	if err != nil {
		t.Fatal(err.Error())
	}

	if judgeEnt.Id != accInDB.Id.Hex() {
		t.Fatalf("expected id to be %s, got %s", accInDB.Id.Hex(), judgeEnt.Id)
	}

	if judgeEnt.Email != accInDB.Email {
		t.Fatalf("expected email address to be %s, got %s", accInDB.Email, judgeEnt.Email)
	}

	if judgeEnt.LastName != accInDB.LastName {
		t.Fatalf("expected last name to be %s, got %s", accInDB.LastName, judgeEnt.LastName)
	}

	if judgeEnt.FirstName != accInDB.FirstName {
		t.Fatalf("expected firstname to be %s, got %s", accInDB.FirstName, judgeEnt.FirstName)
	}

	if judgeEnt.Bio != accInDB.Bio {
		t.Fatalf("expected bio to be %s, got %s", accInDB.Bio, judgeEnt.Bio)
	}

	if judgeEnt.HackathonId != accInDB.HackathonId {
		t.Fatalf("expected hackathon id to be %s, got %s", accInDB.HackathonId, judgeEnt.HackathonId)
	}

	if judgeEnt.PhoneNumber != accInDB.PhoneNumber {
		t.Fatalf("expected phone number to be %s, got %s", accInDB.PhoneNumber, judgeEnt.PhoneNumber)
	}
}

func TestGetJudges(t *testing.T) {

	testsetup.SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := testsetup.GetMongoInstance(cfg)
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	q := testsetup.GetQueryInstance(dbInstance)
	judgeAccountRepository := testrepos.GetJudgeAccountRepositoryWithQueryInstance(q)
	_, _, err := seeders.CreateFakeJudgeAccount(dbInstance)
	if err != nil {
		panic(err)
	}

	ents, err := judgeAccountRepository.GetJudges()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(ents) == 0 {
		t.Fatalf("expected judges, got %d", len(ents))
	}
}
