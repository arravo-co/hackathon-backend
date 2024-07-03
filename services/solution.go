package services

import (
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
)

type SolutionService struct {
	ParticipantRepository *repository.ParticipantRepository
	AccountRepository     *repository.AccountRepository
	SolutionRepository    *repository.SolutionRepository
}

func NewSolutionService() *SolutionService {

	q := query.GetDefaultQuery()
	part := repository.NewParticipantRepository(q)
	acc := repository.NewAccountRepository(q)
	sol := repository.NewSolutionRepository(q)
	return &SolutionService{
		ParticipantRepository: part,
		AccountRepository:     acc,
		SolutionRepository:    sol,
	}
}
func (s *SolutionService) CreateSolution(dataInput *exports.CreateSolutionData) (*entity.Solution, error) {
	return s.SolutionRepository.CreateSolution(dataInput)
}

func (s *SolutionService) GetSolutionDataById(id string) (*entity.Solution, error) {
	return s.SolutionRepository.GetSolutionDataById(id)
}

func (s *SolutionService) UpdateSolutionDataById(id string, updates *exports.UpdateSolutionData) (*entity.Solution, error) {
	return s.SolutionRepository.UpdateSolution(id, updates)
}

func (s *SolutionService) GetSolutionData(dataInput *exports.GetSolutionsQueryData) ([]entity.Solution, error) {
	return s.SolutionRepository.GetSolutionsData(dataInput)
}
