package exports

import "go.mongodb.org/mongo-driver/mongo"

type DBInterface interface {
	GetAccountCollection() (*mongo.Collection, error)
	GetParticipantCollection() (*mongo.Collection, error)
	GetTokenCollection() (*mongo.Collection, error)
	GetScoreCollection() (*mongo.Collection, error)
}
