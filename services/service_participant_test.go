package services

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/publishers"
	"github.com/arravoco/hackathon_backend/seeders"
	testsetup "github.com/arravoco/hackathon_backend/testdbsetup"
	testrepos "github.com/arravoco/hackathon_backend/testdbsetup/test_repos"
	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestRegisterTeamLead(t *testing.T) {

	dbInstance, service := Setup()
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	fake := faker.New()
	internet := fake.Internet()
	person := fake.Person()
	password := internet.Password()
	ent, err := service.RegisterTeamLead(&RegisterNewParticipantDTO{
		Email:            internet.Email(),
		FirstName:        person.FirstName(),
		LastName:         person.LastName(),
		Password:         password,
		ConfirmPassword:  password,
		PhoneNumber:      fake.Phone().E164Number(),
		Skillset:         fake.Lorem().Sentences(2),
		Motivation:       fake.Lorem().Sentence(50),
		State:            "Lagos",
		EmploymentStatus: "STUDENT",
		ExperienceLevel:  "SENIOR",
		Type:             "TEAM",
		Gender:           "MALE",
		DOB:              "2000-08-08",
	})
	if err != nil {
		t.Fatal(err)
	}

	if ent.ParticipantId == "" {
		t.Fatal("particopant id missing")
	}
}

func TestInviteToTeam(t *testing.T) {

	dbInstance, service := Setup()
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	t.Run("should be able to invite new team members to a reviewed participating team as the team lead of the team", func(t *testing.T) {

		account_status := "REVIEWED"
		teamLeadPartAccWithReviewedTeam, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
			Status: &account_status,
		})
		if err != nil {
			panic(err)
		}
		// REVIEWED participant document.
		participant_status := "REVIEWED"
		reviewedPartDoc, err := seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance,
			&seeders.OptsToCreateParticipantRecord{
				Status: participant_status,
				TeamleadInfo: seeders.TeamLeadInfoToCreateTeamParticipant{
					Email:         teamLeadPartAccWithReviewedTeam.Email,
					ParticipantId: teamLeadPartAccWithReviewedTeam.ParticipantId,
				},
			})
		if err != nil {
			panic(err)
		}

		invitee_email := faker.New().Internet().Email()
		res, err := service.InviteToTeam(&AddToTeamInviteListData{
			HackathonId:      teamLeadPartAccWithReviewedTeam.HackathonId,
			ParticipantId:    reviewedPartDoc.ParticipantId,
			InviterEmail:     teamLeadPartAccWithReviewedTeam.Email,
			InviterFirstName: teamLeadPartAccWithReviewedTeam.FirstName,
			Email:            invitee_email,
			Role:             "TEAM_ROLE",
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(fmt.Sprintf("%+v", res))

		partDocInDB := &exports.ParticipantDocument{}
		result := dbInstance.Collection("participants").FindOne(context.Background(),
			bson.M{"participant_id": teamLeadPartAccWithReviewedTeam.ParticipantId})
		err = result.Decode(partDocInDB)
		if err != nil {
			panic(err)
		}
		if len(partDocInDB.InviteList) == 0 {
			t.Fatal("no invite list")
		}
		found_in_list := false
		for _, v := range partDocInDB.InviteList {
			if v.Email == invitee_email {
				found_in_list = true
			}
		}
		if !found_in_list {
			t.Fatal("not found in invite list")
		}
	})

	t.Run("should not be able to invite new team members to an unreviewed participating team", func(t *testing.T) {

		account_status := "EMAIL_UNVERIFIED"
		teamLeadPartAccWithUnReviewedTeam, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
			Status: &account_status,
		})
		if err != nil {
			panic(err)
		}
		// UNREVIEWED participant document.
		participant_status := "UNREVIEWED"
		unReviewedPartDoc, err := seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance, &seeders.OptsToCreateParticipantRecord{
			Status: participant_status,
			TeamleadInfo: seeders.TeamLeadInfoToCreateTeamParticipant{
				Email:         teamLeadPartAccWithUnReviewedTeam.Email,
				ParticipantId: teamLeadPartAccWithUnReviewedTeam.ParticipantId,
			},
		})
		if err != nil {
			panic(err)
		}
		invitee_email := faker.New().Internet().Email()
		_, err = service.InviteToTeam(&AddToTeamInviteListData{
			HackathonId:      teamLeadPartAccWithUnReviewedTeam.HackathonId,
			ParticipantId:    unReviewedPartDoc.ParticipantId,
			InviterEmail:     teamLeadPartAccWithUnReviewedTeam.Email,
			InviterFirstName: teamLeadPartAccWithUnReviewedTeam.FirstName,
			Email:            invitee_email,
			Role:             "TEAM_ROLE",
		})
		if err == nil {
			t.Fatal(fmt.Errorf("failed to throw error"))
		}
	})

}

