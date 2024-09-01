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
)

func TestRegisterTeamLead(t *testing.T) {
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
	publisher := publishers.NewPublisherWithChannel(ch)
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	service := NewService(&ServiceConfig{
		JudgeAccountRepository:       judgeAccountRepository,
		ParticipantAccountRepository: partAccRepo,
		Publisher:                    publisher,
		ParticipantRepository:        partRepo,
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
	publisher := publishers.NewPublisherWithChannel(ch)
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	status := "UNREVIEWED"
	teamLeadPartAcc, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
		Status: &status,
	})
	if err != nil {
		panic(err)
	}
	_, err = seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance, nil, seeders.TeamLeadInfoToCreateTeamParticipant{
		Email:         teamLeadPartAcc.Email,
		ParticipantId: teamLeadPartAcc.ParticipantId,
	}, nil, nil)
	if err != nil {
		panic(err)
	}
	service := NewService(&ServiceConfig{
		JudgeAccountRepository:       judgeAccountRepository,
		ParticipantAccountRepository: partAccRepo,
		Publisher:                    publisher,
		ParticipantRepository:        partRepo,
	})
	invitee_email := faker.New().Internet().Email()
	res, err := service.InviteToTeam(&AddToTeamInviteListData{
		HackathonId:      teamLeadPartAcc.HackathonId,
		ParticipantId:    teamLeadPartAcc.ParticipantId,
		InviterEmail:     teamLeadPartAcc.Email,
		InviterFirstName: teamLeadPartAcc.FirstName,
		Email:            invitee_email,
		Role:             "TEAM_ROLE",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("%#v", res))

	partDocInDB := &exports.ParticipantDocument{}
	result := dbInstance.Collection("participants").FindOne(context.Background(),
		bson.M{"participant_id": teamLeadPartAcc.ParticipantId})
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
}

func TestCompleteNewTeamMemberRegistration(t *testing.T) {

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
	publisher := publishers.NewPublisherWithChannel(ch)
	defer t.Cleanup(func() {
		testsetup.CleanupDB(dbInstance)
	})
	status := "UNREVIEWED"
	teamLeadPartAcc, _, err := seeders.CreateFakeParticipantAccount(dbInstance, &seeders.CreateParticpantAccountOpts{
		Status: &status,
	})
	if err != nil {
		panic(err)
	}
	var invite_list []seeders.InvitelistQueuePayload
	for i := 0; i < 20; i++ {
		invite_list = append(invite_list, seeders.InvitelistQueuePayload{
			InviterId: teamLeadPartAcc.Email,
			Email:     faker.New().Internet().Email(),
			Time:      time.Now(),
		})
	}
	_, err = seeders.CreateAccountLinkedTeamParticipantDocument(dbInstance, &seeders.OptsToCreateParticipantRecord{
		InviteList: invite_list,
	}, seeders.TeamLeadInfoToCreateTeamParticipant{
		Email:         teamLeadPartAcc.Email,
		ParticipantId: teamLeadPartAcc.ParticipantId,
	}, nil, nil)
	if err != nil {
		panic(err)
	}
	service := NewService(&ServiceConfig{
		JudgeAccountRepository:       judgeAccountRepository,
		ParticipantAccountRepository: partAccRepo,
		Publisher:                    publisher,
		ParticipantRepository:        partRepo,
	})

	t.Run("Should create if already in invite list", func(t *testing.T) {
		fake := faker.New()
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
