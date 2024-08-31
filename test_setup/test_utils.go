package testsetup

import (
	"context"
	"fmt"
	"math/rand"
	"path"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
	"github.com/arravoco/hackathon_backend/utils"
	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func GetQueryInstance(dbInstance *mongo.Database) *query.Query {
	var dataSourceInstance exports.DBInterface = data.GetDatasourceWithMongoDBInstance(dbInstance)
	q := query.GetQueryWithConfiguredDatasource(dataSourceInstance)
	return q
}

func GetJudgeAccountRepositoryWithQueryInstance(q *query.Query) exports.JudgeRepositoryInterface {
	var judgeRepoInstance *repository.JudgeAccountRepository = repository.NewJudgeAccountRepository(q)
	var judgeAccountRepository exports.JudgeRepositoryInterface = judgeRepoInstance
	return judgeAccountRepository
}

/*
func GetDefaultDBInterface() exports.DBInterface {
	SetupDefaultTestEnv()
	var mongoInstance exports.DBInterface
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	mongoInstance = data.GetDatasourceWithConfig(cfg)
	return mongoInstance
}
func GetDefaultAdminRepository() exports.JudgeRepositoryInterface {
	fp, err := utils.FindProjectRoot(".test.env")
	if err != nil {
		panic(err)
	}
	config.SetupEnvironment(path.Join(fp, ".test.env"))
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}

	var judgeAccountRepository exports.JudgeRepositoryInterface = GetJudgeRepositoryWithConfig(cfg)
	return judgeAccountRepository
}

func GetDefaultJudgeRepository() exports.JudgeRepositoryInterface {
	fp, err := utils.FindProjectRoot(".test.env")
	if err != nil {
		panic(err)
	}
	config.SetupEnvironment(path.Join(fp, ".test.env"))
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}

	var judgeAccountRepository exports.JudgeRepositoryInterface = GetJudgeRepositoryWithConfig(cfg)
	return judgeAccountRepository
}

func GetJudgeRepositoryWithConfig(cfg *exports.MongoDBConnConfig) exports.JudgeRepositoryInterface {
	clientOpts := options.Client().ApplyURI(cfg.Url)
	mongoInstance, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(err)
	}
	dbInstance := mongoInstance.Database(cfg.DBName)
	var dataSourceInstance exports.DBInterface = data.GetDatasourceWithMongoDBInstance(dbInstance)
	q := query.GetQueryWithConfiguredDatasource(dataSourceInstance)
	var judgeRepoInstance *repository.JudgeAccountRepository = repository.NewJudgeAccountRepository(q)
	var judgeAccountRepository exports.JudgeRepositoryInterface = judgeRepoInstance
	return judgeAccountRepository
}*/

func CleanupDB(dbInstance *mongo.Database) {
	err := dbInstance.Drop(context.Background())
	if err != nil {
		panic(err)
	}
}

func CreateFakeJudgeAccount(dbInstance *mongo.Database) (*exports.AccountDocument, string, error) {
	accountCol := dbInstance.Collection("accounts")
	ctx := context.Context(context.Background())

	fake := faker.New()

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
		Email:           "test@test.com",
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
	fmt.Printf("%#v", acc)
	return acc, password, err
}
