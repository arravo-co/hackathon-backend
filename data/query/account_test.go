package query

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestUpdateAccountInfoByEmail(t *testing.T) {

	SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := GetMongoInstance(cfg)
	defer t.Cleanup(func() {
		CleanupDB(dbInstance)
	})
	q := GetQueryInstance(dbInstance)
	accInDB, err := CreateFakeJudgeAccount(dbInstance)
	if err != nil {
		panic(err)
	}
	t.Run("It should update first_name", func(t *testing.T) {
		first_name_updated := "David"
		accDoc, err := q.UpdateAccountInfoByEmail(&exports.UpdateAccountFilter{
			Email: accInDB.Email,
		}, &exports.UpdateAccountDocument{
			FirstName: first_name_updated,
		})
		if err != nil {
			t.Fatal(err)
		}
		accQuery := &exports.AccountDocument{}
		result := dbInstance.Collection("accounts").FindOne(context.Background(), bson.M{"email": accDoc.Email})
		err = result.Decode(accQuery)
		if err != nil {
			panic(err)
		}
		if accQuery.FirstName != first_name_updated {
			t.Fatalf(fmt.Sprintf("Expected %s, got %s", first_name_updated, accQuery.FirstName))
		}
	})

	t.Run("It should update last_name", func(t *testing.T) {
		last_name_updated := "Jones"
		accDoc, err := q.UpdateAccountInfoByEmail(&exports.UpdateAccountFilter{
			Email: accInDB.Email,
		}, &exports.UpdateAccountDocument{
			LastName: last_name_updated,
		})
		if err != nil {
			t.Fatal(err)
		}
		accQuery := &exports.AccountDocument{}
		result := dbInstance.Collection("accounts").FindOne(context.Background(), bson.M{"email": accDoc.Email})
		err = result.Decode(accQuery)
		if err != nil {
			panic(err)
		}
		if accQuery.LastName != last_name_updated {
			t.Fatalf(fmt.Sprintf("Expected %s, got %s", last_name_updated, accQuery.LastName))
		}
	})

	t.Run("It should update bio", func(t *testing.T) {
		bio_updated := "Bio updated"
		accDoc, err := q.UpdateAccountInfoByEmail(&exports.UpdateAccountFilter{
			Email: accInDB.Email,
		}, &exports.UpdateAccountDocument{
			Bio: bio_updated,
		})
		if err != nil {
			t.Fatal(err)
		}
		accQuery := &exports.AccountDocument{}
		result := dbInstance.Collection("accounts").FindOne(context.Background(), bson.M{"email": accDoc.Email})
		err = result.Decode(accQuery)
		if err != nil {
			panic(err)
		}
		if accQuery.Bio != bio_updated {
			t.Fatalf(fmt.Sprintf("Expected %s, got %s", bio_updated, accQuery.Bio))
		}
	})

}

func CreateFakeJudgeAccount(dbInstance *mongo.Database) (*exports.AccountDocument, error) {
	accountCol := dbInstance.Collection("accounts")
	ctx := context.Context(context.Background())
	acc := &exports.AccountDocument{
		Email:           "test@test.com",
		PasswordHash:    "",
		FirstName:       "john",
		LastName:        "doe",
		Gender:          "MALE",
		HackathonId:     "HACKATHON_ID_001",
		Role:            "JUDGE",
		PhoneNumber:     "+2347068968932",
		IsEmailVerified: false,
		Status:          "EMAIL_UNVERIFIED",
		Bio:             "Short bio",
	}
	result, err := accountCol.InsertOne(ctx, acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.InsertedID)
	acc.Id = result.InsertedID
	return acc, err
}

func SetupDefaultTestEnv() {
	fp, err := utils.FindProjectRoot(".test.env")
	if err != nil {
		panic(err)
	}
	config.SetupEnvironment(path.Join(fp, ".test.env"))
}

func GetMongoInstance(cfg *exports.MongoDBConnConfig) *mongo.Database {
	clientOpts := options.Client().ApplyURI(cfg.Url)
	mongoInstance, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(err)
	}
	dbInstance := mongoInstance.Database(cfg.DBName)
	return dbInstance
}

func GetQueryInstance(dbInstance *mongo.Database) *Query {
	dat := db.GetMongoRepositoryWithDB(dbInstance)
	var dataSourceInstance exports.DBInterface = dat
	q := GetQueryWithConfiguredDatasource(dataSourceInstance)
	return q
}

func CleanupDB(dbInstance *mongo.Database) {
	err := dbInstance.Drop(context.Background())
	if err != nil {
		panic(err)
	}
}
