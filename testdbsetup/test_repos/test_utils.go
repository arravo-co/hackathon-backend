package testrepos

import (
	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/exports"
	"github.com/arravoco/hackathon_backend/repository"
)

func GetJudgeAccountRepositoryWithQueryInstance(q *query.Query) exports.JudgeRepositoryInterface {
	var judgeRepoInstance *repository.JudgeAccountRepository = repository.NewJudgeAccountRepository(q)
	var judgeAccountRepository exports.JudgeRepositoryInterface = judgeRepoInstance
	return judgeAccountRepository
}

func GetParticipantAccountRepositoryWithQueryInstance(q *query.Query) exports.ParticipantAccountRepositoryInterface {
	var partAccRepoInstance *repository.ParticipantAccountRepository = repository.NewParticipantAccountRepository(q)
	var judgeAccountRepository exports.ParticipantAccountRepositoryInterface = partAccRepoInstance
	return judgeAccountRepository
}

func GetParticipantRepositoryWithQueryInstance(q *query.Query) exports.ParticipantRepositoryInterface {
	var partRepoInstance *repository.ParticipantRecordRepository = repository.NewParticipantRecordRepository(q)
	var participantRepository exports.ParticipantRepositoryInterface = partRepoInstance
	return participantRepository
}
