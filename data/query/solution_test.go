package query

import (
	"testing"

	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateSolutionData(t *testing.T) {
	q := SetupDB()
	dat := exports.CreateSolutionData{
		Title:       "Legal solution",
		Description: "Legal solution description",
	}

	doc, err := q.CreateSolutionData(&dat)
	assert.NoError(t, err)
	assert.IsType(t, &exports.SolutionDocument{}, doc)
}

func TestUpdateSolutionData(t *testing.T) {
	q := SetupDB()
	cre := exports.CreateSolutionData{
		Title:       "Test Solution",
		Description: "Test Solution Description",
	}
	doc, err := q.CreateSolutionData(&cre)
	assert.NoError(t, err)

	objId := doc.Id.(primitive.ObjectID).Hex()
	dat := exports.UpdateSolutionData{
		Title:       "Legal solution",
		Description: "Legal solution description",
	}
	docUpd, err := q.UpdateSolutionData(objId, &dat)
	assert.NoError(t, err)
	assert.IsType(t, &exports.SolutionDocument{}, docUpd)
	assert.Equal(t, doc.Id, docUpd.Id)
	assert.Equal(t, "Legal solution", docUpd.Title)
}

func TestGetSolutionDataById(t *testing.T) {
	q := SetupDB()
	cre := exports.CreateSolutionData{
		Title:       "Test Solution",
		Description: "Test Solution Description",
	}
	doc, err := q.CreateSolutionData(&cre)
	assert.NoError(t, err)

	objId := doc.Id.(primitive.ObjectID).Hex()
	docFetched, err := q.GetSolutionDataById(objId)
	assert.NoError(t, err)
	assert.IsType(t, &exports.SolutionDocument{}, docFetched)
	assert.Equal(t, doc.Id, docFetched.Id)
	assert.Equal(t, doc.Title, docFetched.Title)
}

func TestGetManySolutionData(t *testing.T) {
	q := SetupDB()
	cre := exports.CreateSolutionData{
		Title:       "Test Solution",
		Description: "Test Solution Description",
	}
	doc, err := q.CreateSolutionData(&cre)
	assert.NoError(t, err)

	objId := doc.Id.(primitive.ObjectID).Hex()
	docFetched, err := q.GetManySolutionData(objId)
	assert.NoError(t, err)
	assert.IsType(t, []exports.SolutionDocument{}, docFetched)
}

func SetupDB() *Query {
	repo, err := db.GetNewMongoRepository(&exports.MongoDBConnConfig{
		Url: "mongodb://localhost:27017",
	})
	if err != nil {
		panic(err)
	}
	return &Query{
		Datasource: repo,
	}
}
