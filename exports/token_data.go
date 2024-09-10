package exports

import (
	"time"
)

type TokenData struct {
	Id             string    `json:"_id,omitempty"`
	Token          string    `json:"token"`
	TokenType      string    `json:"token_type"`
	TokenTypeValue string    `json:"token_type_value"`
	Scope          string    `json:"scope"`
	TTL            time.Time `json:"ttl"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
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
