package repository

import (
	"os"
	"testing"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/seeders"
	"github.com/arravoco/hackathon_backend/testdbsetup"
	"go.mongodb.org/mongo-driver/mongo"
	//testsetup "github.com/arravoco/hackathon_backend/test_setup"
)

func TestGetSingleParticipantRecordAndMemberAccountsInfo(t *testing.T) {

	dbInstance, accInDB, partInDB, partRepo := SetupDBForRepository()
	defer t.Cleanup(func() {
		testdbsetup.CleanupDB(dbInstance)
	})
	//panic(accInDB.ParticipantId)
	recs, err := partRepo.GetSingleParticipantRecordAndMemberAccountsInfo(accInDB.ParticipantId)
	if err != nil {
		t.Fatal(err)
	}

	if recs.HackathonId != partInDB.HackathonId {
		t.Fatalf("hackathon id does not match. expected %v, got %v", recs.HackathonId, partInDB.HackathonId)
	}
	if recs.TeamName != partInDB.TeamName {
		t.Fatalf("hackathon id does not match. expected %v, got %v", recs.HackathonId, partInDB.HackathonId)
	}

}

func TestGetMultipleParticipantRecordAndMemberAccountsInfo(t *testing.T) {

	testdbsetup.SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := testdbsetup.GetMongoInstance(cfg)
	defer t.Cleanup(func() {
		testdbsetup.CleanupDB(dbInstance)
	})
	q := testdbsetup.GetQueryInstance(dbInstance)
	partRepo := NewParticipantRecordRepository(q)
	status := "UNREVIEWED"
	opts := &seeders.CreateParticipantAccountOpts{
		Status: &status,
	}
	accInDB, _, err := seeders.CreateFakeParticipantAccount(dbInstance, opts)
	if err != nil {
		panic(err)
	}
	teamLeadInfo := seeders.TeamLeadInfoToCreateTeamParticipant{
		TeamName:      "Good team",
		Email:         accInDB.Email,
		ParticipantId: accInDB.ParticipantId,
		HackathonId:   accInDB.HackathonId,
	}
	partInDB, err := seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance, &seeders.OptsToCreateParticipantRecord{
		TeamleadInfo: teamLeadInfo,
	})
	if err != nil {
		panic(err)
	}

	recs, err := partRepo.GetMultipleParticipantRecordAndMemberAccountsInfo(exports.GetParticipantsWithAccountsAggregateFilterOpts{})
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range recs {
		if item.HackathonId != partInDB.HackathonId {
			t.Fatalf("hackathon id does not match. expected %v, got %v", item.HackathonId, partInDB.HackathonId)
		}
		if item.TeamName != partInDB.TeamName {
			t.Fatalf("hackathon id does not match. expected %v, got %v", item.HackathonId, partInDB.HackathonId)
		}
	}
}

func SetupDBForRepository() (*mongo.Database, *exports.AccountDocument, *exports.ParticipantDocument, *ParticipantRecordRepository) {

	testdbsetup.SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := testdbsetup.GetMongoInstance(cfg)
	q := testdbsetup.GetQueryInstance(dbInstance)
	partRepo := NewParticipantRecordRepository(q)
	status := "UNREVIEWED"
	opts := &seeders.CreateParticipantAccountOpts{
		Status: &status,
	}
	accInDB, partsInDB := seeders.SeedMultipleAccountsAndParticipants(dbInstance, seeders.SeedMultipleAccountsAndParticipantsOpts{
		CreateParticipantAccountOpts: *opts,
	})
	return dbInstance, &accInDB[0], &partsInDB[0], partRepo
}
