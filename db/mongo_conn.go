package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client

func GetMongoConn() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetDefaultDB() (*mongo.Database, error) {
	db, err := GetDB("hackathons_db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDB(dbName string) (*mongo.Database, error) {
	client, err := GetMongoConn()
	if err != nil {
		return nil, err
	}
	db := client.Database(dbName)
	return db, nil
}

func GetAccountCollection() (*mongo.Collection, error) {
	db, err := GetDefaultDB()
	return db.Collection("accounts"), err
}
