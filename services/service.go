package services

import (
	"context"
	"fmt"

	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
	"github.com/arravoco/hackathon_backend/resources"
	"github.com/go-playground/validator/v10"
	"github.com/rabbitmq/amqp091-go"
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
	TokenRepository              exports.TokenRepositoryInterface
	AppResources                 *resources.AppResources
	AdminAccountRepository       exports.AdminRepositoryInterface
}

type ServiceConfig struct {
	ParticipantRepository        exports.ParticipantRepositoryInterface
	ParticipantAccountRepository exports.ParticipantAccountRepositoryInterface
	JudgeAccountRepository       exports.JudgeRepositoryInterface
	SolutionRepository           SolutionRepository
	TokenRepository              exports.TokenRepositoryInterface
	AppResources                 *resources.AppResources
	AdminAccountRepository       exports.AdminRepositoryInterface
}

var defaultService *Service

func NewService(cfg *ServiceConfig) *Service {
	validate = validator.New()
	return &Service{
		ParticipantRecordRepository:  cfg.ParticipantRepository,
		JudgeAccountRepository:       cfg.JudgeAccountRepository,
		SolutionRepository:           cfg.SolutionRepository,
		ParticipantAccountRepository: cfg.ParticipantAccountRepository,
		AppResources:                 cfg.AppResources,
		TokenRepository:              cfg.TokenRepository,
		AdminAccountRepository:       cfg.AdminAccountRepository,
	}
}

func GetServiceWithDefaultRepositories() *Service {
	if defaultService != nil {
		return defaultService
	}
	res := resources.GetDefaultResources()
	var dataSourceInstance exports.DBInterface = data.GetDatasourceWithMongoDBInstance(res.Mongo)
	q := query.GetQueryWithConfiguredDatasource(dataSourceInstance)

	var judgeRepoInstance *repository.JudgeAccountRepository = repository.NewJudgeAccountRepository(q)
	var partAccRepoInstance *repository.ParticipantAccountRepository = repository.NewParticipantAccountRepository(q)
	var partRecordRepoInstance *repository.ParticipantRecordRepository = repository.NewParticipantRecordRepository(q)

	var solRepoInstance *repository.SolutionRepository = repository.NewSolutionRepository(q)

	var tokenRepoInstance *repository.TokenDataRepository = repository.NewTokenDataRepository(q)

	var adminAccRepoInstance *repository.AdminAccountRepository = repository.NewAdminAccountRepository(q)

	var cfg *ServiceConfig = &ServiceConfig{
		JudgeAccountRepository:       judgeRepoInstance,
		AppResources:                 res,
		TokenRepository:              tokenRepoInstance,
		ParticipantRepository:        partRecordRepoInstance,
		ParticipantAccountRepository: partAccRepoInstance,
		SolutionRepository:           solRepoInstance,
		AdminAccountRepository:       adminAccRepoInstance,
	}
	defaultService = NewService(cfg)
	return defaultService
}

type PublishOpts struct {
	RMQConn      *amqp091.Connection
	Data         []byte
	ExchangeName string
	KeyName      string
	ExchangeKind string
}

func Publish(opts *PublishOpts) error {
	rmqConn := opts.RMQConn
	by := opts.Data
	ch, err := rmqConn.Channel()
	if err != nil {
		fmt.Println(err)
		return err
	}
	exchange_name := opts.ExchangeName
	key_name := opts.KeyName
	err = ch.ExchangeDeclare(exchange_name, opts.ExchangeKind, true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ch.PublishWithContext(context.Background(),
		exchange_name, key_name, false, false, amqp091.Publishing{
			Body:        by,
			ContentType: "application/json",
		})
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	fmt.Printf("\nPublished details: exchange name: %s; key_name: %s\n", exchange_name, key_name)

	return nil
}

var validate *validator.Validate
