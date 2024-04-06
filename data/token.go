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

func UpsertToken(dataInput *exports.UpsertTokenData) (*exports.TokenData, error) {
	fmt.Printf("\nuuuuuuuuuuuuuuuuuuuuuuuuuu\n%+v\npppppppppppppppppppppppppp\n", dataInput)
	tokenCol, err := Datasource.GetTokenCollection()
	if err != nil {
		return nil, err
	}
	tokenInfo := &exports.TokenData{}
	filter := bson.M{
		"token_type":       dataInput.TokenType,
		"token_type_value": dataInput.TokenTypeValue,
	}
	dataInput.UpdatedAt = time.Now()
	var upsert bool = true
	returnDoc := options.After
	updateDoc := bson.M{"$set": dataInput, "$setOnInsert": bson.M{"created_at": time.Now()}}
	err = tokenCol.FindOneAndUpdate(context.TODO(), filter, updateDoc, &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnDoc,
	}).Decode(tokenInfo)
	return tokenInfo, err
}

func VerifyToken(dataInput *exports.VerifyTokenData) error {
	fmt.Printf("\n%+v\n", dataInput)
	var tokenInfo exports.TokenData
	tokenCol, err := Datasource.GetTokenCollection()
	if err != nil {
		return err
	}
	filter := bson.M{
		"token":            dataInput.Token,
		"token_type":       dataInput.TokenType,
		"token_type_value": dataInput.TokenTypeValue,
	}
	result := tokenCol.FindOneAndUpdate(context.TODO(), filter,
		bson.A{bson.D{{
			Key: "$set", Value: bson.D{{
				Key: "status", Value: bson.D{{
					Key: "$switch", Value: bson.D{{
						Key: "branches", Value: bson.A{
							bson.D{{
								Key: "case", Value: bson.M{
									"$and": bson.A{
										bson.M{
											"$lt": bson.A{"$ttl", "$$NOW"},
										}, bson.M{
											"$or": bson.A{
												bson.M{
													"$eq": bson.A{"$status", "PENDING"}}, bson.M{"$eq": bson.A{"$status", "UNVERIFIED"}},
											},
										},
									},
								}}, {
								Key: "then", Value: "EXPIRED",
							},
								{
									Key: "case", Value: bson.M{
										"$and": bson.A{bson.M{
											"$gte": bson.A{"$ttl", "$$NOW"},
										}, bson.M{
											"$or": bson.A{
												bson.M{
													"$eq": bson.A{"$status", "PENDING"}}, bson.M{"$eq": bson.A{"$status", "UNVERIFIED"}},
											},
										},
										},
									}}, {
									Key: "then", Value: "VERIFIED",
								},
							},
						}}, {
						Key: "default", Value: bson.M{
							"$getField": "status",
						},
					}},
				},
				},
			}}}}},
	) //bson.M{"status": "VERIFIED"}

	if result.Err() != nil {
		exports.MySugarLogger.Error(result.Err())
		return errors.New("unable to verify token")
	}
	err = result.Decode(tokenInfo)
	if result.Err() != nil {
		exports.MySugarLogger.Error(result.Err())
		return errors.New("Error verifying token")
	}
	if tokenInfo.Status == "EXPIRED" {
		return errors.New("token has expired")
	}

	return nil
}
