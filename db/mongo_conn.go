package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func init() {
	client, err := GetMongoConn()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	MongoClient = client
}

func GetMongoConn() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(
		"mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, clientOpts)
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
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	db := client.Database(dbName)
	return db, nil
}

func GetAccountCollection() (*mongo.Collection, error) {
	db, err := GetDefaultDB()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	col := db.Collection("accounts")
	err = CreateIndexes()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	return col, nil
}

func CreateIndexes() error {
	db, err := GetDefaultDB()
	if err != nil {
		return err
	}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", -1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err := db.Collection("accounts").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	fmt.Printf("%s", indexName)
	return nil
}
