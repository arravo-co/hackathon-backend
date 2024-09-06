package repository

import (
	"context"
	"os"
	"testing"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/seeders"
	"github.com/arravoco/hackathon_backend/testdbsetup"
	"github.com/go-faker/faker/v4"
	jaswdrFake "github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateAccountMethod(t *testing.T) {

	dbInstance, _ := SetupDBForAccountRepositoryTesting()
	q := testdbsetup.GetQueryInstance(dbInstance)
	accRepo := NewAccountRepository(q)

	defer t.Cleanup(func() {
		testdbsetup.CleanupDB(dbInstance)
	})

	t.Run("It should create an account", func(t *testing.T) {
		email_to_create := faker.Email()
		repo, err := accRepo.CreateAccount(&exports.CreateAccountData{
			Email:    email_to_create,
			LastName: jaswdrFake.New().Person().FirstName(),
		})
		if err != nil {
			t.Fatal(err.Error())
		}
		if repo.Id == "" {
			t.Fatal("no account id returned")
		}
		var acc exports.AccountDocument
		result := dbInstance.Collection("accounts").FindOne(context.Background(), bson.M{"email": email_to_create})
		if err := result.Decode(&acc); err != nil {
			t.Fatalf("Failed to verify account creation: %v", err)
		}

	})
}

func TestGetAccountByEmailMethod(t *testing.T) {

	dbInstance, accs := SetupDBForAccountRepositoryTesting()
	q := testdbsetup.GetQueryInstance(dbInstance)
	accRepo := NewAccountRepository(q)

	defer t.Cleanup(func() {
		testdbsetup.CleanupDB(dbInstance)
	})

	t.Run("It should get an existing account by email", func(t *testing.T) {
		acc := accs[0]
		repo, err := accRepo.GetAccountByEmail(acc.Email)
		if err != nil {
			t.Fatal(err.Error())
		}
		if repo.Id == "" {
			t.Fatal("no account id returned")
		}

		if acc.Email != repo.Email {
			t.Fatalf("expected email address %s, got %s", acc.Email, repo.Email)
		}
		var accFound exports.AccountDocument
		result := dbInstance.Collection("accounts").FindOne(context.Background(), bson.M{"email": acc.Email})
		if err := result.Decode(&accFound); err != nil {
			t.Fatalf("Failed to verify account existence: %v", err)
		}

	})
}

func SetupDBForAccountRepositoryTesting() (*mongo.Database, []seeders.AccountDocumentFromSeed) {

	testdbsetup.SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := testdbsetup.GetMongoInstance(cfg)
	var accsInDB []seeders.AccountDocumentFromSeed
	/*status := "UNREVIEWED"*/
	accsInDB = seeders.SeedMultipleAccounts(dbInstance, seeders.SeedMultipleAccountsOpts{
		NumberOfAccounts: 2,
	})
	return dbInstance, accsInDB
}
