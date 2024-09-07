package query

import (
	"fmt"
	"os"
	"sort"
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
	//docs, err := q.GetParticipantsRecordsAggregate()
	//assert.NoError(t, err)

	//assert.IsType(t, []exports.ParticipantAccountWithCoParticipantsDocument{}, docs)
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
	opts := &seeders.CreateParticipantAccountOpts{
		Status: &status,
	}
	//var accsInDB []exports.AccountDocument
	var partsInDB []exports.ParticipantDocument
	_, partsInDB = seeders.SeedMultipleAccountsAndParticipants(dbInstance, seeders.SeedMultipleAccountsAndParticipantsOpts{
		CreateParticipantAccountOpts: *opts,
		NumberOfAccounts:             4,
	})
	t.Run("It should get participants by solution title or description", func(t *testing.T) {
		search_sol_like := "solution"
		recs, err := q.GetParticipantsWithAccountsAggregate(exports.GetParticipantsWithAccountsAggregateFilterOpts{
			Solution_Like: &search_sol_like,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(recs) == 0 {
			t.Fatal("failed to get participants with search id")
		}
	})

	t.Run("should get by participant id", func(t *testing.T) {
		participant_id := partsInDB[0].ParticipantId
		recs, err := q.GetParticipantsWithAccountsAggregate(exports.GetParticipantsWithAccountsAggregateFilterOpts{
			ParticipantId: &participant_id,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(recs) != 1 {
			t.Fatal("failed to get participants with search id")
		}
	})

	t.Run("It should get participants by top review ranking", func(t *testing.T) {
		var top_3 []int
		var partsInDBCopy []exports.ParticipantDocument = make([]exports.ParticipantDocument, len(partsInDB))
		num := copy(partsInDBCopy, partsInDB)
		sort.Slice(partsInDBCopy, func(i, j int) bool {
			return partsInDBCopy[i].ReviewRanking > partsInDBCopy[j].ReviewRanking
		})
		for i := 0; i < 3; i++ {
			top_3 = append(top_3, partsInDBCopy[i].ReviewRanking)
		}
		review_ranking_top := 3
		recs, err := q.GetParticipantsWithAccountsAggregate(exports.GetParticipantsWithAccountsAggregateFilterOpts{
			ReviewRanking_Top: &review_ranking_top,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(recs) != 3 {
			t.Fatal("failed to get participants with search id")
		}
		var found_3 []bool
		for i := 0; i < len(recs); i++ {
			for _, v := range top_3 {
				if v == recs[i].ReviewRanking {
					found_3 = append(found_3, true)
				}
			}
		}

		fmt.Println("\n\n", len(partsInDB), "\n\n", num, "\n\n", found_3, "\n\n", top_3)
		for _, v := range found_3 {
			if !v {
				t.Fatal("failed to find top ")
			}
		}
	})

}

func TestUpdateParticipantRecord(t *testing.T) {

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
	//var accsInDB []exports.AccountDocument
	//var partsInDB exports.ParticipantDocument
	accInDB, err := seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance, &seeders.OptsToCreateParticipantRecord{
		Status: status,
		TeamleadInfo: seeders.TeamLeadInfoToCreateTeamParticipant{
			Email: "trinitietp@gmail.com",
		},
	})
	if err != nil {
		panic(err)
	}
	t.Run("It should update the  status.", func(t *testing.T) {
		rec, err := q.UpdateSingleParticipantRecord(&exports.UpdateSingleParticipantRecordFilter{
			HackathonId:   accInDB.HackathonId,
			ParticipantId: accInDB.ParticipantId,
		}, &exports.UpdateParticipantRecordData{
			Status: "REVIEWED",
		})
		if err != nil {
			t.Fatal(err)
		}

		if rec.Status == status {
			t.Fatal("failed to update")
		}

	})

}
