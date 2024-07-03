package query

import (
	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
)

type Query struct {
	Datasource exports.DBInterface
}

func GetQueryWithConfiguredDatasource(config *exports.ConfigQueryWithDatasource) *Query {
	return &Query{Datasource: config.Datasource}
}

func GetDefaultQuery() *Query {
	re, err := db.GetNewDefaultMongoRepository()
	if err != nil {
		panic(err)
	}
	//fmt.Println("\n\n\n", re.DB.Client().Connect(context.Background()), "\n\n\n\n")

	return GetQueryWithConfiguredDatasource(&exports.ConfigQueryWithDatasource{
		Datasource: re,
	})
}
