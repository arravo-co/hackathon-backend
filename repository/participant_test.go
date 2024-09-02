package repository

import (
	"os"
	"testing"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/seeders"
	"github.com/arravoco/hackathon_backend/testdbsetup"
	//testsetup "github.com/arravoco/hackathon_backend/test_setup"
)

func TestRegisterTeamLead(t *testing.T) {
}

func TestCreateTeamMemberAccount(t *testing.T) {

}

func TestGetSingleParticipantRecordAndMemberAccountsInfo(t *testing.T) {

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
	partRepo := NewParticipantRepository(q)
	status := "UNREVIEWED"
	opts := &seeders.CreateParticpantAccountOpts{
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

	recs, err := partRepo.GetMultipleParticipantRecordAndMemberAccountsInfo(FilterGetParticipants{})
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
