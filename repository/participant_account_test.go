package repository

import (
	"os"
	"testing"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/testdbsetup"
	"go.mongodb.org/mongo-driver/mongo"
	//testsetup "github.com/arravoco/hackathon_backend/test_setup"
)

func TestRegisterTeamLead(t *testing.T) {
}

func TestCreateTeamMemberAccount(t *testing.T) {

}

func SetupParticipantAccount() (*mongo.Database, exports.ParticipantAccountRepositoryInterface) {
	testdbsetup.SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := testdbsetup.GetMongoInstance(cfg)
	q := testdbsetup.GetQueryInstance(dbInstance)
	partAccRepo := NewParticipantAccountRepository(q)

	return dbInstance, partAccRepo
}

/**
AdminUpdateParticipantRecord(filterOpts *exports.UpdateSingleParticipantRecordFilter, dataInput *exports.AdminParticipantInfoUpdateDTO)
*/
