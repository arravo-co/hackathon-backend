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
	Status     string
	InviteList []InvitelistQueuePayload
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
	opts *OptsToCreateParticipantRecord,
	team_lead_info TeamLeadInfoToCreateTeamParticipant,
	co_parts []CoParticipantInfoToCreateTeamParticipant,
	solOpts *SolutionInfoToCreateParticipantDocument) (*exports.ParticipantDocument, error) {
	accountCol := dbInstance.Collection("participants")
	ctx := context.Context(context.Background())

	fake := faker.New()

	team_name := fake.Lorem().Sentence(2)

	if team_lead_info.TeamName != "" {
		team_name = team_lead_info.TeamName
	}

	if team_lead_info.Email == "" {
		return nil, fmt.Errorf("team lead email is not set")
	}
	status := "UNREVIEWED"
	var invite_list []exports.ParticipantDocumentTeamInviteInfo
	if opts != nil {
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
	}
	var sol_id string
	sol := exports.ParticipantDocumentParticipantSelectedSolution{}
	if solOpts != nil {
		sol_id = solOpts.Id
		sol.Id = solOpts.Id
		sol.Title = solOpts.Title
		sol.Description = solOpts.Description
		sol.Objective = solOpts.Objective
		sol.SolutionImageUrl = solOpts.SolutionImageUrl
	}

	var co_parts_docs []exports.ParticipantDocumentTeamCoParticipantInfo
	if co_parts != nil {
		co_parts_docs = append(co_parts_docs, exports.ParticipantDocumentTeamCoParticipantInfo{})
	}

	acc := &exports.ParticipantDocument{
		TeamLeadEmail:  team_lead_info.Email,
		TeamName:       team_name,
		HackathonId:    team_lead_info.HackathonId,
		Type:           "TEAM",
		ParticipantId:  team_lead_info.ParticipantId,
		Status:         status,
		CoParticipants: co_parts_docs,
		SolutionId:     sol_id,
		Solution:       sol,
		InviteList:     invite_list,
	}

	result, err := accountCol.InsertOne(ctx, acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.InsertedID)
	acc.Id = result.InsertedID.(primitive.ObjectID)
	//fmt.Printf("%#v", acc)
	return acc, err
}
