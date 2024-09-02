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

func (s *SolutionRepository) CreateSolution(dataInput *exports.CreateSolutionData) (*entity.Solution, error) {
	solDoc, err := s.DB.CreateSolutionData(dataInput)
	if err != nil {
		return nil, err
	}

	return &entity.Solution{
		Id:               solDoc.Id.Hex(),
		Title:            solDoc.Title,
		Description:      solDoc.Description,
		HackathonId:      solDoc.HackathonId,
		Objective:        solDoc.Objective,
		SolutionImageUrl: solDoc.SolutionImageUrl,
		CreatorId:        solDoc.CreatorId,
		CreatedAt:        solDoc.CreatedAt,
		UpdatedAt:        solDoc.UpdatedAt,
	}, nil
}

func (s *SolutionRepository) UpdateSolution(creator_id string, dataInput *exports.UpdateSolutionData) (*entity.Solution, error) {
	sol, err := s.DB.UpdateSolutionData(creator_id, dataInput)
	if err != nil {
		return nil, err
	}

	return &entity.Solution{
		Id:               sol.Id.Hex(),
		Title:            sol.Title,
		Description:      sol.Description,
		Objective:        sol.Objective,
		HackathonId:      sol.HackathonId,
		CreatorId:        sol.CreatorId,
		SolutionImageUrl: sol.SolutionImageUrl,
		CreatedAt:        sol.CreatedAt,
		UpdatedAt:        sol.UpdatedAt,
	}, nil
}

func (s *SolutionRepository) GetSolutionsData(dataInput *exports.GetSolutionsQueryData) ([]entity.Solution, error) {
	solDocs, err := s.DB.GetManySolutionData(dataInput)
	if err != nil {
		return nil, err
	}
	var sols []entity.Solution

	for _, sol := range solDocs {
		sols = append(sols, entity.Solution{
			Id:               sol.Id.Hex(),
			Title:            sol.Title,
			Description:      sol.Description,
			Objective:        sol.Objective,
			HackathonId:      sol.HackathonId,
			CreatorId:        sol.CreatorId,
			SolutionImageUrl: sol.SolutionImageUrl,
			CreatedAt:        sol.CreatedAt,
			UpdatedAt:        sol.UpdatedAt,
		})
	}

	return sols, nil
}

func (s *SolutionRepository) GetSolutionDataById(id string) (*entity.Solution, error) {
	sol, err := s.DB.GetSolutionDataById(id)
	if err != nil {
		return nil, err
	}

	return &entity.Solution{
		Id:               sol.Id.Hex(),
		Title:            sol.Title,
		Description:      sol.Description,
		HackathonId:      sol.HackathonId,
		Objective:        sol.Objective,
		CreatorId:        sol.CreatorId,
		SolutionImageUrl: sol.SolutionImageUrl,
		CreatedAt:        sol.CreatedAt,
		UpdatedAt:        sol.UpdatedAt,
	}, nil
}

type FilterGetParticipants struct {
	Status            *string `validate:"omitempty, oneof UNREVIEWED REVIEWED AI_RANKED "`
	ReviewRanking_Eq  *int
	ReviewRanking_Top *int
	Solution_Like     *string
}
