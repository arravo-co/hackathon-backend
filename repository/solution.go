package repository

import (
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
)

type SolutionRepository struct {
	DB *query.Query
}

func NewSolutionRepository(q *query.Query) *SolutionRepository {

	return &SolutionRepository{
		DB: q,
	}
}

func (s *SolutionRepository) CreateSolution(creator_id string, dataInput exports.CreateSolutionData) (*SolutionRepository, error) {
	return nil, nil
}

func (*SolutionRepository) UpdateSolution(creator_id string, dataInput exports.CreateSolutionData) (*SolutionRepository, error) {
	return nil, nil
}

func (*SolutionRepository) GetSolutions(creator_id string, dataInput exports.CreateSolutionData) ([]*SolutionRepository, error) {
	return nil, nil
}

func (*SolutionRepository) GetSolutionDataById(id string) (*entity.Solution, error) {
	return nil, nil
}
