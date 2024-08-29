package services

import (
	"github.com/arravoco/hackathon_backend/entity"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/go-playground/validator/v10"
)

type ParticipantRepository interface {
	AddMemberToParticipatingTeam(dataToSave *exports.AddMemberToParticipatingTeamData) (*exports.ParticipantDocument, error)
}

type AccountRepository interface {
}

type SolutionRepository interface {
}

type Service struct {
	ParticipantRepository  ParticipantRepository
	JudgeAccountRepository exports.JudgeRepositoryInterface
	SolutionRepository     SolutionRepository
}

type ServiceConfig struct {
	ParticipantRepository  ParticipantRepository
	JudgeAccountRepository exports.JudgeRepositoryInterface
	SolutionRepository     SolutionRepository
}

func NewService(cfg *ServiceConfig) *Service {
	return &Service{
		ParticipantRepository:  cfg.ParticipantRepository,
		JudgeAccountRepository: cfg.JudgeAccountRepository,
		SolutionRepository:     cfg.SolutionRepository,
	}
}

var validate *validator.Validate

func (s *Service) RegisterNewJudge(dataInput *RegisterNewJudgeDTO) (*entity.Judge, error) {
	validate = validator.New()
	err := validate.Struct(dataInput)
	if err != nil {
		return nil, err
	}
	dataToSave := &exports.RegisterNewJudgeDTO{
		FirstName: dataInput.FirstName,
		LastName:  dataInput.LastName,
		Email:     dataInput.Email,
		Gender:    dataInput.Gender,
		Password:  dataInput.Password,
		Bio:       dataInput.Bio,
		State:     dataInput.State,
	}
	created, err := s.JudgeAccountRepository.CreateJudgeAccount(dataToSave)
	if err != nil {
		return nil, err
	}

	judge := entity.Judge{
		Id:                created.Id,
		HackathonId:       created.HackathonId,
		LastName:          created.LastName,
		FirstName:         created.FirstName,
		Email:             created.Email,
		Gender:            created.Gender,
		PhoneNumber:       created.PhoneNumber,
		ProfilePictureUrl: created.ProfilePictureUrl,
		Role:              created.Role,
		Status:            created.Status,
		CreatedAt:         created.CreatedAt,
		UpdatedAt:         created.UpdatedAt,
	}
	return &judge, nil
}
