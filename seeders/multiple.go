package seeders

import (
	"math/rand"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/go-faker/faker/v4"
	jawsFaker "github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeedMultipleAccountsOpts struct {
	NumberOfAccounts              int
	NumberOfCoParticipantAccounts int
	//OptsToCreateParticipantRecord
}

type SeedMultipleAccountsAndParticipantsOpts struct {
	NumberOfAccounts              int
	NumberOfCoParticipantAccounts int
	CreateParticipantAccountOpts
	//OptsToCreateParticipantRecord
}

func SeedMultipleAccountsAndParticipants(dbInstance *mongo.Database, opts SeedMultipleAccountsAndParticipantsOpts) ([]exports.AccountDocument, []exports.ParticipantDocument) {

	var accsInDB []exports.AccountDocument
	var partsInDB []exports.ParticipantDocument
	solArray := []OptsToCreateSolutionRecord{
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
	number_of_co_part_accs := 4
	for i := 0; i < opts.NumberOfAccounts; i++ {
		team_lead_email := faker.Email()
		accInDB, _, err := CreateFakeParticipantAccount(dbInstance, &CreateParticipantAccountOpts{
			Email: &team_lead_email,
		})
		if err != nil {
			panic(err)
		}
		sol, err := CreateFakeSolutionDocument(dbInstance, solArray[i])
		if err != nil {
			panic(err)
		}
		teamLeadInfo := TeamLeadInfoToCreateTeamParticipant{
			TeamName:      "Good team",
			Email:         accInDB.Email,
			ParticipantId: accInDB.ParticipantId,
			HackathonId:   accInDB.HackathonId,
		}
		review_ranking := rand.Intn(1000)
		var co_participants_accs []*exports.AccountDocument
		for i := 0; i < number_of_co_part_accs; i++ {
			email := faker.Email()
			status := "EMAIL_VERIFIED"
			acc, _, err := CreateFakeParticipantAccount(dbInstance, &CreateParticipantAccountOpts{
				ParticipantId: &teamLeadInfo.ParticipantId,
				Status:        &status,
				HackathonId:   &teamLeadInfo.HackathonId,
				Email:         &email,
			})
			if err != nil {
				panic(err)
			}
			co_participants_accs = append(co_participants_accs, acc)
		}
		var co_participants []CoParticipantInfoToCreateTeamParticipant
		for _, v := range co_participants_accs {
			co_participants = append(co_participants, CoParticipantInfoToCreateTeamParticipant{
				Email:         v.Email,
				HackathonId:   v.HackathonId,
				ParticipantId: v.ParticipantId,
			})
		}
		partInDB, err := CreateAccountLinkedTeamParticipantDocument(dbInstance, &OptsToCreateParticipantRecord{
			TeamleadInfo:   teamLeadInfo,
			SolutionId:     sol.Id.Hex(),
			ReviewRanking:  &review_ranking,
			CoParticipants: co_participants,
		})
		if err != nil {
			panic(err)
		}
		accsInDB = append(accsInDB, *accInDB)
		partsInDB = append(partsInDB, *partInDB)
	}

	return accsInDB, partsInDB
}

type AccountDocumentFromSeed struct {
	exports.AccountDocument
	Password string
}

func SeedMultipleAccounts(dbInstance *mongo.Database, opts SeedMultipleAccountsOpts) []AccountDocumentFromSeed {

	var accsInDB []AccountDocumentFromSeed
	for i := 0; i < opts.NumberOfAccounts; i++ {
		var accInDB *exports.AccountDocument
		var password string
		var hackathon_id string = "HACKATHON_ID_001"
		jaswdrFake := jawsFaker.New()
		email_to_create := faker.Email()
		acc_type := jaswdrFake.RandomStringElement([]string{"ADMIN", "PARTICIPANT", "JUDGE"})
		if acc_type == "ADMIN" {
			var err error
			accInDB, password, err = CreateFakeAdminAccount(dbInstance, &CreateFakeAdminAccountOpts{
				HackathonId: hackathon_id,
				Email:       email_to_create,
			})
			if err != nil {
				panic(err)
			}
		}
		if acc_type == "PARTICIPANT" {
			var err error
			accInDB, password, err = CreateFakeParticipantAccount(dbInstance, &CreateParticipantAccountOpts{
				Email:       &email_to_create,
				HackathonId: &hackathon_id,
			})
			if err != nil {
				panic(err)
			}
		} else if acc_type == "JUDGE" {
			var err error
			accInDB, password, err = CreateFakeJudgeAccount(dbInstance, CreateFakeJudgeAccountOpts{
				HackathonId: hackathon_id,
				Email:       email_to_create,
			})
			if err != nil {
				panic(err)
			}
		}
		accsInDB = append(accsInDB, AccountDocumentFromSeed{
			AccountDocument: *accInDB,
			Password:        password,
		})
	}

	return accsInDB
}
