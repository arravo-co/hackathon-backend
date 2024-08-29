package repository

import (
	"testing"
	"time"

	"github.com/arravoco/hackathon_backend/data/query"
	"github.com/arravoco/hackathon_backend/db"
	"github.com/arravoco/hackathon_backend/exports"
)

func TestCreateTeamMemberAccountMethod(t *testing.T) {

	var dtSrc exports.DBInterface
	mongoCfg := &exports.MongoDBConnConfig{}
	dtSrc, _ = db.GetNewMongoRepository(mongoCfg)
	var cfg *exports.ConfigQueryWithDatasource = &exports.ConfigQueryWithDatasource{
		Datasource: dtSrc,
	}
	var q *query.Query = query.GetQueryWithConfiguredDatasource(cfg)
	accRepo := NewAccountRepository(q)
	passwordHash, _ := exports.GenerateHashPassword("password")
	const shortForm = "2006-Jan-02"
	dob, _ := time.Parse(shortForm, "2000-Feb-02")
	_, _ = accRepo.CreateParticipantAccount(&exports.CreateParticipantAccountData{
		ParticipantId: "PARTICIPANT_001",
		Skillset:      []string{"nodejs", "javascript"},
		DOB:           dob,
		CreateAccountData: exports.CreateAccountData{
			Email:        "trinitietp@gmail.com",
			PasswordHash: passwordHash,
			PhoneNumber:  "+2347068968932",
			FirstName:    "temitope",
			LastName:     "alabi",
			Gender:       "MALE",
			State:        "LAGOS",
			Role:         "PARTICIPANT",
			Status:       "EMAIL_UNVERIFIED",
			HackathonId:  "HACKATHON_001",
		},
	})
}
