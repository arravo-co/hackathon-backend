package di

import (
	"context"
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
