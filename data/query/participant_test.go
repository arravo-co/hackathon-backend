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
	var accsInDB []exports.AccountDocument
	var partsInDB []exports.ParticipantDocument

	solArray := []seeders.OptsToCreateSolutionRecord{
		{Title: "Solution 1", Objective: "Objective 1", Description: "Description 1"},
		{Title: "Solution 2", Objective: "Objective 2", Description: "Description 2"},
		{Title: "Solution 3", Objective: "Objective 3", Description: "Description 3"},
		{Title: "Solution 4", Objective: "Objective 4", Description: "Description 4"},
		{Title: "Solution 5", Objective: "Objective 5", Description: "Description 5"},
		{Title: "Solution 6", Objective: "Objective 6", Description: "Description 6"},
		{Title: "Solution 7", Objective: "Objective 7", Description: "Description 7"},
		{Title: "Solution 8", Objective: "Objective 8", Description: "Description 8"},
		{Title: "Solution 9", Objective: "Objective 9", Description: "Description 9"},
		{Title: "Solution 10", Objective: "Objective 10", Description: "Description 10"},
	}

	for i := 0; i < 10; i++ {
		accInDB, _, err := seeders.CreateFakeParticipantAccount(dbInstance, opts)
		if err != nil {
			panic(err)
		}
		sol, err := seeders.CreateFakeSolutionDocument(dbInstance, solArray[i])
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
			SolutionId:   sol.Id.Hex(),
		})
		if err != nil {
			panic(err)
		}
		accsInDB = append(accsInDB, *accInDB)
		partsInDB = append(partsInDB, *partInDB)
	}

	search_sol_like := "Solution"
	recs, err := q.GetParticipantsWithAccountsAggregate(exports.GetParticipantsWithAccountsAggregateFilterOpts{
		Solution_Like: &search_sol_like,
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(recs) == 0 {
		t.Fatal("failed to get participants with search id")
	}
}
