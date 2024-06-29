package data

import (
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

func GetQueue(name string) (rmq.Queue, error) {
	return rmqUtils.GetQueue(name)
}
