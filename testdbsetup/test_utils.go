package testdbsetup

import (
	"context"
	"path"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/utils"
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

func CleanupDB(dbInstance *mongo.Database) {
	err := dbInstance.Drop(context.Background())
	if err != nil {
		panic(err)
	}
}