func TestCompleteNewTeamMemberRegistration(t *testing.T) {

	dbInstance, service := Setup()

	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	t.Run("Should create if already in invite list", func(t *testing.T) {
		fake := faker.New()
		status := "REVIEWED"
		admin_email := fake.Internet().FreeEmailDomain()
		teamLeadPartAcc, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
			Status: &status,
			Email:  &admin_email,
		})
		if err != nil {
			panic(err)
		}
		var invite_list []seeders.InvitelistQueuePayload
		for i := 0; i < 2; i++ {
			func() {
				fake := faker.New()
				email := fake.Internet().FreeEmail()
				invite_list = append(invite_list, seeders.InvitelistQueuePayload{
					InviterId: teamLeadPartAcc.Email,
					Email:     email,
					Time:      time.Now(),
				})
			}()
		}
		_, err = seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance, &seeders.OptsToCreateParticipantRecord{
			InviteList: invite_list,
			TeamleadInfo: seeders.TeamLeadInfoToCreateTeamParticipant{
				Email:         teamLeadPartAcc.Email,
				ParticipantId: teamLeadPartAcc.ParticipantId,
			},
		})
		if err != nil {
			panic(err)
		}
		person := fake.Person()
		password := fake.Internet().Password()
		invite_list_entry := invite_list[0]
		partEnt, err := service.CompleteNewTeamMemberRegistration(&CompleteNewTeamMemberRegistrationDTO{
			ParticipantId:    teamLeadPartAcc.ParticipantId,
			Email:            invite_list_entry.Email,
			FirstName:        person.FirstName(),
			LastName:         person.LastName(),
			Password:         password,
			ConfirmPassword:  password,
			PhoneNumber:      fake.Phone().E164Number(),
			Skillset:         fake.Lorem().Sentences(2),
			Motivation:       fake.Lorem().Sentence(50),
			State:            "Lagos",
			EmploymentStatus: "STUDENT",
			ExperienceLevel:  "SENIOR",
			Gender:           "MALE",
			DOB:              "2000-08-08",
		})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%+v\n", partEnt)
		var found_in_entity bool
		for _, v := range partEnt.CoParticipants {
			if v.Email == invite_list_entry.Email {
				found_in_entity = true
			}
		}

		if !found_in_entity {
			t.Fatal("Team member not found in entity's co participants list")
		}

		partDocInDB := &exports.ParticipantDocument{}
		result := dbInstance.Collection("participants").FindOne(context.Background(),
			bson.M{"participant_id": teamLeadPartAcc.ParticipantId})
		err = result.Decode(partDocInDB)
		if err != nil {
			panic(err)
		}

		var found_team_member bool

		for _, v := range partDocInDB.CoParticipants {
			if v.Email == invite_list_entry.Email {
				found_team_member = true
			}
		}

		if !found_team_member {
			t.Fatal("Team member not found in co participants list")
		}

	})
}

func TestGetTeamMembersInfo(t *testing.T) {

	dbInstance, service := Setup()
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	fake := faker.New()
	status := "REVIEWED"
	admin_email := fake.Internet().FreeEmailDomain()
	teamLeadPartAcc, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
		Status: &status,
		Email:  &admin_email,
	})
	if err != nil {
		panic(err)
	}
	var emails []string = []string{"test1@example.com", "test2@example.com"}
	var teamMems []seeders.CoParticipantInfoToCreateTeamParticipant
	for i := 0; i < 2; i++ {
		func() {
			email := emails[i]
			teamMemAcc, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
				Status:        &status,
				Email:         &email,
				ParticipantId: &teamLeadPartAcc.ParticipantId,
			})
			if err != nil {
				panic(err)
			}
			teamMems = append(teamMems, seeders.CoParticipantInfoToCreateTeamParticipant{
				Email:         teamMemAcc.Email,
				HackathonId:   teamMemAcc.HackathonId,
				ParticipantId: teamMemAcc.ParticipantId,
			})
		}()
	}

	_, err = seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance,
		&seeders.OptsToCreateParticipantRecord{
			Status: "REVIEWED",
			TeamleadInfo: seeders.TeamLeadInfoToCreateTeamParticipant{
				Email:         teamLeadPartAcc.Email,
				ParticipantId: teamLeadPartAcc.ParticipantId,
			},
			CoParticipants: teamMems,
		})
	if err != nil {
		panic(err)
	}

	team_mems_info, err := service.GetTeamMembersInfo(teamLeadPartAcc.ParticipantId)
	if err != nil {
		t.Fatal(err)
	}

	if len(team_mems_info) != len(teamMems) {
		t.Fatalf("expected mem size of %d, got %d", len(teamMems), len(team_mems_info))
	}
	var found_emails []bool
	for _, v := range team_mems_info {
		for _, a := range teamMems {
			if v.Email == a.Email {
				found_emails = append(found_emails, true)
				break
			}
		}
	}

	if len(found_emails) != len(teamMems) {
		t.Fatalf("expected mem size of %d, got %d", len(teamMems), len(found_emails))
	}
}

