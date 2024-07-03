package db

import (
	"context"
	"fmt"
	"time"

	"github.com/arravoco/hackathon_backend/config"
	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

type MongoRepository struct {
	DB *mongo.Database
}

func GetMongoRepository(dbname string) (*MongoRepository, error) {
	db, err := GetDB(dbname)
	if err != nil {
		return nil, err
	}
	return &MongoRepository{
		DB: db,
	}, nil
}

func GetMongoRepositoryWithDB(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		DB: db,
	}
}

func GetNewMongoRepository(conf *exports.MongoDBConnConfig) (*MongoRepository, error) {
	var url string = ""
	var dbName string = ""
	if conf.Url != "" {
		url = conf.Url
	}
	if url == "" {
		url = config.GetMongoDBURL()
	}
	if conf.DBName != "" {
		dbName = conf.DBName
	}
	if dbName == "" {
		dbName = "hackathons_db"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(url)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	fmt.Println("\n\n\nConnected to database\n\n\n")
	db := client.Database(dbName)
	return &MongoRepository{
		DB: db,
	}, nil
}

func GetNewDefaultMongoRepository() (*MongoRepository, error) {
	return GetNewMongoRepository(&exports.MongoDBConnConfig{
		Url:    config.GetMongoDBURL(),
		DBName: "hackathons_db",
	})
}

func GetMongoConn(configs ...*exports.MongoDBConnConfig) (*mongo.Client, error) {
	url := config.GetMongoDBURL()
	for _, config := range configs {
		if config.Url != "" {
			url = config.Url
		}
	}
	if MongoClient != nil {
		return MongoClient, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(url)

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

func GetCollection(colName string) (*mongo.Collection, error) {
	db, err := GetDefaultDB()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	col := db.Collection(colName)
	err = CreateIndexes(db)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	return col, nil
}

func (m MongoRepository) GetView() (*mongo.Collection, error) {
	return nil, nil
}

func (m MongoRepository) GetAccountCollection() (*mongo.Collection, error) {
	var err error
	col := m.DB.Collection("accounts")
	err = CreateIndexes(col.Database())
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		return nil, err
	}
	return col, nil
}

func (m MongoRepository) GetParticipantCollection() (*mongo.Collection, error) {
	var err error
	col := m.DB.Collection("participants")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"participant_id", -1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err := m.DB.Collection("participants").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n%s\n", indexName)

	return col, err
}

func (m MongoRepository) GetTokenCollection() (*mongo.Collection, error) {
	col, err := GetCollection("tokens")
	return col, err
}

func (m MongoRepository) GetSolutionCollection() (*mongo.Collection, error) {
	col := m.DB.Collection("solutions")
	err := CreateIndexes(m.DB)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, err
	}
	return col, nil
}
func (m MongoRepository) GetScoreCollection() (*mongo.Collection, error) {
	col, err := GetCollection("scores")
	return col, err
}

func CreateParticipantColIndexes() error {
	db, err := GetDefaultDB()
	if err != nil {
		return err
	}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"participant_id", -1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err := db.Collection("participants").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	fmt.Printf("\n%s\n", indexName)
	return nil
}

func CreateIndexes(db *mongo.Database) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", -1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err := db.Collection("accounts").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{"token", -1}, {"token_type", 1}, {"token_type_value", 1}},
		Options: options.Index().SetUnique(true),
	}
	indexName, err = db.Collection("tokens").Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	fmt.Printf("%s", indexName)
	return nil
}
