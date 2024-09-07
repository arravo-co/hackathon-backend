package exports

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenData struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	Token          string             `bson:"token"`
	TokenType      string             `bson:"token_type"`
	TokenTypeValue string             `bson:"token_type_value"`
	Scope          string             `bson:"scope"`
	TTL            time.Time          `bson:"ttl"`
	Status         string             `bson:"status"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty"`
}

type UpsertTokenData struct {
	Token          string    `bson:"token"`
	TokenType      string    `bson:"token_type"`
	TokenTypeValue string    `bson:"token_type_value"`
	Scope          string    `bson:"scope"`
	TTL            time.Time `bson:"ttl"`
	Status         string    `bson:"status"`
	UpdatedAt      time.Time `bson:"updated_at"`
}

type VerifyTokenData struct {
	Token          string `bson:"token"`
	TokenType      string `bson:"token_type"`
	TokenTypeValue string `bson:"token_type_value"`
	Scope          string `bson:"scope"`
}
