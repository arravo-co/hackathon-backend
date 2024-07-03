package query

import (
	"testing"

	"github.com/arravoco/hackathon_backend/exports"
	"github.com/stretchr/testify/assert"
)

func TestGetParticipantsInfo(t *testing.T) {
	q := SetupDB()
	_, err := q.CreateParticipantAccount(&exports.CreateParticipantAccountData{
		CreateAccountData: exports.CreateAccountData{
			Email:       "test@test.com",
			FirstName:   "John",
			LastName:    "Doe",
			Gender:      "MALE",
			Role:        "PARTICIPANT",
			HackathonId: "hackathon001",
		},
		Skillset:        []string{"nodejs"},
		ExperienceLevel: "SENIOR",
		Motivation:      "Locum ipsum",
		ParticipantId:   "participant001",
	})
	assert.NoError(t, err)
	_, err = q.CreateParticipantRecord(&exports.CreateParticipantRecordData{
		ParticipantId:    "participant001",
		TeamLeadEmail:    "test@test.com",
		TeamName:         "Team",
		HackathonId:      "hackathon001",
		ParticipantEmail: "test@test.com",
		Type:             "TEAM",
		CoParticipants: []exports.CoParticipant{
			{Email: "test2@test.com", Role: "TEAM_MEMBER"},
		},
	})
	assert.NoError(t, err)
	docs, err := q.GetParticipantsRecordsAggregate()
	assert.NoError(t, err)

	assert.IsType(t, []exports.ParticipantAccountWithCoParticipantsDocument{}, docs)
}
