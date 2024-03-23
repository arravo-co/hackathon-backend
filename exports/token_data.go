package exports

import "time"

type TokenData struct {
	Id             interface{} `bson:"_id"`
	Token          string      `bson:"token"`
	TokenType      string      `bson:"token_type"`
	TokenTypeValue string      `bson:"token_type_value"`
	Scope          string      `bson:"scope"`
	TTL            time.Time   `bson:"ttl"`
	Status         string      `bson:"status"`
	CreatedAt      time.Time   `bson:"created_at"`
	UpdatedAt      time.Time   `bson:"updated_at"`
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
