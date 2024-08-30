package di

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
	"github.com/arravoco/hackathon_backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupDefaultTestEnv() {
	fp, err := utils.FindProjectRoot(".env")
	if err != nil {
		panic(err)
	}
	config.SetupEnvironment(path.Join(fp, ".env"))
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

func GetDefaultMongoInstance() *mongo.Database {
	db_url := os.Getenv("MONGODB_URL")
	clientOpts := options.Client().ApplyURI(db_url)
	mongoInstance, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(err)
	}
	dbInstance := mongoInstance.Database("hackathons_db")
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

func GetDefaultJudgeRepository() exports.JudgeRepositoryInterface {
	SetupDefaultTestEnv()
	mongoDBInstance := GetDefaultMongoInstance()
	q := GetQueryInstance(mongoDBInstance)
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
