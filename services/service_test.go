package services

import (
	"errors"
	"testing"

	"github.com/arravoco/hackathon_backend/data"
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
)

func TestRegisterJudge(t *testing.T) {
	var mongoInstance exports.DBInterface
	mongoInstance = data.GetDatasource(&exports.MongoDBConnConfig{})
	query.GetQueryWithConfiguredDatasource(mongoInstance)
	/*q := exports.QueryWithDatasource{
		Datasource:   mongoInstance,
		QueryMethods: mongoInstance,
	}*/
	var judgeRepoInstance *repository.JudgeAccountRepository = repository.NewJudgeAccountRepository(mongoInstance)
	var judgeAccountRepository exports.JudgeRepositoryInterface = judgeRepoInstance
	service := NewService(&ServiceConfig{JudgeAccountRepository: judgeAccountRepository})
	dataInput := &exports.RegisterNewJudgeDTO{}
	judgeRepo, err := service.JudgeAccountRepository.CreateJudgeAccount(dataInput)
	if err != nil {
		t.Fatal(err)
	}
	if judgeRepo.Id == "" {
		t.Fatal(errors.New("id not returned"))
	}
}
