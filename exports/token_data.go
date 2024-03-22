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
}

type CreateTokenData struct {
	Token          string    `bson:"token"`
	TokenType      string    `bson:"token_type"`
	TokenTypeValue string    `bson:"token_type_value"`
	Scope          string    `bson:"scope"`
	TTL            time.Time `bson:"ttl"`
	Status         string    `bson:"status"`
}

type VerifyTokenData struct {
	Token          string `bson:"token"`
	TokenType      string `bson:"token_type"`
	TokenTypeValue string `bson:"token_type_value"`
	Scope          string `bson:"scope"`
}
