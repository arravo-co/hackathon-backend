package query

import (
	"context"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (q *Query) CreateSolutionData(dataInput *exports.CreateSolutionData) (*exports.SolutionDocument, error) {
	solutionCol, err := q.Datasource.GetSolutionCollection()
	if err != nil {
		return nil, err
	}
	dataFromCol := exports.SolutionDocument{
		Title:       dataInput.Title,
		Description: dataInput.Description,
		HackathonId: dataInput.HackathonId,
		Objective:   dataInput.Objective,
		CreatorId:   dataInput.CreatorId,
		CreatedAt:   time.Now(),
	}
	result, err := solutionCol.InsertOne(context.Background(), dataFromCol)
	dataFromCol.Id = result.InsertedID
	return &dataFromCol, err
}

func (q *Query) UpdateSolutionData(id string, dataInput *exports.UpdateSolutionData) (*exports.SolutionDocument, error) {

	solutionCol, err := q.Datasource.GetSolutionCollection()
	if err != nil {
		return nil, err
	}
	dataFromCol := &exports.SolutionDocument{
		Description: dataInput.Description,
		Title:       dataInput.Title,
		Objective:   dataInput.Objective,
		UpdatedAt:   time.Now(),
	}
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": objId,
	}
	var updates *bson.M

	if dataInput.Description != "" {
		if updates == nil {
			updates = &bson.M{
				"description": dataInput.Description,
			}
		}
	}

	if dataInput.Title != "" {
		if updates == nil {
			updates = &bson.M{
				"title": dataInput.Title,
			}
		}
	}

	if dataInput.Objective != "" {
		if updates == nil {
			updates = &bson.M{
				"objective": dataInput.Objective,
			}
		}
	}
	if updates == nil {
		updates = &bson.M{
			"updated_at": time.Now(),
		}
	}
	after := options.After

	result := solutionCol.FindOneAndUpdate(context.Background(), filter, bson.M{"$set": dataFromCol}, &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	})
	err = result.Decode(dataFromCol)
	if err != nil {
		return nil, err
	}
	return dataFromCol, err
}

func (q *Query) GetSolutionDataById(id string) (*exports.SolutionDocument, error) {
	solutionCol, err := q.Datasource.GetSolutionCollection()
	if err != nil {
		return nil, err
	}
	dataFromCol := &exports.SolutionDocument{}
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id": objId,
	}
	result := solutionCol.FindOne(context.Background(), filter)
	err = result.Decode(dataFromCol)
	if err != nil {
		return nil, err
	}
	return dataFromCol, err
}

func (q *Query) GetManySolutionData(filterInput interface{}) ([]exports.SolutionDocument, error) {
	solutionCol, err := q.Datasource.GetSolutionCollection()
	if err != nil {
		return nil, err
	}
	var dataFromCol []exports.SolutionDocument
	filter := bson.M{}
	result, err := solutionCol.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	err = result.All(context.Background(), &dataFromCol)
	fmt.Println(dataFromCol)
	return dataFromCol, err
}
