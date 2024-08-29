package data

import (
	"fmt"

	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/rmqUtils"
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

func GetDatasource(cfg *exports.MongoDBConnConfig) exports.DBInterface {
	dat, err := db.GetNewMongoRepository(cfg)
	if err != nil {
		panic(err.Error())
	}
	DefaultDatasource = dat
	return DefaultDatasource
}

func GetQueue(name string) (rmq.Queue, error) {
	return rmqUtils.GetQueue(name)
}
