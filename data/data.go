package data

import (
	"github.com/adjust/rmq/v5"
	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/rmqUtils"
)

var Datasource exports.DBInterface

var hackathonId string

func init() {
	Datasource = db.Mongo{}

}

func GetDatasource() exports.DBInterface {
	return Datasource
}

func GetQueue(name string) (rmq.Queue, error) {
	return rmqUtils.GetQueue(name)
}
