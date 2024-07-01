package repository

import (
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SolutionRepository struct {
	DB *query.Query
}

func NewSolutionRepository(q *query.Query) *SolutionRepository {

	return &SolutionRepository{
		DB: q,
	}
}

func (s *SolutionRepository) CreateSolution(dataInput *exports.CreateSolutionData) (*entity.Solution, error) {
	solDoc, err := s.DB.CreateSolutionData(dataInput)
	if err != nil {
		return nil, err
	}

	return &entity.Solution{
		Id:          solDoc.Id.(primitive.ObjectID).Hex(),
		Title:       solDoc.Title,
		Description: solDoc.Description,
		HackathonId: solDoc.HackathonId,
		CreatorId:   solDoc.CreatorId,
		CreatedAt:   solDoc.CreatedAt,
		UpdatedAt:   solDoc.UpdatedAt,
	}, nil
}

func (*SolutionRepository) UpdateSolution(creator_id string, dataInput exports.CreateSolutionData) (*SolutionRepository, error) {
	return nil, nil
}

func (*SolutionRepository) GetSolutions(creator_id string, dataInput exports.CreateSolutionData) ([]*SolutionRepository, error) {

	return nil, nil
}

func (s *SolutionRepository) GetSolutionDataById(id string) (*entity.Solution, error) {
	sol, err := s.DB.GetSolutionDataById(id)
	if err != nil {
		return nil, err
	}

	return &entity.Solution{
		Id:          sol.Id.(primitive.ObjectID).Hex(),
		Title:       sol.Title,
		Description: sol.Description,
		HackathonId: sol.HackathonId,
		CreatorId:   sol.CreatorId,
		CreatedAt:   sol.CreatedAt,
		UpdatedAt:   sol.UpdatedAt,
	}, nil
}
