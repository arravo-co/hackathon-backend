package services

import (
	"errors"

	"github.com/arravoco/hackathon_backend/di"
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
	ParticipantRecordRepository  exports.ParticipantRepositoryInterface
	JudgeAccountRepository       exports.JudgeRepositoryInterface
	ParticipantAccountRepository exports.ParticipantAccountRepositoryInterface
	SolutionRepository           SolutionRepository
}

type ServiceConfig struct {
	ParticipantRepository  ParticipantRepository
	JudgeAccountRepository exports.JudgeRepositoryInterface
	SolutionRepository     SolutionRepository
}

var service *Service

func NewService(cfg *ServiceConfig) *Service {
	return &Service{
		ParticipantRecordRepository: cfg.ParticipantRepository,
		JudgeAccountRepository:      cfg.JudgeAccountRepository,
		SolutionRepository:          cfg.SolutionRepository,
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

func (s *Service) UpdateJudgeInfo(email string, dataInput *UpdateJudgeDTO) (*entity.Judge, error) {
	validate = validator.New()
	err := validate.Struct(dataInput)
	if err != nil {
		return nil, err
	}
	dataToSave := &exports.UpdateJudgeDTO{
		FirstName: dataInput.FirstName,
		LastName:  dataInput.LastName,
		Gender:    dataInput.Gender,
		Bio:       dataInput.Bio,
		State:     dataInput.State,
	}
	err = s.JudgeAccountRepository.UpdateJudgeAccount(email, dataToSave)
	if err != nil {
		return nil, err
	}

	judgeAccRepo, err := s.JudgeAccountRepository.GetJudgeAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	judge := entity.Judge{
		Id:                judgeAccRepo.Id,
		HackathonId:       judgeAccRepo.HackathonId,
		LastName:          judgeAccRepo.LastName,
		FirstName:         judgeAccRepo.FirstName,
		Email:             judgeAccRepo.Email,
		Gender:            judgeAccRepo.Gender,
		PhoneNumber:       judgeAccRepo.PhoneNumber,
		ProfilePictureUrl: judgeAccRepo.ProfilePictureUrl,
		Role:              judgeAccRepo.Role,
		Bio:               judgeAccRepo.Bio,
		Status:            judgeAccRepo.Status,
		CreatedAt:         judgeAccRepo.CreatedAt,
		UpdatedAt:         judgeAccRepo.UpdatedAt,
	}
	return &judge, nil
}

func (s *Service) GetJudgeByEmail(email string) (*entity.Judge, error) {
	if email == "" {
		return nil, errors.New("email must not be empty")
	}
	fetched, err := s.JudgeAccountRepository.GetJudgeAccountByEmail(email)
	if err != nil {
		return nil, err
	}

	judge := entity.Judge{
		Id:                fetched.Id,
		HackathonId:       fetched.HackathonId,
		LastName:          fetched.LastName,
		FirstName:         fetched.FirstName,
		Email:             fetched.Email,
		Gender:            fetched.Gender,
		PhoneNumber:       fetched.PhoneNumber,
		ProfilePictureUrl: fetched.ProfilePictureUrl,
		Role:              fetched.Role,
		Status:            fetched.Status,
		Bio:               fetched.Bio,
		IsEmailVerified:   fetched.IsEmailVerified,
		CreatedAt:         fetched.CreatedAt,
		UpdatedAt:         fetched.UpdatedAt,
	}
	return &judge, nil
}

func (s *Service) GetJudges() ([]*entity.Judge, error) {
	repoAccounts, err := s.JudgeAccountRepository.GetJudges()
	if err != nil {
		return nil, err
	}
	var ents []*entity.Judge
	for _, re := range repoAccounts {
		ents = append(ents, &entity.Judge{
			Id:                re.Id,
			LastName:          re.LastName,
			FirstName:         re.FirstName,
			Email:             re.Email,
			Bio:               re.Bio,
			Gender:            re.Gender,
			State:             re.State,
			Status:            re.Status,
			PhoneNumber:       re.PhoneNumber,
			ProfilePictureUrl: re.ProfilePictureUrl,
			HackathonId:       re.HackathonId,
			IsEmailVerified:   re.IsEmailVerified,
			CreatedAt:         re.CreatedAt,
			UpdatedAt:         re.UpdatedAt,
		})
	}
	return ents, nil
}
