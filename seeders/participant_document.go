package seeders

import (
	"context"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OptsToCreateParticipantRecord struct {
	Status         string
	TeamleadInfo   TeamLeadInfoToCreateTeamParticipant
	CoParticipants []CoParticipantInfoToCreateTeamParticipant
	InviteList     []InvitelistQueuePayload
	SolutionId     string
}
type InvitelistQueuePayload struct {
	InviterId string
	Email     string
	Time      time.Time
}

type TeamLeadInfoToCreateTeamParticipant struct {
	Email         string
	HackathonId   string
	ParticipantId string
	TeamName      string
}

type CoParticipantInfoToCreateTeamParticipant struct {
	Email         string
	HackathonId   string
	ParticipantId string
}

type SolutionInfoToCreateParticipantDocument struct {
	Id               string
	HackathonId      string
	Title            string
	Description      string
	Objective        string
	SolutionImageUrl string
}

func CreateAccountLinkedTeamParticipantDocument(dbInstance *mongo.Database,
	opts *OptsToCreateParticipantRecord) (*exports.ParticipantDocument, error) {
	partCol := dbInstance.Collection("participants")
	ctx := context.Context(context.Background())

	fake := faker.New()

	team_name := fake.Lorem().Sentence(2)

	team_lead_info := opts.TeamleadInfo
	if team_lead_info.TeamName != "" {
		team_name = team_lead_info.TeamName
	}

	if team_lead_info.Email == "" {
		return nil, fmt.Errorf("team lead email is not set")
	}
	status := "UNREVIEWED"
	var co_parts_docs []exports.ParticipantDocumentTeamCoParticipantInfo
	var invite_list []exports.ParticipantDocumentTeamInviteInfo

	if opts.Status != "" {
		status = opts.Status
	}

	if len(opts.InviteList) > 0 {
		for _, v := range opts.InviteList {

			invite_list = append(invite_list, exports.ParticipantDocumentTeamInviteInfo{
				Email:     v.Email,
				InviterId: v.InviterId,
				Time:      v.Time,
			})
		}
	}

	if len(opts.CoParticipants) > 0 {
		for _, v := range opts.CoParticipants {

			co_parts_docs = append(co_parts_docs, exports.ParticipantDocumentTeamCoParticipantInfo{
				Email:         v.Email,
				ParticipantId: v.ParticipantId,
				HackathonId:   v.HackathonId,
				CreatedAt:     time.Now(),
				AddedToTeamAt: time.Now(),
			})
		}
	}

	var sol_id string = opts.SolutionId

	acc := &exports.ParticipantDocument{
		TeamLeadEmail:  team_lead_info.Email,
		TeamName:       team_name,
		HackathonId:    team_lead_info.HackathonId,
		Type:           "TEAM",
		ParticipantId:  team_lead_info.ParticipantId,
		Status:         status,
		CoParticipants: co_parts_docs,
		SolutionId:     sol_id,
		InviteList:     invite_list,
	}

	result, err := partCol.InsertOne(ctx, acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.InsertedID)
	acc.Id = result.InsertedID.(primitive.ObjectID)
	fmt.Printf("%#v", acc.CoParticipants)
	return acc, err
}
