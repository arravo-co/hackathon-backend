package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateParticipantRecord(dataToSave *exports.CreateParticipantRecordData) (*exports.ParticipantDocument, error) {
	participantCol, err := Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	dat := exports.ParticipantDocument{
		ParticipantId:       dataToSave.ParticipantId,
		HackathonId:         dataToSave.HackathonId,
		Type:                dataToSave.Type,
		TeamLeadEmail:       dataToSave.TeamLeadEmail,
		TeamName:            dataToSave.TeamName,
		CoParticipantEmails: dataToSave.CoParticipantEmails,
		GithubAddress:       dataToSave.GithubAddress,
		ParticipantEmail:    dataToSave.ParticipantEmail,
	}
	result, err := participantCol.InsertOne(ctx, dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	dat.Id = result.InsertedID
	return &dat, nil
}

func GetParticipantRecord(dataToSave string) (*exports.ParticipantDocument, error) {
	participantCol, err := Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	dat := exports.ParticipantDocument{}
	result := participantCol.FindOne(ctx, bson.M{"participant_id": dataToSave})
	if result.Err() != nil {
		fmt.Printf("\n%s\n", result.Err())
		return nil, result.Err()
	}
	err = result.Decode(&dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	return &dat, nil
}

func AddToTeamInviteList(dataToSave *exports.AddToTeamInviteListData) (interface{}, error) {
	participantCol, err := Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Account\n")
	fmt.Printf("%#v\n", dataToSave)
	filter := bson.M{
		"participant_id":    dataToSave.ParticipantId,
		"hackathon_id":      dataToSave.HackathonId,
		"invite_list.email": bson.M{"$nin": bson.A{dataToSave.Email}},
	}
	upd := bson.M{
		"$addToSet": bson.M{"invite_list": exports.InviteInfo{Email: dataToSave.Email,
			InviterId: dataToSave.InviterEmail, Time: time.Now()}},
	}
	fmt.Println(upd)
	/*

			"": bson.M{
				"$cond": bson.M{
				"if": bson.M{"$isArray":"$invite_list"},

				},},

		upd1 := bson.M{
			"$set": bson.M{
				"invite_list": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$in": bson.A{dataToSave.Email, "$invite_list.email"}},
						"then": bson.M{"$getField": "invite_list"},
						"else": bson.M{
							"$concatArrays": bson.A{
								bson.M{"$getField": "invite_list"},
								exports.InviteInfo{Email: dataToSave.Email,
									InviterId: dataToSave.InviterEmail, Time: time.Now()}}},
					},
				},
			},
		}
	*/
	result, err := participantCol.UpdateOne(ctx, filter, upd)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result)
	if result.MatchedCount == 0 {
		fmt.Printf("failed to add to invite list")
		return nil, errors.New("failed to add to invite list")
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		fmt.Printf("No changes made")
		return nil, errors.New("failed to add to invite list")
	}
	return result, err
}

func AddMemberToParticipatingTeam(dataToSave *exports.AddMemberToParticipatingTeamData) (interface{}, error) {
	participantCol, err := Datasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Printf("Account\n")
	fmt.Printf("%#v\n", dataToSave)
	filter := bson.M{
		"participant_id":    dataToSave.ParticipantId,
		"hackathon_id":      dataToSave.HackathonId,
		"invite_list.email": dataToSave.Email,
	}
	upd := bson.M{
		"$addToSet": bson.M{"co_participant_emails": dataToSave.Email},
		"$pull":     bson.M{"invite_list": dataToSave.Email},
	}
	result, err := participantCol.UpdateOne(ctx, filter, upd, &options.UpdateOptions{})
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.ModifiedCount)
	return participantCol, err
}
