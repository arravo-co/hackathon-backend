package services

import (
	"github.com/arravoco/hackathon_backend/di"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/publish"
	"github.com/go-playground/validator/v10"
)

/*
type ParticipantRepository interface {
	AddMemberToParticipatingTeam(dataToSave *exports.AddMemberToParticipatingTeamData) (*exports.ParticipantDocument, error)
}*/

type AccountRepository interface {
}

type SolutionRepository interface {
}

type Service struct {
	ParticipantRecordRepository  exports.ParticipantRepositoryInterface
	JudgeAccountRepository       exports.JudgeRepositoryInterface
	ParticipantAccountRepository exports.ParticipantAccountRepositoryInterface
	SolutionRepository           SolutionRepository
	Publisher                    publish.PublisherInterface
}

type ServiceConfig struct {
	ParticipantRepository        exports.ParticipantRepositoryInterface
	ParticipantAccountRepository exports.ParticipantAccountRepositoryInterface
	JudgeAccountRepository       exports.JudgeRepositoryInterface
	SolutionRepository           SolutionRepository
	Publisher                    publish.PublisherInterface
}

var service *Service

func NewService(cfg *ServiceConfig) *Service {
	validate = validator.New()
	return &Service{
		ParticipantRecordRepository:  cfg.ParticipantRepository,
		JudgeAccountRepository:       cfg.JudgeAccountRepository,
		SolutionRepository:           cfg.SolutionRepository,
		ParticipantAccountRepository: cfg.ParticipantAccountRepository,
		Publisher:                    cfg.Publisher,
	}
}

func GetServiceWithDefaultRepositories() *Service {
	if service != nil {
		return service
	}
	var judge exports.JudgeRepositoryInterface = di.GetDefaultJudgeRepository()
	var cfg *ServiceConfig = &ServiceConfig{
		JudgeAccountRepository: judge,
	}
	service = NewService(cfg)
	return service
}

var validate *validator.Validate
