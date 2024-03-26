package data

import (
	"context"
	"fmt"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateParticipantScore(dataToSave *exports.CreateParticipantScoreData) (*exports.ParticipantScoreDocument, error) {
	participantScoreCol, err := Datasource.GetScoreCollection()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	participantScoreDoc := &exports.ParticipantScoreDocument{
		HackathonId:   dataToSave.HackathonId,
		ParticipantId: dataToSave.ParticipantId,
		ScoresInfo:    []exports.ScoreInfo{},
	}
	ctx := context.Context(context.Background())
	result, err := participantScoreCol.InsertOne(ctx, &participantScoreDoc)
	if err != nil {
		return nil, err
	}
	participantScoreDoc.Id = result.InsertedID
	return participantScoreDoc, nil
}

func UpdateParticipantScoreRecordByJudge(updateFilter *exports.UpdateParticipantScoreFilterData, dataToSave *exports.UpdateParticipantScoreDataByJudge) (*exports.ParticipantScoreDocument, error) {
	participantScoreCol, err := Datasource.GetScoreCollection()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	participantScoreDoc := &exports.ParticipantScoreDocument{}
	ctx := context.Context(context.Background())
	filter := bson.M{
		"$and": updateFilter,
	}
	updateDoc := bson.M{
		"$set": bson.M{"scores_info.$[score_info]": dataToSave.ScoreInfo},
	}
	upsert := true
	retDoc := options.After
	result := participantScoreCol.FindOneAndUpdate(ctx, filter, updateDoc, &options.FindOneAndUpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"scores_info.judge_id": dataToSave.ScoreInfo.JudgeId},
			},
		},
		ReturnDocument: &retDoc,
		Upsert:         &upsert,
	})
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	err = result.Decode(&participantScoreDoc)
	if err != nil {
		return nil, err
	}
	return participantScoreDoc, nil
}

func GetParticipantScore(updateFilter *exports.UpdateParticipantScoreFilterData) (*exports.ParticipantScoreDocument, error) {
	participantScoreCol, err := Datasource.GetScoreCollection()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	participantScoreDoc := &exports.ParticipantScoreDocument{}
	ctx := context.Context(context.Background())
	result := participantScoreCol.FindOne(ctx, &participantScoreDoc)
	err = result.Decode(&participantScoreDoc)
	if err != nil {
		return nil, err
	}
	return participantScoreDoc, nil
}

func GetParticipantsScores(updateFilter *exports.UpdateParticipantScoreFilterData) (*[]exports.ParticipantScoreDocument, error) {
	participantScoreCol, err := Datasource.GetScoreCollection()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	participantsScoresDocs := []exports.ParticipantScoreDocument{}
	ctx := context.Context(context.Background())
	cursor, err := participantScoreCol.Find(ctx, &participantsScoresDocs)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &participantsScoresDocs)
	if err != nil {
		return nil, err
	}
	return &participantsScoresDocs, nil
}

func GetJudgeScoring(judgeId string) (*[]exports.ParticipantScoreDocument, error) {
	participantScoreCol, err := Datasource.GetScoreCollection()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	participantsScoresDocs := []exports.ParticipantScoreDocument{}
	ctx := context.Context(context.Background())
	cursor, err := participantScoreCol.Find(ctx, bson.M{"judge_id": judgeId})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &participantsScoresDocs)
	if err != nil {
		return nil, err
	}
	return &participantsScoresDocs, nil
}
