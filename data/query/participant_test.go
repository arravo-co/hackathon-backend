package query

import (
	"os"
	"testing"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/seeders"
	"github.com/stretchr/testify/assert"
)

func TestGetParticipantsInfo(t *testing.T) {
	q := SetupDB()
	_, err := q.CreateParticipantAccount(&exports.CreateParticipantAccountData{
		CreateAccountData: exports.CreateAccountData{
			Email:       "test@test.com",
			FirstName:   "John",
			LastName:    "Doe",
			Gender:      "MALE",
			Role:        "PARTICIPANT",
			HackathonId: "hackathon001",
		},
		Skillset:        []string{"nodejs"},
		ExperienceLevel: "SENIOR",
		Motivation:      "Locum ipsum",
		ParticipantId:   "participant001",
	})
	assert.NoError(t, err)
	_, err = q.CreateParticipantRecord(&exports.CreateParticipantRecordData{
		ParticipantId:    "participant001",
		TeamLeadEmail:    "test@test.com",
		TeamName:         "Team",
		HackathonId:      "hackathon001",
		ParticipantEmail: "test@test.com",
		Type:             "TEAM",
		CoParticipants: []exports.CoParticipant{
			{Email: "test2@test.com", Role: "TEAM_MEMBER"},
		},
	})
	assert.NoError(t, err)
	docs, err := q.GetParticipantsRecordsAggregate()
	assert.NoError(t, err)

	assert.IsType(t, []exports.ParticipantAccountWithCoParticipantsDocument{}, docs)
}

func TestGetParticipantsWithAccountsAggregate(t *testing.T) {

	SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := GetMongoInstance(cfg)
	defer t.Cleanup(func() {
		CleanupDB(dbInstance)
	})
	q := GetQueryInstance(dbInstance)
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
	partInDB, err := seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance, nil, teamLeadInfo, nil, nil)
	if err != nil {
		panic(err)
	}

	recs, err := q.GetParticipantsWithAccountsAggregate(exports.GetParticipantsWithAccountsAggregateFilterOpts{})
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range recs {
		if item.HackathonId != partInDB.HackathonId {
			t.Fatalf("hackathon id does not match. expected %v, got %v", item.HackathonId, partInDB.HackathonId)
		}
	}
}
