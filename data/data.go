package data

import (
	"fmt"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/rmqUtils"
	"go.mongodb.org/mongo-driver/mongo"
)

var DefaultDatasource exports.DBInterface

func SetupDefaultDataSource() {
	dat, err := db.GetMongoRepository("hackathons_db")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("\n\nDefault DataSource initialized\n\n")
	DefaultDatasource = dat
}

func GetDefaultDatasource() exports.DBInterface {
	dat, err := db.GetMongoRepository("hackathons_db")
	if err != nil {
		panic(err.Error())
	}
	DefaultDatasource = dat
	return DefaultDatasource
}

// GetMongoRepositoryWithDB

func GetDatasourceWithMongoDBInstance(monDb *mongo.Database) exports.DBInterface {
	dat := db.GetMongoRepositoryWithDB(monDb)
	return dat
}

func GetDatasourceWithConfig(cfg *exports.MongoDBConnConfig) exports.DBInterface {
	dat, err := db.GetNewMongoRepository(cfg)
	if err != nil {
		panic(err.Error())
	}
	return dat
}

func GetQueue(name string) (rmq.Queue, error) {
	return rmqUtils.GetQueue(name)
}
