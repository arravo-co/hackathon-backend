package repository

import (
	"fmt"

	"github.com/arravoco/hackathon_backend/exports"
)

type TokenDataRepository struct {
	DB exports.TokenDatasourceQueryMethods
}

func NewTokenDataRepository(datasource exports.TokenDatasourceQueryMethods) *TokenDataRepository {
	return &TokenDataRepository{
		DB: datasource,
	}
}

func (t *TokenDataRepository) UpsertToken(dataInput *exports.UpsertTokenData) (*exports.TokenDataRepository, error) {
	tokenDoc, err := t.DB.UpsertToken(dataInput)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\n\n%#v\n\n\n", tokenDoc)

	return &exports.TokenDataRepository{
		Id:             tokenDoc.Id.Hex(),
		Token:          tokenDoc.Token,
		TokenType:      tokenDoc.TokenType,
		TokenTypeValue: tokenDoc.TokenTypeValue,
		Scope:          tokenDoc.Scope,
		TTL:            tokenDoc.TTL,
		Status:         tokenDoc.Status,
		CreatedAt:      tokenDoc.CreatedAt,
		UpdatedAt:      tokenDoc.UpdatedAt,
	}, nil
}

func (t *TokenDataRepository) VerifyToken(dataInput *exports.VerifyTokenData) error {
	err := t.DB.VerifyToken(dataInput)
	return err
}
