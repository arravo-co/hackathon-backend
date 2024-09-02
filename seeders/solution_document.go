package seeders

import (
	"context"
	"fmt"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/jaswdr/faker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OptsToCreateSolutionRecord struct {
	Title       string
	Objective   string
	Description string
}

func CreateFakeSolutionDocument(dbInstance *mongo.Database, opts OptsToCreateSolutionRecord) (*exports.SolutionDocument, error) {
	solCol := dbInstance.Collection("solutions")
	ctx := context.Context(context.Background())

	fake := faker.New()

	title := opts.Title
	obj := fake.Lorem().Sentence(2)
	desc := fake.Lorem().Sentence(50)

	if title != opts.Title {
		title = fake.Lorem().Sentence(2)
	}

	if obj != opts.Objective {
		obj = fake.Lorem().Sentence(2)
	}

	if desc != opts.Description {
		obj = fake.Lorem().Sentence(2)
	}

	acc := &exports.SolutionDocument{
		Title:            title,
		Objective:        obj,
		Description:      desc,
		SolutionImageUrl: "",
	}

	result, err := solCol.InsertOne(ctx, acc)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return nil, err
	}
	fmt.Printf("%#v", result.InsertedID.(primitive.ObjectID).Hex())
	acc.Id = result.InsertedID.(primitive.ObjectID)
	return acc, err
}
