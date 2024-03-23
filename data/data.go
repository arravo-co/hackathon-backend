package data

import (
	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
)

var Datasource exports.DBInterface

var hackathonId string

func init() {
	Datasource = db.Mongo{}
	hackathonId = "hackathon001"
}

func GetDatasource() exports.DBInterface {
	return Datasource
}
