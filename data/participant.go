package data

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//ParticipantDocumentTeamCoParticipantInfo

type ParticipantData struct {
	AccountCollection *mongo.Collection
}

func (p *ParticipantData) CreateParticipantRecord(dataToSave *exports.CreateParticipantRecordData) (*exports.ParticipantDocument, error) {
	participantCol, err := DefaultDatasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	dat := exports.ParticipantDocument{
		ParticipantId: dataToSave.ParticipantId,
		HackathonId:   dataToSave.HackathonId,
		Type:          dataToSave.Type,
		TeamLeadEmail: dataToSave.TeamLeadEmail,
		TeamName:      dataToSave.TeamName,
		//CoParticipants:   dataToSave.CoParticipants,
		GithubAddress:    dataToSave.GithubAddress,
		ParticipantEmail: dataToSave.ParticipantEmail,
		InviteList:       []exports.ParticipantDocumentTeamInviteInfo{},
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Status:           "UNREVIEWED",
	}
	result, err := participantCol.InsertOne(ctx, dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	dat.Id = result.InsertedID.(primitive.ObjectID)
	return &dat, nil
}

func (p *ParticipantData) GetParticipantsRecords() ([]exports.ParticipantDocument, error) {
	participantCol, err := DefaultDatasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	dat := &[]exports.ParticipantDocument{}
	result, err := participantCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		fmt.Printf("\n%s\n", result.Err())
		return nil, result.Err()
	}
	err = result.All(context.Background(), dat)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	return *dat, nil
}

func (p *ParticipantData) GetParticipantRecord(participantId string) (*exports.ParticipantDocument, error) {
	participantCol, err := DefaultDatasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	dat := exports.ParticipantDocument{}
	result := participantCol.FindOne(ctx, bson.M{"participant_id": participantId})
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

func (p *ParticipantData) AddToTeamInviteList(dataToSave *exports.AddToTeamInviteListData) (interface{}, error) {
	participantCol, err := DefaultDatasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"participant_id":    dataToSave.ParticipantId,
		"hackathon_id":      dataToSave.HackathonId,
		"invite_list.email": bson.M{"$nin": bson.A{dataToSave.Email}},
	}
	upd := bson.M{
		"$addToSet": bson.M{"invite_list": exports.ParticipantDocumentTeamInviteInfo{Email: dataToSave.Email,
			InviterId: dataToSave.InviterEmail, Time: time.Now()}},
		"$set": bson.M{"updated_at": time.Now()},
	}
	fmt.Println(upd)

	result, err := participantCol.UpdateOne(ctx, filter, upd)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result)
	if result.MatchedCount == 0 {
		fmt.Printf("failed to add to invite list")
		return nil, errors.New("failed to add to invite list: failed to match")
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		fmt.Printf("No changes made")
		return nil, errors.New("failed to add to invite list: failed to save")
	}
	return result, err
}

func (p *ParticipantData) AddMemberToParticipatingTeam(dataToSave *exports.AddMemberToParticipatingTeamData) (*exports.ParticipantDocument, error) {
	partDoc := &exports.ParticipantDocument{}
	participantCol, err := DefaultDatasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"participant_id":    dataToSave.ParticipantId,
		"hackathon_id":      dataToSave.HackathonId,
		"invite_list.email": dataToSave.Email,
	}

	fmt.Println("\n\n\n", filter, "\n\n\n")

	upd := bson.M{
		"$addToSet": bson.M{"co_participants": bson.M{"email": dataToSave.Email, "role": dataToSave.TeamRole}},
		"$pull":     bson.M{"invite_list": bson.M{"email": dataToSave.Email}},
		"$set":      bson.M{"updated_at": time.Now()},
	}
	retDoc := options.After
	result := participantCol.FindOneAndUpdate(ctx, filter, upd, &options.FindOneAndUpdateOptions{ReturnDocument: &retDoc})
	err = result.Decode(partDoc)
	if err != nil {
		fmt.Printf("%s\n\n\n\n", err.Error())
		return nil, err
	}
	return partDoc, err
}

func (p *ParticipantData) RemoveMemberFromParticipatingTeam(dataToSave *exports.RemoveMemberFromParticipatingTeamData) (interface{}, error) {
	participantCol, err := DefaultDatasource.GetParticipantCollection()
	ctx := context.Context(context.Background())
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"participant_id": dataToSave.ParticipantId,
		"hackathon_id":   dataToSave.HackathonId,
	}

	fmt.Println("\n\n", dataToSave.MemberEmail, "\n\n")

	upd := bson.M{
		"$pull": bson.M{"co_participants.email": dataToSave.MemberEmail},
	}
	result, err := participantCol.UpdateOne(ctx, filter, upd, &options.UpdateOptions{})
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		if strings.Contains(err.Error(), "mongo") {
			return nil, fmt.Errorf("Unexpected error")
		}
		return nil, err
	}
	fmt.Printf("%#v", result.ModifiedCount)
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no records found")
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no changes made.")
	}
	return participantCol, err
}