func TestGetParticipantInfo(t *testing.T) {
	dbInstance, service := Setup()

	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})

	fake := faker.New()
	status := "REVIEWED"
	admin_email := fake.Internet().FreeEmailDomain()
	teamLeadPartAcc, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
		Status: &status,
		Email:  &admin_email,
	})
	if err != nil {
		panic(err)
	}
	var emails []string = []string{"test1@example.com", "test2@example.com"}
	var teamMems []seeders.CoParticipantInfoToCreateTeamParticipant
	for i := 0; i < 2; i++ {
		func() {
			email := emails[i]
			teamMemAcc, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
				Status:        &status,
				Email:         &email,
				ParticipantId: &teamLeadPartAcc.ParticipantId,
			})
			if err != nil {
				panic(err)
			}
			teamMems = append(teamMems, seeders.CoParticipantInfoToCreateTeamParticipant{
				Email:         teamMemAcc.Email,
				HackathonId:   teamMemAcc.HackathonId,
				ParticipantId: teamMemAcc.ParticipantId,
			})
		}()
	}

	_, err = seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance,
		&seeders.OptsToCreateParticipantRecord{
			Status: "REVIEWED",
			TeamleadInfo: seeders.TeamLeadInfoToCreateTeamParticipant{
				Email:         teamLeadPartAcc.Email,
				ParticipantId: teamLeadPartAcc.ParticipantId,
			},
			CoParticipants: teamMems,
		})
	if err != nil {
		panic(err)
	}

	partEnt, err := service.GetParticipantInfo(teamLeadPartAcc.ParticipantId)
	if err != nil {
		panic(err)
	}

	if partEnt.ParticipantId != teamLeadPartAcc.ParticipantId {
		t.Fatalf("Participant id mismatch: expected %s, got %s", teamLeadPartAcc.ParticipantId, partEnt.ParticipantId)
	}
}

func TestSelectTeamSolution(t *testing.T) {
	dbInstance, service := Setup()

	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})

	fake := faker.New()
	status := "REVIEWED"
	admin_email := fake.Internet().FreeEmailDomain()
	teamLeadPartAcc, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
		Status: &status,
		Email:  &admin_email,
	})
	if err != nil {
		panic(err)
	}

	sol, err := seeders.CreateFakeSolutionDocument(dbInstance, seeders.OptsToCreateSolutionRecord{})
	if err != nil {
		panic(err)
	}

	partDoc, err := seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance,
		&seeders.OptsToCreateParticipantRecord{
			Status: "REVIEWED",
			TeamleadInfo: seeders.TeamLeadInfoToCreateTeamParticipant{
				Email:         teamLeadPartAcc.Email,
				ParticipantId: teamLeadPartAcc.ParticipantId,
			},
			//SolutionId: sol.Id.Hex(),
		})
	if err != nil {
		panic(err)
	}

	err = service.SelectTeamSolution(&SelectTeamSolutionData{
		SolutionId:    sol.Id.Hex(),
		ParticipantId: partDoc.ParticipantId,
		HackathonId:   partDoc.HackathonId,
	})
	if err != nil {
		t.Fatal(err)
	}
	partDocInDB := exports.ParticipantDocument{}
	result := dbInstance.Collection("participants").FindOne(context.Background(), bson.M{"participant_id": partDoc.ParticipantId})
	err = result.Decode(&partDocInDB)
	if err != nil {
		panic(err)
	}
	if partDocInDB.SolutionId != sol.Id.Hex() {
		t.Fatalf("failed to select solution: expected solution id %s, got %s", sol.Id, partDocInDB.Id)
	}
}

func Setup() (*mongo.Database, *Service) {
	testsetup.SetupDefaultTestEnv()
	db_url := os.Getenv("MONGODB_URL")
	cfg := &exports.MongoDBConnConfig{
		Url:    db_url,
		DBName: "hackathon_db",
	}
	dbInstance := testsetup.GetMongoInstance(cfg)
	q := testsetup.GetQueryInstance(dbInstance)
	partAccRepo := testrepos.GetParticipantAccountRepositoryWithQueryInstance(q)
	partRepo := testrepos.GetParticipantRepositoryWithQueryInstance(q)
	judgeAccountRepository := testrepos.GetJudgeAccountRepositoryWithQueryInstance(q)
	/**/
	rmq_url := os.Getenv("RABBITMQ_URL")
	fmt.Println(rmq_url)
	ch, err := publishers.GetRMQChannelWithURL(publishers.SetupRMQConfig{
		Url: rmq_url,
	})
	if err != nil {
		panic(err)
	}
	publisher := publishers.NewRMQPublisherWithChannel(ch)

	service := NewService(&ServiceConfig{
		JudgeAccountRepository:       judgeAccountRepository,
		ParticipantAccountRepository: partAccRepo,
		Publisher:                    publisher,
		ParticipantRepository:        partRepo,
	})
	return dbInstance, service
}
